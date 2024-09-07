package controller

import (
	"github.com/weichen-lin/kabaka"
	"github.com/weichen-lin/stargazer/db"
	"github.com/weichen-lin/stargazer/util"
)

type Controller struct {
	db     *db.Database
	kabaka *kabaka.Kabaka
}

func NewController(logger kabaka.Logger) *Controller {

	db := db.NewDatabase()
	kbk := kabaka.NewKabaka(&kabaka.Config{
		Logger: logger,
	})

	starSyncerHandleFunc := func(msg *kabaka.Message) error {
		err := util.GetGithubRepos(db, *msg, kbk)

		if err != nil {
			return err
		}

		return nil
	}

	topicHandlerFunc := func(msg *kabaka.Message) error {
		_ = util.GetRepositoryTopics(db, *msg, kbk)
		return nil
	}

	kbk.CreateTopic("star-syncer")
	kbk.CreateTopic("topic-syncer")

	kbk.Subscribe("star-syncer", starSyncerHandleFunc)
	kbk.Subscribe("topic-syncer", topicHandlerFunc)

	return &Controller{
		db:     db,
		kabaka: kbk,
	}
}
