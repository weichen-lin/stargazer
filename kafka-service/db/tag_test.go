package db

import (
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
	"github.com/weichen-lin/kafka-service/domain"
)

func Test_GetTagsByRepo(t *testing.T) {
	t.Run("Save tag and check exist, finally try to delete", func(t *testing.T) {

		user, ctx := createFakeUser(t)
		repo := createRepositoryAtFakeUser(t, user)

		tag, err := domain.NewTag(faker.Name())
		require.NoError(t, err)
		require.NotEmpty(t, tag)

		err = db.SaveTag(ctx, tag, repo.RepoID())
		require.NoError(t, err)

		fakeName := "NotExist"

		tagfromDB, err := db.GetTagByName(ctx, fakeName)
		require.ErrorIs(t, err, ErrorNotFoundTag)
		require.Empty(t, tagfromDB)

		tagfromDB, err = db.GetTagByName(ctx, tag.Name())
		require.NoError(t, err)
		require.Equal(t, tag.Name(), tagfromDB.Name())

		err = db.RemoveTag(ctx, tag, repo.RepoID())
		require.NoError(t, err)

		tagfromDB, err = db.GetTagByName(ctx, tag.Name())
		require.ErrorIs(t, err, ErrorNotFoundTag)
		require.Empty(t, tagfromDB)
	})

	t.Run("test get tags by repo", func(t *testing.T) {
		user, ctx := createFakeUser(t)
		repo := createRepositoryAtFakeUser(t, user)

		for i := 0; i < 10; i++ {
			newTag, err := domain.NewTag(faker.Name())

			require.NoError(t, err)
			require.NotEmpty(t, newTag)

			err = db.SaveTag(ctx, newTag, repo.RepoID())
			require.NoError(t, err)
		}

		tags, err := db.GetTagsByRepo(ctx, repo.RepoID())
		require.NoError(t, err)
		require.Equal(t, 10, len(tags))

		newRepo := createRepositoryAtFakeUser(t, user)
		tags, err = db.GetTagsByRepo(ctx, newRepo.RepoID())
		require.NoError(t, err)
		require.Equal(t, 0, len(tags))
	})
}
