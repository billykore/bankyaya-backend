package handler

import (
	"github.com/labstack/echo/v4"
	"go.bankyaya.org/app/backend/pkg/framework/http/dto"
	"go.bankyaya.org/app/backend/pkg/framework/http/response"
	"go.bankyaya.org/app/backend/pkg/service"
	"go.bankyaya.org/app/backend/pkg/util/validation"
)

type QRIS struct {
	va  *validation.Validator
	svc *service.QRIS
}

func NewQRIS(svc *service.QRIS) *QRIS {
	return &QRIS{
		svc: svc,
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
func (qris *QRIS) Inquiry(ctx echo.Context) error {
	req := new(dto.QRISInquiryRequest)
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	if err := qris.va.Validate(req); err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	data, err := qris.svc.Inquiry(ctx.Request().Context(), req.SourceAccount, req.QRCode)
	if err != nil {
		return ctx.JSON(response.Error(err))
	}
	resp := dto.NewQRISInquiryResponse(data)
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
func (qris *QRIS) Payment(ctx echo.Context) error {
	req := new(dto.QRISPaymentRequest)
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	if err := qris.va.Validate(req); err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	result, err := qris.svc.Payment(ctx.Request().Context(), req.ToQRISData())
	if err != nil {
		return ctx.JSON(response.Error(err))
	}
	resp := dto.NewQRISPaymentResponse(result)
	return ctx.JSON(response.Success(resp))
}
