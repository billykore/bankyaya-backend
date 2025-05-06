package intrabank

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.bankyaya.org/app/backend/internal/pkg/codes"
	"go.bankyaya.org/app/backend/internal/pkg/ctxt"
	"go.bankyaya.org/app/backend/internal/pkg/logger"
	"go.bankyaya.org/app/backend/internal/pkg/pkgerror"
)

func TestTransferInquirySuccess(t *testing.T) {
	var (
		corebankingMock = NewMockCoreBanking(t)
		repoMock        = NewMockRepository(t)
		mailerMock      = NewMockReceiptMailer(t)
		seqGenMock      = NewMockSequenceGenerator(t)
		notifierMock    = NewMockNotifier(t)
		svc             = NewService(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock, notifierMock)
		ctx             = ctxt.ContextWithUser(context.Background(), &ctxt.User{
			ID:    123,
			CIF:   "1234567",
			Name:  "Olivia Rodrigo",
			Email: "olivia@gmail.com",
		})
	)

	corebankingMock.EXPECT().GetCoreStatus(mock.Anything).Return(&CoreStatus{
		SystemDate:    "25-03-2025",
		Status:        "FINISHED",
		StandInStatus: "N",
	}, nil)
	corebankingMock.EXPECT().GetAccountDetails(mock.Anything, "001001234567891").
		Return(&Account{
			Name:   "Olivia Rodrigo",
			Status: "1",
		}, nil)
	corebankingMock.EXPECT().GetAccountDetails(mock.Anything, "001001234567892").
		Return(&Account{
			Name:   "Destination Account",
			Status: "1",
		}, nil)

	repoMock.EXPECT().GetTransactionLimit(mock.Anything).
		Return(&Limits{
			MinAmount:      1,
			MaxAmount:      50_000_000,
			MaxDailyAmount: 200_000_000,
		}, nil)
	repoMock.EXPECT().InsertSequence(mock.Anything, &Sequence{
		SequenceNumber:     "123456",
		Amount:             100000,
		SourceAccount:      "001001234567891",
		DestinationAccount: "001001234567892",
		DestinationName:    "Destination Account",
		SourceName:         "Olivia Rodrigo",
	}).Return(nil)

	seqGenMock.EXPECT().Generate().
		Return("123456", nil)

	sequence, err := svc.Inquiry(ctx, &Sequence{
		SequenceNumber:     "123456",
		Amount:             100000,
		SourceAccount:      "001001234567891",
		DestinationAccount: "001001234567892",
	})

	assert.Nil(t, err)
	assert.Equal(t, sequence, &Sequence{
		SequenceNumber:     "123456",
		Amount:             100000,
		SourceAccount:      "001001234567891",
		DestinationAccount: "001001234567892",
		DestinationName:    "Destination Account",
		SourceName:         "Olivia Rodrigo",
	})

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}

func TestTransferInquiryFailed_CheckEODFailed(t *testing.T) {
	var (
		corebankingMock = NewMockCoreBanking(t)
		repoMock        = NewMockRepository(t)
		mailerMock      = NewMockReceiptMailer(t)
		seqGenMock      = NewMockSequenceGenerator(t)
		notifierMock    = NewMockNotifier(t)
		svc             = NewService(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock, notifierMock)
		ctx             = ctxt.ContextWithUser(context.Background(), &ctxt.User{
			ID:    123,
			CIF:   "1234567",
			Name:  "Olivia Rodrigo",
			Email: "olivia@gmail.com",
		})
	)

	corebankingMock.EXPECT().GetCoreStatus(mock.Anything).
		Return(nil, errors.New("check EOD failed"))

	sequence, err := svc.Inquiry(ctx, &Sequence{
		SequenceNumber:     "123456",
		Amount:             100000,
		SourceAccount:      "001001234567891",
		DestinationAccount: "001001234567892",
	})
	assert.Nil(t, sequence)
	assert.Equal(t, pkgerror.New(codes.Internal, ErrGeneral), err)

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}

func TestTransferInquiryFailed_EODIsRunning(t *testing.T) {
	var (
		corebankingMock = NewMockCoreBanking(t)
		repoMock        = NewMockRepository(t)
		mailerMock      = NewMockReceiptMailer(t)
		seqGenMock      = NewMockSequenceGenerator(t)
		notifierMock    = NewMockNotifier(t)
		svc             = NewService(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock, notifierMock)
		ctx             = ctxt.ContextWithUser(context.Background(), &ctxt.User{
			ID:    123,
			CIF:   "1234567",
			Name:  "Olivia Rodrigo",
			Email: "olivia@gmail.com",
		})
	)

	corebankingMock.EXPECT().GetCoreStatus(mock.Anything).Return(&CoreStatus{
		SystemDate:    "25-03-2025",
		Status:        "STARTED",
		StandInStatus: "N",
	}, nil)

	sequence, err := svc.Inquiry(ctx, &Sequence{
		SequenceNumber:     "123456",
		Amount:             100000,
		SourceAccount:      "001001234567891",
		DestinationAccount: "001001234567892",
	})
	assert.Nil(t, sequence)
	assert.Equal(t, pkgerror.New(codes.Internal, ErrEODInProgress), err)

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}

