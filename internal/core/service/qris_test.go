package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.bankyaya.org/app/backend/internal/core/entity"
	apimock "go.bankyaya.org/app/backend/internal/core/port/api/mock"
	emailmock "go.bankyaya.org/app/backend/internal/core/port/email/mock"
	"go.bankyaya.org/app/backend/internal/pkg/ctxt"
	"go.bankyaya.org/app/backend/internal/pkg/data"
	"go.bankyaya.org/app/backend/internal/pkg/logger"
)

func TestQRISInquirySuccess(t *testing.T) {
	var (
		corebankingMock = apimock.NewCoreBanking(t)
		qrisAPIMock     = apimock.NewQRIS(t)
		mailerMock      = emailmock.NewQRISReceiptMailer(t)
		svc             = NewQRIS(logger.New(), corebankingMock, qrisAPIMock, mailerMock)
		ctx             = ctxt.ContextWithUser(context.Background(), data.User{
			Id:       123,
			CIF:      "1234567",
			FullName: "Olivia Rodrigo",
			Email:    "olivia@gmail.com",
		})
	)

	corebankingMock.EXPECT().CheckEOD(mock.Anything).
		Return(&entity.EODData{
			Status:        "FINISHED",
			StandInStatus: "N",
		}, nil)
	corebankingMock.EXPECT().GetAccountDetails(mock.Anything, "1234567").
		Return(&entity.AccountDetails{
			Status:        "1",
			AccountNumber: "1234567",
		}, nil)

	qrisAPIMock.EXPECT().GetDetails(mock.Anything, "1234567", "qr-code-001").
		Return(&entity.QRISData{
			SourceAccount: "1234567",
			QRCode:        "qr-code-001",
		}, nil)

	qrisData, err := svc.Inquiry(ctx, "1234567", "qr-code-001")

	assert.NoError(t, err)
	assert.Equal(t, &entity.QRISData{
		SourceAccount: "1234567",
		QRCode:        "qr-code-001",
	}, qrisData)

	corebankingMock.AssertExpectations(t)
	qrisAPIMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
}

func TestQRISPaymentSuccess(t *testing.T) {
	var (
		corebankingMock = apimock.NewCoreBanking(t)
		qrisAPIMock     = apimock.NewQRIS(t)
		mailerMock      = emailmock.NewQRISReceiptMailer(t)
		svc             = NewQRIS(logger.New(), corebankingMock, qrisAPIMock, mailerMock)
		ctx             = ctxt.ContextWithUser(context.Background(), data.User{
			Id:       123,
			CIF:      "1234567",
			FullName: "Olivia Rodrigo",
			Email:    "olivia@gmail.com",
		})
	)

	qrisAPIMock.EXPECT().Pay(mock.Anything, &entity.QRISPaymentData{
		AccountNumber: "1234567",
	}).Return(&entity.QRISPaymentResult{
		TransactionReference: "tx001",
	}, nil)

	mailerMock.EXPECT().SendQRISReceipt(mock.Anything, mock.Anything).
		Return(nil)

	result, err := svc.Payment(ctx, &entity.QRISData{
		SourceAccount: "1234567",
	})

	assert.NoError(t, err)
	assert.Equal(t, &entity.QRISPaymentResult{
		TransactionReference: "tx001",
	}, result)

	corebankingMock.AssertExpectations(t)
	qrisAPIMock.AssertExpectations(t)
	mailerMock.AssertExpectations(t)
}
