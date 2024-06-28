package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/weichen-lin/kafka-service/controller"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("kafka-service")

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	tp := InitTracerHTTP()
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			fmt.Println("Error shutting down tracer provider: ", err)
		}
	}()

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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	c := controller.NewController(service.DB, service.Producer)

	r := gin.Default()

	r.Use(otelgin.Middleware(""))

	r.HEAD("/health", func(c *gin.Context) {

		ctx := c.Request.Context()

		_, span := tracer.Start(ctx, "health")
		defer span.End()
	

		c.JSON(http.StatusOK, gin.H{
			"message": "OK",
		})
	})

	r.GET("/get_user_stars", m.JWTAuth(), c.GetUserStars)

	r.OPTIONS("/sync_user_stars", m.Cors())
	r.GET("/sync_user_stars", m.JWTAuth(), c.HandleConnections)

	r.PATCH("/update_cron_tab_setting", m.Cors(), m.JWTAuth(), c.UpdateCronTabSetting)

	r.Run(":" + port)
}
