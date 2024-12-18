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

// ErrEODInProgress indicates that the End of Day (EOD) process is currently in progress.
var ErrEODInProgress = errors.New("EOD process is running")

// Repository defines methods for managing transfer sequence persistence.
type Repository interface {
	// InsertSequence inserts a transfer sequence into the persistence repository.
	// Requires a context and a Sequence object to execute.
	// Returns an error if the operation fails.
	InsertSequence(ctx context.Context, seq *Sequence) error

	// GetSequence retrieves a transfer sequence based on the sequence number.
	// Requires a context and the sequence number as inputs.
	// Returns a Sequence object and an error if retrieval fails.
	GetSequence(ctx context.Context, sequenceNumber string) (*Sequence, error)

	// GetUserById retrieves a user by their unique ID.
	// Requires a context and an integer ID as input parameters.
	// Returns a User object and an error if retrieval fails.
	GetUserById(ctx context.Context, id int) (*User, error)
}

// CoreBankingService defines methods for core banking operations.
type CoreBankingService interface {
	// CheckEOD verifies the current End-of-Day (EOD) process status in the core banking system.
	CheckEOD(ctx context.Context) (*EODStatus, error)

	// GetAccountDetails retrieves account information for the given account number.
	GetAccountDetails(ctx context.Context, accountNumber string) (*AccountDetails, error)

	// PerformOverbooking executes a transfer between two accounts with the specified amount and remark.
	// It returns an OverbookingResponse and an error if the operation fails.
	PerformOverbooking(ctx context.Context, req OverbookingRequest) (*OverbookingResponse, error)
}

// Email is an interface for sending transfer receipt emails.
type Email interface {
	// SendTransferReceipt sends a transfer receipt email using provided context and email data.
	SendTransferReceipt(context.Context, *EmailData) error
}

type Service struct {
	log         *logger.Logger
	repo        Repository
	corebanking CoreBankingService
	email       Email
}

func NewService(log *logger.Logger, repo Repository, corebanking CoreBankingService, email Email) *Service {
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
	if result.Code != "00" {
		s.log.Usecase("Payment").Errorf("Overbook code (%v): %v", result.Code, result.Description)
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
