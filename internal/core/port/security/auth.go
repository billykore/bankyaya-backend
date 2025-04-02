package security

import (
	"time"

	"go.bankyaya.org/app/backend/internal/core/entity"
)

// TokenService defines methods for creating and validating authorization tokens.
type TokenService interface {
	// Create generates a new token for a given user ID and expiration time.
	Create(user *entity.User, duration time.Duration) (entity.Token, error)
}

// PasswordHasher defines an interface for hashing and verifying passwords.
type PasswordHasher interface {
	// Hash generates a hashed representation of the given password.
	// password: The plain-text password to hash.
	// Returns the hashed password string and an error if hashing fails.
	Hash(password string) (string, error)

	// Compare checks whether the provided plain-text password matches the given hashed password.
	// password: The plain-text password input. hashed: The stored hashed password.
	// Returns true if the password matches, otherwise false.
	Compare(password, hashed string) bool
}
