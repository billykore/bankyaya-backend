package handler

import (
	"github.com/labstack/echo/v4"
	"go.bankyaya.org/app/backend/domain/qris"
	"go.bankyaya.org/app/backend/pkg/entity"
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
//	@Tags			transfer
//	@Accept			json
//	@Produce		json
//	@Param			InquiryRequest	body		qris.InquiryRequest	true	"QRIS inquiry request"
//	@Success		200				{object}	entity.Response
//	@Failure		400				{object}	entity.Response
//	@Failure		404				{object}	entity.Response
//	@Failure		500				{object}	entity.Response
//	@Router			/qris/inquiry [post]
func (h *QRISHandler) Inquiry(ctx echo.Context) error {
	req := new(qris.InquiryRequest)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(entity.ResponseBadRequest(err))
	}
	if err := h.va.Validate(req); err != nil {
		return ctx.JSON(entity.ResponseBadRequest(err))
	}
	resp, err := h.svc.Inquiry(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(entity.ResponseError(err))
	}
	return ctx.JSON(entity.ResponseSuccess(resp))
}
