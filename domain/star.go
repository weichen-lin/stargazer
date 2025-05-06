package domain

import (
	"errors"
	"time"
)

type Star struct {
	userID string
	repoID int32

	isDeleted bool
	createdAt time.Time
	syncedAt  time.Time
	updatedAt time.Time
}

func NewStar(userID string, repoID int32) (*Star, error) {
	if userID == "" {
		return nil, errors.New("star userID cannot be empty")
	}

	if repoID <= 0 {
		return nil, errors.New("star repoID must be positive")
	}

	now := time.Now().UTC()

	star := &Star{
		userID:    userID,
		repoID:    repoID,
		isDeleted: false,
		createdAt: now,
		updatedAt: now,
		syncedAt:  now,
	}

	return star, nil
}

func (s *Star) UserID() string {
	return s.userID
}

func (s *Star) RepoID() int32 {
	return s.repoID
}

func (s *Star) IsDeleted() bool {
	return s.isDeleted
}

func (s *Star) CreatedAt() time.Time {
	return s.createdAt
}

func (s *Star) SyncedAt() time.Time {
	return s.syncedAt
}

func (s *Star) UpdatedAt() time.Time {
	return s.updatedAt
}

func (s *Star) MarkAsDeleted() {
	if !s.isDeleted {
		s.isDeleted = true
		s.updatedAt = time.Now().UTC()
	}
}

func (s *Star) Restore() {
	if s.isDeleted {
		s.isDeleted = false
		s.updatedAt = time.Now().UTC()
	}
}

func (s *Star) UpdateSyncedAt() {
	s.syncedAt = time.Now().UTC()
}
