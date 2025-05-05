package intrabank

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

	// ErrDestinationAccountInactive indicates that the destination account is inactive.
	ErrDestinationAccountInactive = errors.New("destination account is inactive")

	// ErrInvalidSequenceNumber indicates that the sequence number is invalid.
	ErrInvalidSequenceNumber = errors.New("invalid sequence number")

	// ErrUnauthenticatedUser indicates that the user is not authenticated.
	ErrUnauthenticatedUser = errors.New("unauthenticated user")

	// ErrFailedParseMoney is returned when the system fails to parse a money value
	// from a string or an unsupported format.
	ErrFailedParseMoney = errors.New("failed to parse money")

	// ErrInvalidAmount is returned when the request amount is invalid.
	ErrInvalidAmount = errors.New("invalid request amount")

	// ErrNotifyFailed is returned when the notification fails to send.
	ErrNotifyFailed = errors.New("notify failed")
)
