package entity

import "time"

type User struct {
	ID            int       `gorm:"column:ID;primaryKey"`
	CIF           string    `gorm:"column:CIF"`
	AccountNumber string    `gorm:"column:ACCNO"`
	FullName      string    `gorm:"column:FULL_NAME"`
	Email         string    `gorm:"column:EMAIL"`
	PhoneNumber   string    `gorm:"column:PHONE_NUMBER"`
	IdentityNo    string    `gorm:"column:KTP_NUMBER"`
	CreateDate    time.Time `gorm:"column:CREATE_DATE"`
	AuthData      AuthData  `gorm:"foreignKey:AuthDataID"`
}

func (*User) TableName() string {
	return "_users"
}

type AuthData struct {
	Id              int64     `gorm:"column:ID;primaryKey"`
	Pin             string    `gorm:"column:PIN"`
	TempPin         string    `gorm:"column:TEMP_PIN"`
	FingerPrintSign string    `gorm:"column:FINGER_PRINT_SIGN"`
	UserId          int64     `gorm:"column:USER_ID"`
	Status          int       `gorm:"column:STATUS"`
	Password        string    `gorm:"column:PASSWORD"`
	DeviceId        string    `gorm:"column:DEVICE_ID"`
	FirebaseId      string    `gorm:"column:FIREBASE_ID"`
	ActivationKey   string    `gorm:"column:ACTIVATION_KEY"`
	CreateDate      time.Time `gorm:"column:CREATE_DATE"`
	UpdateDate      time.Time `gorm:"column:UPDATE_DATE"`
	DeviceStatus    string    `gorm:"column:DEVICE_STATUS"`
	Device          Device    `gorm:"foreignKey:DeviceId;references:DEVICE_ID"`
}

func (*AuthData) TableName() string {
	return "_user_auth"
}

func (data *AuthData) ValidFirebaseId(firebaseId string) bool {
	return data.FirebaseId == firebaseId
}

func (data *AuthData) ValidDeviceId(deviceId string) bool {
	return data.DeviceId == deviceId
}

type Device struct {
	Id         int64           `gorm:"column:ID;primaryKey"`
	DeviceId   string          `gorm:"column:DEVICE_ID"`
	FirebaseId string          `gorm:"column:FIREBASE_ID"`
	UserId     int64           `gorm:"column:USER_ID"`
	Brand      string          `gorm:"column:BRAND"`
	BrandOs    string          `gorm:"column:BRAND_OS"`
	OsVersion  string          `gorm:"column:OS_VERSION"`
	Model      string          `gorm:"column:MODEL"`
	Status     string          `gorm:"column:STATUS"`
	GpsCity    string          `gorm:"column:GPS_CITY"`
	CreatedAt  time.Time       `gorm:"column:CREATED_AT"`
	UpdatedAt  time.Time       `gorm:"column:UPDATED_AT"`
	UpdatedBy  string          `gorm:"column:UPDATED_BY"`
	LastLogin  time.Time       `gorm:"column:LAST_LOGIN"`
	Blacklist  BlacklistDevice `gorm:"foreignKey:DeviceId;references:DeviceId"`
}

func (*Device) TableName() string {
	return "_user_devices"
}

// IsBlacklisted tells if device blacklist status equals 'active'.
func (device *Device) IsBlacklisted() bool {
	return device.Blacklist.Status == "active"
}

type BlacklistDevice struct {
	Id        int64     `gorm:"column:ID;primaryKey"`
	DeviceId  string    `gorm:"column:DEVICE_ID"`
	Status    string    `gorm:"column:STATUS"`
	CreatedAt time.Time `gorm:"column:CREATED_AT"`
	UpdatedAt time.Time `gorm:"column:UPDATED_AT"`
}

func (*BlacklistDevice) TableName() string {
	return "_user_blacklist_device_id"
}

type Token struct {
	AccessToken string `json:"accessToken"`
	ExpiredTime int64  `json:"expiredTime"`
}
