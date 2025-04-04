package qris

import (
	"context"
	"fmt"

	"go.bankyaya.org/app/backend/internal/core/domain"
	"go.bankyaya.org/app/backend/internal/core/entity"
	"go.bankyaya.org/app/backend/internal/pkg/qris"
	"go.bankyaya.org/app/backend/internal/pkg/types"
)

type QRIS struct {
	client *qris.Client
}

func NewQRIS(client *qris.Client) *QRIS {
	return &QRIS{
		client: client,
	}
}

func (q *QRIS) GetDetails(ctx context.Context, accountNumber, qrCode string) (*entity.QRISData, error) {
	details, err := q.client.Inquiry(ctx, qris.InquiryRequest{
		AccountNumber: accountNumber,
		QrCode:        qrCode,
	})
	if err != nil {
		return nil, err
	}
	amount, err := types.ParseMoney(details.Data.Amount)
	if err != nil {
		return nil, err
	}
	return &entity.QRISData{
		Status:                       details.Data.Status,
		RRN:                          details.Data.RRN,
		CustomerName:                 details.Data.CustomerName,
		CustomerDetail:               details.Data.DetailCustomer,
		FinancialOrganisation:        details.Data.LembagaKeuangan,
		FinancialOrganisationDetails: details.Data.DetailLembagaKeuangan,
		MerchantID:                   details.Data.MerchantId,
		MerchantCriteria:             details.Data.MerchantCriteria,
		NMId:                         details.Data.NMId,
		Amount:                       amount,
		TipIndicator:                 details.Data.TipIndicator,
		TipValueOfFixed:              details.Data.TipValueOfFixed,
		TipValueOfPercentage:         details.Data.TipValueOfPercentage,
	}, nil
}

func (q *QRIS) Pay(ctx context.Context, data *entity.QRISPaymentData) (*entity.QRISPaymentResult, error) {
	res, err := q.client.Payment(ctx, qris.PaymentRequest{
		NoRekNasabah:     data.AccountNumber,
		QRCode:           data.QRCode,
		RRN:              data.RRN,
		AmountPay:        float64(data.Amount),
		Tips:             float64(data.Tip),
		LembagaKeuangan:  data.FinancialOrganisation,
		CustomerName:     data.CustomerName,
		MerchantId:       data.MerchantId,
		MerchantCriteria: data.MerchantCriteria,
		NMId:             data.NMId,
		AccountName:      data.AccountName,
	})
	if err != nil {
		return nil, err
	}
	if res.Code != "0200" {
		return nil, fmt.Errorf("%s: %s", domain.ErrUnsuccessfulPayment, res.Description)
	}
	return &entity.QRISPaymentResult{
		Message:              res.Data.Message,
		RRN:                  res.Data.RRN,
		InvoiceNumber:        res.Data.InvoiceNumber,
		Remark:               res.Data.Remark,
		TransactionReference: res.Data.TransactionRef,
		TransactionDate:      res.Data.Detail.TransactionDate,
		TransactionStatus:    res.Data.Detail.TransactionStatus,
		AcquirerName:         res.Data.Detail.AcquirerName,
		MerchantName:         res.Data.Detail.MerchantName,
		MerchantLocation:     res.Data.Detail.MerchantLocation,
		MerchantPan:          res.Data.Detail.MerchantPan,
		TerminalId:           res.Data.Detail.TerminalId,
		CustomerPan:          res.Data.Detail.CustomerPan,
		ReferenceId:          res.Data.Detail.ReffId,
		Amount:               types.Money(res.Data.Detail.Amount),
		Tip:                  types.Money(res.Data.Detail.TipsAmount),
	}, nil
}
