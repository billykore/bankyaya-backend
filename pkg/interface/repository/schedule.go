package repository

import (
	"context"

	"go.bankyaya.org/app/backend/pkg/entity"
)

// ScheduleRepository defines methods for managing transaction scheduler persistence.
type ScheduleRepository interface {
	// CreateSchedule creates a new schedule.
	// Requires a context and a Schedule object to execute.
	// Returns an error if the operation fails.
	CreateSchedule(ctx context.Context, schedule *entity.Schedule) error

	// GetTodaySchedules retrieves the schedules for the current day based on the provided cron expressions.
	// If no cron expressions are provided, it returns all schedules for today.
	// Requires a context and the ID of the schedule to retrieve.
	// Returns the corresponding Schedule object and an error if not found.
	GetTodaySchedules(ctx context.Context, cronExpr ...string) ([]*entity.Schedule, error)

	// GetSchedulesByUserId retrieves schedules by its user ID.
	// Requires a context and the ID of the user.
	// Returns the corresponding Schedule object and an error if not found.
	GetSchedulesByUserId(ctx context.Context, userId int) ([]*entity.Schedule, error)

	// GetScheduleById retrieves a schedule by its ID.
	// Requires a context and the ID of the schedule to retrieve.
	// Returns the corresponding Schedule object and an error if not found.
	GetScheduleById(ctx context.Context, id int) (*entity.Schedule, error)

	// DeleteScheduleByIdAndUserId deletes a schedule by its ID and user ID.
	// This function takes a context and the ID and user ID of the schedule to be deleted.
	// It returns an error if the deletion fails.
	DeleteScheduleByIdAndUserId(ctx context.Context, id, userId int) error
}