func TestTransferInquiryFailed_GetTransactionLimitFailed(t *testing.T) {
	var (
		corebankingMock = NewMockCoreBanking(t)
		repoMock        = NewMockRepository(t)
		mailerMock      = NewMockReceiptMailer(t)
		seqGenMock      = NewMockSequenceGenerator(t)
		notifierMock    = NewMockNotifier(t)
		svc             = NewService(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock, notifierMock)
		ctx             = ctxt.ContextWithUser(context.Background(), &ctxt.User{
			ID:    123,
			CIF:   "1234567",
			Name:  "Olivia Rodrigo",
			Email: "olivia@gmail.com",
		})
	)

	corebankingMock.EXPECT().GetCoreStatus(mock.Anything).Return(&CoreStatus{
		SystemDate:    "25-03-2025",
		Status:        "FINISHED",
		StandInStatus: "N",
	}, nil)

	repoMock.EXPECT().GetTransactionLimit(mock.Anything).
		Return(nil, errors.New("some error"))

	sequence, err := svc.Inquiry(ctx, &Sequence{
		SequenceNumber:     "123456",
		Amount:             100000,
		SourceAccount:      "001001234567891",
		DestinationAccount: "001001234567892",
	})

	assert.Nil(t, sequence)
	assert.Equal(t, pkgerror.New(codes.Internal, ErrGeneral), err)

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}

func TestTransferInquiryFailed_TransactionLimitCannotTransfer(t *testing.T) {
	var (
		corebankingMock = NewMockCoreBanking(t)
		repoMock        = NewMockRepository(t)
		mailerMock      = NewMockReceiptMailer(t)
		seqGenMock      = NewMockSequenceGenerator(t)
		notifierMock    = NewMockNotifier(t)
		svc             = NewService(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock, notifierMock)
		ctx             = ctxt.ContextWithUser(context.Background(), &ctxt.User{
			ID:    123,
			CIF:   "1234567",
			Name:  "Olivia Rodrigo",
			Email: "olivia@gmail.com",
		})
	)

	corebankingMock.EXPECT().GetCoreStatus(mock.Anything).Return(&CoreStatus{
		SystemDate:    "25-03-2025",
		Status:        "FINISHED",
		StandInStatus: "N",
	}, nil)

	repoMock.EXPECT().GetTransactionLimit(mock.Anything).
		Return(&Limits{
			MinAmount:      1,
			MaxAmount:      100_000_000,
			MaxDailyAmount: 50_000,
		}, nil)

	sequence, err := svc.Inquiry(ctx, &Sequence{
		SequenceNumber:     "123456",
		Amount:             100000,
		SourceAccount:      "001001234567891",
		DestinationAccount: "001001234567892",
	})

	assert.Nil(t, sequence)
	assert.Equal(t, pkgerror.New(codes.BadRequest, ErrInvalidAmount), err)

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}

func TestTransferInquiryFailed_FailedCheckSourceAccount(t *testing.T) {
	var (
		corebankingMock = NewMockCoreBanking(t)
		repoMock        = NewMockRepository(t)
		mailerMock      = NewMockReceiptMailer(t)
		seqGenMock      = NewMockSequenceGenerator(t)
		notifierMock    = NewMockNotifier(t)
		svc             = NewService(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock, notifierMock)
		ctx             = ctxt.ContextWithUser(context.Background(), &ctxt.User{
			ID:    123,
			CIF:   "1234567",
			Name:  "Olivia Rodrigo",
			Email: "olivia@gmail.com",
		})
	)

	corebankingMock.EXPECT().GetCoreStatus(mock.Anything).Return(&CoreStatus{
		SystemDate:    "25-03-2025",
		Status:        "FINISHED",
		StandInStatus: "N",
	}, nil)
	corebankingMock.EXPECT().GetAccountDetails(mock.Anything, "001001234567891").
		Return(nil, errors.New("GetAccountDetails failed"))

	repoMock.EXPECT().GetTransactionLimit(mock.Anything).
		Return(&Limits{
			MinAmount:      1,
			MaxAmount:      50_000_000,
			MaxDailyAmount: 200_000_000,
		}, nil)

	sequence, err := svc.Inquiry(ctx, &Sequence{
		SequenceNumber:     "123456",
		Amount:             100000,
		SourceAccount:      "001001234567891",
		DestinationAccount: "001001234567892",
	})

	assert.Nil(t, sequence)
	assert.Equal(t, pkgerror.New(codes.Internal, ErrGeneral), err)

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}

func TestTransferInquiryFailed_SourceAccountIsInactive(t *testing.T) {
	var (
		corebankingMock = NewMockCoreBanking(t)
		repoMock        = NewMockRepository(t)
		mailerMock      = NewMockReceiptMailer(t)
		seqGenMock      = NewMockSequenceGenerator(t)
		notifierMock    = NewMockNotifier(t)
		svc             = NewService(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock, notifierMock)
		ctx             = ctxt.ContextWithUser(context.Background(), &ctxt.User{
			ID:    123,
			CIF:   "1234567",
			Name:  "Olivia Rodrigo",
			Email: "olivia@gmail.com",
		})
	)

	corebankingMock.EXPECT().GetCoreStatus(mock.Anything).Return(&CoreStatus{
		SystemDate:    "25-03-2025",
		Status:        "FINISHED",
		StandInStatus: "N",
	}, nil)
	corebankingMock.EXPECT().GetAccountDetails(mock.Anything, "001001234567891").
		Return(&Account{
			Status: "9",
		}, nil)

	repoMock.EXPECT().GetTransactionLimit(mock.Anything).
		Return(&Limits{
			MinAmount:      1,
			MaxAmount:      50_000_000,
			MaxDailyAmount: 200_000_000,
		}, nil)

	sequence, err := svc.Inquiry(ctx, &Sequence{
		SequenceNumber:     "123456",
		Amount:             100000,
		SourceAccount:      "001001234567891",
		DestinationAccount: "001001234567892",
	})

	assert.Nil(t, sequence)
	assert.Equal(t, pkgerror.New(codes.BadRequest, ErrSourceAccountInactive), err)

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}

