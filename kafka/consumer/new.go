package consumer

import (
	"fmt"

	"github.com/IBM/sarama"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/weichen-lin/kafka-service/neo4j_kafka"
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
					return
				}

				err = neo4j_kafka.CreateUser(driver, user)
				if err != nil {
					fmt.Println("Error creating user:", err)
				}
			}
		}
	}, nil
}
