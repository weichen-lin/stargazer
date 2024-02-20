package consumer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/IBM/sarama"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	database "github.com/weichen-lin/kafka-service/db"
	"gorm.io/gorm"
)

func GetUserStarredRepos(info *database.GetGithubReposInfo) ([]database.Repository, error) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("GITHUB_TOKEN not set")
	}

	url := fmt.Sprintf("https://api.github.com/users/%s/starred?&page=%d", info.Username, info.Page)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "token "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error code at get user stars: %d", resp.StatusCode)
	}

	var repos []database.Repository

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &repos)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON: %s", err)
	}

	return repos, nil
}

func GetGithubReposConsumer() (func(neo4j.DriverWithContext, *gorm.DB), error) {
	brokers := []string{"localhost:9092"}
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Producer.Return.Successes = true

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		fmt.Println("Error creating consumer:", err)
	}

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		fmt.Println("Error creating producer:", err)
		return nil, err
	}

	consumerPartitionConsumer, err := consumer.ConsumePartition("get_user_stars", 0, sarama.OffsetNewest)
	if err != nil {
		fmt.Println("Error creating partition consumer:", err)
	}

	return func(driver neo4j.DriverWithContext, pool *gorm.DB) {
		for {
			select {
			case err := <-consumerPartitionConsumer.Errors():
				fmt.Println("Error:", err)
			case message := <-consumerPartitionConsumer.Messages():
				fmt.Printf("Received message: Topic - %s, Key - %s, Value - %s\n",
					message.Topic, message.Key, message.Value)

				var info database.GetGithubReposInfo

				err = json.Unmarshal(message.Value, &info)
				if err != nil {
					fmt.Println("Error parsing JSON:", err)
				}

				stars, err := GetUserStarredRepos(&info)
				if err != nil {
					fmt.Println("Error getting user stars:", err)
				}

				if len(stars) == 30 {
					info.Page++

					jsonString, err := json.Marshal(info)
					if err != nil {
						fmt.Println("Error marshalling JSON:", err)
					}

					partition, offset, err := producer.SendMessage(&sarama.ProducerMessage{
						Topic: "get_user_stars",
						Value: sarama.StringEncoder(jsonString),
					})
					if err != nil {
						fmt.Println("Error sending message:", err)
					}

					fmt.Printf("Sent message: Topic - %s, Partition - %d, Offset - %d\n",
						"get_user_stars", partition, offset)
				}

				for _, repo := range stars {
					err = database.CreateRepository(driver, &repo, info.UserId, pool)
				}
			}
		}
	}, nil
}
