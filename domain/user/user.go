package user

import (
	"context"
	"errors"
)

// ErrInvalidDevice is returned when an operation encounters an invalid device.
var ErrInvalidDevice = errors.New("invalid device")

// Repository defines methods for managing user persistence.
type Repository interface {
	// GetUserByPhoneNumber retrieves a user from the database by their phone number.
	// Requires a context and a string phone number as input parameters.
	// Returns a User object and an error if retrieval fails.
	GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (*User, error)

	// GetDeviceById retrieves a device from the database by its unique ID.
	// Requires a context and a strings device ID as input parameters.
	// Returns a Device object and an error if retrieval fails.
	GetDeviceById(ctx context.Context, deviceId string) (*Device, error)
}
