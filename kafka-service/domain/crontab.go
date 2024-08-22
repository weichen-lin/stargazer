package domain

import (
	"errors"
	"time"
)

type Crontab struct {
	*AggregateRoot

	triggeredAt     time.Time
	createdAt       time.Time
	updatedAt       time.Time
	status          string
	lastTriggeredAt time.Time
}

type CrontabEntity struct {
	TriggeredAt     string `json:"triggered_at"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
	Status          string `json:"status"`
	LastTriggeredAt string `json:"last_triggered_at"`
	Version         int64  `json:"version"`
}

func (c *Crontab) TriggeredAt() time.Time {
	return c.triggeredAt
}

func (c *Crontab) CreatedAt() time.Time {
	return c.createdAt
}

func (c *Crontab) UpdatedAt() time.Time {
	return c.updatedAt
}

func (c *Crontab) LastTriggeredAt() time.Time {
	return c.lastTriggeredAt
}

func (c *Crontab) Status() string {
	return c.status
}

func (c *Crontab) SetTriggeredAt(t string) error {
	if t == "" {
		c.triggeredAt = time.Time{}
		return nil
	}

	parsedTime, err := ParseTime(t)
	if err != nil {
		return err
	}

	c.triggeredAt = parsedTime
	return nil
}

func (c *Crontab) SetLastTriggerAt(t string) error {
	if t == "" {
		c.lastTriggeredAt = time.Time{}
		return nil
	}

	parsedTime, err := ParseTime(t)
	if err != nil {
		return err
	}

	if parsedTime.Before(c.createdAt) {
		return errors.New("last triggered time cannot be before created time")
	}

	c.lastTriggeredAt = parsedTime
	return nil
}

func (c *Crontab) SetCreatedAt(t string) error {
	if t == "" {
		return errors.New("created time cannot be empty")
	}

	parsedTime, err := ParseTime(t)
	if err != nil {
		return err
	}

	c.createdAt = parsedTime
	return nil
}

func (c *Crontab) SetUpdatedAt(t string) error {
	if t == "" {
		c.updatedAt = time.Time{}
		return nil
	}

	parsedTime, err := ParseTime(t)
	if err != nil {
		return err
	}

	if parsedTime.Before(c.createdAt) {
		return errors.New("updated time cannot be before created time")
	}

	c.updatedAt = parsedTime
	return nil
}

func (c *Crontab) SetStatus(s string) error {
	if s == "" {
		return errors.New("invalid status: cannot be empty")
	}
	c.status = s
	return nil
}

func NewCrontab() *Crontab {
	Crontab := &Crontab{}

	now := time.Now()

	Crontab.SetTriggeredAt("")
	Crontab.SetCreatedAt(now.Format(time.RFC3339))
	Crontab.SetUpdatedAt("")
	Crontab.SetStatus("new")

	root := NewAggregateRoot()
	Crontab.AggregateRoot = root

	return Crontab
}

func (c *Crontab) ToCrontabEntity() *CrontabEntity {

	var triggeredAt, updatedAt, lastTriggeredAt string

	if c.TriggeredAt().IsZero() {
		triggeredAt = ""
	} else {
		triggeredAt = c.TriggeredAt().Format(time.RFC3339)
	}

	if c.UpdatedAt().IsZero() {
		updatedAt = ""
	} else {
		updatedAt = c.UpdatedAt().Format(time.RFC3339)
	}

	if c.LastTriggeredAt().IsZero() {
		lastTriggeredAt = ""
	} else {
		lastTriggeredAt = c.LastTriggeredAt().Format(time.RFC3339)
	}

	return &CrontabEntity{
		CreatedAt:       c.CreatedAt().Format(time.RFC3339),
		TriggeredAt:     triggeredAt,
		UpdatedAt:       updatedAt,
		LastTriggeredAt: lastTriggeredAt,
		Status:          c.Status(),
		Version:         c.Version(),
	}
}

func FromCrontabEntity(CrontabEntity *CrontabEntity) (*Crontab, error) {
	Crontab := &Crontab{}

	root := NewAggregateRoot()
	root.version = CrontabEntity.Version

	if err := Crontab.SetTriggeredAt(CrontabEntity.TriggeredAt); err != nil {
		return nil, err
	}

	if err := Crontab.SetCreatedAt(CrontabEntity.CreatedAt); err != nil {
		return nil, err
	}

	if err := Crontab.SetUpdatedAt(CrontabEntity.UpdatedAt); err != nil {
		return nil, err
	}

	if err := Crontab.SetStatus(CrontabEntity.Status); err != nil {
		return nil, err
	}

	if err := Crontab.SetLastTriggerAt(CrontabEntity.LastTriggeredAt); err != nil {
		return nil, err
	}

	Crontab.AggregateRoot = root

	return Crontab, nil
}
