package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewRepository(t *testing.T) {
	createdAt := time.Now().Add(-24 * time.Hour)
	updatedAt := time.Now()

	githubRepo := &GithubRepository{
		ID:              12345,
		Name:            "sample-repo",
		FullName:        "user/sample-repo",
		Owner:           Owner{Login: "user", AvatarURL: "https://avatar.url", HTMLURL: "https://github.com/user"},
		HTMLURL:         "https://github.com/user/sample-repo",
		Description:     "This is a sample repository",
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
		Homepage:        "https://sample-repo.com",
		StargazersCount: 10,
		WatchersCount:   20,
		Language:        "Go",
		Archived:        false,
		OpenIssuesCount: 5,
		Topics:          []string{"go", "sample"},
		DefaultBranch:   "main",
	}

	repo, err := NewRepository(githubRepo)
	require.NoError(t, err)

	require.Equal(t, int64(12345), repo.RepoID())
	require.Equal(t, "sample-repo", repo.RepoName())
	require.Equal(t, "user", repo.OwnerName())
	require.Equal(t, "https://avatar.url", repo.AvatarURL())
	require.Equal(t, "https://github.com/user/sample-repo", repo.HTMLURL())
	require.Equal(t, "https://sample-repo.com", repo.Homepage())
	require.Equal(t, "This is a sample repository", repo.Description())
	require.Equal(t, createdAt, repo.CreatedAt())
	require.Equal(t, updatedAt, repo.UpdatedAt())
	require.Equal(t, 10, repo.StargazersCount())
	require.Equal(t, 20, repo.WatchersCount())
	require.Equal(t, 5, repo.OpenIssuesCount())
	require.Equal(t, "Go", repo.Language())
	require.Equal(t, "main", repo.DefaultBranch())
	require.False(t, repo.IsArchived())
}

