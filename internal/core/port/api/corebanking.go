package api

import (
	"context"

	"go.bankyaya.org/app/backend/internal/core/entity"
)

// CoreBanking defines methods for core banking operations.
type CoreBanking interface {
	// CheckEOD verifies the current End-of-Day (EOD) process status in the core banking system.
	CheckEOD(ctx context.Context) (*entity.EODData, error)

	// GetAccountDetails retrieves account information for the given account number.
	GetAccountDetails(ctx context.Context, accountNumber string) (*entity.AccountDetails, error)

	// PerformOverbooking executes a transfer between two accounts with the specified amount and remark.
	// It returns an OverbookingResponse and an error if the operation fails.
	PerformOverbooking(ctx context.Context, req *entity.OverbookingRequest) (*entity.OverbookingResponse, error)
}
