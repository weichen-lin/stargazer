package db

import (
	"context"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
	"github.com/weichen-lin/stargazer/domain"
)

func TestGetRepository(t *testing.T) {
	repo, err := db.GetRepository(context.Background(), 123)
	require.ErrorIs(t, err, ErrNotFoundEmailAtContext)
	require.Empty(t, repo)

	repositoryEntity := &domain.RepositoryEntity{
		RepoID:            123456789,
		RepoName:          "example-repo",
		OwnerName:         "example-owner",
		AvatarURL:         "https://example.com/avatar.png",
		HtmlURL:           "https://github.com/example/repo",
		Homepage:          "https://example.com",
		CreatedAt:         "2024-01-01T00:00:00Z",
		UpdatedAt:         "2024-01-02T00:00:00Z",
		StargazersCount:   100,
		WatchersCount:     50,
		OpenIssuesCount:   10,
		DefaultBranch:     "main",
		Description:       "test-description",
		Language:          "Go",
		Archived:          true,
		Topics:            []string{"TEST", "test-2"},
		ExternalCreatedAt: time.Now().Format(time.RFC3339),
		LastSyncedAt:      time.Now().Format(time.RFC3339),
		LastModifiedAt:    time.Now().Format(time.RFC3339),
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

	ctx, err = WithEmail(context.Background(), user.Email())
	require.NoError(t, err)

	err = db.CreateRepository(ctx, repo)
	require.NoError(t, err)

	repo, err = db.GetRepository(ctx, repositoryEntity.RepoID)
	require.NoError(t, err)

	expectCreated, _ := time.Parse(time.RFC3339, repositoryEntity.CreatedAt)
	expectUpdateAt, _ := time.Parse(time.RFC3339, repositoryEntity.UpdatedAt)

	expectExternalCreatedAt, err := time.Parse(time.RFC3339, repositoryEntity.ExternalCreatedAt)
	require.NoError(t, err)

	expectLastSyncedAt, err := time.Parse(time.RFC3339, repositoryEntity.LastSyncedAt)
	require.NoError(t, err)

	expectLastModifiedAt, err := time.Parse(time.RFC3339, repositoryEntity.LastModifiedAt)
	require.NoError(t, err)

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
	require.Equal(t, repo.Topics(), repositoryEntity.Topics)
	require.WithinDuration(t, repo.ExternalCreateAt(), expectExternalCreatedAt, time.Duration(time.Second*3))
	require.WithinDuration(t, repo.LastSyncedAt(), expectLastSyncedAt, time.Duration(time.Second*3))
	require.WithinDuration(t, repo.LastModifiedAt(), expectLastModifiedAt, time.Duration(time.Second*3))
}

func TestGetRepoLanguageDistribution(t *testing.T) {
	repo, err := db.GetRepository(context.Background(), 123)
	require.ErrorIs(t, err, ErrNotFoundEmailAtContext)
	require.Empty(t, repo)

	repositoryEntityGo := &domain.RepositoryEntity{
		RepoID:            1233211234,
		RepoName:          "example-repo-for-language-distribution",
		OwnerName:         "example-owner",
		AvatarURL:         "https://example.com/avatar.png",
		HtmlURL:           "https://github.com/example/repo",
		Homepage:          "https://example.com",
		CreatedAt:         "2024-01-01T00:00:00Z",
		UpdatedAt:         "2024-01-02T00:00:00Z",
		StargazersCount:   100,
		WatchersCount:     50,
		OpenIssuesCount:   10,
		DefaultBranch:     "main",
		Description:       "test-description",
		Language:          "Go",
		Archived:          true,
		Topics:            []string{"TEST"},
		ExternalCreatedAt: time.Now().Format(time.RFC3339),
		LastSyncedAt:      time.Now().Format(time.RFC3339),
		LastModifiedAt:    time.Now().Format(time.RFC3339),
	}

	repositoryEntityPython := &domain.RepositoryEntity{
		RepoID:            432112345,
		RepoName:          "example-repo-for-language-distribution",
		OwnerName:         "example-owner",
		AvatarURL:         "https://example.com/avatar.png",
		HtmlURL:           "https://github.com/example/repo",
		Homepage:          "https://example.com",
		CreatedAt:         "2024-01-01T00:00:00Z",
		UpdatedAt:         "2024-01-02T00:00:00Z",
		StargazersCount:   100,
		WatchersCount:     50,
		OpenIssuesCount:   10,
		DefaultBranch:     "main",
		Description:       "test-description",
		Language:          "Python",
		Archived:          true,
		Topics:            []string{"TEST-2"},
		ExternalCreatedAt: time.Now().Format(time.RFC3339),
		LastSyncedAt:      time.Now().Format(time.RFC3339),
		LastModifiedAt:    time.Now().Format(time.RFC3339),
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

	result, err = db.SearchRepositoryByLanguage(ctx, &SearchParams{
		Languages: []string{"NotExist"},
		Page:      1,
		Limit:     10,
	})
	require.NoError(t, err)
	require.Equal(t, result.Total, int64(0))
	require.Equal(t, len(result.Data), 0)
}

func TestGetAllRepositoryTopics(t *testing.T) {
	user, ctx := createFakeUser(t)

	repo1 := createRepositoryAtFakeUser(t, user)
	require.NotEmpty(t, repo1)

	repo2 := createRepositoryAtFakeUser(t, user)
	require.NotEmpty(t, repo2)

	results, err := db.GetAllRepositoryTopics(ctx)
	require.NoError(t, err)
	require.Equal(t, 2, len(results))
}

func TestGetRepositoriesOrderBy(t *testing.T) {
	user, ctx := createFakeUser(t)

	for i := 0; i < 10; i++ {
		repo := createRepositoryAtFakeUser(t, user)
		require.NotEmpty(t, repo)
	}

	t.Run("test sort key and sort order validate", func(t *testing.T) {
		_, err := db.GetRepositoriesOrderBy(ctx, &SortParams{
			Key:   "test-invalid-sort-key",
			Order: "DESC",
		})
		require.ErrorIs(t, err, ErrInvalidSortKey)

		_, err = db.GetRepositoriesOrderBy(ctx, &SortParams{
			Key:   "created_at",
			Order: "adsasda",
		})
		require.ErrorIs(t, err, ErrInvalidSortOrder)
	})

	t.Run("test limit repos activate", func(t *testing.T) {
		repos, err := db.GetRepositoriesOrderBy(ctx, &SortParams{
			Key:   "created_at",
			Order: "DESC",
		})
		require.NoError(t, err)
		require.Equal(t, 5, len(repos))
	})
}

func TestFullTextSearch(t *testing.T) {
	user, ctx := createFakeUser(t)
	for i := 0; i < 20; i++ {
		repo := createRepositoryAtFakeUser(t, user)
		require.NotEmpty(t, repo)
	}

	_, err := db.FullTextSearch(ctx, "et")
	require.NoError(t, err)
}
