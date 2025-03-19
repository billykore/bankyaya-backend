package corebanking

import (
	"context"
	"fmt"

	"go.bankyaya.org/app/backend/internal/core/transfer"
	"go.bankyaya.org/app/backend/pkg/corebanking"
)

const transactionType = "sa-ovb-sa"

type Transfer struct {
	client *corebanking.Client
}

func NewTransfer(corebanking *corebanking.Client) *Transfer {
	return &Transfer{client: corebanking}
}

func (tf *Transfer) CheckEOD(ctx context.Context) (*transfer.EODData, error) {
	eod, err := tf.client.EOD(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check EOD: %w", err)
	}
	if eod.Code != "00" {
		return nil, fmt.Errorf("core banking: %v (%v)", eod.Description, eod.Code)
	}
	return &transfer.EODData{
		SystemDate:    eod.Data.SystemDate,
		Status:        eod.Data.EodStatus,
		StandInStatus: eod.Data.StandInStatus,
	}, nil
}

func (tf *Transfer) GetAccountDetails(ctx context.Context, accountNumber string) (*transfer.AccountDetails, error) {
	inquiry, err := tf.client.Inquiry(ctx, accountNumber)
	if err != nil {
		return nil, err
	}
	if inquiry.StatusCode != "00" {
		return nil, fmt.Errorf("core banking: %v (%v)", inquiry.StatusDescription, inquiry.ErrorCode)
	}
	return &transfer.AccountDetails{
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

func (tf *Transfer) PerformOverbooking(ctx context.Context, req *transfer.OverbookingRequest) (*transfer.OverbookingResponse, error) {
	ovb, err := tf.client.Overbook(ctx, corebanking.OverbookRequest{
		TransactionType: transactionType,
		AccNoSrc:        req.SourceAccount,
		Amount:          req.Amount.String(),
		TransactionInfo: req.Remark,
		AccNoCredit:     req.DestinationAccount,
		Fee:             req.Fee.String(),
	})
	if err != nil {
		return nil, err
	}
	if ovb.Code != "00" {
		return nil, fmt.Errorf("core banking overbook failed: %v (%v)", ovb.Description, ovb.Code)
	}
	return &transfer.OverbookingResponse{
		JournalSequence:      ovb.JournalSequence,
		TransactionReference: ovb.TransactionReference,
		ABMsg:                ovb.ABMsg,
	}, nil
}
