package user

import (
	"context"
	"errors"
	"time"
)

var (
	// ErrDeviceIsBlacklisted is returned when an operation is attempted on a blacklisted device.
	ErrDeviceIsBlacklisted = errors.New("device is blacklisted")

	// ErrInvalidDevice is returned when an operation encounters an invalid device.
	ErrInvalidDevice = errors.New("invalid device")

	// ErrDeviceNotFound is returned when the requested device cannot be found in the system.
	ErrDeviceNotFound = errors.New("device not found")

	// ErrNotFound is returned when the requested user cannot be found in the system.
	ErrNotFound = errors.New("user not found")

	// ErrInvalidPassword is returned when the provided password does not match the stored hash.
	ErrInvalidPassword = errors.New("invalid password")
)

// Repository defines methods for managing user persistence.
type Repository interface {
	// GetUserDataByPhoneNumber retrieves a user from the database by their phone number.
	// Requires a context and a string phone number as input parameters.
	// Returns a User object and an error if retrieval fails.
	GetUserDataByPhoneNumber(ctx context.Context, phoneNumber string) (*User, error)

	// GetDeviceById retrieves a device from the database by its unique ID.
	// Requires a context and a strings device ID as input parameters.
	// Returns a Device object and an error if retrieval fails.
	GetDeviceById(ctx context.Context, deviceId string) (*Device, error)
}

// TokenService defines methods for creating and validating authorization tokens.
type TokenService interface {
	// Create generates a new token for a given user ID and expiration time.
	Create(user *User, duration time.Duration) (Token, error)
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
