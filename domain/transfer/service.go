package transfer

import (
	"context"
	"errors"
	"time"

	"go.bankyaya.org/app/backend/pkg/codes"
	"go.bankyaya.org/app/backend/pkg/constant"
	"go.bankyaya.org/app/backend/pkg/ctxt"
	"go.bankyaya.org/app/backend/pkg/logger"
	"go.bankyaya.org/app/backend/pkg/status"
	"go.bankyaya.org/app/backend/pkg/uuid"
)

const (
	transferType           = "internal"
	transferSuccessSubject = "Transfer success"
	transferFee            = 0
)

const (
	messageEODIsRunning       = "EOD process is running"
	messageInquiryFailed      = "Inquiry failed"
	messagePaymentFailed      = "Payment failed"
	messageAccountIsNotActive = "Account is not active"
)

// Service handles intra-bank transfer process.
type Service struct {
	log         *logger.Logger
	repo        Repository
	corebanking CoreBanking
	email       Email
}

func NewService(log *logger.Logger, repo Repository, corebanking CoreBanking, email Email) *Service {
	return &Service{
		log:         log,
		repo:        repo,
		corebanking: corebanking,
		email:       email,
	}
}

func (s *Service) Inquiry(ctx context.Context, req InquiryRequest) (*InquiryResponse, error) {
	eod, err := s.corebanking.CheckEOD(ctx)
	if err != nil {
		s.log.Usecase("Inquiry").Errorf("EOD: %v", err)
		return nil, status.Errorf(codes.Internal, "%s: %v", messageInquiryFailed, err)
	}
	if eod.IsRunning() {
		s.log.Usecase("Inquiry").Error(ErrEODInProgress)
		return nil, status.Error(codes.Internal, messageEODIsRunning)
	}

	srcAccount, err := s.corebanking.GetAccountDetails(ctx, req.SourceAccount)
	if err != nil {
		s.log.Usecase("Inquiry").Errorf("Inquiry failed: %v", err)
		return nil, status.Errorf(codes.Internal, "%s: %v", messageInquiryFailed, err)
	}
	if !srcAccount.IsAccountActive() {
		s.log.Usecase("Inquiry").Errorf("account status is not active (%s)", srcAccount.Status)
		return nil, status.Errorf(codes.Internal, "%s: %v", messageInquiryFailed, messageAccountIsNotActive)
	}

	destAccount, err := s.corebanking.GetAccountDetails(ctx, req.DestinationAccount)
	if err != nil {
		s.log.Usecase("Inquiry").Errorf("Inquiry: %v", err)
		return nil, status.Errorf(codes.Internal, "%s: %v", messageInquiryFailed, err)
	}
	if !destAccount.IsAccountActive() {
		s.log.Usecase("Inquiry").Errorf("account status is not active (%s)", destAccount.Status)
		return nil, status.Errorf(codes.Internal, "%s: %v", messageInquiryFailed, messageAccountIsNotActive)
	}

	sequence, err := uuid.New()
	if err != nil {
		s.log.Usecase("Inquiry").Errorf("UUID: %v", err)
		return nil, status.Error(codes.Internal, messageInquiryFailed)
	}

	err = s.repo.InsertSequence(ctx, &Sequence{
		SeqNo:           sequence,
		Amount:          req.StringAmount(),
		AccNoSrc:        req.SourceAccount,
		AccNoDest:       req.DestinationAccount,
		AccNameSrc:      srcAccount.Name,
		AccNameDest:     destAccount.Name,
		TransactionType: transferType,
		CifDest:         destAccount.CIF,
		CreateDate:      time.Now(),
	})
	if err != nil {
		s.log.Usecase("Inquiry").Errorf("InsertTransferSequence: %v", err)
		return nil, status.Error(codes.Internal, messageInquiryFailed)
	}

	return &InquiryResponse{
		SequenceNumber:     sequence,
		SourceAccount:      req.SourceAccount,
		DestinationAccount: req.DestinationAccount,
		Status:             destAccount.Status,
	}, nil
}

func (s *Service) Payment(ctx context.Context, req PaymentRequest) (*PaymentResponse, error) {
	eod, err := s.corebanking.CheckEOD(ctx)
	if err != nil {
		s.log.Usecase("Payment").Errorf("EOD: %v", err)
		return nil, status.Errorf(codes.Internal, "%s: %v", messagePaymentFailed, err)
	}
	if eod.IsRunning() {
		s.log.Usecase("Payment").Error(ErrEODInProgress)
		return nil, status.Error(codes.Internal, messageEODIsRunning)
	}

	sequence, err := s.repo.GetSequence(ctx, req.Sequence)
	if err != nil {
		s.log.Usecase("Payment").Errorf("GetTransferSequence: %v", err)
		return nil, status.Error(codes.Internal, messagePaymentFailed)
	}
	if sequence.SeqNo == "" {
		s.log.Usecase("Payment").Error(errors.New("empty sequence transfer"))
		return nil, status.Error(codes.Internal, messagePaymentFailed)
	}

	result, err := s.corebanking.PerformOverbooking(ctx, OverbookingRequest{
		SourceAccount:      req.SourceAccount,
		DestinationAccount: req.DestinationAccount,
		Amount:             req.Amount,
		Fee:                transferFee,
		Remark:             req.Remark(),
	})
	if err != nil {
		s.log.Usecase("Payment").Errorf("Overbook: %v", err)
		return nil, status.Error(codes.Internal, messagePaymentFailed)
	}

	user, ok := ctxt.UserFromContext(ctx)
	if !ok {
		s.log.Usecase("Payment").Error(ctxt.ErrUserFromContext)
		return nil, status.Error(codes.Unauthenticated, "user not found")
	}

	userData, err := s.repo.GetUserById(ctx, user.Id)
	if err != nil {
		s.log.Usecase("Payment").Errorf("GetUserById: %v", err)
		return nil, status.Error(codes.Internal, messagePaymentFailed)
	}

	go s.sendTransferReceipt(ctx, &EmailData{
		Subject:            transferSuccessSubject,
		Recipient:          userData.Email,
		Amount:             req.Amount,
		Fee:                transferFee,
		SourceName:         userData.FullName,
		SourceAccount:      req.SourceAccount,
		DestinationName:    sequence.AccNameDest,
		DestinationAccount: req.DestinationAccount,
		DestinationBank:    constant.CompanyName,
		TransactionRef:     result.TransactionReference,
		Note:               req.Remark(),
	})

	return &PaymentResponse{
		ABMsg:                result.ABMsg,
		JournalSequence:      result.JournalSequence,
		DestinationAccount:   req.DestinationAccount,
		SourceAccount:        req.SourceAccount,
		Amount:               req.Amount,
		Notes:                req.Notes,
		BankName:             constant.CompanyName,
		TransactionReference: result.TransactionReference,
		Remark:               req.Remark(),
	}, nil
}

// sendTransferReceipt is a dedicated function to handle the asynchronous sending of transfer receipt emails.
func (s *Service) sendTransferReceipt(ctx context.Context, emailData *EmailData) {
	if err := s.email.SendTransferReceipt(ctx, emailData); err != nil {
		s.log.Usecase("sendTransferReceipt").
			Errorf("Failed to send transfer receipt: %v", err)
	}
}
