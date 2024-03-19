package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/IBM/sarama"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/weichen-lin/kafka-service/consumer"
	"github.com/weichen-lin/kafka-service/controller"
	database "github.com/weichen-lin/kafka-service/db"
	"github.com/weichen-lin/kafka-service/middleware"
	"github.com/weichen-lin/kafka-service/zapper"
)

func main() {

	godotenv.Load(
		"secrets.env",
	)

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

	kafka_url := os.Getenv("KAFKA_URL")
	brokers := []string{kafka_url}
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		fmt.Println("Error creating producer:", err)
		return
	}

	get_repo_consumer, err := consumer.GetGithubReposConsumer()
	if err != nil {
		fmt.Println("Error creating consumer:", err)
		return
	}

	go get_repo_consumer(driver, pool)

	r := gin.Default()

	cors_config := cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders: []string{"Content-Length"},
		MaxAge:        12 * time.Hour,
	}


	r.Use(cors.New(cors_config))

	r.GET("/health", func(c *gin.Context) {
		zapper.Info("Health check")
		c.JSON(http.StatusOK, gin.H{
			"message": "OK",
		})
	})

	r.POST("/get_user_stars", middleware.AuthMiddleware(), middleware.ProducerMiddleware(producer), controller.GetUserStars)

	r.GET("/sync_user_stars", middleware.AuthJWTMiddleware(), middleware.Neo4jDriverMiddleware(driver), controller.HandleConnections)

	r.Run(":8080")
}
