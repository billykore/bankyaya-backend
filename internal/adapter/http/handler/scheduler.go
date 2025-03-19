package handler

import (
	"github.com/labstack/echo/v4"
	"go.bankyaya.org/app/backend/internal/adapter/http/response"
	"go.bankyaya.org/app/backend/internal/application/scheduler"
)

type SchedulerHandler struct {
	uc *scheduler.Usecase
}

func NewSchedulerHandler(uc *scheduler.Usecase) *SchedulerHandler {
	return &SchedulerHandler{
		uc: uc,
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
//	@Router				/scheduler [post]
func (h *SchedulerHandler) CreateSchedule(ctx echo.Context) error {
	req := new(scheduler.CreateScheduleRequest)
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	err := h.uc.CreateSchedule(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(response.Error(err))
	}
	return ctx.JSON(response.SuccessWithoutData())
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
//	@Router				/scheduler/{id} [get]
func (h *SchedulerHandler) GetSchedule(ctx echo.Context) error {
	req := new(scheduler.GetScheduleRequest)
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	resp, err := h.uc.GetSchedule(ctx.Request().Context(), req)
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
//	@Router				/scheduler/{id} [delete]
func (h *SchedulerHandler) DeleteSchedule(ctx echo.Context) error {
	req := new(scheduler.DeleteScheduleRequest)
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	err := h.uc.DeleteSchedule(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(response.Error(err))
	}
	return ctx.JSON(response.SuccessWithoutData())
}
