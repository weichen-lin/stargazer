package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNewFolder(t *testing.T) {
	name := "Test Folder"
	folder, err := NewFolder(name)
	require.NoError(t, err)

	parsedTime, err := ParseTime(time.Now().Format(time.RFC3339))
	require.NoError(t, err)

	require.NotNil(t, folder, "NewFolder should not return nil")
	require.NotNil(t, folder.id, "New folder should have a UUID")
	require.Equal(t, name, folder.Name(), "Folder name should match")
	require.False(t, folder.IsPublic(), "New folder should not be public by default")
	require.False(t, folder.CreatedAt().IsZero(), "CreatedAt should not be zero")
	require.False(t, folder.UpdatedAt().IsZero(), "UpdatedAt should be zero for a new folder")
	require.Equal(t, parsedTime, folder.CreatedAt())
	require.Equal(t, parsedTime, folder.UpdatedAt())

	folder2, err := NewFolder("")
	require.Error(t, err)
	require.Nil(t, folder2)
}

func TestFolder_SetId(t *testing.T) {
	folder, err := NewFolder("Test Folder")
	require.NoError(t, err)

	err = folder.SetId("invalid-id")
	require.Error(t, err)

	err = folder.SetId(uuid.New().String())
	require.NoError(t, err)
}

func TestFolder_SetName(t *testing.T) {
	folder, err := NewFolder("Initial Name")
	require.NoError(t, err)

	tests := []struct {
		name     string
		newName  string
		wantErr  bool
		expected string
	}{
		{"Valid new name", "New Name", false, "New Name"},
		{"Empty name", "", true, "New Name"},
		{"Same name", "New Name", false, "New Name"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := folder.SetName(tt.newName)
			if tt.wantErr {
				require.Error(t, err, "SetName() should return an error")
			} else {
				require.NoError(t, err, "SetName() should not return an error")
			}
			require.Equal(t, tt.expected, folder.Name(), "Folder name should match expected value")
		})
	}
}

func TestFolder_SetIsPublic(t *testing.T) {
	folder, err := NewFolder("Test Folder")
	require.NoError(t, err)

	folder.SetIsPublic(true)
	require.True(t, folder.IsPublic(), "Folder should be public")

	folder.SetIsPublic(false)
	require.False(t, folder.IsPublic(), "Folder should be private")
}

func TestFolder_SetCreatedAt(t *testing.T) {
	folder, err := NewFolder("Test Folder")
	require.NoError(t, err)

	validTime := time.Now().Format(time.RFC3339)
	err = folder.SetCreatedAt(validTime)
	require.NoError(t, err, "SetCreatedAt() should not return an error for valid time")

	err = folder.SetCreatedAt("")
	require.Error(t, err, "SetCreatedAt() should return an error for empty string")

	err = folder.SetCreatedAt("invalid time")
	require.Error(t, err, "SetCreatedAt() should return an error for invalid time")
}

func TestFolder_SetUpdatedAt(t *testing.T) {
	folder, err := NewFolder("Test Folder")
	require.NoError(t, err)

	validTime := time.Now().Format(time.RFC3339)
	err = folder.SetUpdatedAt(validTime)
	require.NoError(t, err, "SetUpdatedAt() should not return an error for valid time")

	err = folder.SetUpdatedAt("")
	require.Error(t, err, "SetUpdatedAt() should return an error for empty string")

	err = folder.SetUpdatedAt("invalid time")
	require.Error(t, err, "SetUpdatedAt() should return an error for invalid time")
}

func TestFromFolderEntity(t *testing.T) {
	now := time.Now().Format(time.RFC3339)
	entity := &FolderEntity{
		Id:        uuid.New().String(),
		Name:      "Test Folder",
		IsPublic:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	folder, err := FromFolderEntity(entity)
	require.NoError(t, err, "FromFolderEntity() should not return an error")
	require.NotNil(t, folder, "FromFolderEntity() should return a non-nil folder")

	require.Equal(t, entity.Name, folder.Name(), "Folder name should match")
	require.Equal(t, entity.IsPublic, folder.IsPublic(), "Folder IsPublic should match")
	require.Equal(t, entity.CreatedAt, folder.CreatedAt().Format(time.RFC3339), "Folder CreatedAt should match")
	require.Equal(t, entity.UpdatedAt, folder.UpdatedAt().Format(time.RFC3339), "Folder UpdatedAt should match")

	entity.Id = "invalid-id"
	_, err = FromFolderEntity(entity)
	require.Error(t, err, "FromFolderEntity() should return an error for invalid uuid")
	entity.Id = uuid.New().String()

	entity.Name = ""
	_, err = FromFolderEntity(entity)
	require.Error(t, err, "FromFolderEntity() should return an error for invalid name")
	entity.Name = "Test Folder"

	entity.CreatedAt = ""
	_, err = FromFolderEntity(entity)
	require.Error(t, err, "FromFolderEntity() should return an error for invalid createAt")
	entity.CreatedAt = now

	entity.UpdatedAt = ""
	_, err = FromFolderEntity(entity)
	require.Error(t, err, "FromFolderEntity() should return an error for invalid updateAt")
	entity.UpdatedAt = now
}

func TestToFolderEntity(t *testing.T) {
	folder, err := NewFolder("Test Folder")
	require.NoError(t, err)

	entity := folder.ToFolderEntity()

	require.Equal(t, folder.Id().String(), entity.Id, "Entity Id should match")
	require.Equal(t, folder.Name(), entity.Name, "Entity name should match")
	require.Equal(t, folder.IsPublic(), entity.IsPublic, "Entity IsPublic should match")
	require.Equal(t, folder.CreatedAt().Format(time.RFC3339), entity.CreatedAt, "Entity CreatedAt should match")
	require.Equal(t, folder.UpdatedAt().Format(time.RFC3339), entity.UpdatedAt, "Entity UpdatedAt should match")
}
