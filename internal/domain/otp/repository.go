package otp

import "context"

// Repository defines methods to persist and retrieve OTP entities.
type Repository interface {
	// Save persists an OTP entity into the data storage.
	// Returns an error if the operation fails.
	Save(context.Context, *OTP) error

	// Get retrieves an OTP entity from the data storage by its ID.
	// It takes a context and an ID as parameters.
	// Returns the OTP entity if found,
	// or an error if the operation fails or the OTP is not found.
	Get(ctx context.Context, id int64) (*OTP, error)

	// Update updates an existing OTP entity in the data storage.
	// Returns an error if the operation fails.
	Update(ctx context.Context, otp *OTP) error
}
