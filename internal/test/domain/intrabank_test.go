package domain

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.bankyaya.org/app/backend/internal/domain/intrabank"
	"go.bankyaya.org/app/backend/internal/pkg/codes"
	"go.bankyaya.org/app/backend/internal/pkg/ctxt"
	"go.bankyaya.org/app/backend/internal/pkg/logger"
	"go.bankyaya.org/app/backend/internal/pkg/status"
	intrabankmock "go.bankyaya.org/app/backend/internal/test/domain/mocks/intrabank"
)

func TestTransferInquirySuccess(t *testing.T) {
	var (
		corebankingMock = intrabankmock.NewCoreBanking(t)
		repoMock        = intrabankmock.NewRepository(t)
		mailerMock      = intrabankmock.NewReceiptMailer(t)
		seqGenMock      = intrabankmock.NewSequenceGenerator(t)
		svc             = intrabank.NewService(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock)
		ctx             = ctxt.ContextWithUser(context.Background(), &ctxt.User{
			Id:       123,
			CIF:      "1234567",
			FullName: "Olivia Rodrigo",
			Email:    "olivia@gmail.com",
		})
	)

	corebankingMock.EXPECT().GetCoreStatus(mock.Anything).Return(&intrabank.CoreStatus{
		SystemDate:    "25-03-2025",
		Status:        "FINISHED",
		StandInStatus: "N",
	}, nil)
	corebankingMock.EXPECT().GetAccountDetails(mock.Anything, "001001234567891").
		Return(&intrabank.Account{
			Name:   "Olivia Rodrigo",
			Status: "1",
		}, nil)
	corebankingMock.EXPECT().GetAccountDetails(mock.Anything, "001001234567892").
		Return(&intrabank.Account{
			Name:   "Destination Account",
			Status: "1",
		}, nil)

	repoMock.EXPECT().InsertSequence(mock.Anything, &intrabank.Sequence{
		SequenceNumber:     "123456",
		Amount:             100000,
		SourceAccount:      "001001234567891",
		DestinationAccount: "001001234567892",
		DestinationName:    "Destination Account",
		SourceName:         "Olivia Rodrigo",
	}).Return(nil)

	seqGenMock.EXPECT().Generate().
		Return("123456", nil)

	sequence, err := svc.Inquiry(ctx, &intrabank.Sequence{
		SequenceNumber:     "123456",
		Amount:             100000,
		SourceAccount:      "001001234567891",
		DestinationAccount: "001001234567892",
	})

	assert.Nil(t, err)
	assert.Equal(t, sequence, &intrabank.Sequence{
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
		corebankingMock = intrabankmock.NewCoreBanking(t)
		repoMock        = intrabankmock.NewRepository(t)
		mailerMock      = intrabankmock.NewReceiptMailer(t)
		seqGenMock      = intrabankmock.NewSequenceGenerator(t)
		svc             = intrabank.NewService(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock)
		ctx             = ctxt.ContextWithUser(context.Background(), &ctxt.User{
			Id:       123,
			CIF:      "1234567",
			FullName: "Olivia Rodrigo",
			Email:    "olivia@gmail.com",
		})
	)

	corebankingMock.EXPECT().GetCoreStatus(mock.Anything).
		Return(nil, errors.New("check EOD failed"))

	sequence, err := svc.Inquiry(ctx, &intrabank.Sequence{
		SequenceNumber:     "123456",
		Amount:             100000,
		SourceAccount:      "001001234567891",
		DestinationAccount: "001001234567892",
	})
	assert.Nil(t, sequence)
	assert.Equal(t, status.Error(codes.Internal, intrabank.ErrGeneral), err)

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}

func TestTransferInquiryFailed_EODIsRunning(t *testing.T) {
	var (
		corebankingMock = intrabankmock.NewCoreBanking(t)
		repoMock        = intrabankmock.NewRepository(t)
		mailerMock      = intrabankmock.NewReceiptMailer(t)
		seqGenMock      = intrabankmock.NewSequenceGenerator(t)
		svc             = intrabank.NewService(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock)
		ctx             = ctxt.ContextWithUser(context.Background(), &ctxt.User{
			Id:       123,
			CIF:      "1234567",
			FullName: "Olivia Rodrigo",
			Email:    "olivia@gmail.com",
		})
	)

	corebankingMock.EXPECT().GetCoreStatus(mock.Anything).Return(&intrabank.CoreStatus{
		SystemDate:    "25-03-2025",
		Status:        "STARTED",
		StandInStatus: "N",
	}, nil)

	sequence, err := svc.Inquiry(ctx, &intrabank.Sequence{
		SequenceNumber:     "123456",
		Amount:             100000,
		SourceAccount:      "001001234567891",
		DestinationAccount: "001001234567892",
	})
	assert.Nil(t, sequence)
	assert.Equal(t, status.Error(codes.Internal, intrabank.ErrEODInProgress), err)

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}

func TestTransferInquiryFailed_FailedCheckSourceAccount(t *testing.T) {
	var (
		corebankingMock = intrabankmock.NewCoreBanking(t)
		repoMock        = intrabankmock.NewRepository(t)
		mailerMock      = intrabankmock.NewReceiptMailer(t)
		seqGenMock      = intrabankmock.NewSequenceGenerator(t)
		svc             = intrabank.NewService(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock)
		ctx             = ctxt.ContextWithUser(context.Background(), &ctxt.User{
			Id:       123,
			CIF:      "1234567",
			FullName: "Olivia Rodrigo",
			Email:    "olivia@gmail.com",
		})
	)

	corebankingMock.EXPECT().GetCoreStatus(mock.Anything).Return(&intrabank.CoreStatus{
		SystemDate:    "25-03-2025",
		Status:        "FINISHED",
		StandInStatus: "N",
	}, nil)
	corebankingMock.EXPECT().GetAccountDetails(mock.Anything, "001001234567891").
		Return(nil, errors.New("GetAccountDetails failed"))

	sequence, err := svc.Inquiry(ctx, &intrabank.Sequence{
		SequenceNumber:     "123456",
		Amount:             100000,
		SourceAccount:      "001001234567891",
		DestinationAccount: "001001234567892",
	})

	assert.Nil(t, sequence)
	assert.Equal(t, status.Error(codes.Internal, intrabank.ErrGeneral), err)

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}

func TestTransferInquiryFailed_SourceAccountIsInactive(t *testing.T) {
	var (
		corebankingMock = intrabankmock.NewCoreBanking(t)
		repoMock        = intrabankmock.NewRepository(t)
		mailerMock      = intrabankmock.NewReceiptMailer(t)
		seqGenMock      = intrabankmock.NewSequenceGenerator(t)
		svc             = intrabank.NewService(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock)
		ctx             = ctxt.ContextWithUser(context.Background(), &ctxt.User{
			Id:       123,
			CIF:      "1234567",
			FullName: "Olivia Rodrigo",
			Email:    "olivia@gmail.com",
		})
	)

	corebankingMock.EXPECT().GetCoreStatus(mock.Anything).Return(&intrabank.CoreStatus{
		SystemDate:    "25-03-2025",
		Status:        "FINISHED",
		StandInStatus: "N",
	}, nil)
	corebankingMock.EXPECT().GetAccountDetails(mock.Anything, "001001234567891").
		Return(&intrabank.Account{
			Status: "9",
		}, nil)

	sequence, err := svc.Inquiry(ctx, &intrabank.Sequence{
		SequenceNumber:     "123456",
		Amount:             100000,
		SourceAccount:      "001001234567891",
		DestinationAccount: "001001234567892",
	})

	assert.Nil(t, sequence)
	assert.Equal(t, status.Error(codes.BadRequest, intrabank.ErrSourceAccountInactive), err)

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}

func TestTransferInquiryFailed_CheckDestinationAccountFailed(t *testing.T) {
	var (
		corebankingMock = intrabankmock.NewCoreBanking(t)
		repoMock        = intrabankmock.NewRepository(t)
		mailerMock      = intrabankmock.NewReceiptMailer(t)
		seqGenMock      = intrabankmock.NewSequenceGenerator(t)
		svc             = intrabank.NewService(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock)
		ctx             = ctxt.ContextWithUser(context.Background(), &ctxt.User{
			Id:       123,
			CIF:      "1234567",
			FullName: "Olivia Rodrigo",
			Email:    "olivia@gmail.com",
		})
	)

	corebankingMock.EXPECT().GetCoreStatus(mock.Anything).Return(&intrabank.CoreStatus{
		SystemDate:    "25-03-2025",
		Status:        "FINISHED",
		StandInStatus: "N",
	}, nil)
	corebankingMock.EXPECT().GetAccountDetails(mock.Anything, "001001234567891").
		Return(&intrabank.Account{
			Name:   "Olivia Rodrigo",
			Status: "1",
		}, nil)
	corebankingMock.EXPECT().GetAccountDetails(mock.Anything, "001001234567892").
		Return(nil, errors.New("GetAccountDetails failed"))

	sequence, err := svc.Inquiry(ctx, &intrabank.Sequence{
		SequenceNumber:     "123456",
		Amount:             100000,
		SourceAccount:      "001001234567891",
		DestinationAccount: "001001234567892",
	})

	assert.Nil(t, sequence)
	assert.Equal(t, status.Error(codes.Internal, intrabank.ErrGeneral), err)

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}

func TestTransferInquiryFailed_DestinationAccountInactive(t *testing.T) {
	var (
		corebankingMock = intrabankmock.NewCoreBanking(t)
		repoMock        = intrabankmock.NewRepository(t)
		mailerMock      = intrabankmock.NewReceiptMailer(t)
		seqGenMock      = intrabankmock.NewSequenceGenerator(t)
		svc             = intrabank.NewService(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock)
		ctx             = ctxt.ContextWithUser(context.Background(), &ctxt.User{
			Id:       123,
			CIF:      "1234567",
			FullName: "Olivia Rodrigo",
			Email:    "olivia@gmail.com",
		})
	)

	corebankingMock.EXPECT().GetCoreStatus(mock.Anything).Return(&intrabank.CoreStatus{
		SystemDate:    "25-03-2025",
		Status:        "FINISHED",
		StandInStatus: "N",
	}, nil)
	corebankingMock.EXPECT().GetAccountDetails(mock.Anything, "001001234567891").
		Return(&intrabank.Account{
			Name:   "Olivia Rodrigo",
			Status: "1",
		}, nil)
	corebankingMock.EXPECT().GetAccountDetails(mock.Anything, "001001234567892").
		Return(&intrabank.Account{
			Status: "9",
		}, nil)

	sequence, err := svc.Inquiry(ctx, &intrabank.Sequence{
		SequenceNumber:     "123456",
		Amount:             100000,
		SourceAccount:      "001001234567891",
		DestinationAccount: "001001234567892",
	})

	assert.Nil(t, sequence)
	assert.Equal(t, status.Error(codes.BadRequest, intrabank.ErrDestinationAccountInactive), err)

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}

func TestTransferInquiryFailed_GenerateSequenceFailed(t *testing.T) {
	var (
		corebankingMock = intrabankmock.NewCoreBanking(t)
		repoMock        = intrabankmock.NewRepository(t)
		mailerMock      = intrabankmock.NewReceiptMailer(t)
		seqGenMock      = intrabankmock.NewSequenceGenerator(t)
		svc             = intrabank.NewService(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock)
		ctx             = ctxt.ContextWithUser(context.Background(), &ctxt.User{
			Id:       123,
			CIF:      "1234567",
			FullName: "Olivia Rodrigo",
			Email:    "olivia@gmail.com",
		})
	)

	corebankingMock.EXPECT().GetCoreStatus(mock.Anything).Return(&intrabank.CoreStatus{
		SystemDate:    "25-03-2025",
		Status:        "FINISHED",
		StandInStatus: "N",
	}, nil)
	corebankingMock.EXPECT().GetAccountDetails(mock.Anything, "001001234567891").
		Return(&intrabank.Account{
			Name:   "Olivia Rodrigo",
			Status: "1",
		}, nil)
	corebankingMock.EXPECT().GetAccountDetails(mock.Anything, "001001234567892").
		Return(&intrabank.Account{
			Name:   "Destination Account",
			Status: "1",
		}, nil)

	seqGenMock.EXPECT().Generate().
		Return("", errors.New("some error"))

	sequence, err := svc.Inquiry(ctx, &intrabank.Sequence{
		SequenceNumber:     "123456",
		Amount:             100000,
		SourceAccount:      "001001234567891",
		DestinationAccount: "001001234567892",
	})

	assert.Nil(t, sequence)
	assert.Equal(t, status.Error(codes.Internal, intrabank.ErrGeneral), err)

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}

func TestTransferInquiryFailed_InsertSequenceFailed(t *testing.T) {
	var (
		corebankingMock = intrabankmock.NewCoreBanking(t)
		repoMock        = intrabankmock.NewRepository(t)
		mailerMock      = intrabankmock.NewReceiptMailer(t)
		seqGenMock      = intrabankmock.NewSequenceGenerator(t)
		svc             = intrabank.NewService(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock)
		ctx             = ctxt.ContextWithUser(context.Background(), &ctxt.User{
			Id:       123,
			CIF:      "1234567",
			FullName: "Olivia Rodrigo",
			Email:    "olivia@gmail.com",
		})
	)

	corebankingMock.EXPECT().GetCoreStatus(mock.Anything).Return(&intrabank.CoreStatus{
		SystemDate:    "25-03-2025",
		Status:        "FINISHED",
		StandInStatus: "N",
	}, nil)
	corebankingMock.EXPECT().GetAccountDetails(mock.Anything, "001001234567891").
		Return(&intrabank.Account{
			Name:   "Olivia Rodrigo",
			Status: "1",
		}, nil)
	corebankingMock.EXPECT().GetAccountDetails(mock.Anything, "001001234567892").
		Return(&intrabank.Account{
			Name:   "Destination Account",
			Status: "1",
		}, nil)

	seqGenMock.EXPECT().Generate().
		Return("123456", nil)

	repoMock.EXPECT().InsertSequence(mock.Anything, &intrabank.Sequence{
		SequenceNumber:     "123456",
		Amount:             100000,
		SourceAccount:      "001001234567891",
		DestinationAccount: "001001234567892",
		DestinationName:    "Destination Account",
		SourceName:         "Olivia Rodrigo",
	}).Return(errors.New("some error"))

	sequence, err := svc.Inquiry(ctx, &intrabank.Sequence{
		SequenceNumber:     "123456",
		Amount:             100000,
		SourceAccount:      "001001234567891",
		DestinationAccount: "001001234567892",
	})

	assert.Nil(t, sequence)
	assert.Equal(t, status.Error(codes.Internal, intrabank.ErrGeneral), err)

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}

func TestTransferDoPaymentSuccess(t *testing.T) {
	var (
		corebankingMock = intrabankmock.NewCoreBanking(t)
		repoMock        = intrabankmock.NewRepository(t)
		mailerMock      = intrabankmock.NewReceiptMailer(t)
		seqGenMock      = intrabankmock.NewSequenceGenerator(t)
		svc             = intrabank.NewService(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock)
		ctx             = ctxt.ContextWithUser(context.Background(), &ctxt.User{
			Id:       123,
			CIF:      "1234567",
			FullName: "Olivia Rodrigo",
			Email:    "olivia@gmail.com",
		})
	)

	corebankingMock.EXPECT().GetCoreStatus(mock.Anything).Return(&intrabank.CoreStatus{
		SystemDate:    "25-03-2025",
		Status:        "FINISHED",
		StandInStatus: "N",
	}, nil)

	repoMock.EXPECT().GetSequence(mock.Anything, "123456").
		Return(&intrabank.Sequence{
			SequenceNumber:     "123456",
			Amount:             100000,
			SourceAccount:      "001001234567891",
			DestinationAccount: "001001234567892",
			DestinationName:    "Destination Account",
			SourceName:         "Olivia Rodrigo",
		}, nil)

	corebankingMock.EXPECT().PerformOverbooking(mock.Anything, &intrabank.OverbookingInput{
		SourceAccount:      "001001234567891",
		DestinationAccount: "001001234567892",
		Amount:             100000,
		Fee:                0,
		Remark:             "TRF 001001234567891 001001234567892 BNKYAYA 123456",
	}).Return(&intrabank.OverbookingResult{
		JournalSequence:      "111111",
		TransactionReference: "222222",
	}, nil)

	repoMock.EXPECT().InsertTransaction(mock.Anything, &intrabank.Transaction{
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

	mailerMock.EXPECT().SendReceipt(mock.Anything, mock.Anything).Return(nil)

	transaction, err := svc.DoPayment(ctx, "123456")

	assert.NoError(t, err)
	assert.Equal(t, &intrabank.Transaction{
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
		corebankingMock = intrabankmock.NewCoreBanking(t)
		repoMock        = intrabankmock.NewRepository(t)
		mailerMock      = intrabankmock.NewReceiptMailer(t)
		seqGenMock      = intrabankmock.NewSequenceGenerator(t)
		svc             = intrabank.NewService(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock)
		ctx             = ctxt.ContextWithUser(context.Background(), &ctxt.User{
			Id:       123,
			CIF:      "1234567",
			FullName: "Olivia Rodrigo",
			Email:    "olivia@gmail.com",
		})
	)

	corebankingMock.EXPECT().GetCoreStatus(mock.Anything).
		Return(nil, errors.New("some error"))

	transaction, err := svc.DoPayment(ctx, "123456")

	assert.Nil(t, transaction)
	assert.Equal(t, status.Error(codes.Internal, intrabank.ErrGeneral), err)

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}

func TestTransferDoPaymentFailed_EODIsRunning(t *testing.T) {
	var (
		corebankingMock = intrabankmock.NewCoreBanking(t)
		repoMock        = intrabankmock.NewRepository(t)
		mailerMock      = intrabankmock.NewReceiptMailer(t)
		seqGenMock      = intrabankmock.NewSequenceGenerator(t)
		svc             = intrabank.NewService(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock)
		ctx             = ctxt.ContextWithUser(context.Background(), &ctxt.User{
			Id:       123,
			CIF:      "1234567",
			FullName: "Olivia Rodrigo",
			Email:    "olivia@gmail.com",
		})
	)

	corebankingMock.EXPECT().GetCoreStatus(mock.Anything).
		Return(&intrabank.CoreStatus{
			SystemDate:    "25-03-2025",
			Status:        "STARTED",
			StandInStatus: "N",
		}, nil)

	transaction, err := svc.DoPayment(ctx, "123456")

	assert.Nil(t, transaction)
	assert.Equal(t, status.Error(codes.Internal, intrabank.ErrEODInProgress), err)

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}

func TestTransferDoPaymentFailed_GetSequenceFailed(t *testing.T) {
	var (
		corebankingMock = intrabankmock.NewCoreBanking(t)
		repoMock        = intrabankmock.NewRepository(t)
		mailerMock      = intrabankmock.NewReceiptMailer(t)
		seqGenMock      = intrabankmock.NewSequenceGenerator(t)
		svc             = intrabank.NewService(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock)
		ctx             = ctxt.ContextWithUser(context.Background(), &ctxt.User{
			Id:       123,
			CIF:      "1234567",
			FullName: "Olivia Rodrigo",
			Email:    "olivia@gmail.com",
		})
	)

	corebankingMock.EXPECT().GetCoreStatus(mock.Anything).Return(&intrabank.CoreStatus{
		SystemDate:    "25-03-2025",
		Status:        "FINISHED",
		StandInStatus: "N",
	}, nil)

	repoMock.EXPECT().GetSequence(mock.Anything, "123456").
		Return(nil, errors.New("some error"))

	transaction, err := svc.DoPayment(ctx, "123456")

	assert.Nil(t, transaction)
	assert.Equal(t, status.Error(codes.Internal, intrabank.ErrGeneral), err)

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}

func TestTransferDoPaymentFailed_InvalidSequenceNumber(t *testing.T) {
	var (
		corebankingMock = intrabankmock.NewCoreBanking(t)
		repoMock        = intrabankmock.NewRepository(t)
		mailerMock      = intrabankmock.NewReceiptMailer(t)
		seqGenMock      = intrabankmock.NewSequenceGenerator(t)
		svc             = intrabank.NewService(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock)
		ctx             = ctxt.ContextWithUser(context.Background(), &ctxt.User{
			Id:       123,
			CIF:      "1234567",
			FullName: "Olivia Rodrigo",
			Email:    "olivia@gmail.com",
		})
	)

	corebankingMock.EXPECT().GetCoreStatus(mock.Anything).Return(&intrabank.CoreStatus{
		SystemDate:    "25-03-2025",
		Status:        "FINISHED",
		StandInStatus: "N",
	}, nil)

	repoMock.EXPECT().GetSequence(mock.Anything, "123456").
		Return(&intrabank.Sequence{
			SequenceNumber:     "111111",
			Amount:             100000,
			SourceAccount:      "001001234567891",
			DestinationAccount: "001001234567892",
			DestinationName:    "Destination Account",
			SourceName:         "Olivia Rodrigo",
		}, nil)

	transaction, err := svc.DoPayment(ctx, "123456")

	assert.Nil(t, transaction)
	assert.Equal(t, status.Error(codes.BadRequest, intrabank.ErrInvalidSequenceNumber), err)

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}

func TestTransferDoPaymentFailed_PerformOverbookingFailed(t *testing.T) {
	var (
		corebankingMock = intrabankmock.NewCoreBanking(t)
		repoMock        = intrabankmock.NewRepository(t)
		mailerMock      = intrabankmock.NewReceiptMailer(t)
		seqGenMock      = intrabankmock.NewSequenceGenerator(t)
		svc             = intrabank.NewService(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock)
		ctx             = ctxt.ContextWithUser(context.Background(), &ctxt.User{
			Id:       123,
			CIF:      "1234567",
			FullName: "Olivia Rodrigo",
			Email:    "olivia@gmail.com",
		})
	)

	corebankingMock.EXPECT().GetCoreStatus(mock.Anything).Return(&intrabank.CoreStatus{
		SystemDate:    "25-03-2025",
		Status:        "FINISHED",
		StandInStatus: "N",
	}, nil)

	repoMock.EXPECT().GetSequence(mock.Anything, "123456").
		Return(&intrabank.Sequence{
			SequenceNumber:     "123456",
			Amount:             100000,
			SourceAccount:      "001001234567891",
			DestinationAccount: "001001234567892",
			DestinationName:    "Destination Account",
			SourceName:         "Olivia Rodrigo",
		}, nil)

	corebankingMock.EXPECT().PerformOverbooking(mock.Anything, &intrabank.OverbookingInput{
		SourceAccount:      "001001234567891",
		DestinationAccount: "001001234567892",
		Amount:             100000,
		Fee:                0,
		Remark:             "TRF 001001234567891 001001234567892 BNKYAYA 123456",
	}).Return(nil, errors.New("some error"))

	transaction, err := svc.DoPayment(ctx, "123456")

	assert.Nil(t, transaction)
	assert.Equal(t, status.Error(codes.Internal, intrabank.ErrGeneral), err)

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}

func TestTransferDoPaymentFailed_GetUserFromContextFailed(t *testing.T) {
	var (
		corebankingMock = intrabankmock.NewCoreBanking(t)
		repoMock        = intrabankmock.NewRepository(t)
		mailerMock      = intrabankmock.NewReceiptMailer(t)
		seqGenMock      = intrabankmock.NewSequenceGenerator(t)
		svc             = intrabank.NewService(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock)
		ctx             = context.Background()
	)

	corebankingMock.EXPECT().GetCoreStatus(mock.Anything).Return(&intrabank.CoreStatus{
		SystemDate:    "25-03-2025",
		Status:        "FINISHED",
		StandInStatus: "N",
	}, nil)

	repoMock.EXPECT().GetSequence(mock.Anything, "123456").
		Return(&intrabank.Sequence{
			SequenceNumber:     "123456",
			Amount:             100000,
			SourceAccount:      "001001234567891",
			DestinationAccount: "001001234567892",
			DestinationName:    "Destination Account",
			SourceName:         "Olivia Rodrigo",
		}, nil)

	corebankingMock.EXPECT().PerformOverbooking(mock.Anything, &intrabank.OverbookingInput{
		SourceAccount:      "001001234567891",
		DestinationAccount: "001001234567892",
		Amount:             100000,
		Fee:                0,
		Remark:             "TRF 001001234567891 001001234567892 BNKYAYA 123456",
	}).Return(&intrabank.OverbookingResult{
		JournalSequence:      "111111",
		TransactionReference: "222222",
	}, nil)

	transaction, err := svc.DoPayment(ctx, "123456")

	assert.Nil(t, transaction)
	assert.Equal(t, status.Error(codes.Unauthenticated, intrabank.ErrUserUnauthenticated), err)

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}

func TestTransferDoPaymentFailed_InsertTransactionFailed(t *testing.T) {
	var (
		corebankingMock = intrabankmock.NewCoreBanking(t)
		repoMock        = intrabankmock.NewRepository(t)
		mailerMock      = intrabankmock.NewReceiptMailer(t)
		seqGenMock      = intrabankmock.NewSequenceGenerator(t)
		svc             = intrabank.NewService(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock)
		ctx             = ctxt.ContextWithUser(context.Background(), &ctxt.User{
			Id:       123,
			CIF:      "1234567",
			FullName: "Olivia Rodrigo",
			Email:    "olivia@gmail.com",
		})
	)

	corebankingMock.EXPECT().GetCoreStatus(mock.Anything).Return(&intrabank.CoreStatus{
		SystemDate:    "25-03-2025",
		Status:        "FINISHED",
		StandInStatus: "N",
	}, nil)

	repoMock.EXPECT().GetSequence(mock.Anything, "123456").
		Return(&intrabank.Sequence{
			SequenceNumber:     "123456",
			Amount:             100000,
			SourceAccount:      "001001234567891",
			DestinationAccount: "001001234567892",
			DestinationName:    "Destination Account",
			SourceName:         "Olivia Rodrigo",
		}, nil)

	corebankingMock.EXPECT().PerformOverbooking(mock.Anything, &intrabank.OverbookingInput{
		SourceAccount:      "001001234567891",
		DestinationAccount: "001001234567892",
		Amount:             100000,
		Fee:                0,
		Remark:             "TRF 001001234567891 001001234567892 BNKYAYA 123456",
	}).Return(&intrabank.OverbookingResult{
		JournalSequence:      "111111",
		TransactionReference: "222222",
	}, nil)

	repoMock.EXPECT().InsertTransaction(mock.Anything, &intrabank.Transaction{
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
	assert.Equal(t, status.Error(codes.Internal, intrabank.ErrGeneral), err)

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}

func TestTransferDoPaymentFailed_SendEmailFailed(t *testing.T) {
	var (
		corebankingMock = intrabankmock.NewCoreBanking(t)
		repoMock        = intrabankmock.NewRepository(t)
		mailerMock      = intrabankmock.NewReceiptMailer(t)
		seqGenMock      = intrabankmock.NewSequenceGenerator(t)
		svc             = intrabank.NewService(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock)
		ctx             = ctxt.ContextWithUser(context.Background(), &ctxt.User{
			Id:       123,
			CIF:      "1234567",
			FullName: "Olivia Rodrigo",
			Email:    "olivia@gmail.com",
		})
	)

	corebankingMock.EXPECT().GetCoreStatus(mock.Anything).Return(&intrabank.CoreStatus{
		SystemDate:    "25-03-2025",
		Status:        "FINISHED",
		StandInStatus: "N",
	}, nil)

	repoMock.EXPECT().GetSequence(mock.Anything, "123456").
		Return(&intrabank.Sequence{
			SequenceNumber:     "123456",
			Amount:             100000,
			SourceAccount:      "001001234567891",
			DestinationAccount: "001001234567892",
			DestinationName:    "Destination Account",
			SourceName:         "Olivia Rodrigo",
		}, nil)

	corebankingMock.EXPECT().PerformOverbooking(mock.Anything, &intrabank.OverbookingInput{
		SourceAccount:      "001001234567891",
		DestinationAccount: "001001234567892",
		Amount:             100000,
		Fee:                0,
		Remark:             "TRF 001001234567891 001001234567892 BNKYAYA 123456",
	}).Return(&intrabank.OverbookingResult{
		JournalSequence:      "111111",
		TransactionReference: "222222",
	}, nil)

	repoMock.EXPECT().InsertTransaction(mock.Anything, &intrabank.Transaction{
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
	assert.Equal(t, status.Error(codes.Internal, intrabank.ErrSendEmailFailed), err)

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}
