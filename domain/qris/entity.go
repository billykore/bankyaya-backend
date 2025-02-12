package qris

import "go.bankyaya.org/app/backend/pkg/types"

type EODData struct {
	SystemDate    string
	Status        string
	StandInStatus string
}

// IsRunning checks if the EOD process is running and stand-in mode is not activated.
func (eod *EODData) IsRunning() bool {
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

// accountStatus defines status of account. true means active.
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

type QRISData struct {
	Status                       bool        `json:"status"`
	RRN                          string      `json:"rrn"`
	CustomerName                 string      `json:"customer_name"`
	DetailCustomer               string      `json:"detail_customer"`
	FinancialOrganisation        string      `json:"financialOrganisation"`
	FinancialOrganisationDetails string      `json:"financialOrganisationDetails"`
	MerchantID                   string      `json:"merchant_id"`
	MerchantCriteria             string      `json:"merchant_criteria"`
	NMId                         string      `json:"nmid"`
	Amount                       types.Money `json:"amount"`
	TipIndicator                 string      `json:"tip_indicator"`
	TipValueOfFixed              string      `json:"tip_value_of_fixed"`
	TipValueOfPercentage         string      `json:"tip_value_of_percentage"`
}

type PaymentData struct {
	AccountNumber         string
	QRCode                string
	RRN                   string
	Amount                types.Money
	Tip                   types.Money
	FinancialOrganisation string
	CustomerName          string
	MerchantId            string
	MerchantCriteria      string
	NMId                  string
	AccountName           string
}

type PaymentResult struct {
	Message              string      `json:"message"`
	RRN                  string      `json:"rrn"`
	InvoiceNumber        string      `json:"invoiceNumber"`
	Remark               string      `json:"remark"`
	TransactionReference string      `json:"transactionReference"`
	TransactionDate      string      `json:"transactionDate"`
	TransactionStatus    string      `json:"transactionStatus"`
	AcquirerName         string      `json:"acquirerName"`
	MerchantName         string      `json:"merchantName"`
	MerchantLocation     string      `json:"merchantLocation"`
	MerchantPan          string      `json:"merchantPan"`
	TerminalId           string      `json:"terminalId"`
	CustomerPan          string      `json:"customerPan"`
	ReferenceId          string      `json:"referenceId"`
	Amount               types.Money `json:"amount"`
	Tip                  types.Money `json:"tip"`
}

func (res *PaymentResult) TotalAmount() types.Money {
	return res.Amount + res.Tip
}
