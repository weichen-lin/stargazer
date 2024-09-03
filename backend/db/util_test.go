package db

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWithEmail(t *testing.T) {
	t.Run("valid email", func(t *testing.T) {
		ctx := context.Background()
		newCtx, err := WithEmail(ctx, "valid@example.com")

		require.NoError(t, err)
		require.NotNil(t, newCtx)

		email, ok := EmailFromContext(newCtx)
		require.True(t, ok)
		require.Equal(t, "valid@example.com", email)
	})

	t.Run("invalid email format", func(t *testing.T) {
		ctx := context.Background()
		_, err := WithEmail(ctx, "invalid_email")

		require.Error(t, err)
		require.EqualError(t, err, "invalid email format")
	})

	t.Run("email too short", func(t *testing.T) {
		ctx := context.Background()
		_, err := WithEmail(ctx, "ab@d")

		require.Error(t, err)
		require.EqualError(t, err, "email length should be between 5 and 254 characters")
	})

	t.Run("email too long", func(t *testing.T) {
		longEmail := "a" + strings.Repeat("a", 254) + "@example.com"
		ctx := context.Background()
		_, err := WithEmail(ctx, longEmail)

		require.Error(t, err)
		require.EqualError(t, err, "email length should be between 5 and 254 characters")
	})
}

func TestEmailFromContext(t *testing.T) {
	t.Run("empty context", func(t *testing.T) {
		email, ok := EmailFromContext(context.Background())

		require.False(t, ok)
		require.Equal(t, "", email)
	})

	t.Run("no email in context", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), EmailKey(""), nil)
		email, ok := EmailFromContext(ctx)

		require.False(t, ok)
		require.Equal(t, "", email)
	})
}
