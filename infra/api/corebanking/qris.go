package corebanking

import (
	"context"
	"fmt"

	"go.bankyaya.org/app/backend/domain/qris"
	"go.bankyaya.org/app/backend/pkg/corebanking"
)

type QRIS struct {
	client *corebanking.Client
}

func NewQRIS(client *corebanking.Client) *QRIS {
	return &QRIS{
		client: client,
	}
}

func (q *QRIS) CheckEOD(ctx context.Context) (*qris.EODStatus, error) {
	eod, err := q.client.EOD(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check EOD: %w", err)
	}
	return &qris.EODStatus{
		Code:          eod.Code,
		Description:   eod.Description,
		SystemDate:    eod.Data.SystemDate,
		Status:        eod.Data.EodStatus,
		StandInStatus: eod.Data.StandInStatus,
	}, nil
}

func (q *QRIS) GetAccountDetails(ctx context.Context, accountNumber string) (*qris.AccountDetails, error) {
	inquiry, err := q.client.Inquiry(ctx, accountNumber)
	if err != nil {
		return nil, err
	}
	if inquiry.StatusCode != "00" {
		return nil, fmt.Errorf("core banking: %v (%v)", inquiry.StatusDescription, inquiry.ErrorCode)
	}
	return &qris.AccountDetails{
		StatusCode:           inquiry.StatusCode,
		StatusDescription:    inquiry.StatusDescription,
		ErrorCode:            inquiry.ErrorCode,
		JournalSequence:      inquiry.JournalSequence,
		TransactionReference: inquiry.TransactionReference,
		AccountNumber:        inquiry.AccountData.AccountNumber,
		AccountType:          inquiry.AccountData.AccountType,
		Name:                 inquiry.AccountData.Name,
		Currency:             inquiry.AccountData.Currency,
		Status:               inquiry.AccountData.Status,
		Blocked:              inquiry.AccountData.Blocked,
		Balance:              inquiry.AccountData.Balance,
		MinBalance:           inquiry.AccountData.MinBalance,
		AvailableBalance:     inquiry.AccountData.AvailableBalance,
		CIF:                  inquiry.AccountData.CIF,
		ProductType:          inquiry.AccountData.ProductType,
	}, nil
}
