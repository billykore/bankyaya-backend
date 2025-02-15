package scheduler

import "context"

// Repository defines methods for managing transaction scheduler persistence.
type Repository interface {
	// CreateSchedule creates a new schedule.
	// Requires a context and a Schedule object to execute.
	// Returns an error if the operation fails.
	CreateSchedule(ctx context.Context, schedule *Schedule) error
}
