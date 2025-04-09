package service

import (
	"context"
	"errors"

	"go.bankyaya.org/app/backend/internal/core/domain"
	"go.bankyaya.org/app/backend/internal/core/entity"
	"go.bankyaya.org/app/backend/internal/core/port/messaging"
	"go.bankyaya.org/app/backend/internal/core/port/repository"
	"go.bankyaya.org/app/backend/internal/pkg/codes"
	"go.bankyaya.org/app/backend/internal/pkg/ctxt"
	"go.bankyaya.org/app/backend/internal/pkg/logger"
	"go.bankyaya.org/app/backend/internal/pkg/status"
	"go.bankyaya.org/app/backend/internal/pkg/types"
)

const schedulerService = "Scheduler"

type Scheduler struct {
	log                        *logger.Logger
	repo                       repository.ScheduleRepository
	scheduledTransferProcessor messaging.ScheduledTransferProcessor
}

func NewScheduler(log *logger.Logger, repo repository.ScheduleRepository, scheduledTransferProcessor messaging.ScheduledTransferProcessor) *Scheduler {
	return &Scheduler{
		log:                        log,
		repo:                       repo,
		scheduledTransferProcessor: scheduledTransferProcessor,
	}
}

func (s *Scheduler) Create(ctx context.Context, schedule *entity.Schedule) error {
	user, ok := ctxt.UserFromContext(ctx)
	if !ok {
		s.log.ServiceUsecase(schedulerService, "Create").Errorf("UserFromContext: %v", ctxt.ErrUserFromContext)
		return status.Error(codes.Unauthenticated, ctxt.ErrUserFromContext)
	}
	schedule.UserId = user.Id
	schedule.Status = "inactive"

	err := s.repo.CreateSchedule(ctx, schedule)
	if err != nil {
		s.log.ServiceUsecase(schedulerService, "Create").Errorf("CreateSchedule: %v", err)
		return status.Error(codes.Internal, domain.ErrGeneral)
	}

	return nil
}

func (s *Scheduler) GetSchedules(ctx context.Context) ([]*entity.Schedule, error) {
	user, ok := ctxt.UserFromContext(ctx)
	if !ok {
		s.log.ServiceUsecase(schedulerService, "GetSchedules").Errorf("UserFromContext: %v", ctxt.ErrUserFromContext)
		return nil, status.Error(codes.Unauthenticated, ctxt.ErrUserFromContext)
	}

	schedules, err := s.repo.GetSchedulesByUserId(ctx, user.Id)
	if err != nil {
		s.log.ServiceUsecase(schedulerService, "GetSchedules").Errorf("GetSchedulesByUserId: %v", err)
		if errors.Is(err, domain.ErrScheduleNotFound) {
			return nil, status.Error(codes.NotFound, domain.ErrScheduleNotFound)
		}
		return nil, status.Error(codes.Internal, domain.ErrGeneral)
	}

	return schedules, nil
}

func (s *Scheduler) GetById(ctx context.Context, scheduleId int) (*entity.Schedule, error) {
	schedule, err := s.repo.GetScheduleById(ctx, scheduleId)
	if err != nil {
		s.log.ServiceUsecase(schedulerService, "GetById").Errorf("GetScheduleById: %v", err)
		if errors.Is(err, domain.ErrScheduleNotFound) {
			return nil, status.Error(codes.NotFound, err)
		}
		return nil, status.Error(codes.Internal, domain.ErrGeneral)
	}
	return schedule, nil
}

func (s *Scheduler) Delete(ctx context.Context, scheduleId int) error {
	user, ok := ctxt.UserFromContext(ctx)
	if !ok {
		s.log.ServiceUsecase(schedulerService, "Delete").Errorf("UserFromContext: %v", ctxt.ErrUserFromContext)
		return status.Error(codes.Unauthenticated, ctxt.ErrUserFromContext)
	}
	err := s.repo.DeleteScheduleByIdAndUserId(ctx, scheduleId, user.Id)
	if err != nil {
		s.log.ServiceUsecase(schedulerService, "Delete").Errorf("DeleteScheduleByIdAndUserId: %v", err)
		return status.Error(codes.Internal, domain.ErrGeneral)
	}
	return nil
}

func (s *Scheduler) ProcessScheduledTransfer(ctx context.Context) error {
	schedules, err := s.repo.GetTodaySchedules(ctx, "")
	if err != nil {
		s.log.ServiceUsecase(schedulerService, "ProcessScheduledTransfer").Errorf("GetTodaySchedules: %v", err)
		if errors.Is(err, domain.ErrNoScheduleForToday) {
			return status.Error(codes.NotFound, err)
		}
		return status.Error(codes.Internal, domain.ErrGeneral)
	}

	var processErr error
	for _, schedule := range schedules {
		amount, err := types.ParseMoney(schedule.Amount)
		if err != nil {
			processErr = errors.Join(processErr, err)
		}

		err = s.scheduledTransferProcessor.Process(ctx, &entity.TransferRequest{
			ScheduleId:    schedule.ID,
			Destination:   schedule.Destination,
			Amount:        amount,
			AccountNumber: schedule.AccountType,
			UserId:        schedule.UserId,
			Notes:         schedule.Note,
			BankCode:      schedule.BankCode,
			Status:        schedule.Status,
			PhoneNumber:   "",
			DeviceId:      "",
		})
		if err != nil {
			processErr = errors.Join(processErr, err)
		}
	}
	if processErr != nil {
		s.log.ServiceUsecase(schedulerService, "ProcessScheduledTransfer").Errorf("ProcessScheduledTransfer: %v", processErr)
		return status.Error(codes.Internal, err)
	}

	return nil
}
