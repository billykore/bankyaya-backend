package handler

import (
	"github.com/labstack/echo/v4"
	"go.bankyaya.org/app/backend/internal/adapter/http/response"
	"go.bankyaya.org/app/backend/internal/application/user"
)

type UserHandler struct {
	uc *user.Usecase
}

func NewUserHandler(uc *user.Usecase) *UserHandler {
	return &UserHandler{
		uc: uc,
	}
}

// Login swaggo annotation.
//
//	@Summary		User login
//	@Description	User login to get access token
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			LoginRequest	body		user.LoginRequest	true	"Login request"
//	@Success		200				{object}	response.Response
//	@Failure		400				{object}	response.Response
//	@Failure		404				{object}	response.Response
//	@Failure		500				{object}	response.Response
//	@Router			/user/login [post]
func (h *UserHandler) Login(ctx echo.Context) error {
	req := new(user.LoginRequest)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	resp, err := h.uc.Login(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(response.Error(err))
	}
	return ctx.JSON(response.Success(resp))
}
