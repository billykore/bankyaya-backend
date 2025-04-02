package entity

import (
	"fmt"
	"time"

	"go.bankyaya.org/app/backend/internal/pkg/types"
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
	return fmt.Sprintf("TRF %v %v BNKYAYA %v",
		seq.AccNoSrc,
		seq.AccNoDest,
		seq.SeqNo,
	)
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
	CoreRequestPayload      string    `gorm:"column:CORE_REQUEST_PAYLOAD"`
	CoreResponsePayload     string    `gorm:"column:CORE_RESPONSE_PAYLOAD"`
	EChannelRequestPayload  string    `gorm:"column:E_CHANNEL_REQUEST_PAYLOAD"`
	EChannelResponsePayload string    `gorm:"column:E_CHANNEL_RESPONSE_PAYLOAD"`
	Status                  string    `gorm:"column:STATUS"`
	CreatedAt               time.Time `gorm:"column:CREATED_AT;default:CURRENT_TIMESTAMP()"`
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

type TransferEmailData struct {
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
