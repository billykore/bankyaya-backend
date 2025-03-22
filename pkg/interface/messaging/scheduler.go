package messaging

import (
	"context"

	"go.bankyaya.org/app/backend/pkg/entity"
)

// AutoDebitEventPublisher defines an interface for publishing auto-debit-related events.
// Implementations of this interface should send events to a message broker (e.g., Kafka, RabbitMQ).
type AutoDebitEventPublisher interface {
	// Publish sends an auto-debit event to a message broker.
	// The method takes a context for handling timeouts and cancellations,
	// and an Event struct containing event details.
	Publish(ctx context.Context, event *entity.Event) error
}