func TestTransferInquiryFailed_CheckDestinationAccountFailed(t *testing.T) {
	var (
		corebankingMock = NewMockCoreBanking(t)
		repoMock        = NewMockRepository(t)
		mailerMock      = NewMockReceiptMailer(t)
		seqGenMock      = NewMockSequenceGenerator(t)
		notifierMock    = NewMockNotifier(t)
		svc             = NewService(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock, notifierMock)
		ctx             = ctxt.ContextWithUser(context.Background(), &ctxt.User{
			ID:    123,
			CIF:   "1234567",
			Name:  "Olivia Rodrigo",
			Email: "olivia@gmail.com",
		})
	)

	corebankingMock.EXPECT().GetCoreStatus(mock.Anything).Return(&CoreStatus{
		SystemDate:    "25-03-2025",
		Status:        "FINISHED",
		StandInStatus: "N",
	}, nil)
	corebankingMock.EXPECT().GetAccountDetails(mock.Anything, "001001234567891").
		Return(&Account{
			Name:   "Olivia Rodrigo",
			Status: "1",
		}, nil)
	corebankingMock.EXPECT().GetAccountDetails(mock.Anything, "001001234567892").
		Return(nil, errors.New("GetAccountDetails failed"))

	repoMock.EXPECT().GetTransactionLimit(mock.Anything).
		Return(&Limits{
			MinAmount:      1,
			MaxAmount:      50_000_000,
			MaxDailyAmount: 200_000_000,
		}, nil)

	sequence, err := svc.Inquiry(ctx, &Sequence{
		SequenceNumber:     "123456",
		Amount:             100000,
		SourceAccount:      "001001234567891",
		DestinationAccount: "001001234567892",
	})

	assert.Nil(t, sequence)
	assert.Equal(t, pkgerror.New(codes.Internal, ErrGeneral), err)

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}

func TestTransferInquiryFailed_DestinationAccountInactive(t *testing.T) {
	var (
		corebankingMock = NewMockCoreBanking(t)
		repoMock        = NewMockRepository(t)
		mailerMock      = NewMockReceiptMailer(t)
		seqGenMock      = NewMockSequenceGenerator(t)
		notifierMock    = NewMockNotifier(t)
		svc             = NewService(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock, notifierMock)
		ctx             = ctxt.ContextWithUser(context.Background(), &ctxt.User{
			ID:    123,
			CIF:   "1234567",
			Name:  "Olivia Rodrigo",
			Email: "olivia@gmail.com",
		})
	)

	repoMock.EXPECT().GetTransactionLimit(mock.Anything).
		Return(&Limits{
			MinAmount:      1,
			MaxAmount:      50_000_000,
			MaxDailyAmount: 200_000_000,
		}, nil)

	corebankingMock.EXPECT().GetCoreStatus(mock.Anything).Return(&CoreStatus{
		SystemDate:    "25-03-2025",
		Status:        "FINISHED",
		StandInStatus: "N",
	}, nil)
	corebankingMock.EXPECT().GetAccountDetails(mock.Anything, "001001234567891").
		Return(&Account{
			Name:   "Olivia Rodrigo",
			Status: "1",
		}, nil)
	corebankingMock.EXPECT().GetAccountDetails(mock.Anything, "001001234567892").
		Return(&Account{
			Status: "9",
		}, nil)

	sequence, err := svc.Inquiry(ctx, &Sequence{
		SequenceNumber:     "123456",
		Amount:             100000,
		SourceAccount:      "001001234567891",
		DestinationAccount: "001001234567892",
	})

	assert.Nil(t, sequence)
	assert.Equal(t, pkgerror.New(codes.BadRequest, ErrDestinationAccountInactive), err)

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}

func TestTransferInquiryFailed_GenerateSequenceFailed(t *testing.T) {
	var (
		corebankingMock = NewMockCoreBanking(t)
		repoMock        = NewMockRepository(t)
		mailerMock      = NewMockReceiptMailer(t)
		seqGenMock      = NewMockSequenceGenerator(t)
		notifierMock    = NewMockNotifier(t)
		svc             = NewService(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock, notifierMock)
		ctx             = ctxt.ContextWithUser(context.Background(), &ctxt.User{
			ID:    123,
			CIF:   "1234567",
			Name:  "Olivia Rodrigo",
			Email: "olivia@gmail.com",
		})
	)

	repoMock.EXPECT().GetTransactionLimit(mock.Anything).
		Return(&Limits{
			MinAmount:      1,
			MaxAmount:      50_000_000,
			MaxDailyAmount: 200_000_000,
		}, nil)

	corebankingMock.EXPECT().GetCoreStatus(mock.Anything).Return(&CoreStatus{
		SystemDate:    "25-03-2025",
		Status:        "FINISHED",
		StandInStatus: "N",
	}, nil)
	corebankingMock.EXPECT().GetAccountDetails(mock.Anything, "001001234567891").
		Return(&Account{
			Name:   "Olivia Rodrigo",
			Status: "1",
		}, nil)
	corebankingMock.EXPECT().GetAccountDetails(mock.Anything, "001001234567892").
		Return(&Account{
			Name:   "Destination Account",
			Status: "1",
		}, nil)

	seqGenMock.EXPECT().Generate().
		Return("", errors.New("some error"))

	sequence, err := svc.Inquiry(ctx, &Sequence{
		SequenceNumber:     "123456",
		Amount:             100000,
		SourceAccount:      "001001234567891",
		DestinationAccount: "001001234567892",
	})

	assert.Nil(t, sequence)
	assert.Equal(t, pkgerror.New(codes.Internal, ErrGeneral), err)

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}

