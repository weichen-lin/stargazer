package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewCronTab(t *testing.T) {
	now := time.Now()

    cronTab, err := NewCronTab()
    require.NoError(t, err)

    require.Equal(t, "new", cronTab.Status())
    require.True(t, cronTab.TriggerAt().IsZero())
    require.False(t, cronTab.CreatedAt().IsZero())
	require.Equal(t, now.Sub(cronTab.createdAt).Milliseconds() > 0, true)
    require.True(t, cronTab.UpdatedAt().IsZero())
}

func Test_ToCronTabEntityEmpty(t *testing.T) {
    cronTab, err := NewCronTab()
    require.NoError(t, err)

    entity := cronTab.ToCronTabEntity()
	require.Equal(t, entity.CreatedAt, cronTab.CreatedAt().Format(time.RFC3339))
	require.Equal(t, entity.TriggerAt, "")
	require.Equal(t, entity.UpdatedAt, "")
	require.Equal(t, entity.Status, cronTab.Status())
}


func TestCronTab_SetTriggerAt(t *testing.T) {
	cronTab := &CronTab{}

	err := cronTab.SetTriggerAt("2023-01-01T00:00:00Z")
	require.NoError(t, err)
	require.Equal(t, "2023-01-01T00:00:00Z", cronTab.TriggerAt().Format(time.RFC3339))

	// Test empty string (should set to zero time)
	err = cronTab.SetTriggerAt("")
	require.NoError(t, err)
	require.Equal(t, cronTab.TriggerAt(), time.Time{})

	// Test invalid time format
	err = cronTab.SetTriggerAt("invalid-time")
	require.Error(t, err)
}

func TestCronTab_SetCreatedAt(t *testing.T) {
	cronTab := &CronTab{}

	// Test valid time
	err := cronTab.SetCreatedAt("2023-01-01T00:00:00Z")
	require.NoError(t, err)
	require.Equal(t, "2023-01-01T00:00:00Z", cronTab.CreatedAt().Format(time.RFC3339))

	// Test empty string (should return error)
	err = cronTab.SetCreatedAt("")
	require.Error(t, err)

	// Test invalid time format
	err = cronTab.SetCreatedAt("invalid-time")
	require.Error(t, err)
}

func TestCronTab_SetUpdatedAt(t *testing.T) {
	cronTab := &CronTab{}
	cronTab.SetCreatedAt("2023-01-01T00:00:00Z")

	// Test valid time
	err := cronTab.SetUpdatedAt("2023-01-02T00:00:00Z")
	require.NoError(t, err)
	require.Equal(t, "2023-01-02T00:00:00Z", cronTab.UpdatedAt().Format(time.RFC3339))

	// Test empty string (should set to zero time)
	err = cronTab.SetUpdatedAt("")
	require.NoError(t, err)
	require.Equal(t, cronTab.UpdatedAt(), time.Time{})

	// Test invalid time format
	err = cronTab.SetUpdatedAt("invalid-time")
	require.Error(t, err)

	// Test time before created time
	err = cronTab.SetUpdatedAt("2022-12-31T23:59:59Z")
	require.Error(t, err)
}

func TestCronTab_SetStatus(t *testing.T) {
	cronTab := &CronTab{}

	// Test valid status
	err := cronTab.SetStatus("active")
	require.NoError(t, err)
	require.Equal(t, "active", cronTab.Status())

	// Test empty string (should return error)
	err = cronTab.SetStatus("")
	require.Error(t, err)
}

func TestCronTab_ToCronTabEntity(t *testing.T) {
	cronTab := &CronTab{}
	cronTab.SetCreatedAt("2023-01-01T00:00:00Z")
	cronTab.SetTriggerAt("2023-01-02T00:00:00Z")
	cronTab.SetUpdatedAt("2023-01-03T00:00:00Z")
	cronTab.SetStatus("active")

	entity := cronTab.ToCronTabEntity()

	require.Equal(t, "2023-01-01T00:00:00Z", entity.CreatedAt)
	require.Equal(t, "2023-01-02T00:00:00Z", entity.TriggerAt)
	require.Equal(t, "2023-01-03T00:00:00Z", entity.UpdatedAt)
	require.Equal(t, "active", entity.Status)
}

func TestFromCronTabEntity(t *testing.T) {
	entity := &CronTabEntity{
		TriggerAt:     "2023-01-02T00:00:00Z",
		CreatedAt:     "2023-01-01T00:00:00Z",
		UpdatedAt:     "2023-01-03T00:00:00Z",
		Status:        "active",
	}

	cronTab, err := FromCronTabEntity(entity)
	require.NoError(t, err)

	require.Equal(t, "2023-01-02T00:00:00Z", cronTab.TriggerAt().Format(time.RFC3339))
	require.Equal(t, "2023-01-01T00:00:00Z", cronTab.CreatedAt().Format(time.RFC3339))
	require.Equal(t, "2023-01-03T00:00:00Z", cronTab.UpdatedAt().Format(time.RFC3339))
	require.Equal(t, "active", cronTab.Status())
}

func TestFromCronTabEntity_ErrorCases(t *testing.T) {

	// Test invalid times
	invalidEntity := &CronTabEntity{
		TriggerAt:     "invalid-time",
		CreatedAt:     "",
		UpdatedAt:     "",
		Status:        "active",
	}

	_, err := FromCronTabEntity(invalidEntity)
	require.Error(t, err)

	invalidEntity.TriggerAt = ""
	invalidEntity.CreatedAt = "invalid_createdAt"

	_, err = FromCronTabEntity(invalidEntity)
	require.Error(t, err)

	invalidEntity.CreatedAt = "2023-01-01T00:00:00Z"
	invalidEntity.UpdatedAt = "invalid_createdAt"

	_, err = FromCronTabEntity(invalidEntity)
	require.Error(t, err)

	invalidEntity.UpdatedAt = ""
	invalidEntity.Status = ""

	_, err = FromCronTabEntity(invalidEntity)
	require.Error(t, err)
}
