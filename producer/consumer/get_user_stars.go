package consumer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/IBM/sarama"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/patrickmn/go-cache"
	neo4jOpeartion "github.com/weichen-lin/kafka-service/neo4j"
	"github.com/weichen-lin/kafka-service/util"
)

var tokenCache = cache.New(20*time.Minute, 10*time.Minute)


type GetGithubReposInfo struct {
	UserId   string `json:"user_id"`
	Username string `json:"user_name"`
	Page     int    `json:"page"`
}

func GetUserStarredRepos(info *GetGithubReposInfo, token string) ([]neo4jOpeartion.Repository, error) {

	url := fmt.Sprintf("https://api.github.com/user/starred?&page=%d", info.Page)

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

	var repos []neo4jOpeartion.Repository

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

func GetGithubReposConsumer() (func(neo4j.DriverWithContext), error) {
	kafka_url := os.Getenv("KAFKA_URL")
	brokers := []string{kafka_url}
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
		return nil, err
	}

	return func(driver neo4j.DriverWithContext) {
		for {
			select {
			case err := <-consumerPartitionConsumer.Errors():
				fmt.Println("Error:", err)
			case message := <-consumerPartitionConsumer.Messages():
				fmt.Printf("Received message: Topic - %s, Key - %s, Value - %s\n",
					message.Topic, message.Key, message.Value)

				var info GetGithubReposInfo

				err = json.Unmarshal(message.Value, &info)
				if err != nil {
					fmt.Println("Error parsing JSON:", err)
					continue
				}

				var token string

				if tokenFormCache, found := tokenCache.Get(info.UserId); !found {
					tokenValue, err := neo4jOpeartion.GetUserGithubToken(driver, info.Username)
					if err != nil {
						fmt.Println("Error getting user token:", err)
						continue
					}
					token = tokenValue
					tokenCache.Set(info.UserId, tokenValue, cache.DefaultExpiration)
				} else {
					tokenValue, ok := tokenFormCache.(string)
					if !ok {
						fmt.Println("Error converting token to string:", token)
						continue
					}
					token = tokenValue
				}

				stars, err := GetUserStarredRepos(&info, token)

				if err != nil {
					fmt.Println("Error getting user stars:", err)
					continue
				}

				if len(stars) == 30 {
					info.Page++

					jsonString, err := json.Marshal(info)
					if err != nil {
						fmt.Println("Error marshalling JSON:", err)
						continue
					}

					_, _, err = producer.SendMessage(&sarama.ProducerMessage{
						Topic: "get_user_stars",
						Value: sarama.StringEncoder(jsonString),
					})

					if err != nil {
						fmt.Println("Error sending message:", err)
						continue
					}
				} else {
					email, _ := neo4jOpeartion.GetUserEmail(driver, info.Username)

					params := &util.SendMailParams{
						Email:      email,
						Name:       info.Username,
						StarsCount: (info.Page-1)*30 + len(stars),
					}

					err := util.SendMail(params)
					if err != nil {
						fmt.Println("Error sending email:", err)
					}
				}

				for _, repo := range stars {
					neo4jOpeartion.CreateRepository(driver, &repo, info.UserId)
				}
			}
		}
	}, nil
}
