package handler

import (
	"github.com/labstack/echo/v4"
	"go.bankyaya.org/app/backend/domain/scheduler"
	"go.bankyaya.org/app/backend/pkg/response"
	"go.bankyaya.org/app/backend/pkg/validation"
)

type SchedulerHandler struct {
	va  *validation.Validator
	svc *scheduler.Service
}

func NewSchedulerHandler(va *validation.Validator, svc *scheduler.Service) *SchedulerHandler {
	return &SchedulerHandler{
		va:  va,
		svc: svc,
	}
}

// CreateSchedule swaggo annotation.
//
//	@Summary		Create a new schedule
//	@Description	Create a new transaction schedule
//	@Tags			transaction
//	@Accept			json
//	@Produce		json
//	@Param			CreateScheduleRequest	body		scheduler.CreateScheduleRequest	true	"Create a new schedule request"
//	@Success		200						{object}	response.Response
//	@Failure		400						{object}	response.Response
//	@Failure		404						{object}	response.Response
//	@Failure		500						{object}	response.Response
//	@Router			/scheduler [post]
func (h *SchedulerHandler) CreateSchedule(ctx echo.Context) error {
	var req scheduler.CreateScheduleRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	if err := h.va.Validate(req); err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	err := h.svc.Create(ctx.Request().Context(), &req)
	if err != nil {
		return ctx.JSON(response.Error(err))
	}
	return ctx.JSON(response.SuccessWithoutData())
}
