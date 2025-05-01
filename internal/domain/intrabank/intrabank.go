// Package intrabank provides domain logic for handling intrabank transfers.
//
// It encapsulates business rules related to money transfers within the same bank,
// such as validating account ownership, checking balance sufficiency, and applying
// any relevant transaction policies.
//
// This package defines core entities, interfaces, and use cases that are independent
// of transport layers (e.g., HTTP, gRPC) and persistence mechanisms (e.g., database, cache).
package intrabank

import "context"

// Repository defines methods for managing transfer sequence persistence.
type Repository interface {
	// GetTransactionLimit retrieves the transaction limit for the current day.
	// Requires a context as an input parameter.
	// Returns a Limits object and an error if retrieval fails.
	GetTransactionLimit(ctx context.Context) (*Limits, error)

	// InsertSequence inserts a transfer sequence into the persistence repository.
	// Requires a context and a Sequence object to execute.
	// Returns an error if the operation fails.
	InsertSequence(ctx context.Context, seq *Sequence) error

	// GetSequence retrieves a transfer sequence based on the sequence number.
	// Requires a context and the sequence number as inputs.
	// Returns a Sequence object and an error if retrieval fails.
	GetSequence(ctx context.Context, sequenceNumber string) (*Sequence, error)

	// InsertTransaction inserts a transaction into the persistence repository.
	// Requires a context and a Transaction object as input parameters.
	// Returns an error if the operation fails.
	InsertTransaction(ctx context.Context, transaction *Transaction) error
}

// CoreBanking defines methods for core banking operations.
type CoreBanking interface {
	// GetCoreStatus gets the current status of the core banking system.
	GetCoreStatus(ctx context.Context) (*CoreStatus, error)

	// GetAccountDetails retrieves account information for the given account number.
	GetAccountDetails(ctx context.Context, accountNumber string) (*Account, error)

	// PerformOverbooking executes a transfer between two accounts with the specified amount and remark.
	// It returns an OverbookingResponse and an error if the operation fails.
	PerformOverbooking(ctx context.Context, req *OverbookingInput) (*OverbookingResult, error)
}

// SequenceGenerator defines an interface for generating unique sequences.
type SequenceGenerator interface {
	// Generate produces the unique sequence as a string
	// and error if the sequence cannot be generated.
	Generate() (string, error)
}

// ReceiptMailer sends transfer receipt emails.
type ReceiptMailer interface {
	// SendReceipt sends a transfer receipt email with the provided context and email data.
	SendReceipt(ctx context.Context, data *EmailData) error
}

// Notifier sends intrabank transfer notifications to users.
type Notifier interface {
	// Notify sends a transfer notification to the specified user.
	Notify(ctx context.Context, notification *Notification) error
}
