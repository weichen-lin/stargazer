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

	utcTime := t.UTC()

	now := time.Now()
	location := now.Location()

	localTime := utcTime.In(location)

	j, err := s.centraller.NewJob(
		gocron.DailyJob(
			1,
			gocron.NewAtTimes(
				gocron.NewAtTime(uint(localTime.Hour()), 0, 0),
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

func (s *Scheduler) Update(email string, triggered_at time.Time, fn func() error) error {
	s.Lock()
	defer s.Unlock()

	id, exists := s.jobs[email]

	if exists {
		s.centraller.RemoveJob(id)
	}

	utcTime := triggered_at.UTC()

	now := time.Now()
	location := now.Location()

	localTime := utcTime.In(location)

	j, err := s.centraller.NewJob(
		gocron.DailyJob(
			1,
			gocron.NewAtTimes(
				gocron.NewAtTime(uint(localTime.Hour()), 0, 0),
			),
		),
		gocron.NewTask(func() {
			fn()
		}),
	)
	if err != nil {
		return err
	}

	s.jobs[email] = j.ID()
	return nil
}

func (s *Scheduler) Remove(email string) {
	s.Lock()
	defer s.Unlock()

	id, exists := s.jobs[email]

	if exists {
		s.centraller.RemoveJob(id)
		delete(s.jobs, email)
	}
}

func (s *Scheduler) GetJobInfo(id uuid.UUID) gocron.Job {
	s.Lock()
	defer s.Unlock()

	jobs := s.centraller.Jobs()

	for _, job := range jobs {
		if job.ID() == id {
			return job
		}
	}

	return nil
}
