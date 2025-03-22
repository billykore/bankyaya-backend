package handler

import (
	"github.com/labstack/echo/v4"
	"go.bankyaya.org/app/backend/pkg/entity"
	"go.bankyaya.org/app/backend/pkg/framework/http/dto"
	"go.bankyaya.org/app/backend/pkg/framework/http/response"
	"go.bankyaya.org/app/backend/pkg/service"
	"go.bankyaya.org/app/backend/pkg/util/validation"
)

type Scheduler struct {
	va  *validation.Validator
	svc *service.Scheduler
}

func NewScheduler(va *validation.Validator, svc *service.Scheduler) *Scheduler {
	return &Scheduler{
		va:  va,
		svc: svc,
	}
}

// CreateSchedule swaggo annotation.
//
//	@Summary			Create a new schedule
//	@Description		Create a new transaction schedule
//	@Tags				transaction
//	@Accept				json
//	@Produce			json
//	@Param				CreateScheduleRequest	body		scheduler.CreateScheduleRequest	true	"Create a new schedule request"
//	@successResponse	200																					{object}	response.Response
//	@Failure			400						{object}	response.Response
//	@Failure			404						{object}	response.Response
//	@Failure			500						{object}	response.Response
//	@Router				/schedules [post]
func (s *Scheduler) CreateSchedule(ctx echo.Context) error {
	req := new(dto.CreateScheduleRequest)
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	if err := ctx.Validate(req); err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	err := s.svc.Create(ctx.Request().Context(), &entity.Schedule{})
	if err != nil {
		return ctx.JSON(response.Error(err))
	}
	return ctx.JSON(response.SuccessWithoutData())
}

// GetSchedules swaggo annotation.
//
//	@Summary			Get a schedule by ID
//	@Description		Retrieve a transaction schedule by the given ID
//	@Tags				transaction
//	@Accept				json
//	@Produce			json
//	@Param				GetScheduleRequest	path		int	true	"Schedule ID"
//	@successResponse	200												{object}	response.Response
//	@Failure			400					{object}	response.Response
//	@Failure			404					{object}	response.Response
//	@Failure			500					{object}	response.Response
//	@Router				/schedules [get]
func (s *Scheduler) GetSchedules(ctx echo.Context) error {
	resp, err := s.svc.GetSchedules(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(response.Error(err))
	}
	return ctx.JSON(response.Success(resp))
}

// GetSchedule swaggo annotation.
//
//	@Summary			Get a schedule by ID
//	@Description		Retrieve a transaction schedule by the given ID
//	@Tags				transaction
//	@Accept				json
//	@Produce			json
//	@Param				GetScheduleRequest	path		int	true	"Schedule ID"
//	@successResponse	200												{object}	response.Response
//	@Failure			400					{object}	response.Response
//	@Failure			404					{object}	response.Response
//	@Failure			500					{object}	response.Response
//	@Router				/schedules/{id} [get]
func (s *Scheduler) GetSchedule(ctx echo.Context) error {
	req := new(dto.GetScheduleRequest)
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	if err := ctx.Validate(req); err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	resp, err := s.svc.GetById(ctx.Request().Context(), req.ScheduleId)
	if err != nil {
		return ctx.JSON(response.Error(err))
	}
	return ctx.JSON(response.Success(resp))
}

// DeleteSchedule swaggo annotation.
//
//	@Summary			Delete a schedule by ID
//	@Description		Delete a transaction schedule by the given ID
//	@Tags				transaction
//	@Accept				json
//	@Produce			json
//	@Param				DeleteScheduleRequest	path		int	true	"Schedule ID"
//	@successResponse	200														{object}	response.Response
//	@Failure			400						{object}	response.Response
//	@Failure			404						{object}	response.Response
//	@Failure			500						{object}	response.Response
//	@Router				/schedules/{id} [delete]
func (s *Scheduler) DeleteSchedule(ctx echo.Context) error {
	req := new(dto.DeleteScheduleRequest)
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	if err := ctx.Validate(req); err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	err := s.svc.Delete(ctx.Request().Context(), req.ScheduleId)
	if err != nil {
		return ctx.JSON(response.Error(err))
	}
	return ctx.JSON(response.SuccessWithoutData())
}
