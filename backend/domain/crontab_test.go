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
	require.False(t, crontab.UpdatedAt().IsZero())
	require.WithinDuration(t, crontab.createdAt, time.Now(), time.Duration(2*time.Second))
	require.True(t, crontab.LastTriggeredAt().IsZero())

	updateAt := now.Add(time.Hour * 2).Format(time.RFC3339)
	crontab.SetUpdatedAt(updateAt)

	require.Equal(t, updateAt, crontab.UpdatedAt().Format(time.RFC3339))
}

func Test_ToCrontabEntityEmpty(t *testing.T) {
	Crontab := NewCrontab()

	entity := Crontab.ToCrontabEntity()
	require.Equal(t, entity.CreatedAt, Crontab.CreatedAt().Format(time.RFC3339))
	require.Equal(t, entity.TriggeredAt, "")
	require.Equal(t, entity.UpdatedAt, Crontab.UpdatedAt().Format(time.RFC3339))
	require.Equal(t, entity.Status, Crontab.Status())
	require.Equal(t, entity.LastTriggeredAt, "")

	crontabEmptyUpdatedAt := NewCrontab()
	crontabEmptyUpdatedAt.SetUpdatedAt("")
	crontabEmptyUpdatedAtEntity := crontabEmptyUpdatedAt.ToCrontabEntity()
	require.Equal(t, crontabEmptyUpdatedAtEntity.UpdatedAt, time.Now().Format(time.RFC3339))
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

	err = Crontab.SetUpdatedAt("")
	require.NoError(t, err)
	require.WithinDuration(t, Crontab.UpdatedAt(), time.Now(), time.Duration(2*time.Second))

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
	}

	Crontab, err := FromCrontabEntity(entity)
	require.NoError(t, err)

	require.Equal(t, "2023-01-02T00:00:00Z", Crontab.TriggeredAt().Format(time.RFC3339))
	require.Equal(t, "2023-01-01T00:00:00Z", Crontab.CreatedAt().Format(time.RFC3339))
	require.Equal(t, "2023-01-03T00:00:00Z", Crontab.UpdatedAt().Format(time.RFC3339))
	require.Equal(t, "2023-01-03T00:00:00Z", Crontab.LastTriggeredAt().Format(time.RFC3339))
	require.Equal(t, "active", Crontab.Status())
}

func TestFromCrontabEntity_ErrorCasess(t *testing.T) {
	validTime := "2023-01-01T00:00:00Z"
	testCases := []struct {
		name     string
		entity   *CrontabEntity
		expected string
	}{
		{
			name: "Invalid TriggeredAt",
			entity: &CrontabEntity{
				TriggeredAt: "invalid-time",
				Status:      "active",
			},
			expected: "parsing time",
		},
		{
			name: "Invalid CreatedAt",
			entity: &CrontabEntity{
				CreatedAt: "invalid_createdAt",
				Status:    "active",
			},
			expected: "parsing time",
		},
		{
			name: "Invalid UpdatedAt",
			entity: &CrontabEntity{
				CreatedAt: validTime,
				UpdatedAt: "invalid_updatedAt",
				Status:    "active",
			},
			expected: "parsing time",
		},
		{
			name: "Empty Status",
			entity: &CrontabEntity{
				CreatedAt: validTime,
				Status:    "",
			},
			expected: "invalid status: cannot be empty",
		},
		{
			name: "Invalid LastTriggeredAt",
			entity: &CrontabEntity{
				CreatedAt:       validTime,
				Status:          "active",
				LastTriggeredAt: "invalid_lastTriggeredAt",
			},
			expected: "parsing time",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := FromCrontabEntity(tc.entity)
			require.Error(t, err)
			require.Contains(t, err.Error(), tc.expected)
		})
	}
}