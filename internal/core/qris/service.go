package qris

import (
	"context"

	"go.bankyaya.org/app/backend/pkg/ctxt"
	"go.bankyaya.org/app/backend/pkg/logger"
)

const (
	messageEODIsRunning       = "EOD process is running"
	messageInquiryFailed      = "QRIS inquiry failed"
	messagePaymentFailed      = "QRIS payment failed"
	messageAccountIsNotActive = "Account is not active"
)

const (
	qrisSuccessSubject = "QRIS payment success"
)

const qrisFee = 0

// Service handles QRIS payment process.
type Service struct {
	corebanking CoreBanking
	qris        QRIS
	mailer      ReceiptMailer
}

func NewService(log *logger.Logger, corebanking CoreBanking, qris QRIS) *Service {
	return &Service{
		corebanking: corebanking,
		qris:        qris,
	}
}

func (s *Service) Inquiry(ctx context.Context, sourceAccount, qrCode string) (*QRISData, error) {
	eod, err := s.corebanking.CheckEOD(ctx)
	if err != nil {
		return nil, err
	}
	if eod.IsRunning() {
		return nil, ErrEODInProgress
	}

	srcAccount, err := s.corebanking.GetAccountDetails(ctx, sourceAccount)
	if err != nil {
		return nil, err
	}
	if !srcAccount.IsAccountActive() {
		return nil, ErrSourceAccountInactive
	}

	details, err := s.qris.GetDetails(ctx, srcAccount.AccountNumber, qrCode)
	if err != nil {
		return nil, err
	}

	return details, nil
}

func (s *Service) Payment(ctx context.Context, data *QRISData) (*PaymentResult, error) {
	payRes, err := s.qris.Pay(ctx, &PaymentData{
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
		return nil, err
	}

	user, ok := ctxt.UserFromContext(ctx)
	if !ok {
		return nil, ctxt.ErrUserFromContext
	}

	sendErr := s.sendTransferReceipt(ctx, EmailData{
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
			return nil, ErrSendEmailFailed
		}
	default:
	}

	return payRes, nil
}

// sendTransferReceipt is a dedicated function to handle the asynchronous sending of QRIS payment receipt emails.
func (s *Service) sendTransferReceipt(ctx context.Context, emailData EmailData) <-chan error {
	errCh := make(chan error, 1)
	go func() {
		err := s.mailer.SendQRISReceipt(ctx, emailData)
		if err != nil {
			errCh <- err
			close(errCh)
		}
	}()
	return errCh
}
