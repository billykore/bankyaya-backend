package qris

import (
	"context"
	"errors"

	"go.bankyaya.org/app/backend/internal/core/qris"
	"go.bankyaya.org/app/backend/pkg/codes"
	"go.bankyaya.org/app/backend/pkg/ctxt"
	"go.bankyaya.org/app/backend/pkg/logger"
	"go.bankyaya.org/app/backend/pkg/status"
	"go.bankyaya.org/app/backend/pkg/validation"
)

type Usecase struct {
	va  *validation.Validator
	log *logger.Logger
	svc *qris.Service
}

func NewUsecase(va *validation.Validator, log *logger.Logger, svc *qris.Service) *Usecase {
	return &Usecase{
		va:  va,
		log: log,
		svc: svc,
	}
}

func (uc *Usecase) Inquiry(ctx context.Context, req *InquiryRequest) (*InquiryResponse, error) {
	if err := uc.va.Validate(req); err != nil {
		uc.log.Usecase("QRIS", "Inquiry").Errorf("Validate error: %v", err)
		return nil, status.Error(codes.BadRequest, "Bad Request")
	}

	data, err := uc.svc.Inquiry(ctx, req.SourceAccount, req.QRCode)
	if err != nil {
		uc.log.Usecase("QRIS", "Inquiry").Errorf("Inquiry error: %v", err)
		if errors.Is(err, qris.ErrEODInProgress) {
			return nil, status.Error(codes.Internal, "EOD process is running")
		}
		return nil, status.Error(codes.Internal, "Internal Error")
	}

	return &InquiryResponse{
		Status:                       data.Status,
		RRN:                          data.RRN,
		CustomerName:                 data.CustomerName,
		CustomerDetail:               data.CustomerDetail,
		FinancialOrganisation:        data.FinancialOrganisation,
		FinancialOrganisationDetails: data.FinancialOrganisationDetails,
		MerchantId:                   data.MerchantID,
		MerchantCriteria:             data.MerchantCriteria,
		NMId:                         data.QRCode,
		Amount:                       data.Amount,
		TipIndicator:                 data.TipIndicator,
		TipValueOfFixed:              data.TipValueOfFixed,
		TipValueOfPercentage:         data.TipValueOfPercentage,
		Fee:                          data.Fee,
	}, nil
}

func (uc *Usecase) Payment(ctx context.Context, req *PaymentRequest) (*PaymentResponse, error) {
	if err := uc.va.Validate(req); err != nil {
		uc.log.Usecase("QRIS", "Inquiry").Errorf("Validate error: %v", err)
		return nil, status.Error(codes.BadRequest, "Bad Request")
	}

	result, err := uc.svc.Payment(ctx, &qris.QRISData{})
	if err != nil {
		uc.log.Usecase("QRIS", "Payment").Errorf("Payment error: %v", err)
		if errors.Is(err, ctxt.ErrUserFromContext) {
			return nil, status.Error(codes.Unauthenticated, "User not authenticated")
		}
		return nil, status.Error(codes.Internal, "Internal Error")
	}

	return &PaymentResponse{
		Amount:           result.Amount,
		Tip:              result.Tip,
		Total:            result.TotalAmount(),
		Message:          result.Message,
		RRN:              result.RRN,
		InvoiceNumber:    result.InvoiceNumber,
		TransactionLogId: result.TransactionReference,
	}, nil
}
