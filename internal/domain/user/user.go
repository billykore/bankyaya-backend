// Package user provides domain logic related to user operations.
//
// It includes business rules for user authentication, registration, profile management,
// and other identity-related behaviors. This package defines core entities and interfaces
// for handling user data and workflows, ensuring separation of concerns from infrastructure
// (e.g., database, HTTP) and application layers.
package user

import (
	"context"
	"time"
)

// Repository defines methods for managing user persistence.
type Repository interface {
	// GetUserByPhoneNumber retrieves a user from the database by their phone number.
	// Requires a context and a string phone number as input parameters.
	// Returns a User object and an error if retrieval fails.
	GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (*User, error)
}

// TokenService defines methods for creating and validating authorization tokens.
type TokenService interface {
	// Create generates a new token for a given user ID and expiration time.
	Create(user *User, duration time.Duration) (*Token, error)
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
