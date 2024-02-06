package consumer

import (
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/weichen-lin/kafka-service/workflow"
)

func GetUserProfileConsumer() (sarama.PartitionConsumer, func(neo4j.DriverWithContext), error) {
	brokers := []string{"localhost:9092"}
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		fmt.Println("Error creating consumer:", err)
	}

	consumerPartitionConsumer, err := consumer.ConsumePartition("get_user_profile", 0, sarama.OffsetNewest)
	if err != nil {
		fmt.Println("Error creating partition consumer:", err)
	}

	return consumerPartitionConsumer, func(driver neo4j.DriverWithContext) {
		for {
			select {
			case err := <-consumerPartitionConsumer.Errors():
				fmt.Println("Error:", err.Err)
			case message := <-consumerPartitionConsumer.Messages():
				fmt.Printf("Received message: Topic - %s, Partition - %d, Offset - %d, Key - %s, Value - %s\n",
					message.Topic, message.Partition, message.Offset, message.Key, message.Value)

				user, err := workflow.GetUserInfo(string(message.Value))
				if err != nil {
					fmt.Println("Error getting user info:", err)
				}

				fmt.Println("User:", user)

				// err = neo4j_kafka.CreateUser(driver, user)
				// if err != nil {
				// 	fmt.Println("Error creating user:", err)
				// }
			}
		}
	}, nil
}

func GetGithubReposConsumer() (sarama.PartitionConsumer, func(neo4j.DriverWithContext), error) {
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
		return nil, nil, err
	}
	

	consumerPartitionConsumer, err := consumer.ConsumePartition("get_user_stars", 0, sarama.OffsetNewest)
	if err != nil {
		fmt.Println("Error creating partition consumer:", err)
	}

	return consumerPartitionConsumer, func(driver neo4j.DriverWithContext) {
		for {
			select {
			case err := <-consumerPartitionConsumer.Errors():
				fmt.Println("Error:", err.Err)
			case message := <-consumerPartitionConsumer.Messages():
				fmt.Printf("Received message: Topic - %s, Partition - %d, Offset - %d, Key - %s, Value - %s\n",
					message.Topic, message.Partition, message.Offset, message.Key, message.Value)

				var info workflow.GetGithubReposInfo

				err = json.Unmarshal(message.Value, &info)
				if err != nil {
					fmt.Println("Error parsing JSON:", err)
				}

				stars, err := workflow.GetUserStarredRepos(&info)
				if err != nil {
					fmt.Println("Error getting user info:", err)
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
					fmt.Println("Repo:", repo.FullName)
				}

				// err = neo4j_kafka.CreateUser(driver, user)
				// if err != nil {
				// 	fmt.Println("Error creating user:", err)
				// }
			}
		}
	}, nil
}

func GetGithubRepoInfoConsumer() (sarama.PartitionConsumer, func(neo4j.DriverWithContext), error) {
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
		return nil, nil, err
	}
	

	consumerPartitionConsumer, err := consumer.ConsumePartition("get_user_stars", 0, sarama.OffsetNewest)
	if err != nil {
		fmt.Println("Error creating partition consumer:", err)
	}

	return consumerPartitionConsumer, func(driver neo4j.DriverWithContext) {
		for {
			select {
			case err := <-consumerPartitionConsumer.Errors():
				fmt.Println("Error:", err.Err)
			case message := <-consumerPartitionConsumer.Messages():
				fmt.Printf("Received message: Topic - %s, Partition - %d, Offset - %d, Key - %s, Value - %s\n",
					message.Topic, message.Partition, message.Offset, message.Key, message.Value)

				var info workflow.GetGithubReposInfo

				err = json.Unmarshal(message.Value, &info)
				if err != nil {
					fmt.Println("Error parsing JSON:", err)
				}

				stars, err := workflow.GetUserStarredRepos(&info)
				if err != nil {
					fmt.Println("Error getting user info:", err)
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
					fmt.Println("Repo:", repo.FullName)
				}

				// err = neo4j_kafka.CreateUser(driver, user)
				// if err != nil {
				// 	fmt.Println("Error creating user:", err)
				// }
			}
		}
	}, nil
}