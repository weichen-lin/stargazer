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

	sharedCollection, err := db.GetCollectionById(ctx, collection.Id().String())
	require.NoError(t, err)
	require.Equal(t, collection.Id().String(), sharedCollection.Collection.Id)

	checkCollectionName, err := db.GetCollectionByName(ctx, collection.Name())
	require.NoError(t, err)
	require.Equal(t, collection.Name(), checkCollectionName.Name())

	total := 25

	repoIds := make([]int64, total)

	for i := 0; i < total; i++ {
		repo := createRepositoryAtFakeUser(t, user)
		repoIds[i] = repo.RepoID()
	}

	err = db.AddRepoToCollection(ctx, collection, repoIds)
	require.NoError(t, err)

	repos, err := db.GetCollectionContainRepos(ctx, collection, 1, 20)
	require.NoError(t, err)
	require.Len(t, repos.Data, 20)
	require.Equal(t, repos.Total, int64(25))

	repos, err = db.GetCollectionContainRepos(ctx, collection, 2, 20)
	require.NoError(t, err)
	require.Len(t, repos.Data, 5)
	require.Equal(t, repos.Total, int64(25))

	err = db.DeleteRepoFromCollection(ctx, collection, repoIds)
	require.NoError(t, err)

	repos, err = db.GetCollectionContainRepos(ctx, collection, 1, 25)
	require.NoError(t, err)
	require.Len(t, repos.Data, 0)
	require.Equal(t, repos.Total, int64(0))

	for i := 0; i < total; i++ {
		repo := createRepositoryAtFakeUser(t, user)
		repoIds[i] = repo.RepoID()
	}

	err = db.AddRepoToCollection(ctx, collection, repoIds)
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

func Test_ShareCollection(t *testing.T) {
	user, ctx := createFakeUser(t)
	user2, ctx2 := createFakeUser(t)

	collection, err := domain.NewCollection(faker.Name())
	require.NoError(t, err)

	err = db.SaveCollection(ctx, collection)
	require.NoError(t, err)

	shared, err := db.GetCollectionById(ctx, collection.Id().String())
	require.NoError(t, err)
	require.Equal(t, collection.Id().String(), shared.Collection.Id)

	total := 25

	repoIds := make([]int64, total)

	for i := 0; i < total; i++ {
		repo := createRepositoryAtFakeUser(t, user)
		repoIds[i] = repo.RepoID()
	}

	err = db.AddRepoToCollection(ctx, collection, repoIds)
	require.NoError(t, err)

	repos, err := db.GetCollectionContainRepos(ctx, collection, 1, 20)
	require.NoError(t, err)
	require.Len(t, repos.Data, 20)
	require.Equal(t, repos.Total, int64(25))

	repos, err = db.GetCollectionContainRepos(ctx, collection, 2, 20)
	require.NoError(t, err)
	require.Len(t, repos.Data, 5)
	require.Equal(t, repos.Total, int64(25))

	err = db.ShareCollection(ctx, collection, user2.Email())
	require.NoError(t, err)

	sharedCollections, err := db.GetSharedCollections(ctx2)
	require.NoError(t, err)
	require.Len(t, sharedCollections, 1)

	for _, sharedCollection := range sharedCollections {
		require.Equal(t, sharedCollection.Collection.Id, shared.Collection.Id)
		require.Equal(t, sharedCollection.SharedFrom.Email, user.Email())
	}
}

func Test_GetCollectionWithShared(t *testing.T) {
	user, ctx := createFakeUser(t)
	user2, ctx2 := createFakeUser(t)
	_, ctx3 := createFakeUser(t)

	collection, err := domain.NewCollection(faker.Name())
	require.NoError(t, err)

	err = db.SaveCollection(ctx, collection)
	require.NoError(t, err)

	checkCollection, err := db.GetCollectionById(ctx, collection.Id().String())
	require.NoError(t, err)
	require.Equal(t, collection.Id().String(), checkCollection.Collection.Id)
	require.Equal(t, user.Email(), checkCollection.Owner)
	require.Nil(t, checkCollection.SharedFrom)

	err = db.ShareCollection(ctx, collection, user2.Email())
	require.NoError(t, err)

	sharedCollection, err := db.GetCollectionById(ctx2, collection.Id().String())
	require.NoError(t, err)
	require.Equal(t, collection.Id().String(), sharedCollection.Collection.Id)
	require.Equal(t, "", sharedCollection.Owner)
	require.NotNil(t, sharedCollection.SharedFrom)
	require.Equal(t, sharedCollection.SharedFrom.Email, user.Email())

	notSharedCollection, err := db.GetCollectionById(ctx3, collection.Id().String())
	require.NoError(t, err)
	require.Equal(t, collection.Id().String(), notSharedCollection.Collection.Id)
	require.Equal(t, "", notSharedCollection.Owner)
	require.Nil(t, notSharedCollection.SharedFrom)
}