func TestTransferInquiryFailed_InsertSequenceFailed(t *testing.T) {
	var (
		corebankingMock = NewMockCoreBanking(t)
		repoMock        = NewMockRepository(t)
		mailerMock      = NewMockReceiptMailer(t)
		seqGenMock      = NewMockSequenceGenerator(t)
		notifierMock    = NewMockNotifier(t)
		svc             = NewService(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock, notifierMock)
		ctx             = ctxt.ContextWithUser(context.Background(), &ctxt.User{
			ID:    123,
			CIF:   "1234567",
			Name:  "Olivia Rodrigo",
			Email: "olivia@gmail.com",
		})
	)

	corebankingMock.EXPECT().GetCoreStatus(mock.Anything).Return(&CoreStatus{
		SystemDate:    "25-03-2025",
		Status:        "FINISHED",
		StandInStatus: "N",
	}, nil)
	corebankingMock.EXPECT().GetAccountDetails(mock.Anything, "001001234567891").
		Return(&Account{
			Name:   "Olivia Rodrigo",
			Status: "1",
		}, nil)
	corebankingMock.EXPECT().GetAccountDetails(mock.Anything, "001001234567892").
		Return(&Account{
			Name:   "Destination Account",
			Status: "1",
		}, nil)

	seqGenMock.EXPECT().Generate().
		Return("123456", nil)

	repoMock.EXPECT().GetTransactionLimit(mock.Anything).
		Return(&Limits{
			MinAmount:      1,
			MaxAmount:      50_000_000,
			MaxDailyAmount: 200_000_000,
		}, nil)
	repoMock.EXPECT().InsertSequence(mock.Anything, &Sequence{
		SequenceNumber:     "123456",
		Amount:             100000,
		SourceAccount:      "001001234567891",
		DestinationAccount: "001001234567892",
		DestinationName:    "Destination Account",
		SourceName:         "Olivia Rodrigo",
	}).Return(errors.New("some error"))

	sequence, err := svc.Inquiry(ctx, &Sequence{
		SequenceNumber:     "123456",
		Amount:             100000,
		SourceAccount:      "001001234567891",
		DestinationAccount: "001001234567892",
	})

	assert.Nil(t, sequence)
	assert.Equal(t, pkgerror.New(codes.Internal, ErrGeneral), err)

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}

func TestTransferDoPaymentSuccess(t *testing.T) {
	var (
		corebankingMock = NewMockCoreBanking(t)
		repoMock        = NewMockRepository(t)
		mailerMock      = NewMockReceiptMailer(t)
		seqGenMock      = NewMockSequenceGenerator(t)
		notifierMock    = NewMockNotifier(t)
		svc             = NewService(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock, notifierMock)
		ctx             = ctxt.ContextWithUser(context.Background(), &ctxt.User{
			ID:    123,
			CIF:   "1234567",
			Name:  "Olivia Rodrigo",
			Email: "olivia@gmail.com",
		})
	)

	corebankingMock.EXPECT().GetCoreStatus(mock.Anything).Return(&CoreStatus{
		SystemDate:    "25-03-2025",
		Status:        "FINISHED",
		StandInStatus: "N",
	}, nil)

	repoMock.EXPECT().GetTransactionLimit(mock.Anything).
		Return(&Limits{
			MinAmount:      1,
			MaxAmount:      50_000_000,
			MaxDailyAmount: 200_000_000,
		}, nil)
	repoMock.EXPECT().GetSequence(mock.Anything, "123456").
		Return(&Sequence{
			SequenceNumber:     "123456",
			Amount:             100000,
			SourceAccount:      "001001234567891",
			DestinationAccount: "001001234567892",
			DestinationName:    "Destination Account",
			SourceName:         "Olivia Rodrigo",
		}, nil)

	corebankingMock.EXPECT().PerformOverbooking(mock.Anything, &OverbookingInput{
		SourceAccount:      "001001234567891",
		DestinationAccount: "001001234567892",
		Amount:             100000,
		Fee:                0,
		Remark:             "TRF 001001234567891 001001234567892 BNKYAYA 123456",
	}).Return(&OverbookingResult{
		JournalSequence:      "111111",
		TransactionReference: "222222",
	}, nil)

	repoMock.EXPECT().InsertTransaction(mock.Anything, &Transaction{
		SequenceNumber:       "123456",
		SequenceJournal:      "111111",
		UserID:               "123",
		Destination:          "001001234567892",
		Amount:               100000,
		TransactionType:      "internal_transfer",
		TransactionReference: "222222",
		Remarks:              "TRF 001001234567891 001001234567892 BNKYAYA 123456",
		Fee:                  "0",
		DestinationName:      "Destination Account",
	}).Return(nil)

	mailerMock.EXPECT().SendReceipt(mock.Anything, mock.Anything).
		Return(nil)

	notifierMock.EXPECT().Notify(mock.Anything, mock.Anything).
		Return(nil)

	transaction, err := svc.DoPayment(ctx, "123456")

	assert.NoError(t, err)
	assert.Equal(t, &Transaction{
		SequenceNumber:       "123456",
		SequenceJournal:      "111111",
		UserID:               "123",
		Destination:          "001001234567892",
		Amount:               100000,
		TransactionType:      "internal_transfer",
		TransactionReference: "222222",
		Remarks:              "TRF 001001234567891 001001234567892 BNKYAYA 123456",
		Fee:                  "0",
		DestinationName:      "Destination Account",
	}, transaction)

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}

