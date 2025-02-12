package transfer

import (
	"fmt"
	"strconv"

	"go.bankyaya.org/app/backend/pkg/types"
)

type InquiryRequest struct {
	Amount             types.Money `json:"amount" validate:"required"`
	SourceAccount      string      `json:"sourceAccount" validate:"required"`
	DestinationAccount string      `json:"destinationAccount" validate:"required"`
}

func (r InquiryRequest) StringAmount() string {
	return strconv.Itoa(int(r.Amount))
}

type InquiryResponse struct {
	SequenceNumber     string `json:"sequenceNumber"`
	SourceAccount      string `json:"sourceAccount"`
	DestinationAccount string `json:"destinationAccount"`
	Status             string `json:"status"`
}

type PaymentRequest struct {
	DestinationAccount string      `json:"destinationAccount" validate:"required"`
	SourceAccount      string      `json:"sourceAccount" validate:"required"`
	Amount             types.Money `json:"amount" validate:"required"`
	Sequence           string      `json:"sequence" validate:"required"`
	Notes              string      `json:"notes"`
}

// Remark creates remark for the transfer payment.
func (r *PaymentRequest) Remark() string {
	return fmt.Sprintf("TRFAG %v %v RAYA %v %v",
		r.SourceAccount,
		r.DestinationAccount,
		r.Sequence,
		r.Notes,
	)
}

type PaymentResponse struct {
	ABMsg                  []string    `json:"abmsg"`
	JournalSequence        string      `json:"journalSequence"`
	DestinationAccount     string      `json:"destinationAccount"`
	DestinationAccountName string      `json:"destinationAccountName"`
	SourceAccount          string      `json:"sourceAccount"`
	Amount                 types.Money `json:"amount"`
	Notes                  string      `json:"notes"`
	BankName               string      `json:"bankName"`
	TransactionReference   string      `json:"transactionReference"`
	Remark                 string      `json:"remark"`
}
