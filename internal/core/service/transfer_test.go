package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.bankyaya.org/app/backend/internal/core/domain"
	"go.bankyaya.org/app/backend/internal/core/entity"
	"go.bankyaya.org/app/backend/internal/core/port/api/mock"
	"go.bankyaya.org/app/backend/internal/core/port/email/mock"
	"go.bankyaya.org/app/backend/internal/core/port/repository/mock"
	"go.bankyaya.org/app/backend/internal/core/port/security/mock"
	"go.bankyaya.org/app/backend/internal/pkg/codes"
	"go.bankyaya.org/app/backend/internal/pkg/ctxt"
	"go.bankyaya.org/app/backend/internal/pkg/data"
	"go.bankyaya.org/app/backend/internal/pkg/logger"
	"go.bankyaya.org/app/backend/internal/pkg/status"
)

func TestTransferInquirySuccess(t *testing.T) {
	var (
		corebankingMock = apimock.NewCoreBanking(t)
		repoMock        = repomock.NewTransferRepository(t)
		mailerMock      = emailmock.NewTransferReceiptMailer(t)
		seqGenMock      = securitymock.NewSequenceGenerator(t)
		svc             = NewTransfer(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock)
		ctx             = ctxt.ContextWithUser(context.Background(), data.User{
			Id:       123,
			CIF:      "1234567",
			FullName: "Olivia Rodrigo",
			Email:    "olivia@gmail.com",
		})
	)

	corebankingMock.EXPECT().CheckEOD(mock.Anything).Return(&entity.EODData{
		SystemDate:    "25-03-2025",
		Status:        "FINISHED",
		StandInStatus: "N",
	}, nil)
	corebankingMock.EXPECT().GetAccountDetails(mock.Anything, "001001234567891").
		Return(&entity.AccountDetails{
			Name:   "Olivia Rodrigo",
			Status: "1",
		}, nil)
	corebankingMock.EXPECT().GetAccountDetails(mock.Anything, "001001234567892").
		Return(&entity.AccountDetails{
			Name:   "Destination Account",
			Status: "1",
		}, nil)

	repoMock.EXPECT().InsertSequence(mock.Anything, &entity.Sequence{
		SeqNo:       "123456",
		Amount:      "100000",
		AccNoSrc:    "001001234567891",
		AccNoDest:   "001001234567892",
		AccNameDest: "Destination Account",
		AccNameSrc:  "Olivia Rodrigo",
	}).Return(nil)

	seqGenMock.EXPECT().Generate().
		Return("123456", nil)

	sequence, err := svc.Inquiry(ctx, &entity.Sequence{
		SeqNo:     "123456",
		Amount:    "100000",
		AccNoSrc:  "001001234567891",
		AccNoDest: "001001234567892",
	})

	assert.Nil(t, err)
	assert.Equal(t, sequence, &entity.Sequence{
		SeqNo:       "123456",
		Amount:      "100000",
		AccNoSrc:    "001001234567891",
		AccNoDest:   "001001234567892",
		AccNameDest: "Destination Account",
		AccNameSrc:  "Olivia Rodrigo",
	})

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}

func TestTransferInquiryFailed_CheckEODFailed(t *testing.T) {
	var (
		corebankingMock = apimock.NewCoreBanking(t)
		repoMock        = repomock.NewTransferRepository(t)
		mailerMock      = emailmock.NewTransferReceiptMailer(t)
		seqGenMock      = securitymock.NewSequenceGenerator(t)
		svc             = NewTransfer(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock)
		ctx             = ctxt.ContextWithUser(context.Background(), data.User{
			Id:       123,
			CIF:      "1234567",
			FullName: "Olivia Rodrigo",
			Email:    "olivia@gmail.com",
		})
	)

	corebankingMock.EXPECT().CheckEOD(mock.Anything).
		Return(nil, errors.New("check EOD failed"))

	sequence, err := svc.Inquiry(ctx, &entity.Sequence{
		SeqNo:     "123456",
		Amount:    "100000",
		AccNoSrc:  "001001234567891",
		AccNoDest: "001001234567892",
	})
	assert.Nil(t, sequence)
	assert.Equal(t, status.Error(codes.Internal, domain.ErrGeneral), err)

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}

