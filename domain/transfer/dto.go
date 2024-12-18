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

type EODStatus struct {
	Code          string
	Description   string
	SystemDate    string
	Status        string
	StandInStatus string
}

// IsRunning checks if the EOD process is running and stand-in mode is not activated.
func (eod *EODStatus) IsRunning() bool {
	return eod.Status == "STARTED" && eod.StandInStatus == "N"
}

type AccountDetails struct {
	StatusCode           string
	StatusDescription    string
	ErrorCode            string
	JournalSequence      string
	TransactionReference string
	AccountNumber        string
	AccountType          string
	Name                 string
	Currency             string
	Status               string
	Blocked              string
	Balance              string
	MinBalance           string
	AvailableBalance     string
	CIF                  string
	ProductType          string
}

var accountStatus = map[string]bool{
	"1": true,
	"4": true,
	"6": true,
	"2": false,
	"7": false,
	"9": false,
	"3": false,
}

func (ir *AccountDetails) IsAccountActive() bool {
	if v, ok := accountStatus[ir.Status]; ok {
		return v
	}
	return false
}

type OverbookingRequest struct {
	SourceAccount      string
	DestinationAccount string
	Amount             types.Money
	Fee                types.Money
	Remark             string
}

type OverbookingResponse struct {
	Code                 string   `json:"statusCode"`
	Description          string   `json:"statusDescription"`
	JournalSequence      string   `json:"journalSequence"`
	TransactionReference string   `json:"transactionReference"`
	ABMsg                []string `json:"abmsg"`
}

type EmailData struct {
	Subject            string
	Recipient          string
	Amount             types.Money
	Fee                types.Money
	SourceName         string
	SourceAccount      string
	DestinationName    string
	DestinationAccount string
	DestinationBank    string
	TransactionRef     string
	Note               string
}
