package corebanking

import (
	"context"
	"fmt"

	"go.bankyaya.org/app/backend/internal/core/entity"
	"go.bankyaya.org/app/backend/internal/pkg/corebanking"
)

const transactionType = "sa-ovb-sa"

type CoreBanking struct {
	client *corebanking.Client
}

func New(corebanking *corebanking.Client) *CoreBanking {
	return &CoreBanking{client: corebanking}
}

func (cb *CoreBanking) CheckEOD(ctx context.Context) (*entity.EODData, error) {
	eod, err := cb.client.EOD(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check EOD: %w", err)
	}
	if eod.Code != "00" {
		return nil, fmt.Errorf("core banking: %v (%v)", eod.Description, eod.Code)
	}
	return &entity.EODData{
		SystemDate:    eod.Data.SystemDate,
		Status:        eod.Data.EodStatus,
		StandInStatus: eod.Data.StandInStatus,
	}, nil
}

func (cb *CoreBanking) GetAccountDetails(ctx context.Context, accountNumber string) (*entity.AccountDetails, error) {
	inquiry, err := cb.client.Inquiry(ctx, accountNumber)
	if err != nil {
		return nil, err
	}
	if inquiry.StatusCode != "00" {
		return nil, fmt.Errorf("core banking: %v (%v)", inquiry.StatusDescription, inquiry.ErrorCode)
	}
	return &entity.AccountDetails{
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

func (cb *CoreBanking) PerformOverbooking(ctx context.Context, req *entity.OverbookingRequest) (*entity.OverbookingResponse, error) {
	ovb, err := cb.client.Overbook(ctx, corebanking.OverbookRequest{
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
	return &entity.OverbookingResponse{
		JournalSequence:      ovb.JournalSequence,
		TransactionReference: ovb.TransactionReference,
		ABMsg:                ovb.ABMsg,
	}, nil
}
