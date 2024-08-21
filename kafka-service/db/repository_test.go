package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/weichen-lin/kafka-service/domain"
)

func TestGetRepository(t *testing.T) {
	repo, err := db.GetRepository(context.Background(), 123)
	require.Error(t, err)
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
	}

	repo, err = domain.FromRepositoryEntity(repositoryEntity)
	require.NoError(t, err)
	require.NotEmpty(t, repo)

	ctx, err := WithEmail(context.Background(), "valid@example.com")
	require.NoError(t, err)
	require.NotEmpty(t, ctx)

	err = db.CreateRepository(ctx, repo)
	require.NoError(t, err)
}
