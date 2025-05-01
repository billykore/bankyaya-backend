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
