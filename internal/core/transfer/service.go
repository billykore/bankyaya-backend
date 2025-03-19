package transfer

import (
	"context"
	"strconv"
	"time"

	"go.bankyaya.org/app/backend/pkg/constant"
	"go.bankyaya.org/app/backend/pkg/ctxt"
	"go.bankyaya.org/app/backend/pkg/types"
	"go.bankyaya.org/app/backend/pkg/uuid"
)

const (
	transferType           = "internal_transfer"
	transferSuccessSubject = "Transfer success"
	transferFee            = 0
	stringTransferFee      = "0"
)

// Service handles intra-bank transfer process.
type Service struct {
	repo        Repository
	corebanking CoreBanking
	mailer      ReceiptMailer
}

func NewService(repo Repository, corebanking CoreBanking, mailer ReceiptMailer) *Service {
	return &Service{
		repo:        repo,
		corebanking: corebanking,
		mailer:      mailer,
	}
}

func (s *Service) Inquiry(ctx context.Context, seq *Sequence) (*Sequence, error) {
	eod, err := s.corebanking.CheckEOD(ctx)
	if err != nil {
		return nil, err
	}
	if eod.IsRunning() {
		return nil, err
	}

	srcAccount, err := s.corebanking.GetAccountDetails(ctx, seq.AccNoSrc)
	if err != nil {
		return nil, err
	}
	if !srcAccount.IsAccountActive() {
		return nil, ErrSourceAccountInactive
	}
	destAccount, err := s.corebanking.GetAccountDetails(ctx, seq.AccNameDest)
	if err != nil {
		return nil, err
	}
	if !destAccount.IsAccountActive() {
		return nil, ErrDestinationAccountInactive
	}
	seq.AccNameDest = destAccount.Name

	sequenceNo, err := uuid.New()
	if err != nil {
		return nil, err
	}
	seq.SeqNo = sequenceNo

	err = s.repo.InsertSequence(ctx, seq)
	if err != nil {
		return nil, err
	}

	return seq, nil
}

func (s *Service) DoPayment(ctx context.Context, sequenceNumber string) (*Transaction, error) {
	eod, err := s.corebanking.CheckEOD(ctx)
	if err != nil {
		return nil, err
	}
	if eod.IsRunning() {
		return nil, err
	}

	sequence, err := s.repo.GetSequence(ctx, sequenceNumber)
	if err != nil {
		return nil, err
	}
	if sequence.SeqNo == "" {
		return nil, ErrInvalidSequenceNumber
	}
	amount, err := types.ParseMoney(sequence.Amount)
	if err != nil {
		return nil, err
	}

	result, err := s.corebanking.PerformOverbooking(ctx, &OverbookingRequest{
		SourceAccount:      sequence.AccNameSrc,
		DestinationAccount: sequence.AccNameDest,
		Amount:             amount,
		Fee:                transferFee,
		Remark:             sequence.Remark(),
	})
	if err != nil {
		return nil, err
	}

	user, ok := ctxt.UserFromContext(ctx)
	if !ok {
		return nil, ctxt.ErrUserFromContext
	}

	transaction := &Transaction{
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

	err = s.repo.InsertTransaction(ctx, transaction)
	if err != nil {
		return nil, err
	}

	err = s.mailer.SendTransferReceipt(ctx, &EmailData{
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
		return nil, ErrSendEmailFailed
	}

	return transaction, nil
}

func (s *Service) ProcessAutoTransferEvent(ctx context.Context, event *Event) (*Transaction, error) {
	sequence, err := s.Inquiry(ctx, &Sequence{
		Amount:    event.Amount.String(),
		AccNoSrc:  event.AccountNumber,
		AccNoDest: event.Destination,
	})
	if err != nil {
		return nil, err
	}
	transaction, err := s.DoPayment(ctx, sequence.SeqNo)
	if err != nil {
		return nil, err
	}
	return transaction, nil
}
