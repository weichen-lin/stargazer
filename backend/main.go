package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/weichen-lin/stargazer/controller"
	"go.uber.org/zap"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	m := NewMiddleware()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	core, _ := zap.NewProduction()

	logger := NewStarGazerLogger(core)

	c := controller.NewController(logger)

	r := gin.Default()

	r.HEAD("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK",
		})
	})

	repo := r.Group("/repository", m.JWTAuth())
	{
		repo.GET("/", c.SearchRepoByLanguages)
		repo.GET("/sync-repository", c.SyncRepository)
		repo.GET("/topics", c.GetTopics)
		repo.GET("/:id", c.GetRepository)
		repo.GET("/language-distribution", c.GetUserLanguageDistribution)
	}

	crontab := r.Group("/crontab", m.JWTAuth())
	{
		crontab.GET("/", c.GetCrontab)
		crontab.POST("/", c.CreateCrontab)
		crontab.PATCH("/", c.UpdateCrontab)
	}

	tag := r.Group("/tag", m.JWTAuth())
	{
		tag.GET("/:id", c.GetTags)
		tag.POST("/", c.CreateTag)
		tag.DELETE("/", c.DeleteTag)
	}

	r.Run(":" + port)
}
