package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/weichen-lin/kafka-service/domain"
)

func Test_CreateUser(t *testing.T) {

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
}

func Test_GetUser(t *testing.T) {
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

	ctx, err := WithEmail(context.Background(), entity.Email)
	require.NoError(t, err)
	require.NotEmpty(t, ctx)

	userFromdb, err := db.GetUser(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, userFromdb)

	require.Equal(t, entity.Name, userFromdb.Name())
	require.Equal(t, entity.Email, userFromdb.Email())
	require.Equal(t, entity.AccessToken, userFromdb.AccessToken())
	require.Equal(t, entity.Provider, userFromdb.Provider())
	require.Equal(t, entity.ProviderAccountId, userFromdb.ProviderAccountId())

	require.Equal(t, entity.Scope, userFromdb.Scope())
	require.Equal(t, entity.AuthType, userFromdb.AuthType())
	require.Equal(t, entity.TokenType, userFromdb.TokenType())
}

func Test_ErrorGetUser(t *testing.T) {
	ctx, err := WithEmail(context.Background(), "invalid@example.com")
	require.NoError(t, err)
	require.NotEmpty(t, ctx)

	_, err = db.GetUser(ctx)
	require.ErrorIs(t, err, ErrNotFoundEmailAtContext)
}
