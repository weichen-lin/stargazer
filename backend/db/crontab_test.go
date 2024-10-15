package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/weichen-lin/stargazer/domain"
)

func Test_SaveCrontab(t *testing.T) {
	_, ctx := createFakeUser(t)

	crontab := domain.NewCrontab()

	err := db.SaveCrontab(context.Background(), crontab)
	require.ErrorIs(t, err, ErrNotFoundEmailAtContext)

	err = db.SaveCrontab(ctx, crontab)
	require.NoError(t, err)

	ctx, err = WithEmail(context.Background(), "test-not-exists@gmail.com")
	require.NoError(t, err)
	err = db.SaveCrontab(ctx, crontab)
	require.Error(t, err)

	_, ctx = createFakeUser(t)
	err = db.SaveCrontab(ctx, crontab)
	require.NoError(t, err)

	_, ctx = createFakeUser(t)
	err = db.SaveCrontab(ctx, crontab)
	require.NoError(t, err)

	crontabs := db.GetAllCrontab()
	require.NotEmpty(t, crontabs)
}

func Test_GetCrontab(t *testing.T) {
	_, ctx := createFakeUser(t)

	newCrontab := domain.NewCrontab()
	require.NotEmpty(t, newCrontab)

	err := db.SaveCrontab(ctx, newCrontab)
	require.NoError(t, err)

	crontab, err := db.GetCrontab(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, crontab)
	require.WithinDuration(t, crontab.CreatedAt(), newCrontab.CreatedAt(), time.Duration(2*time.Second))
	require.WithinDuration(t, crontab.UpdatedAt(), newCrontab.UpdatedAt(), time.Duration(2*time.Second))
	require.Equal(t, crontab.TriggeredAt(), time.Time{})
	require.Equal(t, crontab.LastTriggeredAt(), time.Time{})
	require.Equal(t, crontab.Status(), "new")
}
