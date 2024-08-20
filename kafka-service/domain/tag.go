package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Tag struct {
	id        uuid.UUID
	name      string
	createdAt time.Time
	updatedAt time.Time
}

type TagEntity struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (t *Tag) ID() uuid.UUID {
	return t.id
}

func (t *Tag) Name() string {
	return t.name
}

func (t *Tag) CreatedAt() time.Time {
	return t.createdAt
}

func (t *Tag) UpdatedAt() time.Time {
	return t.updatedAt
}

func (t *Tag) SetId(id string) error {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	t.id = uuid
	return nil
}

func (t *Tag) SetName(name string) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}

	if name == t.name {
		return nil
	}

	t.name = name
	return nil
}

func (t *Tag) SetCreatedAt(s string) error {
	if s == "" {
		return errors.New("created time cannot be empty")
	}

	parsedTime, err := ParseTime(s)
	if err != nil {
		return err
	}

	t.createdAt = parsedTime
	return nil
}

func (t *Tag) SetUpdatedAt(s string) error {
	if s == "" {
		t.updatedAt = time.Time{}
		return nil
	}

	parsedTime, err := ParseTime(s)
	if err != nil {
		return err
	}

	if parsedTime.Before(t.createdAt) {
		return errors.New("updated time cannot be before created time")
	}

	t.updatedAt = parsedTime
	return nil
}

func NewTag(name string) (*Tag, error) {
	tag := &Tag{}

	id, _ := uuid.NewUUID()

	if err := tag.SetName(name); err != nil {
		return nil, err
	}

	now := time.Now()

	tag.SetCreatedAt(now.Format(time.RFC3339))
	tag.SetUpdatedAt("")
	tag.SetId(id.String())

	return tag, nil
}

func FromTagEntity(tagEntity *TagEntity) (*Tag, error) {
	tag := &Tag{}

	if err := tag.SetName(tagEntity.Name); err != nil {
		return nil, err
	}

	if err := tag.SetCreatedAt(tagEntity.CreatedAt); err != nil {
		return nil, err
	}

	if err := tag.SetUpdatedAt(tagEntity.UpdatedAt); err != nil {
		return nil, err
	}

	if err := tag.SetId(tagEntity.ID); err != nil {
		return nil, err
	}

	return tag, nil
}
