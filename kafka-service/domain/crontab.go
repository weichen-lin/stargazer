package domain

import (
	"errors"
	"time"
)

type CronTab struct {
	triggerAt     time.Time
	createdAt     time.Time
	updatedAt     time.Time
	status        string
}

type CronTabEntity struct {
	TriggerAt     string `json:"trigger_at"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	Status        string `json:"status"`
}

func (c *CronTab) TriggerAt() time.Time {
	return c.triggerAt
}

func (c *CronTab) CreatedAt() time.Time {
	return c.createdAt
}

func (c *CronTab) UpdatedAt() time.Time {
	return c.updatedAt
}

func (c *CronTab) Status() string {
	return c.status
}

func (c *CronTab) SetTriggerAt(t string) error {
	if t == "" {
		c.triggerAt = time.Time{}
		return nil
	}

	parsedTime, err := ParseTime(t)
	if err != nil {
		return err
	}

	c.triggerAt = parsedTime
	return nil
}

func (c *CronTab) SetCreatedAt(t string) error {
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

func (c *CronTab) SetUpdatedAt(t string) error {
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

func (c *CronTab) SetStatus(s string) error {
	if s == "" {
		return errors.New("invalid status: cannot be empty")
	}
	c.status = s
	return nil
}

func NewCronTab() (*CronTab, error) {
	cronTab := &CronTab{}

	now := time.Now()

	cronTab.SetTriggerAt("")
	cronTab.SetCreatedAt(now.Format(time.RFC3339))
	cronTab.SetUpdatedAt("")
	cronTab.SetStatus("new")

	return cronTab, nil
}

func (c *CronTab) ToCronTabEntity() *CronTabEntity {

	var triggerAt, updatedAt string

	if c.TriggerAt().IsZero() {
		triggerAt = ""
	} else {
		triggerAt = c.TriggerAt().Format(time.RFC3339)
	}

	if c.UpdatedAt().IsZero() {
		updatedAt = ""
	} else {
		updatedAt = c.UpdatedAt().Format(time.RFC3339)
	}

	return &CronTabEntity{
		CreatedAt:     c.CreatedAt().Format(time.RFC3339),
		TriggerAt:     triggerAt,
		UpdatedAt:     updatedAt,
		Status:        c.Status(),
	}
}

func FromCronTabEntity(cronTabEntity *CronTabEntity) (*CronTab, error) {
	cronTab := &CronTab{}

	if err := cronTab.SetTriggerAt(cronTabEntity.TriggerAt); err != nil {
		return nil, err
	}

	if err := cronTab.SetCreatedAt(cronTabEntity.CreatedAt); err != nil {
		return nil, err
	}

	if err := cronTab.SetUpdatedAt(cronTabEntity.UpdatedAt); err != nil {
		return nil, err
	}

	if err := cronTab.SetStatus(cronTabEntity.Status); err != nil {
		return nil, err
	}

	return cronTab, nil
}