func TestRepository_Setters(t *testing.T) {
	repo := &Repository{}

	err := repo.setRepoID(12345)
	require.NoError(t, err)
	require.Equal(t, int64(12345), repo.RepoID())

	err = repo.setRepoName("sample-repo")
	require.NoError(t, err)
	require.Equal(t, "sample-repo", repo.RepoName())

	err = repo.setOwnerName("user")
	require.NoError(t, err)
	require.Equal(t, "user", repo.OwnerName())

	err = repo.setAvatarURL("https://avatar.url")
	require.NoError(t, err)
	require.Equal(t, "https://avatar.url", repo.AvatarURL())

	err = repo.setHTMLURL("https://github.com/user/sample-repo")
	require.NoError(t, err)
	require.Equal(t, "https://github.com/user/sample-repo", repo.HTMLURL())

	err = repo.setHomepage("https://sample-repo.com")
	require.NoError(t, err)
	require.Equal(t, "https://sample-repo.com", repo.Homepage())

	repo.setDescription("This is a sample repository")
	require.Equal(t, "This is a sample repository", repo.Description())

	createdAt := time.Now().Add(-24 * time.Hour)
	err = repo.setCreatedAt(createdAt)
	require.NoError(t, err)
	require.Equal(t, createdAt, repo.CreatedAt())

	updatedAt := time.Now()
	err = repo.setUpdatedAt(updatedAt)
	require.NoError(t, err)
	require.Equal(t, updatedAt, repo.UpdatedAt())

	err = repo.setStargazersCount(10)
	require.NoError(t, err)
	require.Equal(t, 10, repo.StargazersCount())

	err = repo.setWatchersCount(20)
	require.NoError(t, err)
	require.Equal(t, 20, repo.WatchersCount())

	err = repo.setOpenIssuesCount(5)
	require.NoError(t, err)
	require.Equal(t, 5, repo.OpenIssuesCount())

	repo.setLanguage("Go")
	require.Equal(t, "Go", repo.Language())

	err = repo.setDefaultBranch("main")
	require.NoError(t, err)
	require.Equal(t, "main", repo.DefaultBranch())

	repo.setArchived(false)
	require.False(t, repo.IsArchived())
}
func TestRepository_Setters_ErrorCases(t *testing.T) {
	repo := Repository{}

	// Test negative repository ID
	err := repo.setRepoID(-1)
	require.Error(t, err)
	require.EqualError(t, err, "repository ID must be positive")

	// Test empty repository name
	err = repo.setRepoName("")
	require.Error(t, err)
	require.EqualError(t, err, "repository name cannot be empty")

	// Test empty owner name
	err = repo.setOwnerName("")
	require.Error(t, err)
	require.EqualError(t, err, "owner name cannot be empty")

	// Test invalid URL
	err = repo.setAvatarURL("invalid-url")
	require.Error(t, err)
	require.EqualError(t, err, "URL must start with http:// or https://")

	// Test invalid HTML URL
	err = repo.setHTMLURL("invalid-url")
	require.Error(t, err)
	require.EqualError(t, err, "URL must start with http:// or https://")

	// Test invalid homepage URL
	err = repo.setHomepage("invalid-url")
	require.Error(t, err)
	require.EqualError(t, err, "URL must start with http:// or https://")

	// Test zero created time
	err = repo.setCreatedAt(time.Time{})
	require.Error(t, err)
	require.EqualError(t, err, "created time cannot be empty")

	// Test zero updated time
	err = repo.setUpdatedAt(time.Time{})
	require.Error(t, err)
	require.EqualError(t, err, "updated time cannot be empty")

	// Test updated time before created time
	createdAt := time.Now()
	updatedAt := createdAt.Add(-1 * time.Hour)
	repo.setCreatedAt(createdAt)
	err = repo.setUpdatedAt(updatedAt)
	require.Error(t, err)
	require.EqualError(t, err, "updated time cannot be before created time")

	// Test negative stargazers count
	err = repo.setStargazersCount(-1)
	require.Error(t, err)
	require.EqualError(t, err, "stargazers count cannot be negative")

	// Test negative watchers count
	err = repo.setWatchersCount(-1)
	require.Error(t, err)
	require.EqualError(t, err, "watchers count cannot be negative")

	// Test negative open issues count
	err = repo.setOpenIssuesCount(-1)
	require.Error(t, err)
	require.EqualError(t, err, "open issues count cannot be negative")

	// Test empty default branch
	err = repo.setDefaultBranch("")
	require.Error(t, err)
	require.EqualError(t, err, "default branch cannot be empty")

	repo.setLanguage("")
	require.Equal(t, "Unknown", repo.Language())
}

