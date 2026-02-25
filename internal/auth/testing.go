package auth

import "context"

func ContextWithUserIDForTests(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}
