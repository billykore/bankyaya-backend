package intrabank

import (
	"context"
	"strconv"

	"go.bankyaya.org/app/backend/internal/pkg/codes"
	"go.bankyaya.org/app/backend/internal/pkg/constant"
	"go.bankyaya.org/app/backend/internal/pkg/ctxt"
	"go.bankyaya.org/app/backend/internal/pkg/logger"
	"go.bankyaya.org/app/backend/internal/pkg/status"
)

const (
	domainName             = "transfer"
	transferType           = "internal_transfer"
	transferSuccessSubject = "Transfer Berhasil"
	transferFee            = 0
	stringTransferFee      = "0"
)

// Service handles the intra-bank transfer process.
type Service struct {
	log         *logger.Logger
	corebanking CoreBanking
	repo        Repository
	seqGen      SequenceGenerator
	mailer      ReceiptMailer
	notifier    Notifier
}

func NewService(
	log *logger.Logger,
	repo Repository,
	corebanking CoreBanking,
	seqGen SequenceGenerator,
	mailer ReceiptMailer,
	notifier Notifier,
) *Service {
	return &Service{
		log:         log,
		corebanking: corebanking,
		repo:        repo,
		seqGen:      seqGen,
		mailer:      mailer,
		notifier:    notifier,
	}
}

func (s *Service) Inquiry(ctx context.Context, seq *Sequence) (*Sequence, error) {
	coreStatus, err := s.corebanking.GetCoreStatus(ctx)
	if err != nil {
		s.log.DomainUsecase(domainName, "Inquiry").Errorf("CheckEOD: %v", err)
		return nil, status.Error(codes.Internal, ErrGeneral)
	}
	if coreStatus.IsEODRunning() {
		s.log.DomainUsecase(domainName, "Inquiry").Errorf("CheckEOD: %v", ErrEODInProgress)
		return nil, status.Error(codes.Internal, ErrEODInProgress)
	}

	intrabankLimit, err := s.repo.GetTransactionLimit(ctx)
	if err != nil {
		s.log.DomainUsecase(domainName, "Inquiry").Errorf("GetTransactionLimit: %v", err)
		return nil, status.Error(codes.Internal, ErrGeneral)
	}
	if !intrabankLimit.CanTransfer(seq.Amount) {
		s.log.DomainUsecase(domainName, "Inquiry").Error(ErrInvalidAmount)
		return nil, status.Error(codes.BadRequest, ErrInvalidAmount)
	}

	srcAccount, err := s.corebanking.GetAccountDetails(ctx, seq.SourceAccount)
	if err != nil {
		s.log.DomainUsecase(domainName, "Inquiry").Errorf("GetAccountDetails: %v", err)
		return nil, status.Error(codes.Internal, ErrGeneral)
	}
	if !srcAccount.IsAccountActive() {
		s.log.DomainUsecase(domainName, "Inquiry").Errorf("source account (%v) not active", seq.SourceAccount)
		return nil, status.Error(codes.BadRequest, ErrSourceAccountInactive)
	}
	seq.SourceName = srcAccount.Name

	destAccount, err := s.corebanking.GetAccountDetails(ctx, seq.DestinationAccount)
	if err != nil {
		s.log.DomainUsecase(domainName, "Inquiry").Errorf("GetAccountDetails: %v", err)
		return nil, status.Error(codes.Internal, ErrGeneral)
	}
	if !destAccount.IsAccountActive() {
		s.log.DomainUsecase(domainName, "Inquiry").Errorf("GetAccountDetails: %v", err)
		return nil, status.Error(codes.BadRequest, ErrDestinationAccountInactive)
	}
	seq.DestinationName = destAccount.Name

	sequenceNo, err := s.seqGen.Generate()
	if err != nil {
		s.log.DomainUsecase(domainName, "Inquiry").Errorf("Generate failed: %v", err)
		return nil, status.Error(codes.Internal, ErrGeneral)
	}
	seq.SequenceNumber = sequenceNo

	err = s.repo.InsertSequence(ctx, seq)
	if err != nil {
		s.log.DomainUsecase(domainName, "Inquiry").Errorf("InsertSequence: %v", err)
		return nil, status.Error(codes.Internal, ErrGeneral)
	}

	return seq, nil
}

