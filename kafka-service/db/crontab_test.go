package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/weichen-lin/kafka-service/domain"
)

func Test_CreateCrontab(t *testing.T) {
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

	err := db.CreateUser(user)
	require.NoError(t, err)

	crontab := domain.NewCrontab()

	err = db.CreateCrontab(context.Background(), crontab)
	require.ErrorIs(t, err, ErrNotFoundEmailAtContext)

	ctx, err := WithEmail(context.Background(), entity.Email)
	require.NoError(t, err)
	require.NotEmpty(t, ctx)

	err = db.CreateCrontab(ctx, crontab)
	require.NoError(t, err)

	ctx, err = WithEmail(context.Background(), "test-not-exists@gmail.com")
	err = db.CreateCrontab(ctx, crontab)
	require.Error(t, err)
}

func Test_GetCrontab(t *testing.T) {
	entity := &domain.UserEntity{
		Name:              "Test 123",
		Email:             "john.doe.456@example.com",
		Image:             "https://example.comhaha/avatar.jpg",
		AccessToken:       "abc123123123123",
		Provider:          "github",
		ProviderAccountId: "123412312312356",
		Scope:             "read:user,user:email",
		AuthType:          "oauth",
		TokenType:         "bearer",
	}

	user := domain.FromUserEntity(entity)

	err := db.CreateUser(user)
	require.NoError(t, err)

	ctx, err := WithEmail(context.Background(), entity.Email)
	require.NoError(t, err)
	require.NotEmpty(t, ctx)

	newCrontab := domain.NewCrontab()
	require.NotEmpty(t, newCrontab)

	err = db.CreateCrontab(ctx, newCrontab)
	require.NoError(t, err)

	crontab, err := db.GetCrontab(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, crontab)
	require.Equal(t, crontab.UpdatedAt(), time.Time{})
	require.Equal(t, crontab.TriggeredAt(), time.Time{})
	require.Equal(t, crontab.LastTriggeredAt(), time.Time{})
	require.Equal(t, crontab.Status(), "new")
	require.Equal(t, crontab.Version(), int64(1))
}
