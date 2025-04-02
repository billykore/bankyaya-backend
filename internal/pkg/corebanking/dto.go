package corebanking

type AuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type EODResponse struct {
	Code        string  `json:"statusCode"`
	Description string  `json:"statusDescription"`
	Data        EODData `json:"data"`
}

type EODData struct {
	SystemDate    string `json:"systemDate"`
	EodStatus     string `json:"eodStatus"`
	StandInStatus string `json:"standinStatus"`
}

type InquiryResponse struct {
	StatusCode           string       `json:"statusCode"`
	StatusDescription    string       `json:"statusDescription"`
	ErrorCode            string       `json:"errorCode"`
	JournalSequence      string       `json:"journalSequence"`
	TransactionReference string       `json:"transactionReference"`
	AccountData          *AccountData `json:"data"`
}

type AccountData struct {
	AccountNumber    string `json:"noRekening"`
	AccountType      string `json:"tipeRekening"`
	Name             string `json:"nama"`
	Currency         string `json:"mataUang"`
	Status           string `json:"status"`
	Blocked          string `json:"blokir"`
	Balance          string `json:"saldo"`
	MinBalance       string `json:"saldoMinimum"`
	AvailableBalance string `json:"saldoTersedia"`
	CIF              string `json:"cif"`
	ProductType      string `json:"tipeProduk"`
}

type OverbookRequest struct {
	TransactionType string `json:"tipeTransaksi"`
	EWalletType     string `json:"jenisEwallet"`
	AccNoSrc        string `json:"noRekeningDebet"`
	Amount          string `json:"nominal"`
	TransactionInfo string `json:"keteranganTransaksi"`
	AccNoCredit     string `json:"noRekeningCredit"`
	Fee             string `json:"biaya"`
	Provider        string `json:"provider"`
	FreeFee         string `json:"bebasBiaya"`
}

type OverbookResponse struct {
	Code                 string   `json:"statusCode"`
	Description          string   `json:"statusDescription"`
	JournalSequence      string   `json:"journalSequence"`
	TransactionReference string   `json:"transactionReference"`
	ABMsg                []string `json:"abmsg"`
}
