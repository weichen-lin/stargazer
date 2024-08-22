package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewCrontab(t *testing.T) {
	now := time.Now()

	crontab := NewCrontab()

	require.Equal(t, "new", crontab.Status())
	require.True(t, crontab.TriggeredAt().IsZero())
	require.False(t, crontab.CreatedAt().IsZero())
	require.Equal(t, now.Sub(crontab.createdAt).Milliseconds() > 0, true)
	require.True(t, crontab.UpdatedAt().IsZero())
	require.True(t, crontab.LastTriggeredAt().IsZero())
	require.Equal(t, crontab.Version(), 1)
}

func Test_ToCrontabEntityEmpty(t *testing.T) {
	Crontab := NewCrontab()

	entity := Crontab.ToCrontabEntity()
	require.Equal(t, entity.CreatedAt, Crontab.CreatedAt().Format(time.RFC3339))
	require.Equal(t, entity.TriggeredAt, "")
	require.Equal(t, entity.UpdatedAt, "")
	require.Equal(t, entity.Status, Crontab.Status())
	require.Equal(t, entity.Version, Crontab.Version())
	require.Equal(t, entity.LastTriggeredAt, "")
}

func TestCrontab_SetTriggerAt(t *testing.T) {
	Crontab := &Crontab{}

	err := Crontab.SetTriggeredAt("2023-01-01T00:00:00Z")
	require.NoError(t, err)
	require.Equal(t, "2023-01-01T00:00:00Z", Crontab.TriggeredAt().Format(time.RFC3339))

	// Test empty string (should set to zero time)
	err = Crontab.SetTriggeredAt("")
	require.NoError(t, err)
	require.Equal(t, Crontab.TriggeredAt(), time.Time{})

	// Test invalid time format
	err = Crontab.SetTriggeredAt("invalid-time")
	require.Error(t, err)
}

func TestCrontab_SetCreatedAt(t *testing.T) {
	Crontab := &Crontab{}

	// Test valid time
	err := Crontab.SetCreatedAt("2023-01-01T00:00:00Z")
	require.NoError(t, err)
	require.Equal(t, "2023-01-01T00:00:00Z", Crontab.CreatedAt().Format(time.RFC3339))

	// Test empty string (should return error)
	err = Crontab.SetCreatedAt("")
	require.Error(t, err)

	// Test invalid time format
	err = Crontab.SetCreatedAt("invalid-time")
	require.Error(t, err)
}

func TestCrontab_SetUpdatedAt(t *testing.T) {
	Crontab := &Crontab{}
	Crontab.SetCreatedAt("2023-01-01T00:00:00Z")

	// Test valid time
	err := Crontab.SetUpdatedAt("2023-01-02T00:00:00Z")
	require.NoError(t, err)
	require.Equal(t, "2023-01-02T00:00:00Z", Crontab.UpdatedAt().Format(time.RFC3339))

	// Test empty string (should set to zero time)
	err = Crontab.SetUpdatedAt("")
	require.NoError(t, err)
	require.Equal(t, Crontab.UpdatedAt(), time.Time{})

	// Test invalid time format
	err = Crontab.SetUpdatedAt("invalid-time")
	require.Error(t, err)

	// Test time before created time
	err = Crontab.SetUpdatedAt("2022-12-31T23:59:59Z")
	require.Error(t, err)
}

func TestCrontab_SetLastTriggerAt(t *testing.T) {
	Crontab := &Crontab{}
	Crontab.SetLastTriggerAt("2023-01-01T00:00:00Z")
	Crontab.SetCreatedAt("2023-01-01T00:00:00Z")
	// Test valid time
	err := Crontab.SetLastTriggerAt("2023-01-02T00:00:00Z")
	require.NoError(t, err)
	require.Equal(t, "2023-01-02T00:00:00Z", Crontab.LastTriggeredAt().Format(time.RFC3339))

	// Test empty string (should set to zero time)
	err = Crontab.SetLastTriggerAt("")
	require.NoError(t, err)
	require.Equal(t, Crontab.LastTriggeredAt(), time.Time{})

	// Test invalid time format
	err = Crontab.SetLastTriggerAt("invalid-time")
	require.Error(t, err)

	// Test time before created time
	err = Crontab.SetLastTriggerAt("2022-12-31T23:59:59Z")
	require.Error(t, err)
}

