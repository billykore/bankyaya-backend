package ctxt

import (
	"context"
	"errors"

	"go.bankyaya.org/app/backend/internal/pkg/data"
)

var ErrUserFromContext = errors.New("failed to get user from context")

const UserContextKey = "user"

// ContextWithUser set user data to the ctx context.
func ContextWithUser(ctx context.Context, user data.User) context.Context {
	return context.WithValue(ctx, UserContextKey, user)
}

// UserFromContext gets user data from ctx context.
func UserFromContext(ctx context.Context) (data.User, bool) {
	user, ok := ctx.Value(UserContextKey).(data.User)
	return user, ok
}
