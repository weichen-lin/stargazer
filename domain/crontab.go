package domain

import (
	"time"

	"github.com/google/uuid"
)

type Crontab struct {
	userId     uuid.UUID
	stargazers int32
	createdAt  time.Time
	updatedAt  time.Time
}

type CrontabEntity struct {
	UserId     uuid.UUID `json:"user_id"`
	Stargazers int32     `json:"stargazers"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CrontabDto struct {
	Stargazers int32     `json:"stargazers"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func NewCrontab(userId uuid.UUID, stargazers int32) *Crontab {
	return &Crontab{
		userId:     userId,
		stargazers: stargazers,
		createdAt:  time.Now(),
		updatedAt:  time.Now(),
	}
}

func FromCrontabEntity(entity *CrontabEntity) *Crontab {
	return &Crontab{
		userId:     entity.UserId,
		stargazers: entity.Stargazers,
		createdAt:  entity.CreatedAt,
		updatedAt:  entity.UpdatedAt,
	}
}

func (c *Crontab) UserId() uuid.UUID {
	return c.userId
}

func (c *Crontab) Stargazers() int32 {
	return c.stargazers
}

func (c *Crontab) CreatedAt() time.Time {
	return c.createdAt
}

func (c *Crontab) UpdatedAt() time.Time {
	return c.updatedAt
}

func (c *Crontab) ToCrontabDto() *CrontabDto {
	return &CrontabDto{
		Stargazers: c.stargazers,
		CreatedAt:  c.createdAt,
		UpdatedAt:  c.updatedAt,
	}
}

func (c *Crontab) ToCrontabEntity() *CrontabEntity {
	return &CrontabEntity{
		UserId:     c.userId,
		Stargazers: c.stargazers,
		CreatedAt:  c.createdAt,
		UpdatedAt:  c.updatedAt,
	}
}
