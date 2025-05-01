package intrabank

import "context"

// CoreBanking defines methods for core banking operations.
type CoreBanking interface {
	// GetCoreStatus gets the current status of the core banking system.
	GetCoreStatus(ctx context.Context) (*CoreStatus, error)

	// GetAccountDetails retrieves account information for the given account number.
	GetAccountDetails(ctx context.Context, accountNumber string) (*Account, error)

	// PerformOverbooking executes a transfer between two accounts with the specified amount and remark.
	// It returns an OverbookingResponse and an error if the operation fails.
	PerformOverbooking(ctx context.Context, req *OverbookingInput) (*OverbookingResult, error)
}
