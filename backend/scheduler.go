package main

import (
	"context"
	"sync"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/weichen-lin/stargazer/db"
)

type Scheduler struct {
	sync.Mutex
	centraller gocron.Scheduler
	jobs       map[string]uuid.UUID
	producer   *kafka.Writer
}

func NewScheduler(db *db.Database, producer *kafka.Writer) *Scheduler {

	c, err := gocron.NewScheduler()
	if err != nil {
		panic(err)
	}

	s := &Scheduler{
		centraller: c,
		jobs:       make(map[string]uuid.UUID),
		producer:   producer,
	}

	c.Start()

	return s
}

func (s *Scheduler) GetJob(jobID string) uuid.UUID {
	s.Lock()
	defer s.Unlock()
	return s.jobs[jobID]
}

func (s *Scheduler) Update(email string, hour int) error {
	s.Lock()
	defer s.Unlock()

	id, exists := s.jobs[email]

	if exists {
		s.centraller.RemoveJob(id)
	}

	j, err := s.centraller.NewJob(
		gocron.DailyJob(
			1,
			gocron.NewAtTimes(
				gocron.NewAtTime(uint(hour), 0, 0),
			),
		),
		gocron.NewTask(func() {
			s.producer.WriteMessages(context.Background(), kafka.Message{
				Value: []byte(`{"email":"` + email + `","page":1}`),
			})
		}),
	)
	if err != nil {
		return err
	}

	s.jobs[email] = j.ID()
	return nil
}

func (s *Scheduler) GetAll() map[string]uuid.UUID {
	s.Lock()
	defer s.Unlock()
	return s.jobs
}
