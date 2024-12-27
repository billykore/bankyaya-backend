package qris

import "go.bankyaya.org/app/backend/pkg/types"

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

type GetDetailRequest struct {
	AccountNumber string `json:"accountNumber,omitempty"`
	QrCode        string `json:"qrCode,omitempty"`
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

type InquiryRequest struct {
	SourceAccount string      `json:"sourceAccount,omitempty" validate:"required"`
	QrCode        string      `json:"qrCode,omitempty" validate:"required"`
	Amount        types.Money `json:"amount,omitempty"`
}

type InquiryResponse struct {
	Status                       bool        `json:"status"`
	RRN                          string      `json:"rrn"`
	CustomerName                 string      `json:"customerName"`
	CustomerDetail               string      `json:"customerDetail"`
	FinancialOrganisation        string      `json:"financialOrganisation"`
	FinancialOrganisationDetails string      `json:"financialOrganisationDetails"`
	MerchantId                   string      `json:"merchantId"`
	MerchantCriteria             string      `json:"merchantCriteria"`
	NMId                         string      `json:"nmId"`
	Amount                       types.Money `json:"amount"`
	TipIndicator                 string      `json:"tipIndicator"`
	TipValueOfFixed              string      `json:"tipValueOfFixed"`
	TipValueOfPercentage         string      `json:"tipValuefPercentage"`
	Fee                          types.Money `json:"fee"`
}

type PaymentRequest struct {
	Amount                types.Money `json:"amount" validate:"required"`
	Tip                   types.Money `json:"tip"`
	SourceAccount         string      `json:"accountNoSource" validate:"required"`
	CustomerName          string      `json:"customerName" validate:"required"`
	CustomerDetail        string      `json:"customerDetail" validate:"required"`
	FinancialOrganisation string      `json:"financialOrganisation" validate:"required"`
	MerchantId            string      `json:"merchantId" validate:"required"`
	MerchantCriteria      string      `json:"merchantCriteria" validate:"required"`
	QrCode                string      `json:"qrCode" validate:"required"`
	RRN                   string      `json:"rrn" validate:"required"`
	NMId                  string      `json:"nmid"`
	Note                  string      `json:"note"`
}

type PaymentResponse struct {
	Amount           types.Money `json:"amountPay"`
	Tip              types.Money `json:"tip"`
	Total            types.Money `json:"total"`
	Message          string      `json:"message"`
	RRN              string      `json:"rrn"`
	InvoiceNumber    string      `json:"invoiceNumber"`
	CustomerName     string      `json:"customerName"`
	CustomerCity     string      `json:"customerCity"`
	Note             string      `json:"note"`
	TransactionLogId string      `json:"transactionLogId"`
}
