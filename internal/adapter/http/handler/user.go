package handler

import (
	"github.com/labstack/echo/v4"
	"go.bankyaya.org/app/backend/internal/adapter/http/dto"
	"go.bankyaya.org/app/backend/internal/adapter/http/response"
	"go.bankyaya.org/app/backend/internal/domain/user"
)

type UserHandler struct {
	svc *user.Service
}

func NewUserHandler(svc *user.Service) *UserHandler {
	return &UserHandler{
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
//	@Param			LoginRequest	body		dto.LoginRequest	true	"Login request"
//	@Success		200				{object}	response.Response
//	@Failure		400				{object}	response.Response
//	@Failure		404				{object}	response.Response
//	@Failure		500				{object}	response.Response
//	@Router			/user/login [post]
func (h *UserHandler) Login(ctx echo.Context) error {
	req := new(dto.LoginRequest)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	token, err := h.svc.Login(ctx.Request().Context(), req.ToUser())
	if err != nil {
		return ctx.JSON(response.Error(err))
	}
	resp := dto.NewLoginResponse(token)
	return ctx.JSON(response.Success(resp))
}
