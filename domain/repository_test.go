package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createValidRequest() *CreateRepositoryRequest {
	return &CreateRepositoryRequest{
		Id:          123,
		Name:        "test-repo",
		OwnerName:   "test-owner",
		AvatarURL:   "https://example.com/avatar.png",
		HtmlURL:     "https://example.com/repo",
		Homepage:    "https://example.com",
		Description: "A test repository",
		CreatedAt:   "2023-01-15T10:00:00Z",
		UpdatedAt:   "2023-01-16T11:30:00Z",
		OpenIssues:  2,
		Language:    "Go",
		Archived:    false,
		Topics:      []string{"go", "test", "domain"},
	}
}

func TestNewRepository_Success(t *testing.T) {
	req := createValidRequest()
	repo, err := NewRepository(req)

	require.NoError(t, err, "NewRepository should not return an error with valid data")
	require.NotNil(t, repo, "NewRepository should return a non-nil repository with valid data")

	require.Equal(t, req.Id, repo.Id(), "ID mismatch")
	require.Equal(t, req.Name, repo.Name(), "Name mismatch")
	require.Equal(t, req.OwnerName, repo.OwnerName(), "OwnerName mismatch")
	require.Equal(t, req.AvatarURL, repo.AvatarURL(), "AvatarURL mismatch")
	require.Equal(t, req.Watchers, repo.Watchers(), "Watchers mismatch")
	require.Equal(t, req.HtmlURL, repo.HTMLURL(), "HTMLURL mismatch")
	require.Equal(t, req.Homepage, repo.Homepage(), "Homepage mismatch")
	require.Equal(t, req.Description, repo.Description(), "Description mismatch")
	require.Equal(t, req.OpenIssues, repo.OpenIssues(), "OpenIssues mismatch")
	require.Equal(t, req.Language, repo.Language(), "Language mismatch")
	require.Equal(t, req.Archived, repo.Archived(), "Archived mismatch")
	require.Equal(t, req.Topics, repo.Topics(), "Topics mismatch")

	expectedCreatedAt, _ := time.Parse(time.RFC3339, req.CreatedAt)
	expectedUpdatedAt, _ := time.Parse(time.RFC3339, req.UpdatedAt)
	require.WithinDuration(t, expectedCreatedAt.UTC(), repo.CreatedAt(), 0, "CreatedAt mismatch or not UTC")
	require.WithinDuration(t, expectedUpdatedAt.UTC(), repo.UpdatedAt(), 0, "UpdatedAt mismatch or not UTC")
}

func TestNewRepository_ValidationErrors(t *testing.T) {
	testCases := []struct {
		name        string
		modifier    func(req *CreateRepositoryRequest)
		expectedErr string // Substring of the expected error message
	}{
		{
			name: "Invalid ID (Zero)",
			modifier: func(req *CreateRepositoryRequest) {
				req.Id = 0
			},
			expectedErr: "repository ID must be positive",
		},
		{
			name: "Invalid ID (Negative)",
			modifier: func(req *CreateRepositoryRequest) {
				req.Id = -10
			},
			expectedErr: "repository ID must be positive",
		},
		{
			name: "Empty Name",
			modifier: func(req *CreateRepositoryRequest) {
				req.Name = ""
			},
			expectedErr: "repository name cannot be empty",
		},
		{
			name: "Empty Owner Name",
			modifier: func(req *CreateRepositoryRequest) {
				req.OwnerName = ""
			},
			expectedErr: "owner name cannot be empty",
		},
		{
			name: "Invalid Avatar URL (No Scheme)",
			modifier: func(req *CreateRepositoryRequest) {
				req.AvatarURL = "example.com/avatar.png"
			},
			expectedErr: "URL must start with http:// or https://",
		},
		{
			name: "Invalid Avatar URL (Wrong Scheme)",
			modifier: func(req *CreateRepositoryRequest) {
				req.AvatarURL = "ftp://example.com/avatar.png"
			},
			expectedErr: "URL must start with http:// or https://",
		},
		{
			name: "Invalid HTML URL",
			modifier: func(req *CreateRepositoryRequest) {
				req.HtmlURL = "invalid-url"
			},
			expectedErr: "URL must start with http:// or https://",
		},
		{
			name: "Invalid Homepage URL (Non-empty)",
			modifier: func(req *CreateRepositoryRequest) {
				req.Homepage = "badsite" // Non-empty but invalid
			},
			expectedErr: "URL must start with http:// or https://",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := createValidRequest() // Start with a valid request
			tc.modifier(req)            // Apply the modification to make it invalid

			repo, err := NewRepository(req)

			require.Error(t, err, "NewRepository should return an error for: "+tc.name)
			require.Nil(t, repo, "Repository object should be nil on validation error for: "+tc.name)
			require.Contains(t, err.Error(), tc.expectedErr, "Error message mismatch for: "+tc.name)
		})
	}
}

