package otp

import "errors"

var (
	// ErrGeneral indicates a general error.
	ErrGeneral = errors.New("failed to generate OTP")

	// ErrUnauthenticatedUser indicates that the user is not authenticated.
	ErrUnauthenticatedUser = errors.New("unauthenticated user")

	// ErrInvalidOTP indicates that the provided OTP is invalid or does not match.
	ErrInvalidOTP = errors.New("invalid OTP")

	// ErrOTPAlreadyUsed indicates that the provided OTP has already been used and cannot be reused.
	ErrOTPAlreadyUsed = errors.New("OTP already used")

	// ErrOTPExpired indicates that the one-time password (OTP) has expired and is no longer valid.
	ErrOTPExpired = errors.New("OTP expired")
)
