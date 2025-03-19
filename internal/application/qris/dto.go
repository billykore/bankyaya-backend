package qris

import "go.bankyaya.org/app/backend/pkg/types"

type InquiryRequest struct {
	SourceAccount string      `json:"sourceAccount,omitempty" validate:"required"`
	QRCode        string      `json:"qrCode,omitempty" validate:"required"`
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
	QRCode                string      `json:"qrCode" validate:"required"`
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
	TransactionLogId string      `json:"transactionLogId"`
}
