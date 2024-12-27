package qris

import (
	"context"

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
