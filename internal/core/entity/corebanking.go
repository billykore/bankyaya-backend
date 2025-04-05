package entity

import (
	"go.bankyaya.org/app/backend/internal/pkg/types"
)

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