func TestTransferInquiryFailed_EODIsRunning(t *testing.T) {
	var (
		corebankingMock = apimock.NewCoreBanking(t)
		repoMock        = repomock.NewTransferRepository(t)
		mailerMock      = emailmock.NewTransferReceiptMailer(t)
		seqGenMock      = securitymock.NewSequenceGenerator(t)
		svc             = NewTransfer(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock)
		ctx             = ctxt.ContextWithUser(context.Background(), data.User{
			Id:       123,
			CIF:      "1234567",
			FullName: "Olivia Rodrigo",
			Email:    "olivia@gmail.com",
		})
	)

	corebankingMock.EXPECT().CheckEOD(mock.Anything).Return(&entity.EODData{
		SystemDate:    "25-03-2025",
		Status:        "STARTED",
		StandInStatus: "N",
	}, nil)

	sequence, err := svc.Inquiry(ctx, &entity.Sequence{
		SeqNo:     "123456",
		Amount:    "100000",
		AccNoSrc:  "001001234567891",
		AccNoDest: "001001234567892",
	})
	assert.Nil(t, sequence)
	assert.Equal(t, status.Error(codes.Internal, domain.ErrEODInProgress), err)

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}

func TestTransferInquiryFailed_FailedCheckSourceAccount(t *testing.T) {
	var (
		corebankingMock = apimock.NewCoreBanking(t)
		repoMock        = repomock.NewTransferRepository(t)
		mailerMock      = emailmock.NewTransferReceiptMailer(t)
		seqGenMock      = securitymock.NewSequenceGenerator(t)
		svc             = NewTransfer(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock)
		ctx             = ctxt.ContextWithUser(context.Background(), data.User{
			Id:       123,
			CIF:      "1234567",
			FullName: "Olivia Rodrigo",
			Email:    "olivia@gmail.com",
		})
	)

	corebankingMock.EXPECT().CheckEOD(mock.Anything).Return(&entity.EODData{
		SystemDate:    "25-03-2025",
		Status:        "FINISHED",
		StandInStatus: "N",
	}, nil)
	corebankingMock.EXPECT().GetAccountDetails(mock.Anything, "001001234567891").
		Return(nil, errors.New("GetAccountDetails failed"))

	sequence, err := svc.Inquiry(ctx, &entity.Sequence{
		SeqNo:     "123456",
		Amount:    "100000",
		AccNoSrc:  "001001234567891",
		AccNoDest: "001001234567892",
	})

	assert.Nil(t, sequence)
	assert.Equal(t, status.Error(codes.Internal, domain.ErrGeneral), err)

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}

func TestTransferInquiryFailed_SourceAccountIsInactive(t *testing.T) {
	var (
		corebankingMock = apimock.NewCoreBanking(t)
		repoMock        = repomock.NewTransferRepository(t)
		mailerMock      = emailmock.NewTransferReceiptMailer(t)
		seqGenMock      = securitymock.NewSequenceGenerator(t)
		svc             = NewTransfer(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock)
		ctx             = ctxt.ContextWithUser(context.Background(), data.User{
			Id:       123,
			CIF:      "1234567",
			FullName: "Olivia Rodrigo",
			Email:    "olivia@gmail.com",
		})
	)

	corebankingMock.EXPECT().CheckEOD(mock.Anything).Return(&entity.EODData{
		SystemDate:    "25-03-2025",
		Status:        "FINISHED",
		StandInStatus: "N",
	}, nil)
	corebankingMock.EXPECT().GetAccountDetails(mock.Anything, "001001234567891").
		Return(&entity.AccountDetails{
			Status: "9",
		}, nil)

	sequence, err := svc.Inquiry(ctx, &entity.Sequence{
		SeqNo:     "123456",
		Amount:    "100000",
		AccNoSrc:  "001001234567891",
		AccNoDest: "001001234567892",
	})

	assert.Nil(t, sequence)
	assert.Equal(t, status.Error(codes.BadRequest, domain.ErrSourceAccountInactive), err)

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}

