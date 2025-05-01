package user

import "errors"

var (
	// ErrDeviceIsBlacklisted is returned when an operation is attempted on a blacklisted device.
	ErrDeviceIsBlacklisted = errors.New("device is blacklisted")

	// ErrInvalidDevice is returned when an operation encounters an invalid device.
	ErrInvalidDevice = errors.New("invalid device")

	// ErrCreateTokenFailed is returned when the authentication token is failed to create.
	ErrCreateTokenFailed = errors.New("create token failed")

	// ErrUserNotFound is returned when the requested user cannot be found in the system.
	ErrUserNotFound = errors.New("user not found")

	// ErrInvalidPassword is returned when the provided password does not match the stored hash.
	ErrInvalidPassword = errors.New("invalid password")
)
