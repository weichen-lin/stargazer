package controller

import (
	"github.com/segmentio/kafka-go"
	"github.com/weichen-lin/kafka-service/db"
	"github.com/weichen-lin/kafka-service/scheduler"
)

type Controller struct {
	db        *db.Database
	producer  *kafka.Writer
	scheduler *scheduler.Scheduler
}

func NewController(db *db.Database, producer *kafka.Writer) *Controller {
	return &Controller{
		db:        db,
		producer:  producer,
		scheduler: scheduler.NewScheduler(db, producer),
	}
}