func TestTransferInquiryFailed_CheckDestinationAccountFailed(t *testing.T) {
	var (
		corebankingMock = apimock.NewCoreBanking(t)
		repoMock        = repomock.NewTransferRepository(t)
		mailerMock      = emailmock.NewTransferReceiptMailer(t)
		seqGenMock      = securitymock.NewSequenceGenerator(t)
		svc             = NewTransfer(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock)
		ctx             = ctxt.ContextWithUser(context.Background(), data.User{
			Id:       123,
			CIF:      "1234567",
			FullName: "Olivia Rodrigo",
			Email:    "olivia@gmail.com",
		})
	)

	corebankingMock.EXPECT().CheckEOD(mock.Anything).Return(&entity.EODData{
		SystemDate:    "25-03-2025",
		Status:        "FINISHED",
		StandInStatus: "N",
	}, nil)
	corebankingMock.EXPECT().GetAccountDetails(mock.Anything, "001001234567891").
		Return(&entity.AccountDetails{
			Name:   "Olivia Rodrigo",
			Status: "1",
		}, nil)
	corebankingMock.EXPECT().GetAccountDetails(mock.Anything, "001001234567892").
		Return(nil, errors.New("GetAccountDetails failed"))

	sequence, err := svc.Inquiry(ctx, &entity.Sequence{
		SeqNo:     "123456",
		Amount:    "100000",
		AccNoSrc:  "001001234567891",
		AccNoDest: "001001234567892",
	})

	assert.Nil(t, sequence)
	assert.Equal(t, status.Error(codes.Internal, domain.ErrGeneral), err)

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}

func TestTransferInquiryFailed_DestinationAccountInactive(t *testing.T) {
	var (
		corebankingMock = apimock.NewCoreBanking(t)
		repoMock        = repomock.NewTransferRepository(t)
		mailerMock      = emailmock.NewTransferReceiptMailer(t)
		seqGenMock      = securitymock.NewSequenceGenerator(t)
		svc             = NewTransfer(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock)
		ctx             = ctxt.ContextWithUser(context.Background(), data.User{
			Id:       123,
			CIF:      "1234567",
			FullName: "Olivia Rodrigo",
			Email:    "olivia@gmail.com",
		})
	)

	corebankingMock.EXPECT().CheckEOD(mock.Anything).Return(&entity.EODData{
		SystemDate:    "25-03-2025",
		Status:        "FINISHED",
		StandInStatus: "N",
	}, nil)
	corebankingMock.EXPECT().GetAccountDetails(mock.Anything, "001001234567891").
		Return(&entity.AccountDetails{
			Name:   "Olivia Rodrigo",
			Status: "1",
		}, nil)
	corebankingMock.EXPECT().GetAccountDetails(mock.Anything, "001001234567892").
		Return(&entity.AccountDetails{
			Status: "9",
		}, nil)

	sequence, err := svc.Inquiry(ctx, &entity.Sequence{
		SeqNo:     "123456",
		Amount:    "100000",
		AccNoSrc:  "001001234567891",
		AccNoDest: "001001234567892",
	})

	assert.Nil(t, sequence)
	assert.Equal(t, status.Error(codes.BadRequest, domain.ErrDestinationAccountInactive), err)

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}