func TestNewRepository_ErrorCases(t *testing.T) {
	// Set up a valid base GithubRepository object
	validGithubRepo := &GithubRepository{
		ID:              1,
		Name:            "valid-repo",
		Owner:           Owner{Login: "valid-owner", AvatarURL: "https://valid-url.com/avatar", HTMLURL: "https://valid-url.com/owner"},
		HTMLURL:         "https://valid-url.com/repo",
		Description:     "A valid description",
		CreatedAt:       time.Now().Add(-24 * time.Hour),
		UpdatedAt:       time.Now(),
		Homepage:        "https://valid-url.com",
		StargazersCount: 100,
		WatchersCount:   50,
		Language:        "Go",
		Archived:        false,
		OpenIssuesCount: 10,
		DefaultBranch:   "main",
	}

	t.Run("Error when RepoID is invalid", func(t *testing.T) {
		invalidRepo := *validGithubRepo
		invalidRepo.ID = -1
		repo, err := NewRepository(&invalidRepo)
		require.Nil(t, repo)
		require.Error(t, err)
		require.EqualError(t, err, "repository ID must be positive")
	})

	t.Run("Error when RepoName is empty", func(t *testing.T) {
		invalidRepo := *validGithubRepo
		invalidRepo.Name = ""
		repo, err := NewRepository(&invalidRepo)
		require.Nil(t, repo)
		require.Error(t, err)
		require.EqualError(t, err, "repository name cannot be empty")
	})

	t.Run("Error when OwnerName is empty", func(t *testing.T) {
		invalidRepo := *validGithubRepo
		invalidRepo.Owner.Login = ""
		repo, err := NewRepository(&invalidRepo)
		require.Nil(t, repo)
		require.Error(t, err)
		require.EqualError(t, err, "owner name cannot be empty")
	})

	t.Run("Error when AvatarURL is invalid", func(t *testing.T) {
		invalidRepo := *validGithubRepo
		invalidRepo.Owner.AvatarURL = "invalid-url"
		repo, err := NewRepository(&invalidRepo)
		require.Nil(t, repo)
		require.Error(t, err)
		require.EqualError(t, err, "URL must start with http:// or https://")
	})

	t.Run("Error when HTMLURL is invalid", func(t *testing.T) {
		invalidRepo := *validGithubRepo
		invalidRepo.HTMLURL = "invalid-url"
		repo, err := NewRepository(&invalidRepo)
		require.Nil(t, repo)
		require.Error(t, err)
		require.EqualError(t, err, "URL must start with http:// or https://")
	})

	t.Run("Error when Homepage URL is invalid", func(t *testing.T) {
		invalidRepo := *validGithubRepo
		invalidRepo.Homepage = "invalid-url"
		repo, err := NewRepository(&invalidRepo)
		require.Nil(t, repo)
		require.Error(t, err)
		require.EqualError(t, err, "URL must start with http:// or https://")
	})

	t.Run("Error when CreatedAt is zero", func(t *testing.T) {
		invalidRepo := *validGithubRepo
		invalidRepo.CreatedAt = time.Time{}
		repo, err := NewRepository(&invalidRepo)
		require.Nil(t, repo)
		require.Error(t, err)
		require.EqualError(t, err, "created time cannot be empty")
	})

	t.Run("Error when UpdatedAt is zero", func(t *testing.T) {
		invalidRepo := *validGithubRepo
		invalidRepo.UpdatedAt = time.Time{}
		repo, err := NewRepository(&invalidRepo)
		require.Nil(t, repo)
		require.Error(t, err)
		require.EqualError(t, err, "updated time cannot be empty")
	})

	t.Run("Error when UpdatedAt is before CreatedAt", func(t *testing.T) {
		invalidRepo := *validGithubRepo
		invalidRepo.CreatedAt = time.Now()
		invalidRepo.UpdatedAt = time.Now().Add(-1 * time.Hour)
		repo, err := NewRepository(&invalidRepo)
		require.Nil(t, repo)
		require.Error(t, err)
		require.EqualError(t, err, "updated time cannot be before created time")
	})

	t.Run("Error when StargazersCount is negative", func(t *testing.T) {
		invalidRepo := *validGithubRepo
		invalidRepo.StargazersCount = -1
		repo, err := NewRepository(&invalidRepo)
		require.Nil(t, repo)
		require.Error(t, err)
		require.EqualError(t, err, "stargazers count cannot be negative")
	})

	t.Run("Error when WatchersCount is negative", func(t *testing.T) {
		invalidRepo := *validGithubRepo
		invalidRepo.WatchersCount = -1
		repo, err := NewRepository(&invalidRepo)
		require.Nil(t, repo)
		require.Error(t, err)
		require.EqualError(t, err, "watchers count cannot be negative")
	})

	t.Run("Error when OpenIssuesCount is negative", func(t *testing.T) {
		invalidRepo := *validGithubRepo
		invalidRepo.OpenIssuesCount = -1
		repo, err := NewRepository(&invalidRepo)
		require.Nil(t, repo)
		require.Error(t, err)
		require.EqualError(t, err, "open issues count cannot be negative")
	})

	t.Run("Error when DefaultBranch is empty", func(t *testing.T) {
		invalidRepo := *validGithubRepo
		invalidRepo.DefaultBranch = ""
		repo, err := NewRepository(&invalidRepo)
		require.Nil(t, repo)
		require.Error(t, err)
		require.EqualError(t, err, "default branch cannot be empty")
	})
}
