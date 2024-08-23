package db

import (
	"context"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
	"github.com/weichen-lin/kafka-service/domain"
)

func TestGetRepository(t *testing.T) {
	repo, err := db.GetRepository(context.Background(), 123)
	require.ErrorIs(t, err, ErrNotFoundEmailAtContext)
	require.Empty(t, repo)

	repositoryEntity := &domain.RepositoryEntity{
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
		DefaultBranch:   "main",
		Description:     "test-description",
		Language:        "Go",
		Archived:        true,
	}

	repo, err = domain.FromRepositoryEntity(repositoryEntity)
	require.NoError(t, err)
	require.NotEmpty(t, repo)

	entity := &domain.UserEntity{
		Name:              "Test 123",
		Email:             "john.doe.123@example.com",
		Image:             "https://example.comhaha/avatar.jpg",
		AccessToken:       "abc123123123123",
		Provider:          "github",
		ProviderAccountId: "123412312312356",
		Scope:             "read:user,user:email",
		AuthType:          "oauth",
		TokenType:         "bearer",
	}

	user := domain.FromUserEntity(entity)

	err = db.CreateUser(user)
	require.NoError(t, err)

	ctx, err := WithEmail(context.Background(), "valid@example.com")
	require.NoError(t, err)
	require.NotEmpty(t, ctx)

	err = db.CreateRepository(ctx, repo)
	require.ErrorIs(t, err, ErrRepositoryNotFound)

	ctx, err = WithEmail(context.Background(), "john.doe.123@example.com")
	require.NoError(t, err)

	err = db.CreateRepository(ctx, repo)
	require.NoError(t, err)

	repo, err = db.GetRepository(ctx, repositoryEntity.RepoID)
	require.NoError(t, err)

	expectCreated, _ := time.Parse(time.RFC3339, repositoryEntity.CreatedAt)
	expectUpdateAt, _ := time.Parse(time.RFC3339, repositoryEntity.UpdatedAt)

	require.Equal(t, repo.RepoID(), repositoryEntity.RepoID)
	require.Equal(t, repo.RepoName(), repositoryEntity.RepoName)
	require.Equal(t, repo.OwnerName(), repositoryEntity.OwnerName)
	require.Equal(t, repo.AvatarURL(), repositoryEntity.AvatarURL)
	require.Equal(t, repo.HTMLURL(), repositoryEntity.HtmlURL)
	require.Equal(t, repo.Homepage(), repositoryEntity.Homepage)
	require.Equal(t, repo.CreatedAt(), expectCreated)
	require.Equal(t, repo.UpdatedAt(), expectUpdateAt)
	require.Equal(t, repo.StargazersCount(), repositoryEntity.StargazersCount)
	require.Equal(t, repo.WatchersCount(), repositoryEntity.WatchersCount)
	require.Equal(t, repo.OpenIssuesCount(), repositoryEntity.OpenIssuesCount)
	require.Equal(t, repo.DefaultBranch(), repositoryEntity.DefaultBranch)
	require.Equal(t, repo.Description(), repositoryEntity.Description)
	require.Equal(t, repo.Language(), repositoryEntity.Language)
	require.Equal(t, repo.Archived(), repositoryEntity.Archived)
}

func TestGetRepoLanguageDistribution(t *testing.T) {
	repo, err := db.GetRepository(context.Background(), 123)
	require.ErrorIs(t, err, ErrNotFoundEmailAtContext)
	require.Empty(t, repo)

	repositoryEntityGo := &domain.RepositoryEntity{
		RepoID:          1233211234,
		RepoName:        "example-repo-for-language-distribution",
		OwnerName:       "example-owner",
		AvatarURL:       "https://example.com/avatar.png",
		HtmlURL:         "https://github.com/example/repo",
		Homepage:        "https://example.com",
		CreatedAt:       "2024-01-01T00:00:00Z",
		UpdatedAt:       "2024-01-02T00:00:00Z",
		StargazersCount: 100,
		WatchersCount:   50,
		OpenIssuesCount: 10,
		DefaultBranch:   "main",
		Description:     "test-description",
		Language:        "Go",
		Archived:        true,
	}

	repositoryEntityPython := &domain.RepositoryEntity{
		RepoID:          432112345,
		RepoName:        "example-repo-for-language-distribution",
		OwnerName:       "example-owner",
		AvatarURL:       "https://example.com/avatar.png",
		HtmlURL:         "https://github.com/example/repo",
		Homepage:        "https://example.com",
		CreatedAt:       "2024-01-01T00:00:00Z",
		UpdatedAt:       "2024-01-02T00:00:00Z",
		StargazersCount: 100,
		WatchersCount:   50,
		OpenIssuesCount: 10,
		DefaultBranch:   "main",
		Description:     "test-description",
		Language:        "Python",
		Archived:        true,
	}

	goRepo, err := domain.FromRepositoryEntity(repositoryEntityGo)
	require.NoError(t, err)
	require.NotEmpty(t, goRepo)

	pythonRepo, err := domain.FromRepositoryEntity(repositoryEntityPython)
	require.NoError(t, err)
	require.NotEmpty(t, pythonRepo)

	entity := &domain.UserEntity{
		Name:              "Test 123",
		Email:             faker.Email(),
		Image:             "https://example.comhaha/avatar.jpg",
		AccessToken:       "abc123123123123",
		Provider:          "github",
		ProviderAccountId: "123412312312356",
		Scope:             "read:user,user:email",
		AuthType:          "oauth",
		TokenType:         "bearer",
	}

	user := domain.FromUserEntity(entity)

	err = db.CreateUser(user)
	require.NoError(t, err)

	ctx, err := WithEmail(context.Background(), entity.Email)
	require.NoError(t, err)

	err = db.CreateRepository(ctx, goRepo)
	require.NoError(t, err)

	err = db.CreateRepository(ctx, pythonRepo)
	require.NoError(t, err)

	goRepo_from_db, err := db.GetRepository(ctx, goRepo.RepoID())
	require.NoError(t, err)
	require.NotEmpty(t, goRepo_from_db)

	pythonRepo_from_db, err := db.GetRepository(ctx, goRepo.RepoID())
	require.NoError(t, err)
	require.NotEmpty(t, pythonRepo_from_db)

	languageDistribution, err := db.GetRepoLanguageDistribution(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, languageDistribution)
	require.Equal(t, len(languageDistribution), 2)

	result, err := db.SearchRepositoryByLanguage(ctx, &SearchParams{
		Languages: []string{"Go"},
		Page:      1,
		Limit:     10,
	})
	require.NoError(t, err)
	require.Equal(t, result.Total, int64(1))
	require.Equal(t, len(result.Data), 1)
}
