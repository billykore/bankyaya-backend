package repo

import (
	"context"

	"go.bankyaya.org/app/backend/domain/scheduler"
	"gorm.io/gorm"
)

type SchedulerRepo struct {
	db *gorm.DB
}

func NewSchedulerRepo(db *gorm.DB) *SchedulerRepo {
	return &SchedulerRepo{db: db}
}

func (s *SchedulerRepo) CreateSchedule(ctx context.Context, schedule *scheduler.Schedule) error {
	res := s.db.WithContext(ctx).Create(schedule)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
