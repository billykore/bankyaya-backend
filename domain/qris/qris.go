package qris

import (
	"context"
	"errors"
)

// ErrEODInProgress indicates that the End of Day (EOD) process is currently in progress.
var ErrEODInProgress = errors.New("EOD process is running")

// ErrUnsuccessfulPayment is returned when a QRIS payment attempt fails.
var ErrUnsuccessfulPayment = errors.New("QRIS payment is unsuccessful")

// CoreBanking defines methods for core banking operations.
type CoreBanking interface {
	// CheckEOD verifies the current End-of-Day (EOD) process status in the core banking system.
	CheckEOD(ctx context.Context) (*EODData, error)

	// GetAccountDetails retrieves account information for the given account number.
	GetAccountDetails(ctx context.Context, accountNumber string) (*AccountDetails, error)
}

// QRIS is an interface for QR Code Indonesian Standard operations.
type QRIS interface {
	// GetDetails retrieves QRIS data based on account number and QR code parameters.
	GetDetails(ctx context.Context, accountNumber, qrCode string) (*QRISData, error)

	// Pay processes a payment using the QRIS (Quick Response Code Indonesian Standard) system.
	Pay(ctx context.Context, data *PaymentData) (*PaymentResult, error)
}
