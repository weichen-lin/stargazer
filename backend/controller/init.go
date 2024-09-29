package controller

import (
	"github.com/weichen-lin/kabaka"
	"github.com/weichen-lin/stargazer/db"
	"github.com/weichen-lin/stargazer/util"
)

type Controller struct {
	db        *db.Database
	kabaka    *kabaka.Kabaka
	scheduler *Scheduler
}

func NewController(logger kabaka.Logger) *Controller {

	db := db.NewDatabase()
	scheduler := NewScheduler()

	cronjobs := db.GetAllCrontab()

	bk := kabaka.NewKabaka(&kabaka.Config{
		Logger: logger,
	})

	starSyncerHandleFunc := func(msg *kabaka.Message) error {
		err := util.GetGithubRepos(db, *msg, bk)

		if err != nil {
			return err
		}

		return nil
	}

	topicHandlerFunc := func(msg *kabaka.Message) error {
		util.GetRepositoryTopics(db, *msg, bk)
		return nil
	}

	bk.CreateTopic("star-syncer")
	bk.CreateTopic("topic-syncer")

	bk.Subscribe("star-syncer", starSyncerHandleFunc)
	bk.Subscribe("topic-syncer", topicHandlerFunc)

	for _, cronjob := range cronjobs {
		if cronjob.TriggeredAt != "" {

			fn := func() error {
				bk.Publish("star-syncer", []byte(`{"email":"`+cronjob.Email+`","page":1}`))
				return nil
			}

			scheduler.AddJob(cronjob, fn)
		}
	}

	return &Controller{
		db:        db,
		kabaka:    bk,
		scheduler: scheduler,
	}
}
