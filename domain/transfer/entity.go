package transfer

import (
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

// User represents user table.
type User struct {
	ID            int       `pg:"ID"`
	CIF           string    `pg:"CIF"`
	AccountNumber string    `pg:"ACCNO"`
	FullName      string    `pg:"FULL_NAME"`
	Email         string    `pg:"EMAIL"`
	PhoneNumber   string    `pg:"PHONE_NUMBER"`
	IdentityNo    string    `pg:"KTP_NUMBER"`
	CreateDate    time.Time `pg:"CREATE_DATE"`
}

func (*User) TableName() string {
	return "_users"
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
