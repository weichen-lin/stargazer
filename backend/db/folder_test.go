package db

import (
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
	"github.com/weichen-lin/stargazer/domain"
)

func Test_SaveFolder(t *testing.T) {
	user, ctx := createFakeUser(t)

	folder, err := domain.NewFolder(faker.Name())
	require.NoError(t, err)

	err = db.SaveFolder(ctx, folder)
	require.NoError(t, err)

	checkFolder, err := db.GetFolderById(ctx, folder.Id().String())
	require.NoError(t, err)
	require.Equal(t, folder.Id(), checkFolder.Id())

	checkFolder, err = db.GetFolderByName(ctx, folder.Name())
	require.NoError(t, err)
	require.Equal(t, folder.Name(), checkFolder.Name())

	total := 25

	repoIds := make([]int64, total)

	for i := 0; i < total; i++ {
		repo := createRepositoryAtFakeUser(t, user)
		repoIds[i] = repo.RepoID()
	}

	err = db.AddRepoToFolder(ctx, checkFolder, repoIds)
	require.NoError(t, err)

	repos, err := db.GetFolderContainRepos(ctx, checkFolder, 1, 20)
	require.NoError(t, err)
	require.Len(t, repos.Data, 20)
	require.Equal(t, repos.Total, int64(25))

	repos, err = db.GetFolderContainRepos(ctx, checkFolder, 2, 20)
	require.NoError(t, err)
	require.Len(t, repos.Data, 5)
	require.Equal(t, repos.Total, int64(25))

	err = db.DeleteRepoFromFolder(ctx, checkFolder, repoIds)
	require.NoError(t, err)

	repos, err = db.GetFolderContainRepos(ctx, checkFolder, 1, 25)
	require.NoError(t, err)
	require.Len(t, repos.Data, 0)
	require.Equal(t, repos.Total, int64(0))

	for i := 0; i < total; i++ {
		repo := createRepositoryAtFakeUser(t, user)
		repoIds[i] = repo.RepoID()
	}

	err = db.AddRepoToFolder(ctx, checkFolder, repoIds)
	require.NoError(t, err)

	err = db.DeleteFolder(ctx, folder.Id().String())
	require.NoError(t, err)

	_, err = db.GetFolderById(ctx, folder.Id().String())
	require.Error(t, err)
}

func Test_GetFolderByPage(t *testing.T) {
	_, ctx := createFakeUser(t)

	for i := 0; i < 25; i++ {
		folder, err := domain.NewFolder(faker.Name())
		require.NoError(t, err)
		err = db.SaveFolder(ctx, folder)
		require.NoError(t, err)
	}

	result, err := db.GetFolders(ctx, &PagingParams{
		Page:  1,
		Limit: 20,
	})

	require.NoError(t, err)
	require.Len(t, result.Data, 20)
}