func TestTransferInquiryFailed_GenerateSequenceFailed(t *testing.T) {
	var (
		corebankingMock = apimock.NewCoreBanking(t)
		repoMock        = repomock.NewTransferRepository(t)
		mailerMock      = emailmock.NewTransferReceiptMailer(t)
		seqGenMock      = securitymock.NewSequenceGenerator(t)
		svc             = NewTransfer(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock)
		ctx             = ctxt.ContextWithUser(context.Background(), data.User{
			Id:       123,
			CIF:      "1234567",
			FullName: "Olivia Rodrigo",
			Email:    "olivia@gmail.com",
		})
	)

	corebankingMock.EXPECT().CheckEOD(mock.Anything).Return(&entity.EODData{
		SystemDate:    "25-03-2025",
		Status:        "FINISHED",
		StandInStatus: "N",
	}, nil)
	corebankingMock.EXPECT().GetAccountDetails(mock.Anything, "001001234567891").
		Return(&entity.AccountDetails{
			Name:   "Olivia Rodrigo",
			Status: "1",
		}, nil)
	corebankingMock.EXPECT().GetAccountDetails(mock.Anything, "001001234567892").
		Return(&entity.AccountDetails{
			Name:   "Destination Account",
			Status: "1",
		}, nil)

	seqGenMock.EXPECT().Generate().
		Return("", errors.New("some error"))

	sequence, err := svc.Inquiry(ctx, &entity.Sequence{
		SeqNo:     "123456",
		Amount:    "100000",
		AccNoSrc:  "001001234567891",
		AccNoDest: "001001234567892",
	})

	assert.Nil(t, sequence)
	assert.Equal(t, status.Error(codes.Internal, domain.ErrGeneral), err)

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}

func TestTransferInquiryFailed_InsertSequenceFailed(t *testing.T) {
	var (
		corebankingMock = apimock.NewCoreBanking(t)
		repoMock        = repomock.NewTransferRepository(t)
		mailerMock      = emailmock.NewTransferReceiptMailer(t)
		seqGenMock      = securitymock.NewSequenceGenerator(t)
		svc             = NewTransfer(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock)
		ctx             = ctxt.ContextWithUser(context.Background(), data.User{
			Id:       123,
			CIF:      "1234567",
			FullName: "Olivia Rodrigo",
			Email:    "olivia@gmail.com",
		})
	)

	corebankingMock.EXPECT().CheckEOD(mock.Anything).Return(&entity.EODData{
		SystemDate:    "25-03-2025",
		Status:        "FINISHED",
		StandInStatus: "N",
	}, nil)
	corebankingMock.EXPECT().GetAccountDetails(mock.Anything, "001001234567891").
		Return(&entity.AccountDetails{
			Name:   "Olivia Rodrigo",
			Status: "1",
		}, nil)
	corebankingMock.EXPECT().GetAccountDetails(mock.Anything, "001001234567892").
		Return(&entity.AccountDetails{
			Name:   "Destination Account",
			Status: "1",
		}, nil)

	seqGenMock.EXPECT().Generate().
		Return("123456", nil)

	repoMock.EXPECT().InsertSequence(mock.Anything, &entity.Sequence{
		SeqNo:       "123456",
		Amount:      "100000",
		AccNoSrc:    "001001234567891",
		AccNoDest:   "001001234567892",
		AccNameDest: "Destination Account",
		AccNameSrc:  "Olivia Rodrigo",
	}).Return(errors.New("some error"))

	sequence, err := svc.Inquiry(ctx, &entity.Sequence{
		SeqNo:     "123456",
		Amount:    "100000",
		AccNoSrc:  "001001234567891",
		AccNoDest: "001001234567892",
	})

	assert.Nil(t, sequence)
	assert.Equal(t, status.Error(codes.Internal, domain.ErrGeneral), err)

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}

