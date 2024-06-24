package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/weichen-lin/kafka-service/controller"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	topic := os.Getenv("GET_USER_STAR_TOPIC")
	if topic == "" {
		panic("Kafka topic variables not set")
	}

	m := NewMiddleware()
	service := NewService(
		RegisterConsumer{
			Topic:       topic,
			HandlerFunc: GetGithubRepos,
		},
	)

	c := controller.NewController(service.DB, service.Producer)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK",
		})
	})

	r.GET("/get_user_stars", m.JWTAuth(), c.GetUserStars)

	r.OPTIONS("/sync_user_stars", m.Cors())
	r.GET("/sync_user_stars", m.JWTAuth(), c.HandleConnections)

	r.PATCH("/update_cron_tab_setting", m.Cors(), m.JWTAuth(), c.UpdateCronTabSetting)

	r.Run(":8080")
}
