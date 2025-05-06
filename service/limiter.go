package service

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func (s *Service) GetSyncLimiter(ctx context.Context, clerkId string) bool {
	limiter := fmt.Sprintf("%s-sync-limiter", clerkId)

	_, err := s.rdb.Get(ctx, limiter).Result()
	if err != nil {
		if err == redis.Nil {
			s.rdb.SetNX(ctx, limiter, true, time.Minute*30)

			return true
		}

		return false
	}

	return false
}
