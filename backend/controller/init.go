package controller

import (
	"github.com/weichen-lin/kabaka"
	"github.com/weichen-lin/stargazer/db"
)

type Controller struct {
	db     *db.Database
	kabaka *kabaka.Kabaka
}

func NewController(logger kabaka.Logger) *Controller {

	db := db.NewDatabase()
	kabaka := kabaka.NewKabaka(&kabaka.Config{
		Logger: logger,
	})

	return &Controller{
		db:     db,
		kabaka: kabaka,
	}
}
