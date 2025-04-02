package email

import (
	"context"

	"go.bankyaya.org/app/backend/internal/core/entity"
)

// QRISReceiptMailer sends QRIS payment receipt emails.
type QRISReceiptMailer interface {
	// SendQRISReceipt sends a QRIS payment receipt email with the provided context and email data.
	SendQRISReceipt(ctx context.Context, data entity.QRISEmailData) error
}
