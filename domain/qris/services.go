package qris

import (
	"context"
	"errors"

	"go.bankyaya.org/app/backend/pkg/codes"
	"go.bankyaya.org/app/backend/pkg/logger"
	"go.bankyaya.org/app/backend/pkg/status"
)

const qrisFee = 0

// ErrEODInProgress indicates that the End of Day (EOD) process is currently in progress.
var ErrEODInProgress = errors.New("EOD process is running")

type Repository interface {
}

// CoreBanking defines methods for core banking operations.
type CoreBanking interface {
	// CheckEOD verifies the current End-of-Day (EOD) process status in the core banking system.
	CheckEOD(ctx context.Context) (*EODStatus, error)

	// GetAccountDetails retrieves account information for the given account number.
	GetAccountDetails(ctx context.Context, accountNumber string) (*AccountDetails, error)
}

// QRIS is an interface for QR Code Indonesian Standard operations.
type QRIS interface {
	// GetDetails retrieves QRIS data based on account number and QR code parameters.
	GetDetails(ctx context.Context, accountNumber, qrCode string) (*QRISData, error)
}

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

	details, err := s.qris.GetDetails(ctx, srcAccount.AccountNumber, req.QrCode)
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