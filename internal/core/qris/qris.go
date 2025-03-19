package qris

import (
	"context"
	"errors"
)

var (
	// ErrEODInProgress indicates that the End of Day (EOD) process is currently in progress.
	ErrEODInProgress = errors.New("EOD process is running")

	// ErrSourceAccountInactive indicates that the source account is inactive.
	ErrSourceAccountInactive = errors.New("source account is inactive")

	// ErrSendEmailFailed is returned when an attempt to send an email fails.
	ErrSendEmailFailed = errors.New("send email failed")

	// ErrUnsuccessfulPayment is returned when a QRIS payment attempt fails.
	ErrUnsuccessfulPayment = errors.New("QRIS payment is unsuccessful")
)

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

// ReceiptMailer sends QRIS payment receipt emails.
type ReceiptMailer interface {
	// SendQRISReceipt sends a QRIS payment receipt email with the provided context and email data.
	SendQRISReceipt(ctx context.Context, data EmailData) error
}
