package domain

import "errors"

var (
	// ErrGeneral indicates a general error.
	ErrGeneral = errors.New("something went wrong")

	// ErrEODInProgress indicates that the End of Day (EOD) process is currently in progress.
	ErrEODInProgress = errors.New("EOD process is running")

	// ErrSourceAccountInactive indicates that the source account is inactive.
	ErrSourceAccountInactive = errors.New("source account is inactive")

	// ErrSendEmailFailed is returned when an attempt to send an email fails.
	ErrSendEmailFailed = errors.New("send email failed")

	// ErrUnsuccessfulPayment is returned when a QRIS payment attempt fails.
	ErrUnsuccessfulPayment = errors.New("QRIS payment is unsuccessful")

	// ErrScheduleNotFound is returned when the requested schedule does not exist.
	ErrScheduleNotFound = errors.New("schedule not found")

	// ErrNoScheduleForToday is returned when there is no schedule available for today.
	ErrNoScheduleForToday = errors.New("no schedule for today")

	ErrDestinationAccountInactive = errors.New("destination account is inactive")

	// ErrInvalidSequenceNumber indicates that the sequence number is invalid.
	ErrInvalidSequenceNumber = errors.New("invalid sequence number")

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

	// ErrUserUnauthenticated indicates that the user is not authenticated.
	ErrUserUnauthenticated = errors.New("user is unauthenticated")

	// ErrDeviceNotFound is returned when the requested device cannot be found in the system
	ErrDeviceNotFound = errors.New("device not found")
)
