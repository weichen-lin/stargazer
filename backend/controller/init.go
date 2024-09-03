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

	handleFunc := func(msg *kabaka.Message) error {
		err := util.GetGithubRepos(db, *msg, kbk)

		if err != nil {
			return err
		}

		return nil
	}

	kbk.CreateTopic("star-syncer")

	kbk.Subscribe("star-syncer", handleFunc)

	return &Controller{
		db:     db,
		kabaka: kbk,
	}
}
