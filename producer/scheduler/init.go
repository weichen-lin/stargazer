package scheduler

import (
	"sync"

	"github.com/IBM/sarama"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/weichen-lin/kafka-service/db"
)

type Scheduler struct {
	sync.Mutex
	centraller gocron.Scheduler
	jobs       map[string]uuid.UUID
	producer   sarama.SyncProducer
}

func NewScheduler(db *db.Database, producer sarama.SyncProducer) *Scheduler {

	c, err := gocron.NewScheduler()
	if err != nil {
		panic(err)
	}

	s := &Scheduler{
		centraller: c,
		jobs:       make(map[string]uuid.UUID),
		producer:   producer,
	}

	crontabs, err := db.GetAllUserCrontab()
	if err != nil {
		panic(err)
	}

	for _, crontab := range crontabs {
		j, err := c.NewJob(
			gocron.DailyJob(
				1,
				gocron.NewAtTimes(
					gocron.NewAtTime(uint(crontab.Hour), 0, 0),
				),
			),
			gocron.NewTask(func() {
				producer.SendMessage(&sarama.ProducerMessage{
					Topic: "get_user_stars",
					Value: sarama.StringEncoder(`{"email":"` + crontab.Email + `","page":1}`),
				})
			}),
		)

		if err != nil {
			panic(err)
		}

		s.jobs[crontab.Email] = j.ID()
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
			s.producer.SendMessage(&sarama.ProducerMessage{
				Topic: "get_user_stars",
				Value: sarama.StringEncoder(`{"email":"` + email + `","page":1}`),
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
