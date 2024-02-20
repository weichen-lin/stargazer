package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/weichen-lin/kafka-service/consumer"
	database "github.com/weichen-lin/kafka-service/db"
)

type GetGithubReposInfo struct {
	UserId   string  `form:"user_id" json:"user_id" binding:"required"`
	Username string `form:"username" json:"username" binding:"required"`
	Page     int    `form:"page" json:"page" binding:"required"`
}

func main() {

	if os.Getenv("APP_ENV") == "production" {
		godotenv.Load(
			".env",
		)
	} else {
		godotenv.Load(
			".env.dev",
		)
	}

	neo4j_url := os.Getenv("NEO4J_URL")
	neo4j_password := os.Getenv("NEO4J_PASSWORD")

	driver, err := neo4j.NewDriverWithContext(
		neo4j_url,
		neo4j.BasicAuth("neo4j", neo4j_password, ""),
	)
	if err != nil {
		fmt.Println("Error creating driver:", err)
		return
	}

	postgresql_url := os.Getenv("POSTGRESQL_URL")

	pool, err := database.NewPostgresDB(postgresql_url)
	if err != nil {
		fmt.Println("Error creating postgres connection:", err)
		return
	}

	brokers := []string{"localhost:9092"}
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		fmt.Println("Error creating producer:", err)
		return
	}

	get_repo_consumer, err := consumer.GetGithubReposConsumer()

	go get_repo_consumer(driver, pool)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})

	r.POST("/get_user_stars", func(c *gin.Context) {
		var info GetGithubReposInfo

		if err := c.ShouldBindJSON(&info); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		jsonString, err := json.Marshal(info)

		partition, offset, err := producer.SendMessage(&sarama.ProducerMessage{
			Topic: "get_user_stars",
			Value: sarama.StringEncoder(jsonString),
		})

		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
		}

		c.JSON(200, gin.H{
			"partition": partition,
			"offset":    offset,
			"message":   "OK",
		})
	})

	r.Run(":8080")
}