func TestTransferDoPaymentSuccess(t *testing.T) {
	var (
		corebankingMock = apimock.NewCoreBanking(t)
		repoMock        = repomock.NewTransferRepository(t)
		mailerMock      = emailmock.NewTransferReceiptMailer(t)
		seqGenMock      = securitymock.NewSequenceGenerator(t)
		svc             = NewTransfer(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock)
		ctx             = ctxt.ContextWithUser(context.Background(), data.User{
			Id:       123,
			CIF:      "1234567",
			FullName: "Olivia Rodrigo",
			Email:    "olivia@gmail.com",
		})
	)

	corebankingMock.EXPECT().CheckEOD(mock.Anything).Return(&entity.EODData{
		SystemDate:    "25-03-2025",
		Status:        "FINISHED",
		StandInStatus: "N",
	}, nil)

	repoMock.EXPECT().GetSequence(mock.Anything, "123456").
		Return(&entity.Sequence{
			SeqNo:       "123456",
			Amount:      "100000",
			AccNoSrc:    "001001234567891",
			AccNoDest:   "001001234567892",
			AccNameDest: "Destination Account",
			AccNameSrc:  "Olivia Rodrigo",
		}, nil)

	corebankingMock.EXPECT().PerformOverbooking(mock.Anything, &entity.OverbookingRequest{
		SourceAccount:      "001001234567891",
		DestinationAccount: "001001234567892",
		Amount:             100000,
		Fee:                0,
		Remark:             "TRF 001001234567891 001001234567892 BNKYAYA 123456",
	}).Return(&entity.OverbookingResponse{
		JournalSequence:      "111111",
		TransactionReference: "222222",
	}, nil)

	repoMock.EXPECT().InsertTransaction(mock.Anything, &entity.Transaction{
		SequenceNumber:       "123456",
		SequenceJournal:      "111111",
		UserID:               "123",
		Destination:          "001001234567892",
		Amount:               "100000",
		TransactionType:      "internal_transfer",
		TransactionReference: "222222",
		Remarks:              "TRF 001001234567891 001001234567892 BNKYAYA 123456",
		Fee:                  "0",
		NameDest:             "Destination Account",
	}).Return(nil)

	mailerMock.EXPECT().SendTransferReceipt(mock.Anything, mock.Anything).Return(nil)

	transaction, err := svc.DoPayment(ctx, "123456")

	assert.NoError(t, err)
	assert.Equal(t, &entity.Transaction{
		SequenceNumber:       "123456",
		SequenceJournal:      "111111",
		UserID:               "123",
		Destination:          "001001234567892",
		Amount:               "100000",
		TransactionType:      "internal_transfer",
		TransactionReference: "222222",
		Remarks:              "TRF 001001234567891 001001234567892 BNKYAYA 123456",
		Fee:                  "0",
		NameDest:             "Destination Account",
	}, transaction)

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}

func TestTransferDoPaymentFailed_CheckEODFailed(t *testing.T) {
	var (
		corebankingMock = apimock.NewCoreBanking(t)
		repoMock        = repomock.NewTransferRepository(t)
		mailerMock      = emailmock.NewTransferReceiptMailer(t)
		seqGenMock      = securitymock.NewSequenceGenerator(t)
		svc             = NewTransfer(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock)
		ctx             = ctxt.ContextWithUser(context.Background(), data.User{
			Id:       123,
			CIF:      "1234567",
			FullName: "Olivia Rodrigo",
			Email:    "olivia@gmail.com",
		})
	)

	corebankingMock.EXPECT().CheckEOD(mock.Anything).
		Return(nil, errors.New("some error"))

	transaction, err := svc.DoPayment(ctx, "123456")

	assert.Nil(t, transaction)
	assert.Equal(t, status.Error(codes.Internal, domain.ErrGeneral), err)

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}

func TestTransferDoPaymentFailed_EODIsRunning(t *testing.T) {
	var (
		corebankingMock = apimock.NewCoreBanking(t)
		repoMock        = repomock.NewTransferRepository(t)
		mailerMock      = emailmock.NewTransferReceiptMailer(t)
		seqGenMock      = securitymock.NewSequenceGenerator(t)
		svc             = NewTransfer(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock)
		ctx             = ctxt.ContextWithUser(context.Background(), data.User{
			Id:       123,
			CIF:      "1234567",
			FullName: "Olivia Rodrigo",
			Email:    "olivia@gmail.com",
		})
	)

	corebankingMock.EXPECT().CheckEOD(mock.Anything).
		Return(&entity.EODData{
			SystemDate:    "25-03-2025",
			Status:        "STARTED",
			StandInStatus: "N",
		}, nil)

	transaction, err := svc.DoPayment(ctx, "123456")

	assert.Nil(t, transaction)
	assert.Equal(t, status.Error(codes.Internal, domain.ErrEODInProgress), err)

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}

