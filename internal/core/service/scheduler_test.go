package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.bankyaya.org/app/backend/internal/core/domain"
	"go.bankyaya.org/app/backend/internal/core/entity"
	"go.bankyaya.org/app/backend/internal/core/port/messaging/mock"
	"go.bankyaya.org/app/backend/internal/core/port/repository/mock"
	"go.bankyaya.org/app/backend/internal/pkg/codes"
	"go.bankyaya.org/app/backend/internal/pkg/ctxt"
	"go.bankyaya.org/app/backend/internal/pkg/data"
	"go.bankyaya.org/app/backend/internal/pkg/logger"
	"go.bankyaya.org/app/backend/internal/pkg/status"
)

func TestCreateScheduleSuccess(t *testing.T) {
	var (
		repoMock      = repomock.NewScheduleRepository(t)
		processorMock = messagingmock.NewScheduledTransferProcessor(t)
		svc           = NewScheduler(logger.New(), repoMock, processorMock)
		ctx           = ctxt.ContextWithUser(context.Background(), data.User{
			Id:       123,
			CIF:      "1234567",
			FullName: "Olivia Rodrigo",
			Email:    "olivia@gmail.com",
		})
	)

	repoMock.EXPECT().CreateSchedule(mock.Anything, &entity.Schedule{
		UserId: 123,
		Status: "inactive",
	}).
		Return(nil)

	err := svc.Create(ctx, &entity.Schedule{
		UserId: 123,
	})

	assert.NoError(t, err)

	repoMock.AssertExpectations(t)
	processorMock.AssertExpectations(t)
}

func TestCreateScheduleFailed_GetUserFromContextFailed(t *testing.T) {
	var (
		repoMock      = repomock.NewScheduleRepository(t)
		processorMock = messagingmock.NewScheduledTransferProcessor(t)
		svc           = NewScheduler(logger.New(), repoMock, processorMock)
		ctx           = context.Background()
	)

	err := svc.Create(ctx, &entity.Schedule{
		UserId: 123,
	})

	assert.Equal(t, status.Error(codes.Unauthenticated, ctxt.ErrUserFromContext), err)

	repoMock.AssertExpectations(t)
	processorMock.AssertExpectations(t)
}

func TestCreateScheduleFailed_SaveScheduleFailed(t *testing.T) {
	var (
		repoMock      = repomock.NewScheduleRepository(t)
		processorMock = messagingmock.NewScheduledTransferProcessor(t)
		svc           = NewScheduler(logger.New(), repoMock, processorMock)
		ctx           = ctxt.ContextWithUser(context.Background(), data.User{
			Id:       123,
			CIF:      "1234567",
			FullName: "Olivia Rodrigo",
			Email:    "olivia@gmail.com",
		})
	)

	repoMock.EXPECT().CreateSchedule(mock.Anything, mock.Anything).
		Return(errors.New("some error"))

	err := svc.Create(ctx, &entity.Schedule{
		UserId: 123,
	})

	assert.Equal(t, status.Error(codes.Internal, domain.ErrGeneral), err)

	repoMock.AssertExpectations(t)
	processorMock.AssertExpectations(t)
}

func TestGetSchedulesSuccess(t *testing.T) {
	var (
		repoMock      = repomock.NewScheduleRepository(t)
		processorMock = messagingmock.NewScheduledTransferProcessor(t)
		svc           = NewScheduler(logger.New(), repoMock, processorMock)
		ctx           = ctxt.ContextWithUser(context.Background(), data.User{
			Id:       123,
			CIF:      "1234567",
			FullName: "Olivia Rodrigo",
			Email:    "olivia@gmail.com",
		})
	)

	repoMock.EXPECT().GetSchedulesByUserId(mock.Anything, 123).
		Return(make([]*entity.Schedule, 3), nil)

	schedules, err := svc.GetSchedules(ctx)

	assert.NoError(t, err)
	assert.Len(t, schedules, 3)
	assert.Equal(t, make([]*entity.Schedule, 3), schedules)

	repoMock.AssertExpectations(t)
	processorMock.AssertExpectations(t)
}

