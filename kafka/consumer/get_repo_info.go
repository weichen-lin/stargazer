package consumer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/IBM/sarama"
	database "github.com/weichen-lin/kafka-service/db"
	"gorm.io/gorm"
)

func GetRepoInfo(info *database.AddRepoInfo) error {

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return fmt.Errorf("GITHUB_TOKEN not set")
	}

	req, err := http.NewRequest("GET", "https://raw.githubusercontent.com/"+info.Name+"/"+info.DefaultBranch+"/"+"README.md", nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "token "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error at get repo info: %d", resp.StatusCode)
	}

	repoInfo, err := io.ReadAll(resp.Body)

	info.RepoInfo = repoInfo
	if err != nil {
		return err
	}

	return nil
}

func GetGithubRepoInfoConsumer() (func(*gorm.DB), error) {
	brokers := []string{"localhost:9092"}
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Producer.Return.Successes = true

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		fmt.Println("Error creating consumer:", err)
	}

	consumerPartitionConsumer, err := consumer.ConsumePartition("get_star_info", 0, sarama.OffsetNewest)
	if err != nil {
		fmt.Println("Error creating partition consumer:", err)
	}

	return func(pool *gorm.DB) {
		for {
			select {
			case err := <-consumerPartitionConsumer.Errors():
				fmt.Println("Error:", err.Err)
			case message := <-consumerPartitionConsumer.Messages():
				fmt.Printf("Received message: Topic - %s, Key - %s, Value - %s\n",
					message.Topic, message.Key, message.Value)

				var info database.AddRepoInfo

				err = json.Unmarshal(message.Value, &info)
				if err != nil {
					fmt.Println("Error parsing JSON:", err)
				}

				err = GetRepoInfo(&info)
				if err != nil {
					fmt.Println("Error getting repo info:", err)
				}

				err = database.AddRepoReadMeData(pool, &info)
			}
		}
	}, nil
}
