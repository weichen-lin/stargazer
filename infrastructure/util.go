package infrastructure

import (
	"context"
)

type contextKey string

const ClerkIdKey contextKey = "clerkId"

func GetClerkId(ctx context.Context) (string, bool) {
	userId, ok := ctx.Value("clerkId").(string)
	return userId, ok
}
