package handler

import (
	"github.com/labstack/echo/v4"
	"go.bankyaya.org/app/backend/domain/qris"
	"go.bankyaya.org/app/backend/pkg/response"
	"go.bankyaya.org/app/backend/pkg/validation"
)

type QRISHandler struct {
	va  *validation.Validator
	svc *qris.Service
}

func NewQRISHandler(va *validation.Validator, svc *qris.Service) *QRISHandler {
	return &QRISHandler{
		va:  va,
		svc: svc,
	}
}

// Inquiry swaggo annotation.
//
//	@Summary		QRIS inquiry
//	@Description	Create new QRIS inquiry
//	@Tags			payment
//	@Accept			json
//	@Produce		json
//	@Param			InquiryRequest	body		qris.InquiryRequest	true	"QRIS inquiry request"
//	@Success		200				{object}	response.Response
//	@Failure		400				{object}	response.Response
//	@Failure		404				{object}	response.Response
//	@Failure		500				{object}	response.Response
//	@Router			/qris/inquiry [post]
func (h *QRISHandler) Inquiry(ctx echo.Context) error {
	req := new(qris.InquiryRequest)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	if err := h.va.Validate(req); err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	resp, err := h.svc.Inquiry(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(response.Error(err))
	}
	return ctx.JSON(response.Success(resp))
}

// Payment swaggo annotation.
//
//	@Summary		QRIS payment
//	@Description	Do QRIS payment
//	@Tags			payment
//	@Accept			json
//	@Produce		json
//	@Param			PaymentRequest	body		qris.PaymentRequest	true	"QRIS payment request"
//	@Success		200				{object}	response.Response
//	@Failure		400				{object}	response.Response
//	@Failure		404				{object}	response.Response
//	@Failure		500				{object}	response.Response
//	@Router			/qris/pay [post]
func (h *QRISHandler) Payment(ctx echo.Context) error {
	req := new(qris.PaymentRequest)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	if err := h.va.Validate(req); err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	resp, err := h.svc.Payment(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(response.Error(err))
	}
	return ctx.JSON(response.Success(resp))
}
