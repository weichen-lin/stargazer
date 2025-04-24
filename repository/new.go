package repository

import (
	"github.com/redis/go-redis/v9"
)

type Repository struct {
	rdb *redis.Client
}

func NewRepository(rdb *redis.Client) *Repository {
	return &Repository{
		rdb: rdb,
	}
}