func TestTransferDoPaymentFailed_CheckEODFailed(t *testing.T) {
	var (
		corebankingMock = NewMockCoreBanking(t)
		repoMock        = NewMockRepository(t)
		mailerMock      = NewMockReceiptMailer(t)
		seqGenMock      = NewMockSequenceGenerator(t)
		notifierMock    = NewMockNotifier(t)
		svc             = NewService(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock, notifierMock)
		ctx             = ctxt.ContextWithUser(context.Background(), &ctxt.User{
			ID:    123,
			CIF:   "1234567",
			Name:  "Olivia Rodrigo",
			Email: "olivia@gmail.com",
		})
	)

	corebankingMock.EXPECT().GetCoreStatus(mock.Anything).
		Return(nil, errors.New("some error"))

	transaction, err := svc.DoPayment(ctx, "123456")

	assert.Nil(t, transaction)
	assert.Equal(t, pkgerror.New(codes.Internal, ErrGeneral), err)

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}

func TestTransferDoPaymentFailed_EODIsRunning(t *testing.T) {
	var (
		corebankingMock = NewMockCoreBanking(t)
		repoMock        = NewMockRepository(t)
		mailerMock      = NewMockReceiptMailer(t)
		seqGenMock      = NewMockSequenceGenerator(t)
		notifierMock    = NewMockNotifier(t)
		svc             = NewService(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock, notifierMock)
		ctx             = ctxt.ContextWithUser(context.Background(), &ctxt.User{
			ID:    123,
			CIF:   "1234567",
			Name:  "Olivia Rodrigo",
			Email: "olivia@gmail.com",
		})
	)

	corebankingMock.EXPECT().GetCoreStatus(mock.Anything).
		Return(&CoreStatus{
			SystemDate:    "25-03-2025",
			Status:        "STARTED",
			StandInStatus: "N",
		}, nil)

	transaction, err := svc.DoPayment(ctx, "123456")

	assert.Nil(t, transaction)
	assert.Equal(t, pkgerror.New(codes.Internal, ErrEODInProgress), err)

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}

func TestTransferDoPaymentFailed_GetSequenceFailed(t *testing.T) {
	var (
		corebankingMock = NewMockCoreBanking(t)
		repoMock        = NewMockRepository(t)
		mailerMock      = NewMockReceiptMailer(t)
		seqGenMock      = NewMockSequenceGenerator(t)
		notifierMock    = NewMockNotifier(t)
		svc             = NewService(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock, notifierMock)
		ctx             = ctxt.ContextWithUser(context.Background(), &ctxt.User{
			ID:    123,
			CIF:   "1234567",
			Name:  "Olivia Rodrigo",
			Email: "olivia@gmail.com",
		})
	)

	corebankingMock.EXPECT().GetCoreStatus(mock.Anything).Return(&CoreStatus{
		SystemDate:    "25-03-2025",
		Status:        "FINISHED",
		StandInStatus: "N",
	}, nil)

	repoMock.EXPECT().GetSequence(mock.Anything, "123456").
		Return(nil, errors.New("some error"))

	transaction, err := svc.DoPayment(ctx, "123456")

	assert.Nil(t, transaction)
	assert.Equal(t, pkgerror.New(codes.Internal, ErrGeneral), err)

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}

func TestTransferDoPaymentFailed_InvalidSequenceNumber(t *testing.T) {
	var (
		corebankingMock = NewMockCoreBanking(t)
		repoMock        = NewMockRepository(t)
		mailerMock      = NewMockReceiptMailer(t)
		seqGenMock      = NewMockSequenceGenerator(t)
		notifierMock    = NewMockNotifier(t)
		svc             = NewService(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock, notifierMock)
		ctx             = ctxt.ContextWithUser(context.Background(), &ctxt.User{
			ID:    123,
			CIF:   "1234567",
			Name:  "Olivia Rodrigo",
			Email: "olivia@gmail.com",
		})
	)

	corebankingMock.EXPECT().GetCoreStatus(mock.Anything).Return(&CoreStatus{
		SystemDate:    "25-03-2025",
		Status:        "FINISHED",
		StandInStatus: "N",
	}, nil)

	repoMock.EXPECT().GetSequence(mock.Anything, "123456").
		Return(&Sequence{
			SequenceNumber:     "111111",
			Amount:             100000,
			SourceAccount:      "001001234567891",
			DestinationAccount: "001001234567892",
			DestinationName:    "Destination Account",
			SourceName:         "Olivia Rodrigo",
		}, nil)

	transaction, err := svc.DoPayment(ctx, "123456")

	assert.Nil(t, transaction)
	assert.Equal(t, pkgerror.New(codes.BadRequest, ErrInvalidSequenceNumber), err)

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}

