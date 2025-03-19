package scheduler

import (
	"context"
	"errors"

	"go.bankyaya.org/app/backend/internal/core/scheduler"
	"go.bankyaya.org/app/backend/pkg/codes"
	"go.bankyaya.org/app/backend/pkg/cron"
	"go.bankyaya.org/app/backend/pkg/ctxt"
	"go.bankyaya.org/app/backend/pkg/datetime"
	"go.bankyaya.org/app/backend/pkg/logger"
	"go.bankyaya.org/app/backend/pkg/status"
	"go.bankyaya.org/app/backend/pkg/validation"
)

type Usecase struct {
	va  *validation.Validator
	log *logger.Logger
	svc *scheduler.Service
}

func NewUsecase(va *validation.Validator, log *logger.Logger, svc *scheduler.Service) *Usecase {
	return &Usecase{
		va:  va,
		log: log,
		svc: svc,
	}
}

func (uc *Usecase) CreateSchedule(ctx context.Context, req *CreateScheduleRequest) error {
	if err := uc.va.Validate(req); err != nil {
		uc.log.Usecase("Scheduler", "Inquiry").Errorf("Validate error: %v", err)
		return status.Error(codes.BadRequest, "Bad Request")
	}
	err := uc.svc.Create(ctx, &scheduler.Schedule{
		SakuId:             req.SakuId,
		Destination:        req.Destination,
		DestinationName:    req.DestinationName,
		Amount:             req.Amount.String(),
		Note:               req.Note,
		BankCode:           req.BankCode,
		TransactionType:    req.TransactionType(),
		TransactionMethod:  req.TransactionMethod,
		TransactionPurpose: req.PurposeType,
		Frequency:          req.Frequency,
		StartDate:          req.ParseStartDate(),
		CrontabSchedule:    parseCrontab(cron.Frequency(req.Frequency), req.Day, req.Date),
		AccountType:        req.AccountType,
		BIFastCode:         req.BiFastCode,
	})
	if err != nil {
		uc.log.Usecase("Scheduler", "Create").Errorf("Create schedule error: %v", err)
		return status.Error(codes.Internal, "Create schedule error")
	}
	return nil
}

func parseCrontab(freq cron.Frequency, day string, date int) string {
	return cron.ParseScheduleExpr(freq, datetime.IndonesianWeekdayValue(day), date)
}

func (uc *Usecase) GetSchedule(ctx context.Context, req *GetScheduleRequest) (*GetScheduleResponse, error) {
	if err := uc.va.Validate(req); err != nil {
		uc.log.Usecase("scheduler", "GetSchedule").Errorf("Validate error: %v", err)
		return nil, status.Error(codes.BadRequest, "Bad Request")
	}

	result, err := uc.svc.Get(ctx, req.ScheduleId)
	if err != nil {
		uc.log.Usecase("scheduler", "GetSchedule").Errorf("Get schedule error: %v", err)
		if errors.Is(err, ctxt.ErrUserFromContext) {
			return nil, status.Errorf(codes.Unauthenticated, "User unauthenticated")
		}
		return nil, status.Error(codes.Internal, "Get schedule failed")
	}

	return &GetScheduleResponse{
		UserId:             result.UserId,
		SakuId:             result.SakuId,
		Destination:        result.Destination,
		DestinationName:    result.DestinationName,
		Amount:             result.Amount,
		Note:               result.Note,
		BankCode:           result.BankCode,
		TransactionType:    result.TransactionType,
		TransactionMethod:  result.TransactionMethod,
		TransactionPurpose: result.TransactionPurpose,
		Frequency:          result.Frequency,
		StartDate:          result.StartDate,
		CrontabSchedule:    result.CrontabSchedule,
		Status:             result.Status,
		AccountType:        result.AccountType,
		BIFastCode:         result.BIFastCode,
	}, nil
}

func (uc *Usecase) DeleteSchedule(ctx context.Context, req *DeleteScheduleRequest) error {
	if err := uc.va.Validate(req); err != nil {
		uc.log.Usecase("scheduler", "DeleteSchedule").Errorf("Validate error: %v", err)
		return status.Error(codes.BadRequest, "Bad Request")
	}

	err := uc.svc.Delete(ctx, req.ScheduleId)
	if err != nil {
		uc.log.Usecase("scheduler", "DeleteSchedule").Errorf("Delete schedule error: %v", err)
		if errors.Is(err, ctxt.ErrUserFromContext) {
			return status.Errorf(codes.Unauthenticated, "User unauthenticated")
		}
		return status.Error(codes.Internal, "Delete schedule failed")
	}

	return nil
}
