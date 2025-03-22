package service

import (
	"context"

	"go.bankyaya.org/app/backend/pkg/entity"
	pkgerrors "go.bankyaya.org/app/backend/pkg/errors"
	"go.bankyaya.org/app/backend/pkg/interface/api"
	"go.bankyaya.org/app/backend/pkg/interface/email"
	"go.bankyaya.org/app/backend/pkg/util/codes"
	"go.bankyaya.org/app/backend/pkg/util/ctxt"
	"go.bankyaya.org/app/backend/pkg/util/logger"
	"go.bankyaya.org/app/backend/pkg/util/status"
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
		return nil, status.Error(codes.Internal, pkgerrors.ErrGeneral)
	}
	if eod.IsRunning() {
		qris.log.ServiceUsecase(qrisService, "CheckEOD").Error(pkgerrors.ErrEODInProgress)
		return nil, status.Error(codes.Internal, pkgerrors.ErrEODInProgress)
	}

	srcAccount, err := qris.corebanking.GetAccountDetails(ctx, sourceAccount)
	if err != nil {
		qris.log.ServiceUsecase(qrisService, "GetAccountDetails").Error(err)
		return nil, status.Error(codes.Internal, pkgerrors.ErrGeneral)
	}
	if !srcAccount.IsAccountActive() {
		qris.log.ServiceUsecase(qrisService, "GetAccountDetails").Error(pkgerrors.ErrSourceAccountInactive)
		return nil, status.Error(codes.Internal, pkgerrors.ErrSourceAccountInactive)
	}

	details, err := qris.qrisAPI.GetDetails(ctx, srcAccount.AccountNumber, qrCode)
	if err != nil {
		qris.log.ServiceUsecase(qrisService, "GetDetails").Error(err)
		return nil, status.Error(codes.Internal, pkgerrors.ErrGeneral)
	}

	return details, nil
}

func (qris *QRIS) Payment(ctx context.Context, data *entity.QRISData) (*entity.QRISPaymentResult, error) {
	payRes, err := qris.qrisAPI.Pay(ctx, &entity.PaymentData{
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
		return nil, status.Error(codes.Internal, pkgerrors.ErrGeneral)
	}

	user, ok := ctxt.UserFromContext(ctx)
	if !ok {
		qris.log.ServiceUsecase(qrisService, "Pay").Error(ctxt.ErrUserFromContext)
		return nil, status.Error(codes.Unauthenticated, pkgerrors.ErrUserUnauthenticated)
	}

	sendErr := qris.sendTransferReceipt(ctx, entity.QRISEmailData{
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
	select {
	case err := <-sendErr:
		if err != nil {
			qris.log.ServiceUsecase(qrisService, "sendTransferReceipt").Error(err)
			return nil, status.Error(codes.Unauthenticated, pkgerrors.ErrSendEmailFailed)
		}
	default:
	}

	return payRes, nil
}

// sendTransferReceipt is a dedicated function to handle the asynchronous sending of QRIS payment receipt emails.
func (qris *QRIS) sendTransferReceipt(ctx context.Context, emailData entity.QRISEmailData) <-chan error {
	errCh := make(chan error, 1)
	go func() {
		err := qris.mailer.SendQRISReceipt(ctx, emailData)
		if err != nil {
			errCh <- err
			close(errCh)
		}
	}()
	return errCh
}
