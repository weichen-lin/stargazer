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
)

func GetUserInfo(name string) (*database.User, error) {

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("GITHUB_TOKEN not set")
	}

	var user database.User

	req, err := http.NewRequest("GET", "https://api.github.com/users/"+name, nil)
	if err != nil {
		return &user, err
	}

	req.Header.Set("Authorization", "token "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &user, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &user, fmt.Errorf("error: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &user, err
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return &user, err
	}

	return &user, nil
}

func GetUserProfileConsumer() (func(neo4j.DriverWithContext), error) {
	brokers := []string{"localhost:9092"}
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Producer.Return.Successes = true

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		fmt.Println("Error creating consumer:", err)
		return nil, err
	}

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		fmt.Println("Error creating producer:", err)
		return nil, err
	}

	consumerPartitionConsumer, err := consumer.ConsumePartition("get_user_profile", 0, sarama.OffsetNewest)
	if err != nil {
		fmt.Println("Error creating partition consumer:", err)
		return nil, err
	}

	return func(driver neo4j.DriverWithContext) {
		for {
			select {
			case err := <-consumerPartitionConsumer.Errors():
				fmt.Println("Error:", err.Err)
			case message := <-consumerPartitionConsumer.Messages():
				fmt.Printf("Received message: Topic - %s, Key - %s, Value - %s\n",
					message.Topic, message.Key, message.Value)

				user, err := GetUserInfo(string(message.Value))
				if err != nil {
					fmt.Println("Error getting user info:", err)
				}

				// create user into neo4j
				user_id, err := database.CreateUser(driver, user)
				if err != nil {
					fmt.Println("Error creating user:", err)
				}

				fmt.Println("Get User ID:", user_id)

				info := database.GetGithubReposInfo{
					UserId:   user_id,
					Username: user.Login,
					Page:     1,
				}

				jsonString, err := json.Marshal(info)
				if err != nil {
					fmt.Println("Error marshalling JSON:", err)
				}

				_, _, err = producer.SendMessage(&sarama.ProducerMessage{
					Topic: "get_user_stars",
					Value: sarama.StringEncoder(jsonString),
				})
				if err != nil {
					fmt.Println("Error sending message:", err)
				}

				fmt.Printf("Sent message to Topic - %s, : %s\n", "get_user_stars", jsonString)
			}
		}
	}, nil
}