func (s *Service) DoPayment(ctx context.Context, sequenceNumber string) (*Transaction, error) {
	coreStatus, err := s.corebanking.GetCoreStatus(ctx)
	if err != nil {
		s.log.DomainUsecase(domainName, "DoPayment").Errorf("CheckEOD: %v", err)
		return nil, status.Error(codes.Internal, ErrGeneral)
	}
	if coreStatus.IsEODRunning() {
		s.log.DomainUsecase(domainName, "DoPayment").Errorf("CheckEOD: %v", ErrEODInProgress)
		return nil, status.Error(codes.Internal, ErrEODInProgress)
	}

	sequence, err := s.repo.GetSequence(ctx, sequenceNumber)
	if err != nil {
		s.log.DomainUsecase(domainName, "DoPayment").Errorf("GetSequence: %v", err)
		return nil, status.Error(codes.Internal, ErrGeneral)
	}
	if !sequence.Valid(sequenceNumber) {
		s.log.DomainUsecase(domainName, "DoPayment").Errorf("GetSequence: %v", ErrInvalidSequenceNumber)
		return nil, status.Error(codes.BadRequest, ErrInvalidSequenceNumber)
	}

	intrabankLimit, err := s.repo.GetTransactionLimit(ctx)
	if err != nil {
		s.log.DomainUsecase(domainName, "DoPayment").Errorf("GetTransactionLimit: %v", err)
		return nil, status.Error(codes.Internal, ErrGeneral)
	}
	if !intrabankLimit.CanTransfer(sequence.Amount) {
		s.log.DomainUsecase(domainName, "DoPayment").Error(ErrInvalidAmount)
		return nil, status.Error(codes.BadRequest, ErrInvalidAmount)
	}

	result, err := s.corebanking.PerformOverbooking(ctx, &OverbookingInput{
		SourceAccount:      sequence.SourceAccount,
		DestinationAccount: sequence.DestinationAccount,
		Amount:             sequence.Amount,
		Fee:                transferFee,
		Remark:             sequence.Remark(),
	})
	if err != nil {
		s.log.DomainUsecase(domainName, "DoPayment").Errorf("PerformOverbooking: %v", err)
		return nil, status.Error(codes.Internal, ErrGeneral)
	}

	user, ok := ctxt.UserFromContext(ctx)
	if !ok {
		s.log.DomainUsecase(domainName, "DoPayment").Errorf("GetUserFromContext: %v", ctxt.ErrUserFromContext)
		return nil, status.Error(codes.Unauthenticated, ErrUserUnauthenticated)
	}

	transaction := &Transaction{
		SequenceNumber:       sequence.SequenceNumber,
		SequenceJournal:      result.JournalSequence,
		UserID:               strconv.Itoa(user.Id),
		Destination:          sequence.DestinationAccount,
		Amount:               sequence.Amount,
		TransactionType:      transferType,
		TransactionReference: result.TransactionReference,
		Remarks:              sequence.Remark(),
		Fee:                  stringTransferFee,
		DestinationName:      sequence.DestinationName,
	}

	err = s.repo.InsertTransaction(ctx, transaction)
	if err != nil {
		s.log.DomainUsecase(domainName, "DoPayment").Errorf("InsertTransaction: %v", err)
		return nil, status.Error(codes.Internal, ErrGeneral)
	}

	err = s.mailer.SendReceipt(ctx, &EmailData{
		Subject:            transferSuccessSubject,
		Recipient:          user.Email,
		Amount:             sequence.Amount,
		Fee:                transferFee,
		SourceName:         user.FullName,
		SourceAccount:      sequence.SourceAccount,
		DestinationName:    sequence.DestinationName,
		DestinationAccount: sequence.DestinationAccount,
		DestinationBank:    constant.CompanyName,
		TransactionRef:     result.TransactionReference,
		Note:               sequence.Remark(),
	})
	if err != nil {
		s.log.DomainUsecase(domainName, "DoPayment").Errorf("SendReceipt: %v", err)
		return nil, status.Error(codes.Internal, ErrSendEmailFailed)
	}

	err = s.notifier.Notify(ctx, &Notification{
		Subject:     transferSuccessSubject,
		Amount:      transaction.Amount,
		Destination: transaction.Destination,
		Status:      SuccessStatus,
	})
	if err != nil {
		s.log.DomainUsecase(domainName, "DoPayment").Errorf("Notify: %v", err)
		return nil, status.Error(codes.Internal, ErrNotifyFailed)
	}

	return transaction, nil
}
