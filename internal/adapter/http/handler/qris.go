package handler

import (
	"github.com/labstack/echo/v4"
	"go.bankyaya.org/app/backend/internal/adapter/http/response"
	"go.bankyaya.org/app/backend/internal/application/qris"
)

type QRISHandler struct {
	uc *qris.Usecase
}

func NewQRISHandler(uc *qris.Usecase) *QRISHandler {
	return &QRISHandler{
		uc: uc,
	}
}

// Inquiry swaggo annotation.
//
//	@Summary			QRIS inquiry
//	@Description		Create new QRIS inquiry
//	@Tags				payment
//	@Accept				json
//	@Produce			json
//	@Param				InquiryRequest	body		qris.InquiryRequest	true	"QRIS inquiry request"
//	@successResponse	200														{object}	response.Response
//	@Failure			400				{object}	response.Response
//	@Failure			404				{object}	response.Response
//	@Failure			500				{object}	response.Response
//	@Router				/qris/inquiry [post]
func (h *QRISHandler) Inquiry(ctx echo.Context) error {
	req := new(qris.InquiryRequest)
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
//	@Summary			QRIS payment
//	@Description		Do QRIS payment
//	@Tags				payment
//	@Accept				json
//	@Produce			json
//	@Param				PaymentRequest	body		qris.PaymentRequest	true	"QRIS payment request"
//	@successResponse	200														{object}	response.Response
//	@Failure			400				{object}	response.Response
//	@Failure			404				{object}	response.Response
//	@Failure			500				{object}	response.Response
//	@Router				/qris/pay [post]
func (h *QRISHandler) Payment(ctx echo.Context) error {
	req := new(qris.PaymentRequest)
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	resp, err := h.uc.Payment(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(response.Error(err))
	}
	return ctx.JSON(response.Success(resp))
}
