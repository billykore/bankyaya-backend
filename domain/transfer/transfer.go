package transfer

import (
	"time"
)

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