func TestCrontab_SetStatus(t *testing.T) {
	Crontab := &Crontab{}

	// Test valid status
	err := Crontab.SetStatus("active")
	require.NoError(t, err)
	require.Equal(t, "active", Crontab.Status())

	// Test empty string (should return error)
	err = Crontab.SetStatus("")
	require.Error(t, err)
}

func TestCrontab_ToCrontabEntity(t *testing.T) {
	Crontab := &Crontab{}
	Crontab.SetCreatedAt("2023-01-01T00:00:00Z")
	Crontab.SetTriggeredAt("2023-01-02T00:00:00Z")
	Crontab.SetUpdatedAt("2023-01-03T00:00:00Z")
	Crontab.SetStatus("active")
	Crontab.SetLastTriggerAt("2023-01-03T00:00:00Z")

	entity := Crontab.ToCrontabEntity()

	require.Equal(t, "2023-01-01T00:00:00Z", entity.CreatedAt)
	require.Equal(t, "2023-01-02T00:00:00Z", entity.TriggeredAt)
	require.Equal(t, "2023-01-03T00:00:00Z", entity.UpdatedAt)
	require.Equal(t, "active", entity.Status)
	require.Equal(t, "2023-01-03T00:00:00Z", entity.LastTriggeredAt)
}

func TestFromCrontabEntity(t *testing.T) {
	entity := &CrontabEntity{
		TriggeredAt:     "2023-01-02T00:00:00Z",
		CreatedAt:       "2023-01-01T00:00:00Z",
		UpdatedAt:       "2023-01-03T00:00:00Z",
		LastTriggeredAt: "2023-01-03T00:00:00Z",
		Status:          "active",
		Version:         12,
	}

	Crontab, err := FromCrontabEntity(entity)
	require.NoError(t, err)

	require.Equal(t, "2023-01-02T00:00:00Z", Crontab.TriggeredAt().Format(time.RFC3339))
	require.Equal(t, "2023-01-01T00:00:00Z", Crontab.CreatedAt().Format(time.RFC3339))
	require.Equal(t, "2023-01-03T00:00:00Z", Crontab.UpdatedAt().Format(time.RFC3339))
	require.Equal(t, "2023-01-03T00:00:00Z", Crontab.LastTriggeredAt().Format(time.RFC3339))
	require.Equal(t, "active", Crontab.Status())
	require.Equal(t, 12, Crontab.Version())
}

func TestFromCrontabEntity_ErrorCases(t *testing.T) {

	// Test invalid times
	invalidEntity := &CrontabEntity{
		TriggeredAt:     "invalid-time",
		CreatedAt:       "",
		UpdatedAt:       "",
		Status:          "active",
		LastTriggeredAt: "",
	}

	_, err := FromCrontabEntity(invalidEntity)
	require.Error(t, err)

	invalidEntity.TriggeredAt = ""
	invalidEntity.CreatedAt = "invalid_createdAt"

	_, err = FromCrontabEntity(invalidEntity)
	require.Error(t, err)

	invalidEntity.CreatedAt = "2023-01-01T00:00:00Z"
	invalidEntity.UpdatedAt = "invalid_createdAt"

	_, err = FromCrontabEntity(invalidEntity)
	require.Error(t, err)

	invalidEntity.UpdatedAt = ""
	invalidEntity.Status = ""

	_, err = FromCrontabEntity(invalidEntity)
	require.Error(t, err)

	invalidEntity.Status = "status"
	invalidEntity.LastTriggeredAt = "invalid_createdAt"

	_, err = FromCrontabEntity(invalidEntity)
	require.Error(t, err)
}
