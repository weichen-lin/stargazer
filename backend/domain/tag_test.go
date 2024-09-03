package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTag_SetName(t *testing.T) {
	tag := &Tag{}

	// Test empty name
	err := tag.SetName("")
	require.Error(t, err)
	require.EqualError(t, err, "name cannot be empty")

	// Test valid name
	err = tag.SetName("Golang")
	require.NoError(t, err)
	require.Equal(t, "Golang", tag.Name())

	// Test setting the same name
	err = tag.SetName("Golang")
	require.NoError(t, err)
}

func TestTag_SetCreatedAt(t *testing.T) {
	tag := &Tag{}

	// Test empty created time
	err := tag.SetCreatedAt("")
	require.Error(t, err)
	require.EqualError(t, err, "created time cannot be empty")

	// Test invalid time format
	err = tag.SetCreatedAt("invalid-time")
	require.Error(t, err)

	// Test valid created time
	validTime := time.Now().Format(time.RFC3339)
	err = tag.SetCreatedAt(validTime)
	require.NoError(t, err)

	parsedTime, _ := time.Parse(time.RFC3339, validTime)
	require.Equal(t, parsedTime, tag.CreatedAt())
}

func TestTag_SetUpdatedAt(t *testing.T) {
	tag := &Tag{}

	// Test setting updated time before created time
	createdTime := time.Now().Add(1 * time.Hour).Format(time.RFC3339)
	err := tag.SetCreatedAt(createdTime)
	require.NoError(t, err)

	updatedTime := time.Now().Format(time.RFC3339)
	err = tag.SetUpdatedAt(updatedTime)
	require.Error(t, err)
	require.EqualError(t, err, "updated time cannot be before created time")

	// Test valid updated time
	err = tag.SetUpdatedAt(createdTime)
	require.NoError(t, err)

	parsedTime, _ := time.Parse(time.RFC3339, createdTime)
	require.Equal(t, parsedTime, tag.UpdatedAt())

	err = tag.SetUpdatedAt("")
	require.NoError(t, err)
}

func TestNewTag(t *testing.T) {
	// Test empty name
	_, err := NewTag("")
	require.Error(t, err)
	require.EqualError(t, err, "name cannot be empty")

	// Test valid creation
	tag, err := NewTag("Golang")
	require.NoError(t, err)
	require.Equal(t, "Golang", tag.Name())
	require.NotZero(t, tag.CreatedAt())
	require.Zero(t, tag.UpdatedAt())
}

func TestFromTagEntity(t *testing.T) {
	nowString := time.Now().Format(time.RFC3339)

	tagEntity := &TagEntity{
		Name:      "Golang",
		CreatedAt: nowString,
		UpdatedAt: time.Now().Add(1 * time.Hour).Format(time.RFC3339),
	}

	// Test valid TagEntity
	tag, err := FromTagEntity(tagEntity)
	require.NoError(t, err)
	require.Equal(t, tagEntity.Name, tag.Name())
	require.Equal(t, tagEntity.CreatedAt, tag.CreatedAt().Format(time.RFC3339))
	require.Equal(t, tagEntity.UpdatedAt, tag.UpdatedAt().Format(time.RFC3339))

	// Test invalid TagEntity (empty name)
	tagEntity.Name = ""
	_, err = FromTagEntity(tagEntity)
	require.Error(t, err)
	require.EqualError(t, err, "name cannot be empty")
	tagEntity.Name = "test"

	tagEntity.CreatedAt = ""
	_, err = FromTagEntity(tagEntity)
	require.Error(t, err)
	tagEntity.CreatedAt = nowString

	tagEntity.UpdatedAt = time.Now().Add(-2 * time.Hour).Format(time.RFC3339)
	_, err = FromTagEntity(tagEntity)
	require.Error(t, err)

	tagEntity.UpdatedAt = "invalid-date"
	_, err = FromTagEntity(tagEntity)
	require.Error(t, err)
}

func TestToTagEntity(t *testing.T) {
	tag, err := NewTag("test")
	require.NoError(t, err)

	entity := tag.ToTagEntity()

	require.Equal(t, entity.Name, tag.Name())

	now := time.Now()
	now.Add(time.Hour * 1)
	updatedAt := now.Format(time.RFC3339)
	tag.SetUpdatedAt(updatedAt)

	entity2 := tag.ToTagEntity()
	require.Equal(t, entity2.UpdatedAt, updatedAt)
}
