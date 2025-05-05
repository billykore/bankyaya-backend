package otp

import "context"

// Sender defines the contract for sending OTP messages across various channels.
type Sender interface {
	// Send sends an OTP message using the specified channel. Returns an error if sending fails.
	Send(context.Context, *OTP) error
}
