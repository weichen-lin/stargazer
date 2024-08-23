package db

import (
	"context"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
	"github.com/weichen-lin/kafka-service/domain"
)

func Test_CreateAndRemoveTag(t *testing.T) {
	entity := &domain.UserEntity{
		Name:              "John Doe",
		Email:             "john.doe@example.com",
		Image:             "https://example.com/avatar.jpg",
		AccessToken:       "abc123",
		Provider:          "github",
		ProviderAccountId: "123456",
		Scope:             "read:user,user:email",
		AuthType:          "oauth",
		TokenType:         "bearer",
	}

	user := domain.FromUserEntity(entity)

	err := db.CreateUser(user)

	require.NoError(t, err)

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

	repo, err := domain.FromRepositoryEntity(repositoryEntity)
	require.NoError(t, err)
	require.NotEmpty(t, repo)

	ctx, err := WithEmail(context.Background(), entity.Email)
	require.NoError(t, err)

	err = db.CreateRepository(ctx, repo)
	require.NoError(t, err)

	name := faker.Name()
	tag, err := domain.NewTag(name)
	require.NoError(t, err)
	require.NotEmpty(t, tag)

	err = db.SaveTag(ctx, tag, repo.RepoID())
	require.NoError(t, err)

	savedTag, err := db.GetTagByName(ctx, tag.Name())
	require.NoError(t, err)
	require.Equal(t, savedTag.Name(), tag.Name())
	require.Equal(t, savedTag.CreatedAt(), tag.CreatedAt())
	require.Equal(t, savedTag.UpdatedAt(), tag.UpdatedAt())

	err = db.RemoveTag(ctx, savedTag, repo.RepoID())
	require.NoError(t, err)
}
