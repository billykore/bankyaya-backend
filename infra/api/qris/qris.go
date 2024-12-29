package qris

import (
	"context"
	"fmt"

	"go.bankyaya.org/app/backend/domain/qris"
	qrisclient "go.bankyaya.org/app/backend/pkg/qris"
	"go.bankyaya.org/app/backend/pkg/types"
)

type QRIS struct {
	client *qrisclient.Client
}

func NewQRIS(client *qrisclient.Client) *QRIS {
	return &QRIS{
		client: client,
	}
}

func (q *QRIS) GetDetails(ctx context.Context, accountNumber, qrCode string) (*qris.QRISData, error) {
	details, err := q.client.Inquiry(ctx, qrisclient.InquiryRequest{
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
	return &qris.QRISData{
		Status:                       details.Data.Status,
		RRN:                          details.Data.RRN,
		CustomerName:                 details.Data.CustomerName,
		DetailCustomer:               details.Data.DetailCustomer,
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

func (q *QRIS) Pay(ctx context.Context, data *qris.PaymentData) (*qris.PaymentResult, error) {
	res, err := q.client.Payment(ctx, qrisclient.PaymentRequest{
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
		return nil, fmt.Errorf("%s: %s", qris.ErrUnsuccessfulPayment, res.Description)
	}
	return &qris.PaymentResult{
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
