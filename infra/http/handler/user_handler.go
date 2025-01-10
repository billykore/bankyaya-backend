package handler

import (
	"github.com/labstack/echo/v4"
	"go.bankyaya.org/app/backend/domain/user"
	"go.bankyaya.org/app/backend/pkg/entity"
	"go.bankyaya.org/app/backend/pkg/validation"
)

type UserHandler struct {
	va  *validation.Validator
	svc *user.Service
}

func NewUserHandler(va *validation.Validator, svc *user.Service) *UserHandler {
	return &UserHandler{
		va:  va,
		svc: svc,
	}
}

func (h *UserHandler) Login(ctx echo.Context) error {
	var req user.LoginRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(entity.ResponseBadRequest(err))
	}
	if err := h.va.Validate(req); err != nil {
		return ctx.JSON(entity.ResponseBadRequest(err))
	}
	resp, err := h.svc.Login(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(entity.ResponseError(err))
	}
	return ctx.JSON(entity.ResponseSuccess(resp))
}
