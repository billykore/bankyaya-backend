package handler

import (
	"github.com/labstack/echo/v4"
	"go.bankyaya.org/app/backend/pkg/framework/http/dto"
	"go.bankyaya.org/app/backend/pkg/framework/http/response"
	"go.bankyaya.org/app/backend/pkg/service"
)

type Transfer struct {
	svc *service.Transfer
}

func NewTransfer(svc *service.Transfer) *Transfer {
	return &Transfer{
		svc: svc,
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
func (t *Transfer) Inquiry(ctx echo.Context) error {
	req := new(dto.TransferInquiryRequest)
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	sequence, err := t.svc.Inquiry(ctx.Request().Context(), req.ToSequence())
	if err != nil {
		return ctx.JSON(response.Error(err))
	}
	resp := dto.NewTransferInquiryResponse(sequence)
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
func (t *Transfer) Payment(ctx echo.Context) error {
	req := new(dto.TransferPaymentRequest)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	transaction, err := t.svc.DoPayment(ctx.Request().Context(), req.Sequence)
	if err != nil {
		return ctx.JSON(response.Error(err))
	}
	resp := dto.NewTransferPaymentResponse(transaction)
	return ctx.JSON(response.Success(resp))
}
