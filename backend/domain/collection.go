package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Collection struct {
	id        uuid.UUID
	name      string
	isPublic  bool
	createdAt time.Time
	updatedAt time.Time
}

type CollectionEntity struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	IsPublic  bool   `json:"is_public"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (f *Collection) Id() uuid.UUID {
	return f.id
}

func (f *Collection) Name() string {
	return f.name
}

func (f *Collection) IsPublic() bool {
	return f.isPublic
}

func (f *Collection) CreatedAt() time.Time {
	return f.createdAt
}

func (f *Collection) UpdatedAt() time.Time {
	return f.updatedAt
}

func (f *Collection) SetId(id string) error {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	f.id = uuid
	return nil
}

func (f *Collection) SetName(name string) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}

	if name == f.name {
		return nil
	}

	f.name = name
	return nil
}

func (f *Collection) SetIsPublic(isPublic bool) {
	f.isPublic = isPublic
}

func (f *Collection) SetCreatedAt(s string) error {
	if s == "" {
		return errors.New("created time cannot be empty")
	}

	parsedTime, err := ParseTime(s)
	if err != nil {
		return err
	}

	f.createdAt = parsedTime
	return nil
}

func (f *Collection) SetUpdatedAt(s string) error {
	if s == "" {
		return errors.New("updated time cannot be empty")
	}

	parsedTime, err := ParseTime(s)
	if err != nil {
		return err
	}

	f.updatedAt = parsedTime
	return nil
}

func (f *Collection) ToCollectionEntity() *CollectionEntity {
	return &CollectionEntity{
		Id:        f.id.String(),
		Name:      f.Name(),
		IsPublic:  f.IsPublic(),
		CreatedAt: f.CreatedAt().Format(time.RFC3339),
		UpdatedAt: f.UpdatedAt().Format(time.RFC3339),
	}
}

func FromCollectionEntity(entity *CollectionEntity) (*Collection, error) {
	collection := &Collection{}

	err := collection.SetId(entity.Id)
	if err != nil {
		return nil, err
	}

	err = collection.SetName(entity.Name)
	if err != nil {
		return nil, err
	}

	err = collection.SetCreatedAt(entity.CreatedAt)
	if err != nil {
		return nil, err
	}

	err = collection.SetUpdatedAt(entity.UpdatedAt)
	if err != nil {
		return nil, err
	}

	collection.SetIsPublic(entity.IsPublic)

	return collection, nil
}

func NewCollection(name string) (*Collection, error) {
	collection := &Collection{}
	collection.id = uuid.New()

	err := collection.SetName(name)
	if err != nil {
		return nil, err
	}

	collection.SetIsPublic(false)

	now := time.Now()
	collection.SetCreatedAt(now.Format(time.RFC3339))
	collection.SetUpdatedAt(now.Format(time.RFC3339))

	return collection, nil
}
