package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id        uuid.UUID `json:"id"`
	ClerkId   string    `json:"clerk_id"`
	CreatedAt time.Time `json:"created_at"`
}

func NewUser(id, clerkId string, createdAt time.Time) (*User, error) {
	if id == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}

	if clerkId == "" {
		return nil, fmt.Errorf("clerkId cannot be empty")
	}

	if createdAt.IsZero() {
		return nil, fmt.Errorf("createdAt cannot be zero")
	}

	// Parse the UUID from the string
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID format: %v", err)
	}

	return &User{
		Id:        parsedId,
		ClerkId:   clerkId,
		CreatedAt: createdAt,
	}, nil
}
