package repository

import (
	"context"

	"go.bankyaya.org/app/backend/pkg/entity"
)

// TransferRepository defines methods for managing transfer sequence persistence.
type TransferRepository interface {
	// InsertSequence inserts a transfer sequence into the persistence repository.
	// Requires a context and a Sequence object to execute.
	// Returns an error if the operation fails.
	InsertSequence(ctx context.Context, seq *entity.Sequence) error

	// GetSequence retrieves a transfer sequence based on the sequence number.
	// Requires a context and the sequence number as inputs.
	// Returns a Sequence object and an error if retrieval fails.
	GetSequence(ctx context.Context, sequenceNumber string) (*entity.Sequence, error)

	// GetUserById retrieves a user by their unique ID.
	// Requires a context and an integer ID as input parameters.
	// Returns a User object and an error if retrieval fails.
	GetUserById(ctx context.Context, id int) (*entity.User, error)

	// InsertTransaction inserts a transaction into the persistence repository.
	// Requires a context and a Transaction object as input parameters.
	// Returns an error if the operation fails.
	InsertTransaction(ctx context.Context, transaction *entity.Transaction) error
}
