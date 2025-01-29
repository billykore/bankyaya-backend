package handler

import (
	"github.com/labstack/echo/v4"
	"go.bankyaya.org/app/backend/domain/transfer"
	"go.bankyaya.org/app/backend/pkg/response"
	"go.bankyaya.org/app/backend/pkg/validation"
)

type TransferHandler struct {
	va  *validation.Validator
	svc *transfer.Service
}

func NewTransferHandler(va *validation.Validator, svc *transfer.Service) *TransferHandler {
	return &TransferHandler{
		va:  va,
		svc: svc,
	}
}

// Inquiry swaggo annotation.
//
//	@Summary		Transfer inquiry
//	@Description	Create new inquiry transfer
//	@Tags			transfer
//	@Accept			json
//	@Produce		json
//	@Param			InquiryRequest	body		transfer.InquiryRequest	true	"Inquiry request"
//	@Success		200				{object}	response.Response
//	@Failure		400				{object}	response.Response
//	@Failure		404				{object}	response.Response
//	@Failure		500				{object}	response.Response
//	@Router			/transfer/inquiry [post]
func (h *TransferHandler) Inquiry(ctx echo.Context) error {
	var req transfer.InquiryRequest
	if err := ctx.Bind(&req); err != nil {
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
//	@Summary		Transfer payment
//	@Description	Performs transfer payment
//	@Tags			transfer
//	@Accept			json
//	@Produce		json
//	@Param			PaymentRequest	body		transfer.PaymentRequest	true	"Payment request"
//	@Success		200				{object}	response.Response
//	@Failure		400				{object}	response.Response
//	@Failure		404				{object}	response.Response
//	@Failure		500				{object}	response.Response
//	@Router			/transfer/payment [post]
func (h *TransferHandler) Payment(ctx echo.Context) error {
	var req transfer.PaymentRequest
	if err := ctx.Bind(&req); err != nil {
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
