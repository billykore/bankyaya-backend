package transfer

import (
	"context"
	"errors"
)

var (
	// ErrEODInProgress indicates that the End of Day (EOD) process is currently in progress.
	ErrEODInProgress = errors.New("EOD process is running")

	// ErrSourceAccountInactive indicates that the source account is inactive.
	ErrSourceAccountInactive = errors.New("source account is inactive")

	// ErrDestinationAccountInactive indicates that the destination account is inactive.
	ErrDestinationAccountInactive = errors.New("destination account is inactive")

	// ErrInvalidSequenceNumber indicates that the sequence number is invalid.
	ErrInvalidSequenceNumber = errors.New("invalid sequence number")

	// ErrSendEmailFailed is returned when an attempt to send an email fails.
	ErrSendEmailFailed = errors.New("send email failed")
)

// Repository defines methods for managing transfer sequence persistence.
type Repository interface {
	// InsertSequence inserts a transfer sequence into the persistence repository.
	// Requires a context and a Sequence object to execute.
	// Returns an error if the operation fails.
	InsertSequence(ctx context.Context, seq *Sequence) error

	// GetSequence retrieves a transfer sequence based on the sequence number.
	// Requires a context and the sequence number as inputs.
	// Returns a Sequence object and an error if retrieval fails.
	GetSequence(ctx context.Context, sequenceNumber string) (*Sequence, error)

	// GetUserById retrieves a user by their unique ID.
	// Requires a context and an integer ID as input parameters.
	// Returns a User object and an error if retrieval fails.
	GetUserById(ctx context.Context, id int) (*User, error)

	// InsertTransaction inserts a transaction into the persistence repository.
	// Requires a context and a Transaction object as input parameters.
	// Returns an error if the operation fails.
	InsertTransaction(ctx context.Context, transaction *Transaction) error
}

// CoreBanking defines methods for core banking operations.
type CoreBanking interface {
	// CheckEOD verifies the current End-of-Day (EOD) process status in the core banking system.
	CheckEOD(ctx context.Context) (*EODData, error)

	// GetAccountDetails retrieves account information for the given account number.
	GetAccountDetails(ctx context.Context, accountNumber string) (*AccountDetails, error)

	// PerformOverbooking executes a transfer between two accounts with the specified amount and remark.
	// It returns an OverbookingResponse and an error if the operation fails.
	PerformOverbooking(ctx context.Context, req *OverbookingRequest) (*OverbookingResponse, error)
}

// ReceiptMailer sends transfer receipt emails.
type ReceiptMailer interface {
	// SendTransferReceipt sends a transfer receipt email with the provided context and email data.
	SendTransferReceipt(ctx context.Context, data *EmailData) error
}
