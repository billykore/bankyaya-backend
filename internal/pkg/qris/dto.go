package qris

type Token struct {
	RefreshTokenExpiresIn string   `json:"refresh_token_expires_in"`
	ApiProductList        string   `json:"api_product_list"`
	ApiProductListJson    []string `json:"api_product_list_json"`
	OrganizationName      string   `json:"organization_name"`
	DeveloperEmail        string   `json:"developer.email"`
	TokenType             string   `json:"token_type"`
	IssuedAt              string   `json:"issued_at"`
	ClientId              string   `json:"client_id"`
	AccessToken           string   `json:"access_token"`
	ApplicationName       string   `json:"application_name"`
	Scope                 string   `json:"scope"`
	ExpiresIn             string   `json:"expires_in"`
	RefreshCount          string   `json:"refresh_count"`
	Status                string   `json:"status"`
}

type InquiryRequest struct {
	AccountNumber string `json:"accountNumber,omitempty"`
	QrCode        string `json:"qrCode,omitempty"`
}

type InquiryResponse struct {
	ResponseCode        string       `json:"responseCode"`
	ResponseDescription string       `json:"responseDescription"`
	ErrorDescription    string       `json:"errorDescription"`
	Data                *InquiryData `json:"data"`
}

type InquiryData struct {
	Status                bool   `json:"status"`
	RRN                   string `json:"rrn"`
	CustomerName          string `json:"customer_name"`
	DetailCustomer        string `json:"detail_customer"`
	LembagaKeuangan       string `json:"lembaga_keuangan"`
	DetailLembagaKeuangan string `json:"detail_lembaga_keuangan"`
	MerchantId            string `json:"merchant_id"`
	MerchantCriteria      string `json:"merchant_criteria"`
	NMId                  string `json:"nmid"`
	Amount                string `json:"amount"`
	TipIndicator          string `json:"tip_indicator"`
	TipValueOfFixed       string `json:"tip_value_of_fixed"`
	TipValueOfPercentage  string `json:"tip_value_of_percentage"`
}

type PaymentRequest struct {
	NoRekNasabah     string  `json:"norek_nasabah"`
	QRCode           string  `json:"qrCode"`
	RRN              string  `json:"rrn"`
	AmountPay        float64 `json:"amount_pay"`
	Tips             float64 `json:"tips"`
	LembagaKeuangan  string  `json:"lembaga_keuangan"`
	CustomerName     string  `json:"customer_name"`
	MerchantId       string  `json:"merchant_id"`
	MerchantCriteria string  `json:"merchant_criteria"`
	NMId             string  `json:"nmid"`
	AccountName      string  `json:"account_name"`
}

type PaymentResponse struct {
	Code             string       `json:"responseCode"`
	Description      string       `json:"responseDescription"`
	ErrorDescription string       `json:"errorDescription"`
	Data             *PaymentData `json:"data"`
}

type PaymentData struct {
	Message        string         `json:"message"`
	RRN            string         `json:"rrn"`
	InvoiceNumber  string         `json:"invoiceNumber"`
	Remark         string         `json:"remark"`
	TransactionRef string         `json:"trfRef"`
	Detail         *PaymentDetail `json:"detail"`
}

type PaymentDetail struct {
	TransactionDate   string `json:"transactionDate"`
	TransactionStatus string `json:"transactionStatus"`
	AcquirerName      string `json:"acquirerName"`
	MerchantName      string `json:"merchantName"`
	MerchantLocation  string `json:"merchantLocation"`
	MerchantPan       string `json:"merchantPan"`
	TerminalId        string `json:"terminalId"`
	CustomerPan       string `json:"customerPan"`
	ReffId            string `json:"reffId"`
	Amount            int    `json:"amount"`
	TipsAmount        int    `json:"tipsAmount"`
	SourceOfFund      string `json:"sourceOfFund"`
}
