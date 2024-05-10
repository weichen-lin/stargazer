package main

import (
	"fmt"
	"os"

	"github.com/IBM/sarama"
	db "github.com/weichen-lin/kafka-service/db"
)

type Service struct {
	DB       db.Database
	Producer sarama.SyncProducer
}

func NewService(registries ...RegisterConsumer) *Service {
	database := db.NewDatabase()

	kafka_url := os.Getenv("KAFKA_URL")
	brokers := []string{kafka_url}

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Consumer.Return.Errors = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		panic(err)
	}

	for _, registry := range registries {

		consumer, err := sarama.NewConsumer(brokers, config)
		if err != nil {
			panic(err)
		}

		consumerPartitionConsumer, err := consumer.ConsumePartition(registry.Topic, 0, sarama.OffsetNewest)
		if err != nil {
			panic(err)
		}

		go func(registry RegisterConsumer) {
			for {
				select {
				case err := <-consumerPartitionConsumer.Errors():
					panic(err)
				case message := <-consumerPartitionConsumer.Messages():
					err := registry.HandlerFunc(database, message, producer)
					if err != nil {
						fmt.Println(err)
					}
				}
			}
		}(registry)

	}

	return &Service{
		DB:       database,
		Producer: producer,
	}
}
