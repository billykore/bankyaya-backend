package intrabank

import "context"

// ReceiptMailer sends transfer receipt emails.
type ReceiptMailer interface {
	// SendReceipt sends a transfer receipt email with the provided context and email data.
	SendReceipt(ctx context.Context, data *EmailData) error
}
