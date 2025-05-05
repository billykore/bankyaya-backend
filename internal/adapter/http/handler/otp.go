package handler

import (
	"github.com/labstack/echo/v4"
	"go.bankyaya.org/app/backend/internal/adapter/http/dto"
	"go.bankyaya.org/app/backend/internal/adapter/http/response"
	"go.bankyaya.org/app/backend/internal/domain/otp"
	"go.bankyaya.org/app/backend/internal/pkg/validation"
)

type OTPHandler struct {
	va  *validation.Validator
	svc *otp.Service
}

func NewOTPHandler(va *validation.Validator, svc *otp.Service) *OTPHandler {
	return &OTPHandler{
		va:  va,
		svc: svc,
	}
}

// SendOTP swaggo annotation.
//
//	@Summary		Send new OTP
//	@Description	Send new OTP to user
//	@Tags			otp
//	@Accept			json
//	@Produce		json
//	@Param			OTPRequest	body		dto.OTPRequest	true	"OTP request"
//	@Success		200				{object}	response.Response
//	@Failure		400				{object}	response.Response
//	@Failure		404				{object}	response.Response
//	@Failure		500				{object}	response.Response
//	@Router			/otp/send [post]
func (h *OTPHandler) SendOTP(ctx echo.Context) error {
	req := new(dto.OTPRequest)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	if err := h.va.Validate(req); err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	otpRes, err := h.svc.Send(ctx.Request().Context(), otp.NewPurpose(req.Purpose), otp.NewChannel(req.Channel))
	if err != nil {
		return ctx.JSON(response.Error(err))
	}
	resp := dto.NewOTPResponse(otpRes)
	return ctx.JSON(response.Success(resp))
}

// VerifyOTP swaggo annotation.
//
//	@Summary		Verify OTP
//	@Description	Verify user OTP
//	@Tags			otp
//	@Accept			json
//	@Produce		json
//	@Param			 VerifyOTPRequest	body		dto.VerifyOTPRequest	true	"Verify OTP request"
//	@Success		200				{object}	response.Response
//	@Failure		400				{object}	response.Response
//	@Failure		404				{object}	response.Response
//	@Failure		500				{object}	response.Response
//	@Router			/otp/verify [post]
func (h *OTPHandler) VerifyOTP(ctx echo.Context) error {
	req := new(dto.VerifyOTPRequest)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	if err := h.va.Validate(req); err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	err := h.svc.Verify(ctx.Request().Context(), req.ToOTP())
	if err != nil {
		return ctx.JSON(response.Error(err))
	}
	return ctx.JSON(response.Success(nil))
}
