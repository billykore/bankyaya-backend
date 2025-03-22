package dto

import (
	"go.bankyaya.org/app/backend/pkg/entity"
	"go.bankyaya.org/app/backend/pkg/util/types"
)

type QRISInquiryRequest struct {
	SourceAccount string      `json:"sourceAccount,omitempty" validate:"required"`
	QRCode        string      `json:"qrCode,omitempty" validate:"required"`
	Amount        types.Money `json:"amount,omitempty"`
}

type QRISInquiryResponse struct {
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

func NewQRISInquiryResponse(data *entity.QRISData) *QRISInquiryResponse {
	return &QRISInquiryResponse{
		Status:                       data.Status,
		RRN:                          data.RRN,
		CustomerName:                 data.CustomerName,
		CustomerDetail:               data.CustomerDetail,
		FinancialOrganisation:        data.FinancialOrganisation,
		FinancialOrganisationDetails: data.FinancialOrganisationDetails,
		MerchantId:                   data.MerchantID,
		MerchantCriteria:             data.MerchantCriteria,
		NMId:                         data.SourceAccount,
		Amount:                       data.Amount,
		TipIndicator:                 data.TipIndicator,
		TipValueOfFixed:              data.TipValueOfFixed,
		TipValueOfPercentage:         data.TipValueOfPercentage,
		Fee:                          data.Fee,
	}
}

type QRISPaymentRequest struct {
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

func (r *QRISPaymentRequest) ToQRISData() *entity.QRISData {
	return &entity.QRISData{
		Amount:                r.Amount,
		Tip:                   r.Tip,
		SourceAccount:         r.SourceAccount,
		CustomerName:          r.CustomerName,
		CustomerDetail:        r.CustomerDetail,
		FinancialOrganisation: r.FinancialOrganisation,
		MerchantID:            r.MerchantId,
		MerchantCriteria:      r.MerchantCriteria,
		QRCode:                r.QRCode,
		RRN:                   r.RRN,
		NMId:                  r.NMId,
		Note:                  r.Note,
	}
}

type QRISPaymentResponse struct {
	Amount           types.Money `json:"amountPay"`
	Tip              types.Money `json:"tip"`
	Total            types.Money `json:"total"`
	Message          string      `json:"message"`
	RRN              string      `json:"rrn"`
	InvoiceNumber    string      `json:"invoiceNumber"`
	TransactionLogId string      `json:"transactionLogId"`
}

func NewQRISPaymentResponse(data *entity.QRISPaymentResult) *QRISPaymentResponse {
	return &QRISPaymentResponse{
		Amount:           data.Amount,
		Tip:              data.Tip,
		Total:            data.TotalAmount(),
		Message:          data.Message,
		RRN:              data.RRN,
		InvoiceNumber:    data.InvoiceNumber,
		TransactionLogId: data.TransactionReference,
	}
}
