package scheduler

import (
	"context"
	"errors"

	"go.bankyaya.org/app/backend/pkg/ctxt"
	"go.bankyaya.org/app/backend/pkg/logger"
)

type Service struct {
	log       *logger.Logger
	repo      Repository
	publisher AutoDebitEventPublisher
}

func NewService(log *logger.Logger, repo Repository, publisher AutoDebitEventPublisher) *Service {
	return &Service{
		log:       log,
		repo:      repo,
		publisher: publisher,
	}
}

func (s *Service) Create(ctx context.Context, schedule *Schedule) error {
	user, ok := ctxt.UserFromContext(ctx)
	if !ok {
		return ctxt.ErrUserFromContext
	}
	schedule.UserId = user.Id
	schedule.Status = "inactive"

	err := s.repo.CreateSchedule(ctx, schedule)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Get(ctx context.Context, scheduleId int) (*Schedule, error) {
	schedule, err := s.repo.GetById(ctx, scheduleId)
	if err != nil {
		return nil, err
	}
	return schedule, nil
}

func (s *Service) Delete(ctx context.Context, scheduleId int) error {
	user, ok := ctxt.UserFromContext(ctx)
	if !ok {
		return ctxt.ErrUserFromContext
	}
	err := s.repo.DeleteScheduleByIdAndUserId(ctx, scheduleId, user.Id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) PublishAutoDebit(ctx context.Context) error {
	schedules, err := s.repo.GetTodaySchedules(ctx, "")
	if err != nil && errors.Is(err, ErrNoScheduleForToday) {
		return err
	}
	if err != nil {
		return err
	}

	var publishErr error
	for _, schedule := range schedules {
		err := s.publisher.Publish(ctx, &Event{
			ScheduleId:    schedule.ID,
			Destination:   schedule.Destination,
			Amount:        schedule.Amount,
			AccountNumber: schedule.AccountType,
			UserId:        schedule.UserId,
			Notes:         schedule.Note,
			BankCode:      schedule.BankCode,
			Status:        schedule.Status,
			PhoneNumber:   "",
			DeviceId:      "",
		})
		if err != nil {
			publishErr = errors.Join(publishErr, err)
		}
	}
	if publishErr != nil {
		return err
	}

	return nil
}