func TestTransferDoPaymentFailed_GetSequenceFailed(t *testing.T) {
	var (
		corebankingMock = apimock.NewCoreBanking(t)
		repoMock        = repomock.NewTransferRepository(t)
		mailerMock      = emailmock.NewTransferReceiptMailer(t)
		seqGenMock      = securitymock.NewSequenceGenerator(t)
		svc             = NewTransfer(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock)
		ctx             = ctxt.ContextWithUser(context.Background(), data.User{
			Id:       123,
			CIF:      "1234567",
			FullName: "Olivia Rodrigo",
			Email:    "olivia@gmail.com",
		})
	)

	corebankingMock.EXPECT().CheckEOD(mock.Anything).Return(&entity.EODData{
		SystemDate:    "25-03-2025",
		Status:        "FINISHED",
		StandInStatus: "N",
	}, nil)

	repoMock.EXPECT().GetSequence(mock.Anything, "123456").
		Return(nil, errors.New("some error"))

	transaction, err := svc.DoPayment(ctx, "123456")

	assert.Nil(t, transaction)
	assert.Equal(t, status.Error(codes.Internal, domain.ErrGeneral), err)

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}

func TestTransferDoPaymentFailed_InvalidSequenceAmount(t *testing.T) {
	var (
		corebankingMock = apimock.NewCoreBanking(t)
		repoMock        = repomock.NewTransferRepository(t)
		mailerMock      = emailmock.NewTransferReceiptMailer(t)
		seqGenMock      = securitymock.NewSequenceGenerator(t)
		svc             = NewTransfer(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock)
		ctx             = ctxt.ContextWithUser(context.Background(), data.User{
			Id:       123,
			CIF:      "1234567",
			FullName: "Olivia Rodrigo",
			Email:    "olivia@gmail.com",
		})
	)

	corebankingMock.EXPECT().CheckEOD(mock.Anything).Return(&entity.EODData{
		SystemDate:    "25-03-2025",
		Status:        "FINISHED",
		StandInStatus: "N",
	}, nil)

	repoMock.EXPECT().GetSequence(mock.Anything, "123456").
		Return(&entity.Sequence{
			SeqNo:       "123456",
			Amount:      "100000m",
			AccNoSrc:    "001001234567891",
			AccNoDest:   "001001234567892",
			AccNameDest: "Destination Account",
			AccNameSrc:  "Olivia Rodrigo",
		}, nil)

	transaction, err := svc.DoPayment(ctx, "123456")

	assert.Nil(t, transaction)
	assert.Equal(t, status.Error(codes.BadRequest, domain.ErrGeneral), err)

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}

