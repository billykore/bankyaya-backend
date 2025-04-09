package messaging

import (
	"context"

	"go.bankyaya.org/app/backend/internal/core/entity"
)

// ScheduledTransferProcessor processes scheduled transfer requests.
type ScheduledTransferProcessor interface {
	// Process executes a scheduled transfer.
	// Returns an error if the transfer fails.
	Process(ctx context.Context, req *entity.TransferRequest) error
}
