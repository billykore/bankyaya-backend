package service

import (
	"context"
	"strconv"
	"time"

	"go.bankyaya.org/app/backend/pkg/entity"
	pkgerrors "go.bankyaya.org/app/backend/pkg/errors"
	"go.bankyaya.org/app/backend/pkg/interface/api"
	"go.bankyaya.org/app/backend/pkg/interface/email"
	"go.bankyaya.org/app/backend/pkg/interface/repository"
	"go.bankyaya.org/app/backend/pkg/util/codes"
	"go.bankyaya.org/app/backend/pkg/util/constant"
	"go.bankyaya.org/app/backend/pkg/util/ctxt"
	"go.bankyaya.org/app/backend/pkg/util/logger"
	"go.bankyaya.org/app/backend/pkg/util/status"
	"go.bankyaya.org/app/backend/pkg/util/types"
	"go.bankyaya.org/app/backend/pkg/util/uuid"
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
	mailer      email.TransferReceiptMailer
}

func NewTransfer(log *logger.Logger, repo repository.TransferRepository, corebanking api.CoreBanking, mailer email.TransferReceiptMailer) *Transfer {
	return &Transfer{
		log:         log,
		repo:        repo,
		corebanking: corebanking,
		mailer:      mailer,
	}
}

func (t *Transfer) Inquiry(ctx context.Context, seq *entity.Sequence) (*entity.Sequence, error) {
	eod, err := t.corebanking.CheckEOD(ctx)
	if err != nil {
		t.log.ServiceUsecase(transferService, "Inquiry").Errorf("CheckEOD: %v", err)
		return nil, status.Error(codes.Internal, pkgerrors.ErrGeneral)
	}
	if eod.IsRunning() {
		t.log.ServiceUsecase(transferService, "Inquiry").Errorf("CheckEOD: %v", pkgerrors.ErrEODInProgress)
		return nil, status.Error(codes.Internal, pkgerrors.ErrEODInProgress)
	}

	srcAccount, err := t.corebanking.GetAccountDetails(ctx, seq.AccNoSrc)
	if err != nil {
		t.log.ServiceUsecase(transferService, "Inquiry").Errorf("GetAccountDetails: %v", err)
		return nil, status.Error(codes.Internal, pkgerrors.ErrGeneral)
	}
	if !srcAccount.IsAccountActive() {
		t.log.ServiceUsecase(transferService, "Inquiry").Errorf("GetAccountDetails: %v", err)
		return nil, status.Error(codes.BadRequest, pkgerrors.ErrSourceAccountInactive)
	}

	destAccount, err := t.corebanking.GetAccountDetails(ctx, seq.AccNameDest)
	if err != nil {
		t.log.ServiceUsecase(transferService, "Inquiry").Errorf("GetAccountDetails: %v", err)
		return nil, status.Error(codes.Internal, pkgerrors.ErrGeneral)
	}
	if !destAccount.IsAccountActive() {
		t.log.ServiceUsecase(transferService, "Inquiry").Errorf("GetAccountDetails: %v", err)
		return nil, status.Error(codes.BadRequest, pkgerrors.ErrSourceAccountInactive)
	}
	seq.AccNameDest = destAccount.Name

	sequenceNo, err := uuid.New()
	if err != nil {
		t.log.ServiceUsecase(transferService, "Inquiry").Errorf("New uuid failed: %v", err)
		return nil, status.Error(codes.Internal, pkgerrors.ErrGeneral)
	}
	seq.SeqNo = sequenceNo

	err = t.repo.InsertSequence(ctx, seq)
	if err != nil {
		t.log.ServiceUsecase(transferService, "Inquiry").Errorf("InsertSequence: %v", err)
		return nil, status.Error(codes.Internal, pkgerrors.ErrGeneral)
	}

	return seq, nil
}

func (t *Transfer) DoPayment(ctx context.Context, sequenceNumber string) (*entity.Transaction, error) {
	eod, err := t.corebanking.CheckEOD(ctx)
	if err != nil {
		t.log.ServiceUsecase(transferService, "DoPayment").Errorf("CheckEOD: %v", err)
		return nil, status.Error(codes.Internal, pkgerrors.ErrGeneral)
	}
	if eod.IsRunning() {
		t.log.ServiceUsecase(transferService, "DoPayment").Errorf("CheckEOD: %v", pkgerrors.ErrEODInProgress)
		return nil, status.Error(codes.Internal, pkgerrors.ErrEODInProgress)
	}

	sequence, err := t.repo.GetSequence(ctx, sequenceNumber)
	if err != nil {
		t.log.ServiceUsecase(transferService, "DoPayment").Errorf("GetSequence: %v", err)
		return nil, status.Error(codes.Internal, pkgerrors.ErrGeneral)
	}
	if sequence.SeqNo == "" {
		t.log.ServiceUsecase(transferService, "DoPayment").Errorf("GetSequence: %v", err)
		return nil, status.Error(codes.BadRequest, pkgerrors.ErrInvalidSequenceNumber)
	}
	amount, err := types.ParseMoney(sequence.Amount)
	if err != nil {
		t.log.ServiceUsecase(transferService, "DoPayment").Errorf("ParseMoney: %v", err)
		return nil, status.Error(codes.BadRequest, pkgerrors.ErrGeneral)
	}

	result, err := t.corebanking.PerformOverbooking(ctx, &entity.OverbookingRequest{
		SourceAccount:      sequence.AccNameSrc,
		DestinationAccount: sequence.AccNameDest,
		Amount:             amount,
		Fee:                transferFee,
		Remark:             sequence.Remark(),
	})
	if err != nil {
		t.log.ServiceUsecase(transferService, "DoPayment").Errorf("PerformOverbooking: %v", err)
		return nil, status.Error(codes.Internal, pkgerrors.ErrGeneral)
	}

	user, ok := ctxt.UserFromContext(ctx)
	if !ok {
		t.log.ServiceUsecase(transferService, "DoPayment").Errorf("GetUserFromContext: %v", ctxt.ErrUserFromContext)
		return nil, status.Error(codes.Unauthenticated, pkgerrors.ErrUserUnauthenticated)
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
		CreatedAt:            time.Now(),
		Fee:                  stringTransferFee,
		NameDest:             sequence.AccNameDest,
	}

	err = t.repo.InsertTransaction(ctx, transaction)
	if err != nil {
		t.log.ServiceUsecase(transferService, "DoPayment").Errorf("InsertTransaction: %v", err)
		return nil, status.Error(codes.Internal, pkgerrors.ErrGeneral)
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
		return nil, status.Error(codes.Internal, pkgerrors.ErrSendEmailFailed)
	}

	return transaction, nil
}

func (t *Transfer) ProcessAutoTransferEvent(ctx context.Context, event *entity.Event) (*entity.Transaction, error) {
	sequence, err := t.Inquiry(ctx, &entity.Sequence{
		Amount:    event.Amount.String(),
		AccNoSrc:  event.AccountNumber,
		AccNoDest: event.Destination,
	})
	if err != nil {
		t.log.ServiceUsecase(transferService, "ProcessAutoTransferEvent").Errorf("Inquiry: %v", err)
		return nil, status.Error(codes.Internal, pkgerrors.ErrGeneral)
	}
	transaction, err := t.DoPayment(ctx, sequence.SeqNo)
	if err != nil {
		t.log.ServiceUsecase(transferService, "ProcessAutoTransferEvent").Errorf("Payment: %v", err)
		return nil, status.Error(codes.Internal, pkgerrors.ErrGeneral)
	}
	return transaction, nil
}
