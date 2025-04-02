package messaging

import (
	"context"

	"go.bankyaya.org/app/backend/internal/core/entity"
)

// AutoDebitEventPublisher defines an interface for publishing auto-debit-related events.
// Implementations of this interface should send events to a message broker (e.g., Kafka, RabbitMQ).
type AutoDebitEventPublisher interface {
	// Publish sends an auto-debit event to a message broker.
	// The method takes a context for handling timeouts and cancellations,
	// and an TransferRequest struct containing event details.
	Publish(ctx context.Context, event *entity.TransferRequest) error
}
