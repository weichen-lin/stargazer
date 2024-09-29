package controller

import (
	"sync"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/weichen-lin/stargazer/db"
)

type Scheduler struct {
	sync.Mutex
	centraller gocron.Scheduler
	jobs       map[string]uuid.UUID
}

func NewScheduler() *Scheduler {

	c, err := gocron.NewScheduler()
	if err != nil {
		panic(err)
	}

	s := &Scheduler{
		centraller: c,
		jobs:       make(map[string]uuid.UUID),
	}

	c.Start()

	return s
}

func (s *Scheduler) GetJob(email string) uuid.UUID {
	s.Lock()
	defer s.Unlock()

	jobId, ok := s.jobs[email]

	if !ok {
		return uuid.Nil
	}

	return jobId

}

func (s *Scheduler) AddJob(info *db.CrontabInfo, fn func() error) error {
	s.Lock()
	defer s.Unlock()

	t, err := time.Parse(time.RFC3339, info.TriggeredAt)
	if err != nil {
		return err
	}

	j, err := s.centraller.NewJob(
		gocron.DailyJob(
			1,
			gocron.NewAtTimes(
				gocron.NewAtTime(uint(t.Hour()), 0, 0),
			),
		),
		gocron.NewTask(func() {
			fn()
		}),
	)

	if err != nil {
		return err
	}

	s.jobs[info.Email] = j.ID()
	return nil
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
