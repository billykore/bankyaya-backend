package dto

import (
	"strconv"

	"go.bankyaya.org/app/backend/internal/domain/intrabank"
)

type IntrabankInquiryRequest struct {
	Amount             int64  `json:"amount" validate:"required"`
	SourceAccount      string `json:"sourceAccount" validate:"required"`
	DestinationAccount string `json:"destinationAccount" validate:"required"`
}

func (r *IntrabankInquiryRequest) StringAmount() string {
	return strconv.Itoa(int(r.Amount))
}

func (r *IntrabankInquiryRequest) ToSequence() *intrabank.Sequence {
	return &intrabank.Sequence{
		Amount:             intrabank.Money(r.Amount),
		SourceAccount:      r.SourceAccount,
		DestinationAccount: r.DestinationAccount,
	}
}

type IntrabankInquiryResponse struct {
	SequenceNumber     string `json:"sequenceNumber"`
	SourceAccount      string `json:"sourceAccount"`
	DestinationAccount string `json:"destinationAccount"`
	Status             string `json:"status"`
}

func NewIntrabankInquiryResponse(sequence *intrabank.Sequence) *IntrabankInquiryResponse {
	return &IntrabankInquiryResponse{
		SequenceNumber:     sequence.SequenceNumber,
		SourceAccount:      sequence.SourceAccount,
		DestinationAccount: sequence.DestinationAccount,
	}
}

type IntrabankPaymentRequest struct {
	DestinationAccount string `json:"destinationAccount" validate:"required"`
	SourceAccount      string `json:"sourceAccount" validate:"required"`
	Amount             int64  `json:"amount" validate:"required"`
	Sequence           string `json:"sequence" validate:"required"`
	Notes              string `json:"notes"`
}

type IntrabankPaymentResponse struct {
	ABMsg                  []string `json:"abmsg"`
	JournalSequence        string   `json:"journalSequence"`
	DestinationAccount     string   `json:"destinationAccount"`
	DestinationAccountName string   `json:"destinationAccountName"`
	SourceAccount          string   `json:"sourceAccount"`
	Amount                 int64    `json:"amount"`
	Notes                  string   `json:"notes"`
	BankName               string   `json:"bankName"`
	TransactionReference   string   `json:"transactionReference"`
	Remark                 string   `json:"remark"`
}

func NewIntrabankPaymentResponse(transaction *intrabank.Transaction) *IntrabankPaymentResponse {
	return &IntrabankPaymentResponse{
		ABMsg:                  nil,
		JournalSequence:        transaction.SequenceJournal,
		DestinationAccount:     transaction.Destination,
		DestinationAccountName: transaction.DestinationName,
		SourceAccount:          "",
		Amount:                 int64(transaction.Amount),
		Notes:                  transaction.Note,
		BankName:               transaction.BankCode,
		TransactionReference:   transaction.TransactionReference,
		Remark:                 transaction.Remarks,
	}
}

type IntrabankGetAccountRequest struct {
	AccountNumber string `json:"accountNumber" validate:"required"`
}
