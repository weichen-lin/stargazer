package main

import (
	"fmt"

	"github.com/IBM/sarama"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/weichen-lin/kafka-service/consumer"
)

func main() {
	driver, err := neo4j.NewDriverWithContext(
		"",
		neo4j.BasicAuth("neo4j", "", ""),
	)
	if err != nil {
		fmt.Println("Error creating driver:", err)
		return
	}

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	brokers := []string{"localhost:9092"}

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		fmt.Println("Error creating producer:", err)
		return
	}
	defer func() {
		if err := producer.Close(); err != nil {
			fmt.Println("Error closing producer:", err)
		}
	}()

	con, get_star_consumer, _ := consumer.GetUserProfileConsumer()
	_, get_repo_consumer, _ := consumer.GetGithubReposConsumer()

	defer func() {
		if err := con.Close(); err != nil {
			fmt.Println("Error closing consumer:", err)
		}
	}()

	go get_star_consumer(driver)
	go get_repo_consumer(driver)

	select {}
}
