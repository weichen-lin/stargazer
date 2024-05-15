package controller

import (
	"github.com/IBM/sarama"
	"github.com/weichen-lin/kafka-service/db"
	"github.com/weichen-lin/kafka-service/scheduler"
)

type Controller struct {
	db       *db.Database
	producer sarama.SyncProducer
	scheduler *scheduler.Scheduler
}

func NewController(db *db.Database, producer sarama.SyncProducer) *Controller {
	return &Controller{
		db:       db,
		producer: producer,
		scheduler: scheduler.NewScheduler(db, producer),
	}
}
