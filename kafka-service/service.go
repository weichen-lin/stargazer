package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
	db "github.com/weichen-lin/kafka-service/db"
)

type Service struct {
	DB       *db.Database
	Producer *kafka.Writer
}

func NewService(registries ...RegisterConsumer) *Service {
	database := db.NewDatabase()

	endpoint := os.Getenv("KAFKA_ENDPOINT")
	user_name := os.Getenv("KAFKA_USER_NAME")
	password := os.Getenv("KAFKA_USER_PASSWORD")
	topic := os.Getenv("GET_USER_STAR_TOPIC")

	if endpoint == "" || user_name == "" || password == "" || topic == "" {
		panic("Kafka environment variables not set")
	}

	mechanism, _ := scram.Mechanism(scram.SHA256, user_name, password)

	producer := &kafka.Writer{
		Addr:  kafka.TCP(endpoint),
		Topic: topic,
		Transport: &kafka.Transport{
			SASL: mechanism,
			TLS:  &tls.Config{},
		},
	}

	for _, registry := range registries {

		go func(registry RegisterConsumer) {

			consumer := kafka.NewReader(kafka.ReaderConfig{
				Brokers: []string{endpoint},
				GroupID: "1",
				Topic:   registry.Topic,
				Dialer: &kafka.Dialer{
					SASLMechanism: mechanism,
					TLS:           &tls.Config{},
				},
			})

			for {
				message, err := consumer.ReadMessage(context.Background())
				if err != nil {
					panic(err)
				}

				err = registry.HandlerFunc(database, message, producer)
				if err != nil {
					fmt.Println(err)
				}

			}
		}(registry)

	}

	return &Service{
		DB:       database,
		Producer: producer,
	}
}