func TestTransferDoPaymentFailed_PerformOverbookingFailed(t *testing.T) {
	var (
		corebankingMock = apimock.NewCoreBanking(t)
		repoMock        = repomock.NewTransferRepository(t)
		mailerMock      = emailmock.NewTransferReceiptMailer(t)
		seqGenMock      = securitymock.NewSequenceGenerator(t)
		svc             = NewTransfer(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock)
		ctx             = ctxt.ContextWithUser(context.Background(), data.User{
			Id:       123,
			CIF:      "1234567",
			FullName: "Olivia Rodrigo",
			Email:    "olivia@gmail.com",
		})
	)

	corebankingMock.EXPECT().CheckEOD(mock.Anything).Return(&entity.EODData{
		SystemDate:    "25-03-2025",
		Status:        "FINISHED",
		StandInStatus: "N",
	}, nil)

	repoMock.EXPECT().GetSequence(mock.Anything, "123456").
		Return(&entity.Sequence{
			SeqNo:       "123456",
			Amount:      "100000",
			AccNoSrc:    "001001234567891",
			AccNoDest:   "001001234567892",
			AccNameDest: "Destination Account",
			AccNameSrc:  "Olivia Rodrigo",
		}, nil)

	corebankingMock.EXPECT().PerformOverbooking(mock.Anything, &entity.OverbookingRequest{
		SourceAccount:      "001001234567891",
		DestinationAccount: "001001234567892",
		Amount:             100000,
		Fee:                0,
		Remark:             "TRF 001001234567891 001001234567892 BNKYAYA 123456",
	}).Return(nil, errors.New("some error"))

	transaction, err := svc.DoPayment(ctx, "123456")

	assert.Nil(t, transaction)
	assert.Equal(t, status.Error(codes.Internal, domain.ErrGeneral), err)

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}

func TestTransferDoPaymentFailed_GetUserFromContextFailed(t *testing.T) {
	var (
		corebankingMock = apimock.NewCoreBanking(t)
		repoMock        = repomock.NewTransferRepository(t)
		mailerMock      = emailmock.NewTransferReceiptMailer(t)
		seqGenMock      = securitymock.NewSequenceGenerator(t)
		svc             = NewTransfer(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock)
		ctx             = context.Background()
	)

	corebankingMock.EXPECT().CheckEOD(mock.Anything).Return(&entity.EODData{
		SystemDate:    "25-03-2025",
		Status:        "FINISHED",
		StandInStatus: "N",
	}, nil)

	repoMock.EXPECT().GetSequence(mock.Anything, "123456").
		Return(&entity.Sequence{
			SeqNo:       "123456",
			Amount:      "100000",
			AccNoSrc:    "001001234567891",
			AccNoDest:   "001001234567892",
			AccNameDest: "Destination Account",
			AccNameSrc:  "Olivia Rodrigo",
		}, nil)

	corebankingMock.EXPECT().PerformOverbooking(mock.Anything, &entity.OverbookingRequest{
		SourceAccount:      "001001234567891",
		DestinationAccount: "001001234567892",
		Amount:             100000,
		Fee:                0,
		Remark:             "TRF 001001234567891 001001234567892 BNKYAYA 123456",
	}).Return(&entity.OverbookingResponse{
		JournalSequence:      "111111",
		TransactionReference: "222222",
	}, nil)

	transaction, err := svc.DoPayment(ctx, "123456")

	assert.Nil(t, transaction)
	assert.Equal(t, status.Error(codes.Unauthenticated, domain.ErrUserUnauthenticated), err)

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}

