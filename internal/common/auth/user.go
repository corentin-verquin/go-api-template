package auth

import "context"

const usernameKey = "username"

func WithUsername(ctx context.Context, username string) context.Context {
	return context.WithValue(ctx, usernameKey, username)
}

func Username(ctx context.Context) string {
	username, ok := ctx.Value(usernameKey).(string)
	if !ok {
		return "unknown"
	}
	return username
}
