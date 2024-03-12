package consumer

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/IBM/sarama"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	database "github.com/weichen-lin/kafka-service/db"
	"github.com/weichen-lin/kafka-service/workflow"
)

func SyncStarVectorConsumer() (func(neo4j.DriverWithContext), error) {
	kafka_url := os.Getenv("KAFKA_URL")
	brokers := []string{kafka_url}
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Producer.Return.Successes = true

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		fmt.Println("Error creating consumer:", err)
	}

	consumerPartitionConsumer, err := consumer.ConsumePartition("sync_star_vector", 0, sarama.OffsetNewest)
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

				var info workflow.SyncUserStarMsg

				err := json.Unmarshal(message.Value, &info)
				if err != nil {
					fmt.Println("Error unmarshalling message:", err)
					continue
				}

				code, err := workflow.VectorizeStar(&info)
				if err != nil {
					fmt.Println("Error vectorizing star:", err)
					continue
				}

				if code != 200 && code != 201 {
					fmt.Println("Error vectorizing star:", err)
					continue
				}

				err = database.ConfirmVectorize(driver, &info)
				if err != nil {
					fmt.Println("Error confirming vectorize:", err)
					continue
				}
			}
		}
	}, nil
}
