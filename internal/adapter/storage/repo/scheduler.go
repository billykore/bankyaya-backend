package repo

import (
	"context"
	"errors"

	pkgerrors "go.bankyaya.org/app/backend/internal/core/domain"
	"go.bankyaya.org/app/backend/internal/core/entity"
	"gorm.io/gorm"
)

const statusActive = "active"

type SchedulerRepo struct {
	db *gorm.DB
}

func NewSchedulerRepo(db *gorm.DB) *SchedulerRepo {
	return &SchedulerRepo{db: db}
}

func (r *SchedulerRepo) CreateSchedule(ctx context.Context, schedule *entity.Schedule) error {
	res := r.db.WithContext(ctx).Create(schedule)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *SchedulerRepo) GetTodaySchedules(ctx context.Context, cronExpr ...string) ([]*entity.Schedule, error) {
	schedules := make([]*entity.Schedule, 0)
	res := r.db.WithContext(ctx).
		Where(`"STATUS" = ?`, statusActive).
		Where(`"CRON_TAB" IN (?)`, cronExpr).
		Find(&schedules)
	if res.Error != nil {
		return nil, res.Error
	}
	if len(schedules) == 0 {
		return nil, pkgerrors.ErrNoScheduleForToday
	}
	return schedules, nil
}

func (r *SchedulerRepo) GetSchedulesByUserId(ctx context.Context, userId int) ([]*entity.Schedule, error) {
	schedules := make([]*entity.Schedule, 0)
	res := r.db.WithContext(ctx).
		Where(`"USER_ID" = (?)`, userId).
		Where(`"ACCOUNT_TYPE"= ?`, "personal").
		Find(&schedules)
	if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, pkgerrors.ErrScheduleNotFound
	}
	if res.Error != nil {
		return nil, res.Error
	}
	return schedules, nil
}

func (r *SchedulerRepo) GetScheduleById(ctx context.Context, id int) (*entity.Schedule, error) {
	schedule := new(entity.Schedule)
	res := r.db.WithContext(ctx).
		Where(`"ID" = ?`, id).
		Where(`"ACCOUNT_TYPE"= ?`, "personal").
		First(schedule)
	if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, pkgerrors.ErrScheduleNotFound
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
		Delete(&entity.Schedule{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return pkgerrors.ErrScheduleNotFound
	}
	return nil
}
