package controller

import (
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/weichen-lin/kabaka"
	"github.com/weichen-lin/stargazer/db"
	"github.com/weichen-lin/stargazer/util"
)

type Controller struct {
	db     *db.Database
	kabaka *kabaka.Kabaka
	rdb    *redis.Client
}

func NewController(logger kabaka.Logger) *Controller {

	db := db.NewDatabase()
	rdb := NewRedisCli()

	bk := kabaka.NewKabaka(&kabaka.Options{
		BufferSize:            24,
		DefaultMaxRetries:     1,
		DefaultRetryDelay:     time.Duration(5 * time.Second),
		DefaultProcessTimeout: time.Duration(10 * time.Second),
		Logger:                logger,
	})

	starSyncerHandleFunc := func(msg *kabaka.Message) error {
		err := util.GetGithubRepos(db, *msg, bk)

		if err != nil {
			return err
		}

		return nil
	}

	bk.CreateTopic("star-syncer")

	bk.Subscribe("star-syncer", starSyncerHandleFunc)

	return &Controller{
		db:     db,
		kabaka: bk,
		rdb:    rdb,
	}
}

func NewRedisCli() *redis.Client {
	endpoint := os.Getenv("REDIS_ENDPOINT")
	opt, err := redis.ParseURL(endpoint)
	if err != nil {
		panic(err)
	}

	return redis.NewClient(opt)
}
