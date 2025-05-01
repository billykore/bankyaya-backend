package ctxt

import (
	"context"
	"errors"
)

var ErrUserFromContext = errors.New("failed to get user from context")

const UserContextKey = "user"

type User struct {
	Id       int
	CIF      string
	FullName string
	Email    string
}

// ContextWithUser set user data to the ctx context.
func ContextWithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, UserContextKey, user)
}

// UserFromContext gets user data from ctx context.
func UserFromContext(ctx context.Context) (*User, bool) {
	user, ok := ctx.Value(UserContextKey).(*User)
	return user, ok
}