func TestNewRepository_OptionalFields(t *testing.T) {
	t.Run("Empty Homepage", func(t *testing.T) {
		req := createValidRequest()
		req.Homepage = "" // Empty homepage is allowed

		repo, err := NewRepository(req)
		require.NoError(t, err)
		require.NotNil(t, repo)
		require.Equal(t, "", repo.Homepage())
	})

	t.Run("Empty Description", func(t *testing.T) {
		req := createValidRequest()
		req.Description = ""

		repo, err := NewRepository(req)
		require.NoError(t, err)
		require.NotNil(t, repo)
		require.Equal(t, "", repo.Description())
	})

	t.Run("Empty Language", func(t *testing.T) {
		req := createValidRequest()
		req.Language = ""

		repo, err := NewRepository(req)
		require.NoError(t, err)
		require.NotNil(t, repo)
		require.Equal(t, "Unknown", repo.Language()) // Should default to "Unknown"
	})

	t.Run("Nil Topics", func(t *testing.T) {
		req := createValidRequest()
		req.Topics = nil // Explicitly nil

		repo, err := NewRepository(req)
		require.NoError(t, err)
		require.NotNil(t, repo)
		require.Nil(t, repo.Topics()) // Should be nil, not an empty slice
	})

	t.Run("Empty Topics", func(t *testing.T) {
		req := createValidRequest()
		req.Topics = []string{} // Empty slice

		repo, err := NewRepository(req)
		require.NoError(t, err)
		require.NotNil(t, repo)
		require.NotNil(t, repo.Topics()) // Should be an empty slice, not nil
		require.Len(t, repo.Topics(), 0)
	})
}

func TestRepository_ToRepositoryEntity(t *testing.T) {
	req := createValidRequest()
	repo, err := NewRepository(req)
	require.NoError(t, err)
	require.NotNil(t, repo)

	entity := repo.ToRepositoryEntity()
	require.NotNil(t, entity)

	// Verify mapping from Repository to RepositoryEntity
	require.Equal(t, repo.Id(), entity.Id)
	require.Equal(t, repo.Name(), entity.Name)
	require.Equal(t, repo.OwnerName(), entity.OwnerName)
	require.Equal(t, repo.AvatarURL(), entity.AvatarURL)
	require.Equal(t, repo.HTMLURL(), entity.HtmlURL)
	require.Equal(t, repo.Homepage(), entity.Homepage)
	require.Equal(t, repo.Description(), entity.Description)
	require.Equal(t, repo.OpenIssues(), entity.OpenIssues)
	require.Equal(t, repo.Language(), entity.Language)
	require.Equal(t, repo.Archived(), entity.Archived)
	require.Equal(t, repo.Topics(), entity.Topics)

	// Verify time zones specifically (should be UTC)
	require.Equal(t, repo.CreatedAt().Location(), time.UTC)
	require.Equal(t, repo.UpdatedAt().Location(), time.UTC)
	require.Equal(t, entity.CreatedAt.Location(), time.UTC)
	require.Equal(t, entity.UpdatedAt.Location(), time.UTC)
	require.True(t, repo.CreatedAt().Equal(entity.CreatedAt))
	require.True(t, repo.UpdatedAt().Equal(entity.UpdatedAt))

}
