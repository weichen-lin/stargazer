package controller

import (
	"github.com/IBM/sarama"
	"github.com/weichen-lin/kafka-service/db"
)

type Controller struct {
	db       *db.Database
	producer sarama.SyncProducer
}

func NewController(db *db.Database, producer sarama.SyncProducer) *Controller {
	return &Controller{
		db:       db,
		producer: producer,
	}
}
