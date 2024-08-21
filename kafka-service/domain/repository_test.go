package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewRepository(t *testing.T) {
	createdAt := "2014-03-24T16:04:04Z"
	updatedAt := "2014-03-24T17:04:04Z"

	expectCreatedAt, _ := time.Parse(time.RFC3339, createdAt)
	expectUpdatedAt, _ := time.Parse(time.RFC3339, updatedAt)

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
	require.Equal(t, expectCreatedAt, repo.CreatedAt())
	require.Equal(t, expectUpdatedAt, repo.UpdatedAt())
	require.Equal(t, int64(10), repo.StargazersCount())
	require.Equal(t, int64(20), repo.WatchersCount())
	require.Equal(t, int64(5), repo.OpenIssuesCount())
	require.Equal(t, "Go", repo.Language())
	require.Equal(t, "main", repo.DefaultBranch())
	require.False(t, repo.Archived())
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

	createdAt := "2014-03-24T16:04:04Z"
	expectedCreatedAt, _ := time.Parse(time.RFC3339, createdAt)

	err = repo.setCreatedAt(createdAt)
	require.NoError(t, err)
	require.Equal(t, expectedCreatedAt, repo.CreatedAt())

	updatedAt := "2014-03-24T17:04:04Z"
	expectedUpdatedAt, _ := time.Parse(time.RFC3339, updatedAt)
	err = repo.setUpdatedAt(updatedAt)
	require.NoError(t, err)
	require.Equal(t, expectedUpdatedAt, repo.UpdatedAt())

	err = repo.setStargazersCount(10)
	require.NoError(t, err)
	require.Equal(t, int64(10), repo.StargazersCount())

	err = repo.setWatchersCount(20)
	require.NoError(t, err)
	require.Equal(t, int64(20), repo.WatchersCount())

	err = repo.setOpenIssuesCount(5)
	require.NoError(t, err)
	require.Equal(t, int64(5), repo.OpenIssuesCount())

	repo.setLanguage("Go")
	require.Equal(t, "Go", repo.Language())

	err = repo.setDefaultBranch("main")
	require.NoError(t, err)
	require.Equal(t, "main", repo.DefaultBranch())

	repo.setArchived(false)
	require.False(t, repo.Archived())
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
	err = repo.setCreatedAt("")
	require.Error(t, err)
	require.EqualError(t, err, "created time cannot be empty")

	// Test zero updated time
	err = repo.setUpdatedAt("")
	require.Error(t, err)
	require.EqualError(t, err, "updated time cannot be empty")

	// Test updated time before created time
	createdAt := "2014-03-24T16:04:04Z"
	updatedAt := "2014-03-24T14:04:04Z"
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
		CreatedAt:       "2014-03-23T16:04:04Z",
		UpdatedAt:       "2014-03-24T16:04:04Z",
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

	t.Run("Error when CreatedAt is empty", func(t *testing.T) {
		invalidRepo := *validGithubRepo
		invalidRepo.CreatedAt = ""
		repo, err := NewRepository(&invalidRepo)
		require.Nil(t, repo)
		require.Error(t, err)
		require.EqualError(t, err, "created time cannot be empty")
	})

	t.Run("Error when CreatedAt is invalid year", func(t *testing.T) {
		invalidRepo := *validGithubRepo
		invalidRepo.CreatedAt = "0000-03-24T16:04:04Z"
		repo, err := NewRepository(&invalidRepo)
		require.Nil(t, repo)
		require.Error(t, err)
	})

	t.Run("Error when CreatedAt is invalid format", func(t *testing.T) {
		invalidRepo := *validGithubRepo
		invalidRepo.CreatedAt = "2023-02-29T16:04:04Z"
		repo, err := NewRepository(&invalidRepo)
		require.Nil(t, repo)
		require.Error(t, err)
	})

	t.Run("Error when UpdatedAt is empty", func(t *testing.T) {
		invalidRepo := *validGithubRepo
		invalidRepo.UpdatedAt = ""
		repo, err := NewRepository(&invalidRepo)
		require.Nil(t, repo)
		require.Error(t, err)
		require.EqualError(t, err, "updated time cannot be empty")
	})

	t.Run("Error when UpdatedAt is before CreatedAt", func(t *testing.T) {
		invalidRepo := *validGithubRepo
		invalidRepo.CreatedAt = "2014-03-24T16:04:04Z"
		invalidRepo.UpdatedAt = "2014-03-23T16:04:04Z"
		repo, err := NewRepository(&invalidRepo)
		require.Nil(t, repo)
		require.Error(t, err)
		require.EqualError(t, err, "updated time cannot be before created time")
	})

	t.Run("Error when UpdatedAt is invalid year", func(t *testing.T) {
		invalidRepo := *validGithubRepo
		invalidRepo.UpdatedAt = "0000-03-24T16:04:04Z"
		repo, err := NewRepository(&invalidRepo)
		require.Nil(t, repo)
		require.Error(t, err)
	})

	t.Run("Error when UpdatedAt is invalid format", func(t *testing.T) {
		invalidRepo := *validGithubRepo
		invalidRepo.UpdatedAt = "2023-02-29T16:04:04Z"
		repo, err := NewRepository(&invalidRepo)
		require.Nil(t, repo)
		require.Error(t, err)
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

func TestToRepositoryEntity(t *testing.T) {
	createdAt := "2014-03-24T16:04:04Z"
	updatedAt := "2014-03-24T17:04:04Z"

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

	repositoryEntity := repo.ToRepositoryEntity()
	require.Equal(t, int64(12345), repositoryEntity.RepoID)
	require.Equal(t, "sample-repo", repositoryEntity.RepoName)
	require.Equal(t, "user", repositoryEntity.OwnerName)
	require.Equal(t, "https://avatar.url", repositoryEntity.AvatarURL)
	require.Equal(t, "https://github.com/user/sample-repo", repositoryEntity.HtmlURL)
	require.Equal(t, "https://sample-repo.com", repositoryEntity.Homepage)
	require.Equal(t, "This is a sample repository", repositoryEntity.Description)
	require.Equal(t, createdAt, repositoryEntity.CreatedAt)
	require.Equal(t, updatedAt, repositoryEntity.UpdatedAt)
	require.Equal(t, int64(10), repositoryEntity.StargazersCount)
	require.Equal(t, int64(20), repositoryEntity.WatchersCount)
	require.Equal(t, int64(5), repositoryEntity.OpenIssuesCount)
	require.Equal(t, "Go", repositoryEntity.Language)
	require.Equal(t, "main", repositoryEntity.DefaultBranch)
	require.False(t, repositoryEntity.Archived)
}

func TestFromRepositoryEntity(t *testing.T) {
	repositoryEntity := &RepositoryEntity{
		RepoID:          123456789,
		RepoName:        "example-repo",
		OwnerName:       "example-owner",
		AvatarURL:       "https://example.com/avatar.png",
		HtmlURL:         "https://github.com/example/repo",
		Homepage:        "https://example.com",
		Description:     "An example repository",
		CreatedAt:       "2024-01-01T00:00:00Z",
		UpdatedAt:       "2024-01-02T00:00:00Z",
		StargazersCount: 100,
		WatchersCount:   50,
		OpenIssuesCount: 5,
		Language:        "Go",
		DefaultBranch:   "main",
		Archived:        false,
	}

	parsedCreatedAt, _ := time.Parse(time.RFC3339, repositoryEntity.CreatedAt)
	parsedUpdatedAt, _ := time.Parse(time.RFC3339, repositoryEntity.UpdatedAt)

	repoFromEntity, err := FromRepositoryEntity(repositoryEntity)
	require.NoError(t, err)
	require.Equal(t, repoFromEntity.RepoID(), repositoryEntity.RepoID)
	require.Equal(t, repoFromEntity.RepoName(), repositoryEntity.RepoName)
	require.Equal(t, repoFromEntity.OwnerName(), repositoryEntity.OwnerName)
	require.Equal(t, repoFromEntity.HTMLURL(), repositoryEntity.HtmlURL)
	require.Equal(t, repoFromEntity.Description(), repositoryEntity.Description)
	require.Equal(t, repoFromEntity.CreatedAt(), parsedCreatedAt)
	require.Equal(t, repoFromEntity.UpdatedAt(), parsedUpdatedAt)
	require.Equal(t, repoFromEntity.Homepage(), repositoryEntity.Homepage)
	require.Equal(t, repoFromEntity.StargazersCount(), repositoryEntity.StargazersCount)
	require.Equal(t, repoFromEntity.WatchersCount(), repositoryEntity.WatchersCount)
	require.Equal(t, repoFromEntity.Language(), repositoryEntity.Language)
	require.Equal(t, repoFromEntity.Archived(), repositoryEntity.Archived)
	require.Equal(t, repoFromEntity.OpenIssuesCount(), repositoryEntity.OpenIssuesCount)
	require.Equal(t, repoFromEntity.DefaultBranch(), repositoryEntity.DefaultBranch)
}

func TestFromRepositoryEntitySetterErrors(t *testing.T) {
	repositoryEntity := &RepositoryEntity{
		RepoID: -1,
	}

	_, err := FromRepositoryEntity(repositoryEntity)
	require.Error(t, err, "expected error for invalid RepoID")

	repositoryEntity = &RepositoryEntity{
		RepoID:   123456789,
		RepoName: "",
	}

	_, err = FromRepositoryEntity(repositoryEntity)
	require.Error(t, err, "expected error for invalid RepoName")

	repositoryEntity = &RepositoryEntity{
		RepoID:    123456789,
		RepoName:  "example-repo",
		OwnerName: "",
	}

	_, err = FromRepositoryEntity(repositoryEntity)
	require.Error(t, err, "expected error for invalid OwnerName")

	repositoryEntity = &RepositoryEntity{
		RepoID:    123456789,
		RepoName:  "example-repo",
		OwnerName: "example-owner",
		AvatarURL: "invalid-url",
	}

	_, err = FromRepositoryEntity(repositoryEntity)
	require.Error(t, err, "expected error for invalid AvatarURL")

	repositoryEntity = &RepositoryEntity{
		RepoID:    123456789,
		RepoName:  "example-repo",
		OwnerName: "example-owner",
		AvatarURL: "https://example.com/avatar.png",
		HtmlURL:   "invalid-url",
	}

	_, err = FromRepositoryEntity(repositoryEntity)
	require.Error(t, err, "expected error for invalid HTMLURL")

	repositoryEntity = &RepositoryEntity{
		RepoID:    123456789,
		RepoName:  "example-repo",
		OwnerName: "example-owner",
		AvatarURL: "https://example.com/avatar.png",
		HtmlURL:   "https://github.com/example/repo",
		Homepage:  "invalid-url",
	}

	_, err = FromRepositoryEntity(repositoryEntity)
	require.Error(t, err, "expected error for invalid Homepage")

	repositoryEntity = &RepositoryEntity{
		RepoID:    123456789,
		RepoName:  "example-repo",
		OwnerName: "example-owner",
		AvatarURL: "https://example.com/avatar.png",
		HtmlURL:   "https://github.com/example/repo",
		Homepage:  "https://example.com",
		CreatedAt: "invalid-date",
	}

	_, err = FromRepositoryEntity(repositoryEntity)
	require.Error(t, err, "expected error for invalid CreatedAt")

	repositoryEntity = &RepositoryEntity{
		RepoID:    123456789,
		RepoName:  "example-repo",
		OwnerName: "example-owner",
		AvatarURL: "https://example.com/avatar.png",
		HtmlURL:   "https://github.com/example/repo",
		Homepage:  "https://example.com",
		CreatedAt: "2024-01-01T00:00:00Z",
		UpdatedAt: "invalid-date",
	}

	_, err = FromRepositoryEntity(repositoryEntity)
	require.Error(t, err, "expected error for invalid UpdatedAt")

	repositoryEntity = &RepositoryEntity{
		RepoID:          123456789,
		RepoName:        "example-repo",
		OwnerName:       "example-owner",
		AvatarURL:       "https://example.com/avatar.png",
		HtmlURL:         "https://github.com/example/repo",
		Homepage:        "https://example.com",
		CreatedAt:       "2024-01-01T00:00:00Z",
		UpdatedAt:       "2024-01-02T00:00:00Z",
		StargazersCount: -100,
	}

	_, err = FromRepositoryEntity(repositoryEntity)
	require.Error(t, err, "expected error for invalid StargazersCount")

	repositoryEntity = &RepositoryEntity{
		RepoID:          123456789,
		RepoName:        "example-repo",
		OwnerName:       "example-owner",
		AvatarURL:       "https://example.com/avatar.png",
		HtmlURL:         "https://github.com/example/repo",
		Homepage:        "https://example.com",
		CreatedAt:       "2024-01-01T00:00:00Z",
		UpdatedAt:       "2024-01-02T00:00:00Z",
		StargazersCount: 100,
		WatchersCount:   -50,
	}

	_, err = FromRepositoryEntity(repositoryEntity)
	require.Error(t, err, "expected error for invalid WatchersCount")

	repositoryEntity = &RepositoryEntity{
		RepoID:          123456789,
		RepoName:        "example-repo",
		OwnerName:       "example-owner",
		AvatarURL:       "https://example.com/avatar.png",
		HtmlURL:         "https://github.com/example/repo",
		Homepage:        "https://example.com",
		CreatedAt:       "2024-01-01T00:00:00Z",
		UpdatedAt:       "2024-01-02T00:00:00Z",
		StargazersCount: 100,
		WatchersCount:   50,
		OpenIssuesCount: -10,
	}

	_, err = FromRepositoryEntity(repositoryEntity)
	require.Error(t, err, "expected error for invalid OpenIssuesCount")

	repositoryEntity = &RepositoryEntity{
		RepoID:          123456789,
		RepoName:        "example-repo",
		OwnerName:       "example-owner",
		AvatarURL:       "https://example.com/avatar.png",
		HtmlURL:         "https://github.com/example/repo",
		Homepage:        "https://example.com",
		CreatedAt:       "2024-01-01T00:00:00Z",
		UpdatedAt:       "2024-01-02T00:00:00Z",
		StargazersCount: 100,
		WatchersCount:   50,
		OpenIssuesCount: 10,
		DefaultBranch:   "",
	}

	_, err = FromRepositoryEntity(repositoryEntity)
	require.Error(t, err, "expected error for invalid DefaultBranch")
}
