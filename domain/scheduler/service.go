package scheduler

import (
	"context"

	"go.bankyaya.org/app/backend/pkg/codes"
	"go.bankyaya.org/app/backend/pkg/cron"
	"go.bankyaya.org/app/backend/pkg/ctxt"
	"go.bankyaya.org/app/backend/pkg/datetime"
	"go.bankyaya.org/app/backend/pkg/logger"
	"go.bankyaya.org/app/backend/pkg/status"
)

type Service struct {
	log  *logger.Logger
	repo Repository
}

func NewService(log *logger.Logger, repo Repository) *Service {
	return &Service{
		log:  log,
		repo: repo,
	}
}

func (s *Service) Create(ctx context.Context, req *CreateScheduleRequest) error {
	user, ok := ctxt.UserFromContext(ctx)
	if !ok {
		s.log.DomainUsecase("scheduler", "Create").Error(ctxt.ErrUserFromContext)
		return status.Error(codes.Unauthenticated, "user not found")
	}
	startDate, err := req.ParseStartDate()
	if err != nil {
		s.log.DomainUsecase("scheduler", "Create").Error(err)
		return status.Error(codes.Internal, "failed to create schedule")
	}
	err = s.repo.CreateSchedule(ctx, &Schedule{
		UserId:             user.Id,
		SakuId:             req.SakuId,
		Destination:        req.Destination,
		DestinationName:    req.DestinationName,
		Amount:             req.Amount.String(),
		Note:               req.Note,
		BankCode:           req.BankCode,
		TransactionType:    req.TransactionMethod,
		TransactionMethod:  req.TransactionMethod,
		TransactionPurpose: req.PurposeType,
		Frequency:          req.Frequency,
		StartDate:          startDate,
		AutoDebet:          false,
		CrontabSchedule:    parseCrontab(req.CronFrequency(), req.Day, req.Date),
		Status:             "active",
		AccountType:        req.AccountType,
		BiFastCode:         req.BiFastCode,
	})
	if err != nil {
		s.log.DomainUsecase("scheduler", "Create").Error(err)
		return status.Error(codes.Internal, "failed to create schedule")
	}
	return nil
}

func parseCrontab(freq cron.Frequency, day string, date int) string {
	return cron.ParseScheduleExpr(freq, datetime.IndonesianWeekdayValue(day), date)
}