func TestTransferDoPaymentFailed_GetTransactionLimitFailed(t *testing.T) {
	var (
		corebankingMock = NewMockCoreBanking(t)
		repoMock        = NewMockRepository(t)
		mailerMock      = NewMockReceiptMailer(t)
		seqGenMock      = NewMockSequenceGenerator(t)
		notifierMock    = NewMockNotifier(t)
		svc             = NewService(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock, notifierMock)
		ctx             = ctxt.ContextWithUser(context.Background(), &ctxt.User{
			ID:    123,
			CIF:   "1234567",
			Name:  "Olivia Rodrigo",
			Email: "olivia@gmail.com",
		})
	)

	corebankingMock.EXPECT().GetCoreStatus(mock.Anything).Return(&CoreStatus{
		SystemDate:    "25-03-2025",
		Status:        "FINISHED",
		StandInStatus: "N",
	}, nil)

	repoMock.EXPECT().GetSequence(mock.Anything, "123456").
		Return(&Sequence{
			SequenceNumber:     "123456",
			Amount:             100000,
			SourceAccount:      "001001234567891",
			DestinationAccount: "001001234567892",
			DestinationName:    "Destination Account",
			SourceName:         "Olivia Rodrigo",
		}, nil)
	repoMock.EXPECT().GetTransactionLimit(mock.Anything).
		Return(nil, errors.New("some error"))

	transaction, err := svc.DoPayment(ctx, "123456")

	assert.Nil(t, transaction)
	assert.Equal(t, pkgerror.New(codes.Internal, ErrGeneral), err)

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}

func TestTransferDoPaymentFailed_TransactionLimitCannotTransfer(t *testing.T) {
	var (
		corebankingMock = NewMockCoreBanking(t)
		repoMock        = NewMockRepository(t)
		mailerMock      = NewMockReceiptMailer(t)
		seqGenMock      = NewMockSequenceGenerator(t)
		notifierMock    = NewMockNotifier(t)
		svc             = NewService(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock, notifierMock)
		ctx             = ctxt.ContextWithUser(context.Background(), &ctxt.User{
			ID:    123,
			CIF:   "1234567",
			Name:  "Olivia Rodrigo",
			Email: "olivia@gmail.com",
		})
	)

	corebankingMock.EXPECT().GetCoreStatus(mock.Anything).Return(&CoreStatus{
		SystemDate:    "25-03-2025",
		Status:        "FINISHED",
		StandInStatus: "N",
	}, nil)

	repoMock.EXPECT().GetSequence(mock.Anything, "123456").
		Return(&Sequence{
			SequenceNumber:     "123456",
			Amount:             100000,
			SourceAccount:      "001001234567891",
			DestinationAccount: "001001234567892",
			DestinationName:    "Destination Account",
			SourceName:         "Olivia Rodrigo",
		}, nil)
	repoMock.EXPECT().GetTransactionLimit(mock.Anything).
		Return(&Limits{
			MinAmount:      1,
			MaxAmount:      100_000_000,
			MaxDailyAmount: 50_000,
		}, nil)

	transaction, err := svc.DoPayment(ctx, "123456")

	assert.Nil(t, transaction)
	assert.Equal(t, pkgerror.New(codes.BadRequest, ErrInvalidAmount), err)

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}

func TestTransferDoPaymentFailed_PerformOverbookingFailed(t *testing.T) {
	var (
		corebankingMock = NewMockCoreBanking(t)
		repoMock        = NewMockRepository(t)
		mailerMock      = NewMockReceiptMailer(t)
		seqGenMock      = NewMockSequenceGenerator(t)
		notifierMock    = NewMockNotifier(t)
		svc             = NewService(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock, notifierMock)
		ctx             = ctxt.ContextWithUser(context.Background(), &ctxt.User{
			ID:    123,
			CIF:   "1234567",
			Name:  "Olivia Rodrigo",
			Email: "olivia@gmail.com",
		})
	)

	corebankingMock.EXPECT().GetCoreStatus(mock.Anything).Return(&CoreStatus{
		SystemDate:    "25-03-2025",
		Status:        "FINISHED",
		StandInStatus: "N",
	}, nil)

	repoMock.EXPECT().GetTransactionLimit(mock.Anything).
		Return(&Limits{
			MinAmount:      1,
			MaxAmount:      50_000_000,
			MaxDailyAmount: 200_000_000,
		}, nil)
	repoMock.EXPECT().GetSequence(mock.Anything, "123456").
		Return(&Sequence{
			SequenceNumber:     "123456",
			Amount:             100000,
			SourceAccount:      "001001234567891",
			DestinationAccount: "001001234567892",
			DestinationName:    "Destination Account",
			SourceName:         "Olivia Rodrigo",
		}, nil)

	corebankingMock.EXPECT().PerformOverbooking(mock.Anything, &OverbookingInput{
		SourceAccount:      "001001234567891",
		DestinationAccount: "001001234567892",
		Amount:             100000,
		Fee:                0,
		Remark:             "TRF 001001234567891 001001234567892 BNKYAYA 123456",
	}).Return(nil, errors.New("some error"))

	transaction, err := svc.DoPayment(ctx, "123456")

	assert.Nil(t, transaction)
	assert.Equal(t, pkgerror.New(codes.Internal, ErrGeneral), err)

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}

