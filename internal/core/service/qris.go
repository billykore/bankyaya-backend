package service

import (
	"context"

	"go.bankyaya.org/app/backend/internal/core/domain"
	"go.bankyaya.org/app/backend/internal/core/entity"
	"go.bankyaya.org/app/backend/internal/core/port/api"
	"go.bankyaya.org/app/backend/internal/core/port/email"
	"go.bankyaya.org/app/backend/internal/pkg/codes"
	"go.bankyaya.org/app/backend/internal/pkg/ctxt"
	"go.bankyaya.org/app/backend/internal/pkg/logger"
	"go.bankyaya.org/app/backend/internal/pkg/status"
)

const (
	qrisService        = "QRIService"
	qrisSuccessSubject = "QRIS payment success"
	qrisFee            = 0
)

// QRIS handles QRIS payment process.
type QRIS struct {
	log         *logger.Logger
	corebanking api.CoreBanking
	qrisAPI     api.QRIS
	mailer      email.QRISReceiptMailer
}

func NewQRIS(log *logger.Logger, corebanking api.CoreBanking, qris api.QRIS, mailer email.QRISReceiptMailer) *QRIS {
	return &QRIS{
		log:         log,
		corebanking: corebanking,
		qrisAPI:     qris,
		mailer:      mailer,
	}
}

func (qris *QRIS) Inquiry(ctx context.Context, sourceAccount, qrCode string) (*entity.QRISData, error) {
	eod, err := qris.corebanking.CheckEOD(ctx)
	if err != nil {
		qris.log.ServiceUsecase(qrisService, "CheckEOD").Error(err)
		return nil, status.Error(codes.Internal, domain.ErrGeneral)
	}
	if eod.IsRunning() {
		qris.log.ServiceUsecase(qrisService, "CheckEOD").Error(domain.ErrEODInProgress)
		return nil, status.Error(codes.Internal, domain.ErrEODInProgress)
	}

	srcAccount, err := qris.corebanking.GetAccountDetails(ctx, sourceAccount)
	if err != nil {
		qris.log.ServiceUsecase(qrisService, "GetAccountDetails").Error(err)
		return nil, status.Error(codes.Internal, domain.ErrGeneral)
	}
	if !srcAccount.IsAccountActive() {
		qris.log.ServiceUsecase(qrisService, "GetAccountDetails").Error(domain.ErrSourceAccountInactive)
		return nil, status.Error(codes.Internal, domain.ErrSourceAccountInactive)
	}

	details, err := qris.qrisAPI.GetDetails(ctx, srcAccount.AccountNumber, qrCode)
	if err != nil {
		qris.log.ServiceUsecase(qrisService, "GetDetails").Error(err)
		return nil, status.Error(codes.Internal, domain.ErrGeneral)
	}

	return details, nil
}

func (qris *QRIS) Payment(ctx context.Context, data *entity.QRISData) (*entity.QRISPaymentResult, error) {
	payRes, err := qris.qrisAPI.Pay(ctx, &entity.QRISPaymentData{
		AccountNumber:         data.SourceAccount,
		QRCode:                data.QRCode,
		RRN:                   data.RRN,
		Amount:                data.Amount,
		Tip:                   data.Tip,
		FinancialOrganisation: data.FinancialOrganisation,
		CustomerName:          data.CustomerName,
		MerchantId:            data.MerchantID,
		MerchantCriteria:      data.MerchantCriteria,
		NMId:                  data.NMId,
		AccountName:           data.CustomerName,
	})
	if err != nil {
		qris.log.ServiceUsecase(qrisService, "Pay").Error(err)
		return nil, status.Error(codes.Internal, domain.ErrGeneral)
	}

	user, ok := ctxt.UserFromContext(ctx)
	if !ok {
		qris.log.ServiceUsecase(qrisService, "Pay").Error(ctxt.ErrUserFromContext)
		return nil, status.Error(codes.Unauthenticated, domain.ErrUserUnauthenticated)
	}

	err = qris.mailer.SendQRISReceipt(ctx, entity.QRISEmailData{
		Subject:        qrisSuccessSubject,
		Recipient:      user.Email,
		Amount:         data.Amount,
		Fee:            qrisFee,
		SourceName:     user.FullName,
		SourceAccount:  data.SourceAccount,
		MerchantName:   payRes.MerchantName,
		MerchantPan:    payRes.MerchantPan,
		TransactionRef: payRes.TransactionReference,
		Note:           data.Note,
	})
	if err != nil {
		qris.log.ServiceUsecase(qrisService, "sendTransferReceipt").Error(err)
		return nil, status.Error(codes.Unauthenticated, domain.ErrSendEmailFailed)
	}

	return payRes, nil
}
