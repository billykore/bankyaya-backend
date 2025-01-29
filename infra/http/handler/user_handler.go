package handler

import (
	"github.com/labstack/echo/v4"
	"go.bankyaya.org/app/backend/domain/user"
	"go.bankyaya.org/app/backend/pkg/response"
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

// Login swaggo annotation.
//
//	@Summary		User login
//	@Description	User login to get access token
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			InquiryRequest	body		transfer.InquiryRequest	true	"Inquiry request"
//	@Success		200				{object}	response.Response
//	@Failure		400				{object}	response.Response
//	@Failure		404				{object}	response.Response
//	@Failure		500				{object}	response.Response
//	@Router			/user/login [post]
func (h *UserHandler) Login(ctx echo.Context) error {
	var req user.LoginRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	if err := h.va.Validate(req); err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	resp, err := h.svc.Login(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(response.Error(err))
	}
	return ctx.JSON(response.Success(resp))
}
