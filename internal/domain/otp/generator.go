package otp

// Generator defines an interface for generating one-time passwords (OTPs).
type Generator interface {
	// Generate creates a one-time password (OTP) of the specified length and returns it.
	Generate(length int) (string, error)
}
