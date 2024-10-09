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

	initOtel()

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
		repo.DELETE("/:id", c.DeleteRepository)
		repo.GET("/language-distribution", c.GetUserLanguageDistribution)
		repo.GET("/sort", c.GetRepositoriesByKey)
		repo.GET("/full-text-search", c.FullTextSearchWithQuery)
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

	collection := r.Group("/collection", m.JWTAuth())
	{
		collection.GET("/", c.GetCollections)
		collection.GET("/:id", c.GetCollection)
		collection.PATCH("/:id", c.UpdateCollection)
		collection.POST("/", c.CreateCollection)
		collection.DELETE("/", c.DeleteCollection)
		collection.POST("/repo", c.AddRepoIntoCollection)
		collection.DELETE("/repo", c.RemoveRepoFromCollection)
	}

	r.Run(":" + port)
}
