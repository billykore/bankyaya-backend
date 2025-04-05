package service

import (
	"context"
	"strconv"

	"go.bankyaya.org/app/backend/internal/core/domain"
	"go.bankyaya.org/app/backend/internal/core/entity"
	"go.bankyaya.org/app/backend/internal/core/port/api"
	"go.bankyaya.org/app/backend/internal/core/port/email"
	"go.bankyaya.org/app/backend/internal/core/port/repository"
	"go.bankyaya.org/app/backend/internal/core/port/security"
	"go.bankyaya.org/app/backend/internal/pkg/codes"
	"go.bankyaya.org/app/backend/internal/pkg/constant"
	"go.bankyaya.org/app/backend/internal/pkg/ctxt"
	"go.bankyaya.org/app/backend/internal/pkg/logger"
	"go.bankyaya.org/app/backend/internal/pkg/status"
	"go.bankyaya.org/app/backend/internal/pkg/types"
)

const (
	transferService        = "Transfer"
	transferType           = "internal_transfer"
	transferSuccessSubject = "Transfer success"
	transferFee            = 0
	stringTransferFee      = "0"
)

// Transfer handles intra-bank transfer process.
type Transfer struct {
	log         *logger.Logger
	corebanking api.CoreBanking
	repo        repository.TransferRepository
	seqGen      security.SequenceGenerator
	mailer      email.TransferReceiptMailer
}

func NewTransfer(
	log *logger.Logger,
	repo repository.TransferRepository,
	corebanking api.CoreBanking,
	seqGen security.SequenceGenerator,
	mailer email.TransferReceiptMailer,
) *Transfer {
	return &Transfer{
		log:         log,
		corebanking: corebanking,
		repo:        repo,
		seqGen:      seqGen,
		mailer:      mailer,
	}
}

func (t *Transfer) Inquiry(ctx context.Context, seq *entity.Sequence) (*entity.Sequence, error) {
	eod, err := t.corebanking.CheckEOD(ctx)
	if err != nil {
		t.log.ServiceUsecase(transferService, "Inquiry").Errorf("CheckEOD: %v", err)
		return nil, status.Error(codes.Internal, domain.ErrGeneral)
	}
	if eod.IsRunning() {
		t.log.ServiceUsecase(transferService, "Inquiry").Errorf("CheckEOD: %v", domain.ErrEODInProgress)
		return nil, status.Error(codes.Internal, domain.ErrEODInProgress)
	}

	srcAccount, err := t.corebanking.GetAccountDetails(ctx, seq.AccNoSrc)
	if err != nil {
		t.log.ServiceUsecase(transferService, "Inquiry").Errorf("GetAccountDetails: %v", err)
		return nil, status.Error(codes.Internal, domain.ErrGeneral)
	}
	if !srcAccount.IsAccountActive() {
		t.log.ServiceUsecase(transferService, "Inquiry").Errorf("source account (%v) not active", seq.AccNoSrc)
		return nil, status.Error(codes.BadRequest, domain.ErrSourceAccountInactive)
	}
	seq.AccNameSrc = srcAccount.Name

	destAccount, err := t.corebanking.GetAccountDetails(ctx, seq.AccNoDest)
	if err != nil {
		t.log.ServiceUsecase(transferService, "Inquiry").Errorf("GetAccountDetails: %v", err)
		return nil, status.Error(codes.Internal, domain.ErrGeneral)
	}
	if !destAccount.IsAccountActive() {
		t.log.ServiceUsecase(transferService, "Inquiry").Errorf("GetAccountDetails: %v", err)
		return nil, status.Error(codes.BadRequest, domain.ErrDestinationAccountInactive)
	}
	seq.AccNameDest = destAccount.Name

	sequenceNo, err := t.seqGen.Generate()
	if err != nil {
		t.log.ServiceUsecase(transferService, "Inquiry").Errorf("New uuid failed: %v", err)
		return nil, status.Error(codes.Internal, domain.ErrGeneral)
	}
	seq.SeqNo = sequenceNo

	err = t.repo.InsertSequence(ctx, seq)
	if err != nil {
		t.log.ServiceUsecase(transferService, "Inquiry").Errorf("InsertSequence: %v", err)
		return nil, status.Error(codes.Internal, domain.ErrGeneral)
	}

	return seq, nil
}

