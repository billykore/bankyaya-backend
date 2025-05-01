package user

import (
	"context"
)

// Repository defines methods for managing user persistence.
type Repository interface {
	// GetUserByPhoneNumber retrieves a user from the database by their phone number.
	// Requires a context and a string phone number as input parameters.
	// Returns a User object and an error if retrieval fails.
	GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (*User, error)
}
