package handler

import (
	"github.com/labstack/echo/v4"
	"go.bankyaya.org/app/backend/internal/adapter/http/dto"
	"go.bankyaya.org/app/backend/internal/adapter/http/response"
	"go.bankyaya.org/app/backend/internal/domain/intrabank"
)

type Intrabank struct {
	svc *intrabank.Service
}

func NewIntrabankHandler(svc *intrabank.Service) *Intrabank {
	return &Intrabank{
		svc: svc,
	}
}

// Inquiry swaggo annotation.
//
//	@Summary			Intrabank transfer inquiry
//	@Description		Create new inquiry intrabank transfer
//	@Tags				transfer
//	@Accept				json
//	@Produce			json
//	@Param				InquiryRequest	body		dto.IntrabankInquiryRequest	true	"Inquiry request"
//	@successResponse	200																										{object}	response.Response
//	@Failure			400				{object}	response.Response
//	@Failure			404				{object}	response.Response
//	@Failure			500				{object}	response.Response
//	@Router				/transfer/intrabank/inquiry [post]
func (h *Intrabank) Inquiry(ctx echo.Context) error {
	req := new(dto.IntrabankInquiryRequest)
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	sequence, err := h.svc.Inquiry(ctx.Request().Context(), req.ToSequence())
	if err != nil {
		return ctx.JSON(response.Error(err))
	}
	resp := dto.NewIntrabankInquiryResponse(sequence)
	return ctx.JSON(response.Success(resp))
}

// Payment swaggo annotation.
//
//	@Summary			Intrabank payment
//	@Description		Performs transfer payment
//	@Tags				transfer
//	@Accept				json
//	@Produce			json
//	@Param				PaymentRequest	body		dto.IntrabankPaymentRequest	true	"Payment request"
//	@successResponse	200																										{object}	response.Response
//	@Failure			400				{object}	response.Response
//	@Failure			404				{object}	response.Response
//	@Failure			500				{object}	response.Response
//	@Router				/transfer/intrabank/payment [post]
func (h *Intrabank) Payment(ctx echo.Context) error {
	req := new(dto.IntrabankPaymentRequest)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	transaction, err := h.svc.DoPayment(ctx.Request().Context(), req.Sequence)
	if err != nil {
		return ctx.JSON(response.Error(err))
	}
	resp := dto.NewIntrabankPaymentResponse(transaction)
	return ctx.JSON(response.Success(resp))
}
