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

	m := NewMiddleware()
	service := NewService(
		RegisterConsumer{
			Topic:       "",
			HandlerFunc: GetGithubRepos,
		},
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	c := controller.NewController(service.DB, service.Producer)

	r := gin.Default()

	r.HEAD("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK",
		})
	})

	repo := r.Group("/repository", m.JWTAuth())
	{
		repo.GET("/:id", c.GetRepository)
	}

	crontab := r.Group("/crontab", m.JWTAuth())
	{
		crontab.GET("/", c.GetCronTab)
		crontab.POST("/", c.GetCronTab)
	}

	r.Run(":" + port)
}
