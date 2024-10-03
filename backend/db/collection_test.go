package db

import (
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
	"github.com/weichen-lin/stargazer/domain"
)

func Test_SaveCollection(t *testing.T) {
	user, ctx := createFakeUser(t)

	collection, err := domain.NewCollection(faker.Name())
	require.NoError(t, err)

	err = db.SaveCollection(ctx, collection)
	require.NoError(t, err)

	checkCollection, err := db.GetCollectionById(ctx, collection.Id().String())
	require.NoError(t, err)
	require.Equal(t, collection.Id(), checkCollection.Id())

	checkCollection, err = db.GetCollectionByName(ctx, collection.Name())
	require.NoError(t, err)
	require.Equal(t, collection.Name(), checkCollection.Name())

	total := 25

	repoIds := make([]int64, total)

	for i := 0; i < total; i++ {
		repo := createRepositoryAtFakeUser(t, user)
		repoIds[i] = repo.RepoID()
	}

	err = db.AddRepoToCollection(ctx, checkCollection, repoIds)
	require.NoError(t, err)

	repos, err := db.GetCollectionContainRepos(ctx, checkCollection, 1, 20)
	require.NoError(t, err)
	require.Len(t, repos.Data, 20)
	require.Equal(t, repos.Total, int64(25))

	repos, err = db.GetCollectionContainRepos(ctx, checkCollection, 2, 20)
	require.NoError(t, err)
	require.Len(t, repos.Data, 5)
	require.Equal(t, repos.Total, int64(25))

	err = db.DeleteRepoFromCollection(ctx, checkCollection, repoIds)
	require.NoError(t, err)

	repos, err = db.GetCollectionContainRepos(ctx, checkCollection, 1, 25)
	require.NoError(t, err)
	require.Len(t, repos.Data, 0)
	require.Equal(t, repos.Total, int64(0))

	for i := 0; i < total; i++ {
		repo := createRepositoryAtFakeUser(t, user)
		repoIds[i] = repo.RepoID()
	}

	err = db.AddRepoToCollection(ctx, checkCollection, repoIds)
	require.NoError(t, err)

	err = db.DeleteCollection(ctx, collection.Id().String())
	require.NoError(t, err)

	_, err = db.GetCollectionById(ctx, collection.Id().String())
	require.Error(t, err)
}

func Test_GetCollectionByPage(t *testing.T) {
	_, ctx := createFakeUser(t)

	for i := 0; i < 25; i++ {
		collection, err := domain.NewCollection(faker.Name())
		require.NoError(t, err)
		err = db.SaveCollection(ctx, collection)
		require.NoError(t, err)
	}

	result, err := db.GetCollections(ctx, &PagingParams{
		Page:  1,
		Limit: 20,
	})

	require.NoError(t, err)
	require.Len(t, result.Data, 20)
}
