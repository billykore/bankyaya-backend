package entity

import (
	"go.bankyaya.org/app/backend/internal/pkg/types"
)

type QRISData struct {
	SourceAccount                string
	Status                       bool
	RRN                          string
	CustomerName                 string
	CustomerDetail               string
	FinancialOrganisation        string
	FinancialOrganisationDetails string
	MerchantID                   string
	MerchantCriteria             string
	NMId                         string
	Amount                       types.Money
	TipIndicator                 string
	TipValueOfFixed              string
	TipValueOfPercentage         string
	Fee                          types.Money
	Note                         string
	QRCode                       string
	Tip                          types.Money
}

type QRISPaymentData struct {
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

type QRISPaymentResult struct {
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

func (res *QRISPaymentResult) TotalAmount() types.Money {
	return res.Amount + res.Tip
}

type QRISEmailData struct {
	Subject        string
	Recipient      string
	Amount         types.Money
	Fee            types.Money
	SourceName     string
	SourceAccount  string
	MerchantName   string
	MerchantPan    string
	TransactionRef string
	Note           string
}
