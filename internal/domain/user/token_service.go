package user

import "time"

// TokenService defines methods for creating and validating authorization tokens.
type TokenService interface {
	// Create generates a new token for a given user ID and expiration time.
	Create(user *User, duration time.Duration) (*Token, error)
}
