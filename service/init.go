package service

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/weichen-lin/stargazer/infrastructure"
	"github.com/weichen-lin/stargazer/repository"
)

type Service struct {
	db         *infrastructure.Database
	rdb        *redis.Client
	repository *repository.Repository
	config     *infrastructure.Config
}

func NewService(config *infrastructure.Config) *Service {

	database, err := infrastructure.NewDatabase(context.Background(), config.DatabaseURL, nil)
	if err != nil {
		panic(err)
	}

	opt, err := redis.ParseURL(config.RedisURL)
	if err != nil {
		panic(err)
	}

	cli := redis.NewClient(opt)

	repository := repository.NewRepository(cli)

	return &Service{
		db:         database,
		rdb:        cli,
		repository: repository,
		config:     config,
	}
}
