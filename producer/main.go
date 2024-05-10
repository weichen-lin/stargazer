package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/weichen-lin/kafka-service/controller"
)

func main() {

	m := NewMiddleware()
	service := NewService(
		RegisterConsumer{
			Topic:       "get_user_stars",
			HandlerFunc: GetGithubRepos,
		},
	)

	port := os.Getenv("PRODUCER_PORT")

	r := gin.Default()

	cors_config := cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders: []string{"Content-Length"},
		MaxAge:        12 * time.Hour,
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK",
		})
	})

	r.POST("/get_user_stars", m.BasicAuth(), m.Producer(service.Producer), controller.GetUserStars)

	r.GET("/sync_user_stars", cors.New(cors_config), m.JWTAuth(), m.DatabaseDriver(service.DB), controller.HandleConnections)

	r.PATCH("/update_cron_tab_setting", m.JWTAuth(), m.DatabaseDriver(service.DB), controller.UpdateCronTabSetting)

	r.Run(port)
}
