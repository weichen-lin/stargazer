package main

import (
	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/gin-gonic/gin"
	"github.com/weichen-lin/stargazer/infrastructure"
	"go.uber.org/zap/zapcore"
)

func main() {

	config, err := infrastructure.NewConfig()
	if err != nil {
		panic(err)
	}

	clerk.SetKey(config.ClerkSecret)
	infrastructure.InitLogger(zapcore.DebugLevel, false)

	r := gin.Default()

	c, b := initStarGazer(config)

	api := r.Group("/api", ClerkAuth())

	user := api.Group("/user", ClerkAuth())
	{
		user.GET("/crontab", c.GetCrontab)
	}

	background := api.Group("/background", ClerkAuth())
	{
		background.GET("/sync-user-repositories", b.PublishSyncUserRepositoriesEvent)
	}

	repository := api.Group("/repository", ClerkAuth())
	{
		repository.GET("/language-distribution", c.GetLanguageDistribution)
		repository.GET("/latest-updated", c.GetLatestUserUpdatedRepositories)
		repository.GET("/latest-starred", c.GetLatestUserSyncedRepositories)
	}

	r.Run()
}
