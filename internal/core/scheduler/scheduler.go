package scheduler

import (
	"context"
	"errors"
)

// ErrScheduleNotFound is returned when the requested schedule does not exist.
var ErrScheduleNotFound = errors.New("schedule not found")

// ErrNoScheduleForToday is returned when there is no schedule available for today.
var ErrNoScheduleForToday = errors.New("no schedule for today")

// Repository defines methods for managing transaction scheduler persistence.
type Repository interface {
	// CreateSchedule creates a new schedule.
	// Requires a context and a Schedule object to execute.
	// Returns an error if the operation fails.
	CreateSchedule(ctx context.Context, schedule *Schedule) error

	// GetTodaySchedules retrieves the schedules for the current day based on the provided cron expressions.
	// If no cron expressions are provided, it returns all schedules for today.
	// Requires a context and the ID of the schedule to retrieve.
	// Returns the corresponding Schedule object and an error if not found.
	GetTodaySchedules(ctx context.Context, cronExpr ...string) ([]*Schedule, error)

	// GetById retrieves a schedule by its ID.
	// Requires a context and the ID of the schedule to retrieve.
	// Returns the corresponding Schedule object and an error if not found.
	GetById(ctx context.Context, id int) (*Schedule, error)

	// DeleteScheduleByIdAndUserId deletes a schedule by its ID and user ID.
	// This function takes a context and the ID and user ID of the schedule to be deleted.
	// It returns an error if the deletion fails.
	DeleteScheduleByIdAndUserId(ctx context.Context, id, userId int) error
}

// AutoDebitEventPublisher defines an interface for publishing auto-debit-related events.
// Implementations of this interface should send events to a message broker (e.g., Kafka, RabbitMQ).
type AutoDebitEventPublisher interface {
	// Publish sends an auto-debit event to a message broker.
	// The method takes a context for handling timeouts and cancellations,
	// and an Event struct containing event details.
	Publish(ctx context.Context, event *Event) error
}
