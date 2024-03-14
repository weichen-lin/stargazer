package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/IBM/sarama"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/patrickmn/go-cache"
	"github.com/weichen-lin/kafka-service/consumer"
	database "github.com/weichen-lin/kafka-service/db"
	"github.com/weichen-lin/kafka-service/util"
)

type GetGithubReposInfo struct {
	UserId   string `form:"user_id" json:"user_id" binding:"required"`
	Username string `form:"username" json:"username" binding:"required"`
	Page     int    `form:"page" json:"page" binding:"required"`
}

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

	sync_star_vector, err := consumer.SyncStarVectorConsumer()
	if err != nil {
		fmt.Println("Error creating consumer:", err)
		return
	}

	go get_repo_consumer(driver, pool)
	go sync_star_vector(driver)

	server_cache := cache.New(20*time.Minute, 10*time.Minute)

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
		c.JSON(http.StatusOK, gin.H{
			"message": "OK",
		})
	})

	r.POST("/get_user_stars", AuthMiddleware(), func(c *gin.Context) {
		var info GetGithubReposInfo

		if err := c.ShouldBindJSON(&info); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if _, found := server_cache.Get(info.UserId); found {

			_, expired, _ := server_cache.GetWithExpiration(info.UserId)
			remain := time.Until(expired)
			mins := int(remain.Minutes())

			c.JSON(http.StatusConflict, gin.H{
				"message": "This user is already being processed. Please try again later.",
				"expires": fmt.Sprintf("%d minutes", mins),
			})
			return
		}

		server_cache.Set(info.UserId, true, time.Minute*30)

		jsonString, err := json.Marshal(info)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		partition, offset, err := producer.SendMessage(&sarama.ProducerMessage{
			Topic: "get_user_stars",
			Value: sarama.StringEncoder(jsonString),
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		c.JSON(200, gin.H{
			"partition": partition,
			"offset":    offset,
			"message":   "OK",
		})
	})

	r.GET("/sync_user_stars", AuthJWTMiddleware(), Neo4jDriverMiddleware(driver), HandleConnections)

	r.Run(":8080")
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		expectedToken := os.Getenv("AUTHENTICATION_TOKEN")

		if token != "Bearer "+expectedToken {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func AuthJWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		expectedToken := os.Getenv("AUTHENTICATION_TOKEN")

		jwtMaker, err := util.NewJWTMaker(expectedToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		payload, err := jwtMaker.VerifyToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Set("userName", payload.UserName)

		c.Next()
	}
}

func Neo4jDriverMiddleware(driver neo4j.DriverWithContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("neo4jDriver", driver)
		c.Next()
	}
}
