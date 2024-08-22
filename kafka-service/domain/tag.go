package domain

import (
	"errors"
	"time"
)

type Tag struct {
	name      string
	createdAt time.Time
	updatedAt time.Time
}

type TagEntity struct {
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
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

func (t *Tag) ToTagEntity() *TagEntity {

	var updatedAt string

	if t.UpdatedAt().IsZero() {
		updatedAt = ""
	} else {
		updatedAt = t.UpdatedAt().Format(time.RFC3339)
	}
	
	return &TagEntity{
		Name:      t.name,
		CreatedAt: t.createdAt.Format(time.RFC3339),
		UpdatedAt: updatedAt,
	}
}

func NewTag(name string) (*Tag, error) {
	tag := &Tag{}

	if err := tag.SetName(name); err != nil {
		return nil, err
	}

	now := time.Now()

	tag.SetCreatedAt(now.Format(time.RFC3339))
	tag.SetUpdatedAt("")

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

	return tag, nil
}