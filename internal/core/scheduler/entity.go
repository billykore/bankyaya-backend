package scheduler

import "time"

// Schedule represents transaction schedule table entity.
type Schedule struct {
	ID                 int       `gorm:"column:ID;primayKey"`
	UserId             int       `gorm:"column:USER_ID"`
	SakuId             int       `gorm:"column:WALLET_ID_SOURCE"`
	Destination        string    `gorm:"column:DESTINATION"`
	DestinationName    string    `gorm:"column:DESTINATION_NAME"`
	Amount             string    `gorm:"column:AMOUNT"`
	Note               string    `gorm:"column:NOTE"`
	BankCode           string    `gorm:"column:BANK_CODE"`
	TransactionType    string    `gorm:"column:TRANSACTION_TYPE"`
	TransactionMethod  string    `gorm:"column:TRANSACTION_METHOD"`
	TransactionPurpose string    `gorm:"column:TRANSACTION_PURPOSE"`
	Frequency          string    `gorm:"column:FREQUENCY"`
	StartDate          time.Time `gorm:"column:DATE_START"`
	AutoDebet          bool      `gorm:"column:AUTO_DEBET"`
	CrontabSchedule    string    `gorm:"column:CRONTAB_SCHEDULE"`
	Status             string    `gorm:"column:STATUS"`
	AccountType        string    `gorm:"column:ACCOUNT_TYPE"`
	BIFastCode         string    `gorm:"column:BI_FAST_CODE"`
}

func (*Schedule) TableName() string {
	return "_transactions_scheduled"
}

// StringStartDate returns string representation of the StartDate field.
func (s *Schedule) StringStartDate() string {
	return s.StartDate.Format(time.DateOnly)
}

func (s *Schedule) IsActive() bool {
	return s.Status == "active"
}

type Event struct {
	ScheduleId    int    `json:"scheduleId"`
	Destination   string `json:"destination"`
	Amount        string `json:"amount"`
	AccountNumber string `json:"accountNumber"`
	UserId        int    `json:"userId"`
	Notes         string `json:"notes"`
	BankCode      string `json:"bankCode"`
	Status        string `json:"status"`
	PhoneNumber   string `json:"phoneNumber"`
	DeviceId      string `json:"deviceId"`
}
