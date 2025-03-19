package repo

import (
	"context"
	"errors"

	scheduler2 "go.bankyaya.org/app/backend/internal/core/scheduler"
	"gorm.io/gorm"
)

const statusActive = "active"

type SchedulerRepo struct {
	db *gorm.DB
}

func NewSchedulerRepo(db *gorm.DB) *SchedulerRepo {
	return &SchedulerRepo{db: db}
}

func (r *SchedulerRepo) CreateSchedule(ctx context.Context, schedule *scheduler2.Schedule) error {
	res := r.db.WithContext(ctx).Create(schedule)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *SchedulerRepo) GetTodaySchedules(ctx context.Context, cronExpr ...string) ([]*scheduler2.Schedule, error) {
	schedules := make([]*scheduler2.Schedule, 0)
	res := r.db.WithContext(ctx).
		Where(`"STATUS" = ?`, statusActive).
		Where(`"CRON_TAB" IN (?)`, cronExpr).
		Find(&schedules)
	if res.Error != nil {
		return nil, res.Error
	}
	if len(schedules) == 0 {
		return nil, scheduler2.ErrNoScheduleForToday
	}
	return schedules, nil
}

func (r *SchedulerRepo) GetById(ctx context.Context, id int) (*scheduler2.Schedule, error) {
	schedule := new(scheduler2.Schedule)
	res := r.db.WithContext(ctx).
		Where(`"ID" = ?`, id).
		First(schedule)
	if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, scheduler2.ErrScheduleNotFound
	}
	if res.Error != nil {
		return nil, res.Error
	}
	return schedule, nil
}

func (r *SchedulerRepo) DeleteScheduleByIdAndUserId(ctx context.Context, id, userId int) error {
	res := r.db.WithContext(ctx).
		Where(`"ID" = ?`, id).
		Where(`"USER_ID" = ?`, userId).
		Delete(&scheduler2.Schedule{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return scheduler2.ErrScheduleNotFound
	}
	return nil
}