func TestTransferDoPaymentFailed_InsertTransactionFailed(t *testing.T) {
	var (
		corebankingMock = apimock.NewCoreBanking(t)
		repoMock        = repomock.NewTransferRepository(t)
		mailerMock      = emailmock.NewTransferReceiptMailer(t)
		seqGenMock      = securitymock.NewSequenceGenerator(t)
		svc             = NewTransfer(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock)
		ctx             = ctxt.ContextWithUser(context.Background(), data.User{
			Id:       123,
			CIF:      "1234567",
			FullName: "Olivia Rodrigo",
			Email:    "olivia@gmail.com",
		})
	)

	corebankingMock.EXPECT().CheckEOD(mock.Anything).Return(&entity.EODData{
		SystemDate:    "25-03-2025",
		Status:        "FINISHED",
		StandInStatus: "N",
	}, nil)

	repoMock.EXPECT().GetSequence(mock.Anything, "123456").
		Return(&entity.Sequence{
			SeqNo:       "123456",
			Amount:      "100000",
			AccNoSrc:    "001001234567891",
			AccNoDest:   "001001234567892",
			AccNameDest: "Destination Account",
			AccNameSrc:  "Olivia Rodrigo",
		}, nil)

	corebankingMock.EXPECT().PerformOverbooking(mock.Anything, &entity.OverbookingRequest{
		SourceAccount:      "001001234567891",
		DestinationAccount: "001001234567892",
		Amount:             100000,
		Fee:                0,
		Remark:             "TRF 001001234567891 001001234567892 BNKYAYA 123456",
	}).Return(&entity.OverbookingResponse{
		JournalSequence:      "111111",
		TransactionReference: "222222",
	}, nil)

	repoMock.EXPECT().InsertTransaction(mock.Anything, &entity.Transaction{
		SequenceNumber:       "123456",
		SequenceJournal:      "111111",
		UserID:               "123",
		Destination:          "001001234567892",
		Amount:               "100000",
		TransactionType:      "internal_transfer",
		TransactionReference: "222222",
		Remarks:              "TRF 001001234567891 001001234567892 BNKYAYA 123456",
		Fee:                  "0",
		NameDest:             "Destination Account",
	}).Return(errors.New("some error"))

	transaction, err := svc.DoPayment(ctx, "123456")

	assert.Nil(t, transaction)
	assert.Equal(t, status.Error(codes.Internal, domain.ErrGeneral), err)

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}

func TestTransferDoPaymentFailed_SendEmailFailed(t *testing.T) {
	var (
		corebankingMock = apimock.NewCoreBanking(t)
		repoMock        = repomock.NewTransferRepository(t)
		mailerMock      = emailmock.NewTransferReceiptMailer(t)
		seqGenMock      = securitymock.NewSequenceGenerator(t)
		svc             = NewTransfer(logger.New(), repoMock, corebankingMock, seqGenMock, mailerMock)
		ctx             = ctxt.ContextWithUser(context.Background(), data.User{
			Id:       123,
			CIF:      "1234567",
			FullName: "Olivia Rodrigo",
			Email:    "olivia@gmail.com",
		})
	)

	corebankingMock.EXPECT().CheckEOD(mock.Anything).Return(&entity.EODData{
		SystemDate:    "25-03-2025",
		Status:        "FINISHED",
		StandInStatus: "N",
	}, nil)

	repoMock.EXPECT().GetSequence(mock.Anything, "123456").
		Return(&entity.Sequence{
			SeqNo:       "123456",
			Amount:      "100000",
			AccNoSrc:    "001001234567891",
			AccNoDest:   "001001234567892",
			AccNameDest: "Destination Account",
			AccNameSrc:  "Olivia Rodrigo",
		}, nil)

	corebankingMock.EXPECT().PerformOverbooking(mock.Anything, &entity.OverbookingRequest{
		SourceAccount:      "001001234567891",
		DestinationAccount: "001001234567892",
		Amount:             100000,
		Fee:                0,
		Remark:             "TRF 001001234567891 001001234567892 BNKYAYA 123456",
	}).Return(&entity.OverbookingResponse{
		JournalSequence:      "111111",
		TransactionReference: "222222",
	}, nil)

	repoMock.EXPECT().InsertTransaction(mock.Anything, &entity.Transaction{
		SequenceNumber:       "123456",
		SequenceJournal:      "111111",
		UserID:               "123",
		Destination:          "001001234567892",
		Amount:               "100000",
		TransactionType:      "internal_transfer",
		TransactionReference: "222222",
		Remarks:              "TRF 001001234567891 001001234567892 BNKYAYA 123456",
		Fee:                  "0",
		NameDest:             "Destination Account",
	}).Return(nil)

	mailerMock.EXPECT().SendTransferReceipt(mock.Anything, mock.Anything).
		Return(errors.New("some error"))

	transaction, err := svc.DoPayment(ctx, "123456")

	assert.Nil(t, transaction)
	assert.Equal(t, status.Error(codes.Internal, domain.ErrSendEmailFailed), err)

	corebankingMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
	seqGenMock.AssertExpectations(t)
}
