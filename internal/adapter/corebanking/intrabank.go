package corebanking

import (
	"context"
	"fmt"

	"go.bankyaya.org/app/backend/internal/domain/intrabank"
	"go.bankyaya.org/app/backend/internal/pkg/corebanking"
)

const (
	transactionType = "sa-ovb-sa"
	successCode     = "00"
)

type IntrabankCoreBanking struct {
	client *corebanking.Client
}

func NewIntrabankCoreBanking(corebanking *corebanking.Client) *IntrabankCoreBanking {
	return &IntrabankCoreBanking{client: corebanking}
}

func (cb *IntrabankCoreBanking) GetCoreStatus(ctx context.Context) (*intrabank.CoreStatus, error) {
	eod, err := cb.client.EOD(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check EOD: %w", err)
	}
	if eod.Code != "00" {
		return nil, fmt.Errorf("core banking: %s (%s)", eod.Description, eod.Code)
	}
	return &intrabank.CoreStatus{
		SystemDate:    eod.Data.SystemDate,
		Status:        eod.Data.EodStatus,
		StandInStatus: eod.Data.StandInStatus,
	}, nil
}

func (cb *IntrabankCoreBanking) GetAccountDetails(ctx context.Context, accountNumber string) (*intrabank.Account, error) {
	resp, err := cb.client.Inquiry(ctx, accountNumber)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != successCode {
		return nil, fmt.Errorf("core banking: %s (%s)", resp.StatusDescription, resp.ErrorCode)
	}

	balance, err := intrabank.ParseMoney(resp.AccountData.Balance)
	if err != nil {
		return nil, err
	}

	minBalance, err := intrabank.ParseMoney(resp.AccountData.MinBalance)
	if err != nil {
		return nil, err
	}

	availableBalance, err := intrabank.ParseMoney(resp.AccountData.AvailableBalance)
	if err != nil {
		return nil, err
	}

	return &intrabank.Account{
		JournalSequence:      resp.JournalSequence,
		TransactionReference: resp.TransactionReference,
		AccountNumber:        resp.AccountData.AccountNumber,
		AccountType:          resp.AccountData.AccountType,
		Name:                 resp.AccountData.Name,
		Currency:             resp.AccountData.Currency,
		Status:               resp.AccountData.Status,
		Blocked:              resp.AccountData.Blocked,
		Balance:              balance,
		MinBalance:           minBalance,
		AvailableBalance:     availableBalance,
		CIF:                  resp.AccountData.CIF,
		ProductType:          resp.AccountData.ProductType,
	}, nil
}

func (cb *IntrabankCoreBanking) PerformOverbooking(ctx context.Context, req *intrabank.OverbookingInput) (*intrabank.OverbookingResult, error) {
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
		return nil, fmt.Errorf("core banking overbook failed: %s (%s)", ovb.Description, ovb.Code)
	}
	return &intrabank.OverbookingResult{
		JournalSequence:      ovb.JournalSequence,
		TransactionReference: ovb.TransactionReference,
		ABMsg:                ovb.ABMsg,
	}, nil
}