func (t *Transfer) DoPayment(ctx context.Context, sequenceNumber string) (*entity.Transaction, error) {
	eod, err := t.corebanking.CheckEOD(ctx)
	if err != nil {
		t.log.ServiceUsecase(transferService, "DoPayment").Errorf("CheckEOD: %v", err)
		return nil, status.Error(codes.Internal, domain.ErrGeneral)
	}
	if eod.IsRunning() {
		t.log.ServiceUsecase(transferService, "DoPayment").Errorf("CheckEOD: %v", domain.ErrEODInProgress)
		return nil, status.Error(codes.Internal, domain.ErrEODInProgress)
	}

	sequence, err := t.repo.GetSequence(ctx, sequenceNumber)
	if err != nil {
		t.log.ServiceUsecase(transferService, "DoPayment").Errorf("GetSequence: %v", err)
		return nil, status.Error(codes.Internal, domain.ErrGeneral)
	}
	if sequence.SeqNo == "" {
		t.log.ServiceUsecase(transferService, "DoPayment").Errorf("GetSequence: %v", err)
		return nil, status.Error(codes.BadRequest, domain.ErrInvalidSequenceNumber)
	}
	amount, err := types.ParseMoney(sequence.Amount)
	if err != nil {
		t.log.ServiceUsecase(transferService, "DoPayment").Errorf("ParseMoney: %v", err)
		return nil, status.Error(codes.BadRequest, domain.ErrGeneral)
	}

	result, err := t.corebanking.PerformOverbooking(ctx, &entity.OverbookingRequest{
		SourceAccount:      sequence.AccNoSrc,
		DestinationAccount: sequence.AccNoDest,
		Amount:             amount,
		Fee:                transferFee,
		Remark:             sequence.Remark(),
	})
	if err != nil {
		t.log.ServiceUsecase(transferService, "DoPayment").Errorf("PerformOverbooking: %v", err)
		return nil, status.Error(codes.Internal, domain.ErrGeneral)
	}

	user, ok := ctxt.UserFromContext(ctx)
	if !ok {
		t.log.ServiceUsecase(transferService, "DoPayment").Errorf("GetUserFromContext: %v", ctxt.ErrUserFromContext)
		return nil, status.Error(codes.Unauthenticated, domain.ErrUserUnauthenticated)
	}

	transaction := &entity.Transaction{
		SequenceNumber:       sequence.SeqNo,
		SequenceJournal:      result.JournalSequence,
		UserID:               strconv.Itoa(user.Id),
		Destination:          sequence.AccNoDest,
		Amount:               sequence.Amount,
		TransactionType:      transferType,
		TransactionReference: result.TransactionReference,
		Remarks:              sequence.Remark(),
		Fee:                  stringTransferFee,
		NameDest:             sequence.AccNameDest,
	}

	err = t.repo.InsertTransaction(ctx, transaction)
	if err != nil {
		t.log.ServiceUsecase(transferService, "DoPayment").Errorf("InsertTransaction: %v", err)
		return nil, status.Error(codes.Internal, domain.ErrGeneral)
	}

	err = t.mailer.SendTransferReceipt(ctx, &entity.TransferEmailData{
		Subject:            transferSuccessSubject,
		Recipient:          user.Email,
		Amount:             amount,
		Fee:                transferFee,
		SourceName:         user.FullName,
		SourceAccount:      sequence.AccNoSrc,
		DestinationName:    sequence.AccNameDest,
		DestinationAccount: sequence.AccNoDest,
		DestinationBank:    constant.CompanyName,
		TransactionRef:     result.TransactionReference,
		Note:               sequence.Remark(),
	})
	if err != nil {
		t.log.ServiceUsecase(transferService, "DoPayment").Errorf("SendTransferReceipt: %v", err)
		return nil, status.Error(codes.Internal, domain.ErrSendEmailFailed)
	}

	return transaction, nil
}

func (t *Transfer) ProcessTransfer(ctx context.Context, event *entity.TransferRequest) (*entity.Transaction, error) {
	sequence, err := t.Inquiry(ctx, &entity.Sequence{
		Amount:    event.Amount.String(),
		AccNoSrc:  event.AccountNumber,
		AccNoDest: event.Destination,
	})
	if err != nil {
		t.log.ServiceUsecase(transferService, "ProcessTransfer").Errorf("Inquiry: %v", err)
		return nil, status.Error(codes.Internal, domain.ErrGeneral)
	}
	transaction, err := t.DoPayment(ctx, sequence.SeqNo)
	if err != nil {
		t.log.ServiceUsecase(transferService, "ProcessTransfer").Errorf("Payment: %v", err)
		return nil, status.Error(codes.Internal, domain.ErrGeneral)
	}
	return transaction, nil
}
