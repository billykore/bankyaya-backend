package corebanking

import (
	"context"
	"fmt"

	"go.bankyaya.org/app/backend/domain/transfer"
	"go.bankyaya.org/app/backend/pkg/corebanking"
)

const transactionType = "sa-ovb-sa"

type Transfer struct {
	client *corebanking.Client
}

func NewTransfer(corebanking *corebanking.Client) *Transfer {
	return &Transfer{client: corebanking}
}

func (tf *Transfer) CheckEOD(ctx context.Context) (*transfer.EODStatus, error) {
	eod, err := tf.client.EOD(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check EOD: %w", err)
	}
	return &transfer.EODStatus{
		Code:          eod.Code,
		Description:   eod.Description,
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
	return &transfer.AccountDetails{
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

func (tf *Transfer) PerformOverbooking(ctx context.Context, req transfer.OverbookingRequest) (*transfer.OverbookingResponse, error) {
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
	return &transfer.OverbookingResponse{
		Code:                 ovb.Code,
		Description:          ovb.Description,
		JournalSequence:      ovb.JournalSequence,
		TransactionReference: ovb.TransactionReference,
		ABMsg:                ovb.ABMsg,
	}, nil
}
