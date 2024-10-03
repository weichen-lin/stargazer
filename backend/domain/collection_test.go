package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNewCollection(t *testing.T) {
	name := "Test Collection"
	collection, err := NewCollection(name)
	require.NoError(t, err)

	parsedTime, err := ParseTime(time.Now().Format(time.RFC3339))
	require.NoError(t, err)

	require.NotNil(t, collection, "NewCollection should not return nil")
	require.NotNil(t, collection.id, "New Collection should have a UUID")
	require.Equal(t, name, collection.Name(), "Collection name should match")
	require.False(t, collection.IsPublic(), "New Collection should not be public by default")
	require.False(t, collection.CreatedAt().IsZero(), "CreatedAt should not be zero")
	require.False(t, collection.UpdatedAt().IsZero(), "UpdatedAt should be zero for a new collection")
	require.Equal(t, parsedTime, collection.CreatedAt())
	require.Equal(t, parsedTime, collection.UpdatedAt())

	collection2, err := NewCollection("")
	require.Error(t, err)
	require.Nil(t, collection2)
}

func TestCollection_SetId(t *testing.T) {
	collection, err := NewCollection("Test Collection")
	require.NoError(t, err)

	err = collection.SetId("invalid-id")
	require.Error(t, err)

	err = collection.SetId(uuid.New().String())
	require.NoError(t, err)
}

func TestCollection_SetName(t *testing.T) {
	collection, err := NewCollection("Initial Name")
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
			err := collection.SetName(tt.newName)
			if tt.wantErr {
				require.Error(t, err, "SetName() should return an error")
			} else {
				require.NoError(t, err, "SetName() should not return an error")
			}
			require.Equal(t, tt.expected, collection.Name(), "Collection name should match expected value")
		})
	}
}

func TestCollection_SetIsPublic(t *testing.T) {
	collection, err := NewCollection("Test Collection")
	require.NoError(t, err)

	collection.SetIsPublic(true)
	require.True(t, collection.IsPublic(), "Collection should be public")

	collection.SetIsPublic(false)
	require.False(t, collection.IsPublic(), "Collection should be private")
}

func TestCollection_SetCreatedAt(t *testing.T) {
	collection, err := NewCollection("Test Collection")
	require.NoError(t, err)

	validTime := time.Now().Format(time.RFC3339)
	err = collection.SetCreatedAt(validTime)
	require.NoError(t, err, "SetCreatedAt() should not return an error for valid time")

	err = collection.SetCreatedAt("")
	require.Error(t, err, "SetCreatedAt() should return an error for empty string")

	err = collection.SetCreatedAt("invalid time")
	require.Error(t, err, "SetCreatedAt() should return an error for invalid time")
}

func TestCollection_SetUpdatedAt(t *testing.T) {
	collection, err := NewCollection("Test Collection")
	require.NoError(t, err)

	validTime := time.Now().Format(time.RFC3339)
	err = collection.SetUpdatedAt(validTime)
	require.NoError(t, err, "SetUpdatedAt() should not return an error for valid time")

	err = collection.SetUpdatedAt("")
	require.Error(t, err, "SetUpdatedAt() should return an error for empty string")

	err = collection.SetUpdatedAt("invalid time")
	require.Error(t, err, "SetUpdatedAt() should return an error for invalid time")
}

func TestFromCollectionEntity(t *testing.T) {
	now := time.Now().Format(time.RFC3339)
	entity := &CollectionEntity{
		Id:        uuid.New().String(),
		Name:      "Test Collection",
		IsPublic:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	collection, err := FromCollectionEntity(entity)
	require.NoError(t, err, "FromCollectionEntity() should not return an error")
	require.NotNil(t, collection, "FromCollectionEntity() should return a non-nil collection")

	require.Equal(t, entity.Name, collection.Name(), "Collection name should match")
	require.Equal(t, entity.IsPublic, collection.IsPublic(), "Collection IsPublic should match")
	require.Equal(t, entity.CreatedAt, collection.CreatedAt().Format(time.RFC3339), "Collection CreatedAt should match")
	require.Equal(t, entity.UpdatedAt, collection.UpdatedAt().Format(time.RFC3339), "Collection UpdatedAt should match")

	entity.Id = "invalid-id"
	_, err = FromCollectionEntity(entity)
	require.Error(t, err, "FromCollectionEntity() should return an error for invalid uuid")
	entity.Id = uuid.New().String()

	entity.Name = ""
	_, err = FromCollectionEntity(entity)
	require.Error(t, err, "FromCollectionEntity() should return an error for invalid name")
	entity.Name = "Test Collection"

	entity.CreatedAt = ""
	_, err = FromCollectionEntity(entity)
	require.Error(t, err, "FromCollectionEntity() should return an error for invalid createAt")
	entity.CreatedAt = now

	entity.UpdatedAt = ""
	_, err = FromCollectionEntity(entity)
	require.Error(t, err, "FromCollectionEntity() should return an error for invalid updateAt")
	entity.UpdatedAt = now
}

func TestToCollectionrEntity(t *testing.T) {
	collection, err := NewCollection("Test Collection")
	require.NoError(t, err)

	entity := collection.ToCollectionEntity()

	require.Equal(t, collection.Id().String(), entity.Id, "Entity Id should match")
	require.Equal(t, collection.Name(), entity.Name, "Entity name should match")
	require.Equal(t, collection.IsPublic(), entity.IsPublic, "Entity IsPublic should match")
	require.Equal(t, collection.CreatedAt().Format(time.RFC3339), entity.CreatedAt, "Entity CreatedAt should match")
	require.Equal(t, collection.UpdatedAt().Format(time.RFC3339), entity.UpdatedAt, "Entity UpdatedAt should match")
}