func TestGetSchedulesFailed_GetUserFromContextFailed(t *testing.T) {
	var (
		repoMock      = repomock.NewScheduleRepository(t)
		processorMock = messagingmock.NewScheduledTransferProcessor(t)
		svc           = NewScheduler(logger.New(), repoMock, processorMock)
		ctx           = context.Background()
	)

	schedules, err := svc.GetSchedules(ctx)

	assert.Nil(t, schedules)
	assert.Equal(t, status.Error(codes.Unauthenticated, ctxt.ErrUserFromContext), err)

	repoMock.AssertExpectations(t)
	processorMock.AssertExpectations(t)
}

func TestGetSchedulesFailed_GetSchedulesByUserIdNotFound(t *testing.T) {
	var (
		repoMock      = repomock.NewScheduleRepository(t)
		processorMock = messagingmock.NewScheduledTransferProcessor(t)
		svc           = NewScheduler(logger.New(), repoMock, processorMock)
		ctx           = ctxt.ContextWithUser(context.Background(), data.User{
			Id:       123,
			CIF:      "1234567",
			FullName: "Olivia Rodrigo",
			Email:    "olivia@gmail.com",
		})
	)

	repoMock.EXPECT().GetSchedulesByUserId(mock.Anything, 123).
		Return(nil, domain.ErrScheduleNotFound)

	schedules, err := svc.GetSchedules(ctx)

	assert.Nil(t, schedules)
	assert.Equal(t, status.Error(codes.NotFound, domain.ErrScheduleNotFound), err)

	repoMock.AssertExpectations(t)
	processorMock.AssertExpectations(t)
}

func TestGetSchedulesFailed_GetSchedulesByUserIdFailed(t *testing.T) {
	var (
		repoMock      = repomock.NewScheduleRepository(t)
		processorMock = messagingmock.NewScheduledTransferProcessor(t)
		svc           = NewScheduler(logger.New(), repoMock, processorMock)
		ctx           = ctxt.ContextWithUser(context.Background(), data.User{
			Id:       123,
			CIF:      "1234567",
			FullName: "Olivia Rodrigo",
			Email:    "olivia@gmail.com",
		})
	)

	repoMock.EXPECT().GetSchedulesByUserId(mock.Anything, 123).
		Return(nil, errors.New("some error"))

	schedules, err := svc.GetSchedules(ctx)

	assert.Nil(t, schedules)
	assert.Equal(t, status.Error(codes.Internal, domain.ErrGeneral), err)

	repoMock.AssertExpectations(t)
	processorMock.AssertExpectations(t)
}

func TestGetScheduleByIdSuccess(t *testing.T) {
	var (
		repoMock      = repomock.NewScheduleRepository(t)
		processorMock = messagingmock.NewScheduledTransferProcessor(t)
		svc           = NewScheduler(logger.New(), repoMock, processorMock)
		ctx           = ctxt.ContextWithUser(context.Background(), data.User{
			Id:       123,
			CIF:      "1234567",
			FullName: "Olivia Rodrigo",
			Email:    "olivia@gmail.com",
		})
	)

	repoMock.EXPECT().GetScheduleById(mock.Anything, 1).
		Return(&entity.Schedule{ID: 1}, nil)

	schedule, err := svc.GetById(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, &entity.Schedule{ID: 1}, schedule)

	repoMock.AssertExpectations(t)
	processorMock.AssertExpectations(t)
}

func TestGetScheduleByIdFailed_GetScheduleByIdNotFound(t *testing.T) {
	var (
		repoMock      = repomock.NewScheduleRepository(t)
		processorMock = messagingmock.NewScheduledTransferProcessor(t)
		svc           = NewScheduler(logger.New(), repoMock, processorMock)
		ctx           = ctxt.ContextWithUser(context.Background(), data.User{
			Id:       123,
			CIF:      "1234567",
			FullName: "Olivia Rodrigo",
			Email:    "olivia@gmail.com",
		})
	)

	repoMock.EXPECT().GetScheduleById(mock.Anything, 0).
		Return(nil, domain.ErrScheduleNotFound)

	schedule, err := svc.GetById(ctx, 0)

	assert.Nil(t, schedule)
	assert.Equal(t, status.Error(codes.NotFound, domain.ErrScheduleNotFound), err)

	repoMock.AssertExpectations(t)
	processorMock.AssertExpectations(t)
}

