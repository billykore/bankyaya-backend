package api

import (
	"context"

	"go.bankyaya.org/app/backend/internal/core/entity"
)

// QRIS is an interface for QR Code Indonesian Standard operations.
type QRIS interface {
	// GetDetails retrieves QRIS data based on account number and QR code parameters.
	GetDetails(ctx context.Context, accountNumber, qrCode string) (*entity.QRISData, error)

	// Pay processes a payment using the QRIS (Quick Response Code Indonesian Standard) system.
	Pay(ctx context.Context, data *entity.QRISPaymentData) (*entity.QRISPaymentResult, error)
}
