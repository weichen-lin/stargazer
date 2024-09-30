package db

import (
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
	"github.com/weichen-lin/stargazer/domain"
)

func Test_SaveFolder(t *testing.T) {
	user, ctx := createFakeUser(t)

	folder := domain.NewFolder(faker.Name())

	err := db.SaveFolder(ctx, folder)
	require.NoError(t, err)

	getFolder, err := db.GetFolder(ctx, folder.Name())
	require.NoError(t, err)
	require.Equal(t, folder.Name(), getFolder.Name())

	_, err = db.GetFolder(ctx, "test invalid")
	require.Error(t, err)

	repo := createRepositoryAtFakeUser(t, user)
	err = db.AddRepoToFolder(ctx, getFolder, repo.RepoID())
	require.NoError(t, err)

	err = db.DeleteRepoFromFolder(ctx, getFolder, repo.RepoID())
	require.NoError(t, err)
}

func Test_GetFolderByPage(t *testing.T) {
	_, ctx := createFakeUser(t)

	for i := 0; i < 25; i++ {
		folder := domain.NewFolder(faker.Name())
		err := db.SaveFolder(ctx, folder)
		require.NoError(t, err)
	}

	result, err := db.GetFolders(ctx, &GetFolderParams{
		Page:  1,
		Limit: 20,
	})

	require.NoError(t, err)
	require.Len(t, result.Data, 20)
}