func TestGetScheduleByIdFailed_GetScheduleByIdFailed(t *testing.T) {
	var (
		repoMock      = repomock.NewScheduleRepository(t)
		processorMock = messagingmock.NewScheduledTransferProcessor(t)
		svc           = NewScheduler(logger.New(), repoMock, processorMock)
		ctx           = ctxt.ContextWithUser(context.Background(), data.User{
			Id:       123,
			CIF:      "1234567",
			FullName: "Olivia Rodrigo",
			Email:    "olivia@gmail.com",
		})
	)

	repoMock.EXPECT().GetScheduleById(mock.Anything, 0).
		Return(nil, errors.New("some error"))

	schedule, err := svc.GetById(ctx, 0)

	assert.Nil(t, schedule)
	assert.Equal(t, status.Error(codes.Internal, domain.ErrGeneral), err)

	repoMock.AssertExpectations(t)
	processorMock.AssertExpectations(t)
}

func TestDeleteScheduleSuccess(t *testing.T) {
	var (
		repoMock      = repomock.NewScheduleRepository(t)
		processorMock = messagingmock.NewScheduledTransferProcessor(t)
		svc           = NewScheduler(logger.New(), repoMock, processorMock)
		ctx           = ctxt.ContextWithUser(context.Background(), data.User{
			Id:       123,
			CIF:      "1234567",
			FullName: "Olivia Rodrigo",
			Email:    "olivia@gmail.com",
		})
	)

	repoMock.EXPECT().DeleteScheduleByIdAndUserId(mock.Anything, 1, 123).
		Return(nil)

	err := svc.Delete(ctx, 1)

	assert.NoError(t, err)

	repoMock.AssertExpectations(t)
	processorMock.AssertExpectations(t)
}

func TestDeleteScheduleFailed_GetUserFromContextFailed(t *testing.T) {
	var (
		repoMock      = repomock.NewScheduleRepository(t)
		processorMock = messagingmock.NewScheduledTransferProcessor(t)
		svc           = NewScheduler(logger.New(), repoMock, processorMock)
		ctx           = context.Background()
	)

	err := svc.Delete(ctx, 1)

	assert.Equal(t, status.Error(codes.Unauthenticated, ctxt.ErrUserFromContext), err)

	repoMock.AssertExpectations(t)
	processorMock.AssertExpectations(t)
}

func TestDeleteScheduleFailed_GetScheduleByIdFailed(t *testing.T) {
	var (
		repoMock      = repomock.NewScheduleRepository(t)
		processorMock = messagingmock.NewScheduledTransferProcessor(t)
		svc           = NewScheduler(logger.New(), repoMock, processorMock)
		ctx           = ctxt.ContextWithUser(context.Background(), data.User{
			Id:       123,
			CIF:      "1234567",
			FullName: "Olivia Rodrigo",
			Email:    "olivia@gmail.com",
		})
	)

	repoMock.EXPECT().DeleteScheduleByIdAndUserId(mock.Anything, 0, 123).
		Return(errors.New("some error"))

	err := svc.Delete(ctx, 0)

	assert.Equal(t, status.Error(codes.Internal, domain.ErrGeneral), err)

	repoMock.AssertExpectations(t)
	processorMock.AssertExpectations(t)
}

func TestProcessScheduledTransferSuccess(t *testing.T) {
	var (
		repoMock      = repomock.NewScheduleRepository(t)
		processorMock = messagingmock.NewScheduledTransferProcessor(t)
		svc           = NewScheduler(logger.New(), repoMock, processorMock)
		ctx           = ctxt.ContextWithUser(context.Background(), data.User{
			Id:       123,
			CIF:      "1234567",
			FullName: "Olivia Rodrigo",
			Email:    "olivia@gmail.com",
		})
	)

	repoMock.EXPECT().GetTodaySchedules(mock.Anything, mock.Anything).
		Return([]*entity.Schedule{
			{Amount: "10000"},
		}, nil)

	processorMock.EXPECT().Process(ctx, mock.Anything).
		Return(nil)

	err := svc.ProcessScheduledTransfer(ctx)

	assert.NoError(t, err)

	repoMock.AssertExpectations(t)
}
