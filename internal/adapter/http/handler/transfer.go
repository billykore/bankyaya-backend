package handler

import (
	"github.com/labstack/echo/v4"
	"go.bankyaya.org/app/backend/internal/adapter/http/response"
	"go.bankyaya.org/app/backend/internal/application/transfer"
)

type TransferHandler struct {
	uc *transfer.Usecase
}

func NewTransferHandler(uc *transfer.Usecase) *TransferHandler {
	return &TransferHandler{
		uc: uc,
	}
}

// Inquiry swaggo annotation.
//
//	@Summary			Transfer inquiry
//	@Description		Create new inquiry transfer
//	@Tags				transfer
//	@Accept				json
//	@Produce			json
//	@Param				InquiryRequest	body		transfer.InquiryRequest	true	"Inquiry request"
//	@successResponse	200															{object}	response.Response
//	@Failure			400				{object}	response.Response
//	@Failure			404				{object}	response.Response
//	@Failure			500				{object}	response.Response
//	@Router				/transfer/inquiry [post]
func (h *TransferHandler) Inquiry(ctx echo.Context) error {
	req := new(transfer.InquiryRequest)
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	resp, err := h.uc.Inquiry(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(response.Error(err))
	}
	return ctx.JSON(response.Success(resp))
}

// Payment swaggo annotation.
//
//	@Summary			Transfer payment
//	@Description		Performs transfer payment
//	@Tags				transfer
//	@Accept				json
//	@Produce			json
//	@Param				PaymentRequest	body		transfer.PaymentRequest	true	"Payment request"
//	@successResponse	200															{object}	response.Response
//	@Failure			400				{object}	response.Response
//	@Failure			404				{object}	response.Response
//	@Failure			500				{object}	response.Response
//	@Router				/transfer/payment [post]
func (h *TransferHandler) Payment(ctx echo.Context) error {
	req := new(transfer.PaymentRequest)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	resp, err := h.uc.DoPayment(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(response.Error(err))
	}
	return ctx.JSON(response.Success(resp))
}
