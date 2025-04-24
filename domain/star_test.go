package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	testUserID = "user-123"
	testRepoID = int32(987)
)

func TestStar_NewStar(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		before := time.Now().UTC()
		star, err := NewStar(testUserID, testRepoID)
		after := time.Now().UTC()

		require.NoError(t, err)
		require.NotNil(t, star)

		require.Equal(t, testUserID, star.userID) // Access internal for verification in test
		require.Equal(t, testRepoID, star.repoID) // Access internal for verification in test
		require.False(t, star.isDeleted)

		// Check timestamps are set and within the expected range
		require.False(t, star.createdAt.IsZero())
		require.False(t, star.updatedAt.IsZero())
		require.False(t, star.syncedAt.IsZero())

		require.True(t, !star.createdAt.Before(before) && !star.createdAt.After(after))
		require.True(t, !star.updatedAt.Before(before) && !star.updatedAt.After(after))
		require.True(t, !star.syncedAt.Before(before) && !star.syncedAt.After(after))

		// Check initial timestamps are the same
		require.Equal(t, star.createdAt, star.updatedAt)
		require.Equal(t, star.createdAt, star.syncedAt)

		// Verify getters
		require.Equal(t, testUserID, star.UserID())
		require.Equal(t, testRepoID, star.RepoID())
		require.False(t, star.IsDeleted())
		require.Equal(t, star.createdAt, star.CreatedAt())
		require.Equal(t, star.updatedAt, star.UpdatedAt())
		require.Equal(t, star.syncedAt, star.SyncedAt())
	})

	t.Run("Error_EmptyUserID", func(t *testing.T) {
		star, err := NewStar("", testRepoID)
		require.Error(t, err)
		require.Nil(t, star)
		require.Contains(t, err.Error(), "star userID cannot be empty")
	})

	t.Run("Error_ZeroRepoID", func(t *testing.T) {
		star, err := NewStar(testUserID, 0)
		require.Error(t, err)
		require.Nil(t, star)
		require.Contains(t, err.Error(), "star repoID must be positive")
	})

	t.Run("Error_NegativeRepoID", func(t *testing.T) {
		star, err := NewStar(testUserID, -10)
		require.Error(t, err)
		require.Nil(t, star)
		require.Contains(t, err.Error(), "star repoID must be positive")
	})
}

func TestStar_MarkAsDeleted(t *testing.T) {
	t.Run("DeleteActiveStar", func(t *testing.T) {
		star, _ := NewStar(testUserID, testRepoID)
		require.False(t, star.IsDeleted(), "Star should initially be active")
		initialUpdatedAt := star.UpdatedAt()

		// Introduce a slight delay to ensure time changes
		time.Sleep(1 * time.Millisecond)
		star.MarkAsDeleted()

		require.True(t, star.IsDeleted(), "Star should be marked as deleted")
		require.True(t, star.UpdatedAt().After(initialUpdatedAt), "UpdatedAt should be updated")
	})

	t.Run("DeleteAlreadyDeletedStar_Idempotency", func(t *testing.T) {
		star, _ := NewStar(testUserID, testRepoID)
		star.MarkAsDeleted() // First deletion
		require.True(t, star.IsDeleted(), "Star should be marked as deleted")
		updatedAtAfterFirstDelete := star.UpdatedAt()

		// Introduce a slight delay
		time.Sleep(1 * time.Millisecond)
		star.MarkAsDeleted() // Second deletion (should have no effect)

		require.True(t, star.IsDeleted(), "Star should still be marked as deleted")
		// Crucially, UpdatedAt should NOT have changed on the second call
		require.Equal(t, updatedAtAfterFirstDelete, star.UpdatedAt(), "UpdatedAt should not change on redundant delete")
	})
}

func TestStar_Restore(t *testing.T) {
	t.Run("RestoreDeletedStar", func(t *testing.T) {
		star, _ := NewStar(testUserID, testRepoID)
		star.MarkAsDeleted() // Mark as deleted first
		require.True(t, star.IsDeleted(), "Star should be deleted before restoring")
		updatedAtAfterDelete := star.UpdatedAt()

		// Introduce a slight delay
		time.Sleep(1 * time.Millisecond)
		star.Restore()

		require.False(t, star.IsDeleted(), "Star should be restored (not deleted)")
		require.True(t, star.UpdatedAt().After(updatedAtAfterDelete), "UpdatedAt should be updated on restore")
	})

	t.Run("RestoreActiveStar_Idempotency", func(t *testing.T) {
		star, _ := NewStar(testUserID, testRepoID)
		require.False(t, star.IsDeleted(), "Star should initially be active")
		initialUpdatedAt := star.UpdatedAt()

		// Introduce a slight delay
		time.Sleep(1 * time.Millisecond)
		star.Restore() // Restore an already active star (should have no effect)

		require.False(t, star.IsDeleted(), "Star should still be active")
		// Crucially, UpdatedAt should NOT have changed
		require.Equal(t, initialUpdatedAt, star.UpdatedAt(), "UpdatedAt should not change on redundant restore")
	})
}

func TestStar_UpdateSyncedAt(t *testing.T) {
	star, _ := NewStar(testUserID, testRepoID)
	initialSyncedAt := star.SyncedAt()
	initialUpdatedAt := star.UpdatedAt() // Capture initial UpdatedAt

	// Introduce a slight delay to ensure time progresses
	time.Sleep(1 * time.Millisecond)
	star.UpdateSyncedAt()
	newSyncedAt := star.SyncedAt()
	newUpdatedAt := star.UpdatedAt() // Capture UpdatedAt after the call

	require.True(t, newSyncedAt.After(initialSyncedAt), "SyncedAt should be updated to a later time")

	require.Equal(t, initialUpdatedAt, newUpdatedAt, "UpdatedAt should NOT be changed by UpdateSyncedAt")
}
