package email

import (
	"context"

	"go.bankyaya.org/app/backend/pkg/entity"
)

// TransferReceiptMailer sends transfer receipt emails.
type TransferReceiptMailer interface {
	// SendTransferReceipt sends a transfer receipt email with the provided context and email data.
	SendTransferReceipt(ctx context.Context, data *entity.TransferEmailData) error
}
