package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFromUserEntity(t *testing.T) {
	entity := &UserEntity{
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

	expectedUser := &User{
		name:              "John Doe",
		email:             "john.doe@example.com",
		image:             "https://example.com/avatar.jpg",
		accessToken:       "abc123",
		provider:          "github",
		providerAccountId: "123456",
		scope:             "read:user,user:email",
		authType:          "oauth",
		tokenType:         "bearer",
	}

	user := FromUserEntity(entity)

	require.Equal(t, expectedUser.name, user.name, "Name should be equal")
	require.Equal(t, expectedUser.email, user.email, "Email should be equal")
	require.Equal(t, expectedUser.image, user.image, "Image should be equal")
	require.Equal(t, expectedUser.accessToken, user.accessToken, "AccessToken should be equal")
	require.Equal(t, expectedUser.provider, user.provider, "Provider should be equal")
	require.Equal(t, expectedUser.providerAccountId, user.providerAccountId, "ProviderAccountId should be equal")
	require.Equal(t, expectedUser.scope, user.scope, "Scope should be equal")
	require.Equal(t, expectedUser.authType, user.authType, "AuthType should be equal")
	require.Equal(t, expectedUser.tokenType, user.tokenType, "TokenType should be equal")

	require.Equal(t, entity.Name, user.Name())
	require.Equal(t, entity.Email, user.Email())
	require.Equal(t, entity.Image, user.Image())
	require.Equal(t, entity.AccessToken, user.AccessToken())
	require.Equal(t, entity.Provider, user.Provider())
	require.Equal(t, entity.ProviderAccountId, user.ProviderAccountId())
	require.Equal(t, entity.Scope, user.Scope())
	require.Equal(t, entity.AuthType, user.AuthType())
	require.Equal(t, entity.TokenType, user.TokenType())

	userEntity := user.ToUserEntity()

	require.Equal(t, userEntity.Name, user.Name())
	require.Equal(t, userEntity.Email, user.Email())
	require.Equal(t, userEntity.Image, user.Image())
	require.Equal(t, userEntity.AccessToken, user.AccessToken())
	require.Equal(t, userEntity.Provider, user.Provider())
	require.Equal(t, userEntity.ProviderAccountId, user.ProviderAccountId())
	require.Equal(t, userEntity.Scope, user.Scope())
	require.Equal(t, userEntity.AuthType, user.AuthType())
	require.Equal(t, userEntity.TokenType, user.TokenType())
}
