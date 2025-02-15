package scheduler

import (
	"errors"
	"time"

	"go.bankyaya.org/app/backend/pkg/cron"
	"go.bankyaya.org/app/backend/pkg/types"
)

type CreateScheduleRequest struct {
	SakuId            int         `json:"sakuId" validate:"required"`
	Destination       string      `json:"destination" validate:"required"`
	DestinationName   string      `json:"destinationName" validate:"required"`
	Amount            types.Money `json:"amount" validate:"required"`
	Note              string      `json:"note"`
	TransactionMethod string      `json:"transactionMethod" validate:"required"`
	BankCode          string      `json:"bankCode" validate:"required"`
	BiFastCode        string      `json:"biFastCode"`
	PurposeType       string      `json:"purposeType"`
	Frequency         string      `json:"frequency" validate:"required"`
	StartDate         string      `json:"startDate" validate:"required"`
	Date              int         `json:"date"`
	Day               string      `json:"day"`
	AccountType       string      `json:"accountType" validate:"required"`
}

// ParseStartDate parses StartDate to time.Time.
func (r *CreateScheduleRequest) ParseStartDate() (time.Time, error) {
	if r.StartDate == "" {
		return time.Time{}, errors.New("empty StartDate")
	}
	date, err := time.Parse(time.DateOnly, r.StartDate)
	if err != nil {
		return time.Time{}, err
	}
	return date, nil
}

// CronFrequency gets Frequency with type of cron.Frequency.
func (r *CreateScheduleRequest) CronFrequency() cron.Frequency {
	return cron.Frequency(r.Frequency)
}

// TransactionType gets transaction type based on transaction method.
func (r *CreateScheduleRequest) TransactionType() string {
	switch r.TransactionMethod {
	case "INTERNAL":
		return "internal_transfer"
	}
	return ""
}