func TestTransferDoPaymentFailed_GetUserFromContextFailed(t *testing.T) {
	var (
		corebankingMock = NewMockCoreBanking(t)
		repoMock        = NewMockRepository(t)
		mailerMock      = NewMockReceiptMailer(t)
		seqGenMock      = NewMockSequenceGenerator(t)
		notifierMock    = NewMockNotifier(t)
		svc             = NewService(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock, notifierMock)
		ctx             = context.Background()
	)

	corebankingMock.EXPECT().GetCoreStatus(mock.Anything).Return(&CoreStatus{
		SystemDate:    "25-03-2025",
		Status:        "FINISHED",
		StandInStatus: "N",
	}, nil)

	repoMock.EXPECT().GetTransactionLimit(mock.Anything).
		Return(&Limits{
			MinAmount:      1,
			MaxAmount:      50_000_000,
			MaxDailyAmount: 200_000_000,
		}, nil)
	repoMock.EXPECT().GetSequence(mock.Anything, "123456").
		Return(&Sequence{
			SequenceNumber:     "123456",
			Amount:             100000,
			SourceAccount:      "001001234567891",
			DestinationAccount: "001001234567892",
			DestinationName:    "Destination Account",
			SourceName:         "Olivia Rodrigo",
		}, nil)

	corebankingMock.EXPECT().PerformOverbooking(mock.Anything, &OverbookingInput{
		SourceAccount:      "001001234567891",
		DestinationAccount: "001001234567892",
		Amount:             100000,
		Fee:                0,
		Remark:             "TRF 001001234567891 001001234567892 BNKYAYA 123456",
	}).Return(&OverbookingResult{
		JournalSequence:      "111111",
		TransactionReference: "222222",
	}, nil)

	transaction, err := svc.DoPayment(ctx, "123456")

	assert.Nil(t, transaction)
	assert.Equal(t, pkgerror.New(codes.Unauthenticated, ErrUnauthenticatedUser), err)

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}

func TestTransferDoPaymentFailed_InsertTransactionFailed(t *testing.T) {
	var (
		corebankingMock = NewMockCoreBanking(t)
		repoMock        = NewMockRepository(t)
		mailerMock      = NewMockReceiptMailer(t)
		seqGenMock      = NewMockSequenceGenerator(t)
		notifierMock    = NewMockNotifier(t)
		svc             = NewService(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock, notifierMock)
		ctx             = ctxt.ContextWithUser(context.Background(), &ctxt.User{
			ID:    123,
			CIF:   "1234567",
			Name:  "Olivia Rodrigo",
			Email: "olivia@gmail.com",
		})
	)

	corebankingMock.EXPECT().GetCoreStatus(mock.Anything).Return(&CoreStatus{
		SystemDate:    "25-03-2025",
		Status:        "FINISHED",
		StandInStatus: "N",
	}, nil)

	repoMock.EXPECT().GetTransactionLimit(mock.Anything).
		Return(&Limits{
			MinAmount:      1,
			MaxAmount:      50_000_000,
			MaxDailyAmount: 200_000_000,
		}, nil)
	repoMock.EXPECT().GetSequence(mock.Anything, "123456").
		Return(&Sequence{
			SequenceNumber:     "123456",
			Amount:             100000,
			SourceAccount:      "001001234567891",
			DestinationAccount: "001001234567892",
			DestinationName:    "Destination Account",
			SourceName:         "Olivia Rodrigo",
		}, nil)

	corebankingMock.EXPECT().PerformOverbooking(mock.Anything, &OverbookingInput{
		SourceAccount:      "001001234567891",
		DestinationAccount: "001001234567892",
		Amount:             100000,
		Fee:                0,
		Remark:             "TRF 001001234567891 001001234567892 BNKYAYA 123456",
	}).Return(&OverbookingResult{
		JournalSequence:      "111111",
		TransactionReference: "222222",
	}, nil)

	repoMock.EXPECT().InsertTransaction(mock.Anything, &Transaction{
		SequenceNumber:       "123456",
		SequenceJournal:      "111111",
		UserID:               "123",
		Destination:          "001001234567892",
		Amount:               100000,
		TransactionType:      "internal_transfer",
		TransactionReference: "222222",
		Remarks:              "TRF 001001234567891 001001234567892 BNKYAYA 123456",
		Fee:                  "0",
		DestinationName:      "Destination Account",
	}).Return(errors.New("some error"))

	transaction, err := svc.DoPayment(ctx, "123456")

	assert.Nil(t, transaction)
	assert.Equal(t, pkgerror.New(codes.Internal, ErrGeneral), err)

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}

