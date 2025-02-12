package qris

import (
	"context"

	"go.bankyaya.org/app/backend/pkg/codes"
	"go.bankyaya.org/app/backend/pkg/logger"
	"go.bankyaya.org/app/backend/pkg/status"
)

const (
	messageEODIsRunning       = "EOD process is running"
	messageInquiryFailed      = "QRIS inquiry failed"
	messagePaymentFailed      = "QRIS payment failed"
	messageAccountIsNotActive = "Account is not active"
)

const qrisFee = 0

// Service handles QRIS payment process.
type Service struct {
	log         *logger.Logger
	corebanking CoreBanking
	qris        QRIS
}

func NewService(log *logger.Logger, corebanking CoreBanking, qris QRIS) *Service {
	return &Service{
		log:         log,
		corebanking: corebanking,
		qris:        qris,
	}
}

func (s *Service) Inquiry(ctx context.Context, req *InquiryRequest) (*InquiryResponse, error) {
	eod, err := s.corebanking.CheckEOD(ctx)
	if err != nil {
		s.log.DomainUsecase("qris", "Inquiry").Errorf("EOD: %v", err)
		return nil, status.Errorf(codes.Internal, "%s: %v", messageInquiryFailed, err)
	}
	if eod.IsRunning() {
		s.log.DomainUsecase("qris", "Inquiry").Error(ErrEODInProgress)
		return nil, status.Error(codes.Internal, messageEODIsRunning)
	}

	srcAccount, err := s.corebanking.GetAccountDetails(ctx, req.SourceAccount)
	if err != nil {
		s.log.DomainUsecase("qris", "Inquiry").Errorf("Inquiry failed: %v", err)
		return nil, status.Errorf(codes.Internal, "%s: %v", messageInquiryFailed, err)
	}
	if !srcAccount.IsAccountActive() {
		s.log.DomainUsecase("qris", "Inquiry").Errorf("account status is not active (%s)", srcAccount.Status)
		return nil, status.Errorf(codes.Internal, "%s: %v", messageInquiryFailed, messageAccountIsNotActive)
	}

	details, err := s.qris.GetDetails(ctx, srcAccount.AccountNumber, req.QRCode)
	if err != nil {
		s.log.DomainUsecase("qris", "Inquiry").Errorf("Inquiry: %v", err)
		return nil, status.Errorf(codes.Internal, "%s: %v", messageInquiryFailed, err)
	}

	return &InquiryResponse{
		Status:                       details.Status,
		RRN:                          details.RRN,
		CustomerName:                 details.CustomerName,
		CustomerDetail:               details.DetailCustomer,
		FinancialOrganisation:        details.FinancialOrganisation,
		FinancialOrganisationDetails: details.FinancialOrganisationDetails,
		MerchantId:                   details.MerchantID,
		MerchantCriteria:             details.MerchantCriteria,
		NMId:                         details.NMId,
		Amount:                       details.Amount,
		TipIndicator:                 details.TipIndicator,
		TipValueOfFixed:              details.TipValueOfFixed,
		TipValueOfPercentage:         details.TipValueOfPercentage,
		Fee:                          qrisFee,
	}, nil
}

func (s *Service) Payment(ctx context.Context, req *PaymentRequest) (*PaymentResponse, error) {
	payRes, err := s.qris.Pay(ctx, &PaymentData{
		AccountNumber:         req.SourceAccount,
		QRCode:                req.QRCode,
		RRN:                   req.RRN,
		Amount:                req.Amount,
		Tip:                   req.Tip,
		FinancialOrganisation: req.FinancialOrganisation,
		CustomerName:          req.CustomerName,
		MerchantId:            req.MerchantId,
		MerchantCriteria:      req.MerchantCriteria,
		NMId:                  req.NMId,
		AccountName:           req.CustomerName,
	})
	if err != nil {
		s.log.DomainUsecase("qris", "Payment").Errorf("Payment: %v", err)
		return nil, status.Error(codes.Internal, messagePaymentFailed)
	}
	return &PaymentResponse{
		Amount:           payRes.Amount,
		Tip:              payRes.Tip,
		Total:            payRes.TotalAmount(),
		Message:          payRes.Message,
		RRN:              payRes.RRN,
		InvoiceNumber:    payRes.InvoiceNumber,
		TransactionLogId: payRes.TransactionReference,
	}, nil
}
