package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Folder struct {
	id        uuid.UUID
	name      string
	isPublic  bool
	createdAt time.Time
	updatedAt time.Time
}

type FolderEntity struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	IsPublic  bool   `json:"is_public"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (f *Folder) Id() uuid.UUID {
	return f.id
}

func (f *Folder) Name() string {
	return f.name
}

func (f *Folder) IsPublic() bool {
	return f.isPublic
}

func (f *Folder) CreatedAt() time.Time {
	return f.createdAt
}

func (f *Folder) UpdatedAt() time.Time {
	return f.updatedAt
}

func (f *Folder) SetId(id string) error {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	f.id = uuid
	return nil
}

func (f *Folder) SetName(name string) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}

	if name == f.name {
		return nil
	}

	f.name = name
	return nil
}

func (f *Folder) SetIsPublic(isPublic bool) {
	f.isPublic = isPublic
}

func (f *Folder) SetCreatedAt(s string) error {
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

func (f *Folder) SetUpdatedAt(s string) error {
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

func (f *Folder) ToFolderEntity() *FolderEntity {
	return &FolderEntity{
		Id:        f.id.String(),
		Name:      f.Name(),
		IsPublic:  f.IsPublic(),
		CreatedAt: f.CreatedAt().Format(time.RFC3339),
		UpdatedAt: f.UpdatedAt().Format(time.RFC3339),
	}
}

func FromFolderEntity(folderEntity *FolderEntity) (*Folder, error) {
	folder := &Folder{}

	err := folder.SetId(folderEntity.Id)
	if err != nil {
		return nil, err
	}

	err = folder.SetName(folderEntity.Name)
	if err != nil {
		return nil, err
	}

	err = folder.SetCreatedAt(folderEntity.CreatedAt)
	if err != nil {
		return nil, err
	}

	err = folder.SetUpdatedAt(folderEntity.UpdatedAt)
	if err != nil {
		return nil, err
	}

	folder.SetIsPublic(folderEntity.IsPublic)

	return folder, nil
}

func NewFolder(name string) (*Folder, error) {
	folder := &Folder{}
	folder.id = uuid.New()

	err := folder.SetName(name)
	if err != nil {
		return nil, err
	}

	folder.SetIsPublic(false)

	now := time.Now()
	folder.SetCreatedAt(now.Format(time.RFC3339))
	folder.SetUpdatedAt(now.Format(time.RFC3339))

	return folder, nil
}
