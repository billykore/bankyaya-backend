package ctxt

import (
	"context"
	"errors"
)

type ContextKey string

const UserContextKey ContextKey = "user"

func (c ContextKey) String() string {
	return string(c)
}

var ErrUserFromContext = errors.New("failed to get user from context")

type User struct {
	ID    int
	CIF   string
	Name  string
	Email string
	Phone string
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
