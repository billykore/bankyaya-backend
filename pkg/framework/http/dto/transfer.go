package dto

import (
	"strconv"

	"go.bankyaya.org/app/backend/pkg/entity"
	"go.bankyaya.org/app/backend/pkg/util/types"
)

type TransferInquiryRequest struct {
	Amount             types.Money `json:"amount" validate:"required"`
	SourceAccount      string      `json:"sourceAccount" validate:"required"`
	DestinationAccount string      `json:"destinationAccount" validate:"required"`
}

func (r *TransferInquiryRequest) StringAmount() string {
	return strconv.Itoa(int(r.Amount))
}

func (r *TransferInquiryRequest) ToSequence() *entity.Sequence {
	return &entity.Sequence{
		Amount:    r.StringAmount(),
		AccNoSrc:  r.SourceAccount,
		AccNoDest: r.DestinationAccount,
	}
}

type TransferInquiryResponse struct {
	SequenceNumber     string `json:"sequenceNumber"`
	SourceAccount      string `json:"sourceAccount"`
	DestinationAccount string `json:"destinationAccount"`
	Status             string `json:"status"`
}

func NewTransferInquiryResponse(sequence *entity.Sequence) *TransferInquiryResponse {
	return &TransferInquiryResponse{
		SequenceNumber:     sequence.SeqNo,
		SourceAccount:      sequence.AccNoSrc,
		DestinationAccount: sequence.AccNoDest,
	}
}

type TransferPaymentRequest struct {
	DestinationAccount string      `json:"destinationAccount" validate:"required"`
	SourceAccount      string      `json:"sourceAccount" validate:"required"`
	Amount             types.Money `json:"amount" validate:"required"`
	Sequence           string      `json:"sequence" validate:"required"`
	Notes              string      `json:"notes"`
}

type TransferPaymentResponse struct {
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

func NewTransferPaymentResponse(transaction *entity.Transaction) *TransferPaymentResponse {
	amount, err := types.ParseMoney(transaction.Amount)
	if err != nil {
		amount = 0
	}
	return &TransferPaymentResponse{
		ABMsg:                  nil,
		JournalSequence:        transaction.SequenceJournal,
		DestinationAccount:     transaction.Destination,
		DestinationAccountName: transaction.NameDest,
		SourceAccount:          "",
		Amount:                 amount,
		Notes:                  transaction.Note,
		BankName:               transaction.BankCode,
		TransactionReference:   transaction.TransactionReference,
		Remark:                 transaction.Remarks,
	}
}
