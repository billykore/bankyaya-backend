package model

type TransactionMethod struct {
	ID                  int    `gorm:"column:ID"`
	Name                string `gorm:"column:NAME"`
	Status              string `gorm:"column:STATUS"`
	Fee                 int    `gorm:"column:FEE"`
	TransactionLimit    int    `gorm:"column:TRANSACTION_LIMIT"`
	DailyLimit          int    `gorm:"column:DAILY_LIMIT"`
	Type                string `gorm:"column:TYPE"`
	NewFeature          string `gorm:"column:NEW_FEATURE"`
	Desc1               string `gorm:"column:DESC1"`
	Desc2               string `gorm:"column:DESC2"`
	Desc3               string `gorm:"column:DESC3"`
	OpenHour            int    `gorm:"column:OPEN_HOUR"`
	CloseHour           int    `gorm:"column:CLOSE_HOUR"`
	TransactionMinLimit string `gorm:"column:TRANSACTION_MIN_LIMIT"`
	Version             int    `gorm:"column:VERSION"`
	Priority            int    `gorm:"column:PRIORITY"`
}

func (TransactionMethod) TableName() string {
	return "_transfer_methods"
}
