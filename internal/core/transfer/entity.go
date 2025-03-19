package transfer

import (
	"fmt"
	"time"

	"go.bankyaya.org/app/backend/pkg/types"
)

// Sequence represents transfer sequence table.
type Sequence struct {
	ID              int       `gorm:"column:ID;primaryKey"`
	SeqNo           string    `gorm:"column:SEQ_NO"`
	Amount          string    `gorm:"column:AMOUNT"`
	AccNoSrc        string    `gorm:"column:ACC_NO_SRC"`
	AccNoDest       string    `gorm:"column:ACC_NO_DEST"`
	AccNameSrc      string    `gorm:"column:ACC_NAME_SRC"`
	AccNameDest     string    `gorm:"column:ACC_NAME_DEST"`
	TransactionType string    `gorm:"column:TRANSACTION_TYPE"`
	CifDest         string    `gorm:"column:CIF_DEST"`
	CreateDate      time.Time `gorm:"column:CREATE_DATE"`
}

func (*Sequence) TableName() string {
	return "_sequence_trf"
}

// Remark returns the remark for the transfer sequence.
func (seq *Sequence) Remark() string {
	return fmt.Sprintf("TRFAG %v %v RAYA %v",
		seq.AccNameSrc,
		seq.AccNameDest,
		seq.SeqNo,
	)
}

// User represents user table.
type User struct {
	ID            int       `gorm:"column:ID"`
	CIF           string    `gorm:"column:CIF"`
	AccountNumber string    `gorm:"column:ACCNO"`
	FullName      string    `gorm:"column:FULL_NAME"`
	Email         string    `gorm:"column:EMAIL"`
	PhoneNumber   string    `gorm:"column:PHONE_NUMBER"`
	IdentityNo    string    `gorm:"column:KTP_NUMBER"`
	CreateDate    time.Time `gorm:"column:CREATE_DATE"`
}

func (*User) TableName() string {
	return "_users"
}

type Transaction struct {
	ID                      int64     `gorm:"column:ID,pk"`
	UUID                    string    `gorm:"column:UUID"`
	UserID                  string    `gorm:"column:USER_ID"`
	WalletIDSource          int64     `gorm:"column:WALLET_ID_SOURCE"`
	Destination             string    `gorm:"column:DESTINATION"`
	Amount                  string    `gorm:"column:AMOUNT"`
	TransactionType         string    `gorm:"column:TRANSACTION_TYPE"`
	TransactionReference    string    `gorm:"column:TRREFN"`
	SequenceJournal         string    `gorm:"column:SEQUENCE_JOURNAL"`
	Remarks                 string    `gorm:"column:REMARKS"`
	Note                    string    `gorm:"column:NOTE"`
	CoreRequestPayload      string    `gorm:"column:CORE_REQUEST_PAYLOAD" json:"-"`
	CoreResponsePayload     string    `gorm:"column:CORE_RESPONSE_PAYLOAD" json:"-"`
	EChannelRequestPayload  string    `gorm:"column:E_CHANNEL_REQUEST_PAYLOAD" json:"-"`
	EChannelResponsePayload string    `gorm:"column:E_CHANNEL_RESPONSE_PAYLOAD" json:"-"`
	Status                  string    `gorm:"column:STATUS"`
	CreatedAt               time.Time `gorm:"column:CREATED_AT"`
	Fee                     string    `gorm:"column:FEE"`
	NameDest                string    `gorm:"column:DEST_NAME"`
	InitialSourceBalance    float64   `gorm:"column:INITIAL_SOURCE_BALANCE"`
	StatusCode              string    `gorm:"column:STATUS_CODE"`
	SequenceNumber          string    `gorm:"column:SEQ_NO"`
	BankCode                string    `gorm:"column:BANK_CODE"`
	SuccessTrxDate          string    `gorm:"column:SUCCESS_TRX_DATE"`
}

func (*Transaction) TableName() string {
	return "_transactions"
}

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

type Event struct {
	ScheduleId    int         `json:"scheduleId"`
	Destination   string      `json:"destination"`
	Amount        types.Money `json:"amount"`
	AccountNumber string      `json:"accountNumber"`
	UserId        int         `json:"userId"`
	Notes         string      `json:"notes"`
	BankCode      string      `json:"bankCode"`
	Status        string      `json:"status"`
	PhoneNumber   string      `json:"phoneNumber"`
	DeviceId      string      `json:"deviceId"`
}