func TestTransferDoPaymentFailed_SendEmailFailed(t *testing.T) {
	var (
		corebankingMock = NewMockCoreBanking(t)
		repoMock        = NewMockRepository(t)
		mailerMock      = NewMockReceiptMailer(t)
		seqGenMock      = NewMockSequenceGenerator(t)
		notifierMock    = NewMockNotifier(t)
		svc             = NewService(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock, notifierMock)
		ctx             = ctxt.ContextWithUser(context.Background(), &ctxt.User{
			ID:    123,
			CIF:   "1234567",
			Name:  "Olivia Rodrigo",
			Email: "olivia@gmail.com",
		})
	)

	corebankingMock.EXPECT().GetCoreStatus(mock.Anything).Return(&CoreStatus{
		SystemDate:    "25-03-2025",
		Status:        "FINISHED",
		StandInStatus: "N",
	}, nil)

	repoMock.EXPECT().GetTransactionLimit(mock.Anything).
		Return(&Limits{
			MinAmount:      1,
			MaxAmount:      50_000_000,
			MaxDailyAmount: 200_000_000,
		}, nil)
	repoMock.EXPECT().GetSequence(mock.Anything, "123456").
		Return(&Sequence{
			SequenceNumber:     "123456",
			Amount:             100000,
			SourceAccount:      "001001234567891",
			DestinationAccount: "001001234567892",
			DestinationName:    "Destination Account",
			SourceName:         "Olivia Rodrigo",
		}, nil)

	corebankingMock.EXPECT().PerformOverbooking(mock.Anything, &OverbookingInput{
		SourceAccount:      "001001234567891",
		DestinationAccount: "001001234567892",
		Amount:             100000,
		Fee:                0,
		Remark:             "TRF 001001234567891 001001234567892 BNKYAYA 123456",
	}).Return(&OverbookingResult{
		JournalSequence:      "111111",
		TransactionReference: "222222",
	}, nil)

	repoMock.EXPECT().InsertTransaction(mock.Anything, &Transaction{
		SequenceNumber:       "123456",
		SequenceJournal:      "111111",
		UserID:               "123",
		Destination:          "001001234567892",
		Amount:               100000,
		TransactionType:      "internal_transfer",
		TransactionReference: "222222",
		Remarks:              "TRF 001001234567891 001001234567892 BNKYAYA 123456",
		Fee:                  "0",
		DestinationName:      "Destination Account",
	}).Return(nil)

	mailerMock.EXPECT().SendReceipt(mock.Anything, mock.Anything).
		Return(errors.New("some error"))

	transaction, err := svc.DoPayment(ctx, "123456")

	assert.Nil(t, transaction)
	assert.Equal(t, pkgerror.New(codes.Internal, ErrSendEmailFailed), err)

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}

func TestTransferDoPaymentFailed_NotifyFailed(t *testing.T) {
	var (
		corebankingMock = NewMockCoreBanking(t)
		repoMock        = NewMockRepository(t)
		mailerMock      = NewMockReceiptMailer(t)
		seqGenMock      = NewMockSequenceGenerator(t)
		notifierMock    = NewMockNotifier(t)
		svc             = NewService(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock, notifierMock)
		ctx             = ctxt.ContextWithUser(context.Background(), &ctxt.User{
			ID:    123,
			CIF:   "1234567",
			Name:  "Olivia Rodrigo",
			Email: "olivia@gmail.com",
		})
	)

	corebankingMock.EXPECT().GetCoreStatus(mock.Anything).Return(&CoreStatus{
		SystemDate:    "25-03-2025",
		Status:        "FINISHED",
		StandInStatus: "N",
	}, nil)

	repoMock.EXPECT().GetTransactionLimit(mock.Anything).
		Return(&Limits{
			MinAmount:      1,
			MaxAmount:      50_000_000,
			MaxDailyAmount: 200_000_000,
		}, nil)
	repoMock.EXPECT().GetSequence(mock.Anything, "123456").
		Return(&Sequence{
			SequenceNumber:     "123456",
			Amount:             100000,
			SourceAccount:      "001001234567891",
			DestinationAccount: "001001234567892",
			DestinationName:    "Destination Account",
			SourceName:         "Olivia Rodrigo",
		}, nil)

	corebankingMock.EXPECT().PerformOverbooking(mock.Anything, &OverbookingInput{
		SourceAccount:      "001001234567891",
		DestinationAccount: "001001234567892",
		Amount:             100000,
		Fee:                0,
		Remark:             "TRF 001001234567891 001001234567892 BNKYAYA 123456",
	}).Return(&OverbookingResult{
		JournalSequence:      "111111",
		TransactionReference: "222222",
	}, nil)

	repoMock.EXPECT().InsertTransaction(mock.Anything, &Transaction{
		SequenceNumber:       "123456",
		SequenceJournal:      "111111",
		UserID:               "123",
		Destination:          "001001234567892",
		Amount:               100000,
		TransactionType:      "internal_transfer",
		TransactionReference: "222222",
		Remarks:              "TRF 001001234567891 001001234567892 BNKYAYA 123456",
		Fee:                  "0",
		DestinationName:      "Destination Account",
	}).Return(nil)

	mailerMock.EXPECT().SendReceipt(mock.Anything, mock.Anything).
		Return(nil)

	notifierMock.EXPECT().Notify(mock.Anything, mock.Anything).
		Return(errors.New("some error"))

	transaction, err := svc.DoPayment(ctx, "123456")

	assert.Nil(t, transaction)
	assert.Equal(t, pkgerror.New(codes.Internal, ErrNotifyFailed), err)

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}
