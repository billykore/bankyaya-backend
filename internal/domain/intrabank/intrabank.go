// Package intrabank provides domain logic for handling intrabank transfers.
//
// It encapsulates business rules related to money transfers within the same bank,
// such as validating account ownership, checking balance sufficiency, and applying
// any relevant transaction policies.
//
// This package defines core entities, interfaces, and use cases that are independent
// of transport layers (e.g., HTTP, gRPC) and persistence mechanisms (e.g., database, cache).
package intrabank

import (
	"fmt"
	"strconv"
	"time"
)

const (
	SuccessStatus = "Success"
	FailedStatus  = "Failed"
)

// Money represents a monetary value stored as a 64-bit integer.
type Money int64

// String converts the Money object to its string representation as an integer.
func (m Money) String() string {
	return strconv.Itoa(int(m))
}

// Rupiah converts the Money object to its string representation as a Rupiah value.
// The currency value is formatted as a string with the currency symbol and the amount.
// For example, "Rp10.000.000" for 10000000.00.
func (m Money) Rupiah() string {
	rp := m.String()
	return fmt.Sprintf("Rp%v", rp)
}

func ParseMoney(s string) (Money, error) {
	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, ErrFailedParseMoney
	}
	return Money(i), nil
}

// Limits represent the minimum and maximum amount and daily amount limits for a transfer.
// The daily amount limit is the maximum amount that can be transferred within a 24-hour period.
// The minimum and maximum amount limits are the minimum and maximum amount that can be transferred.
type Limits struct {
	MinAmount      Money
	MaxAmount      Money
	MaxDailyAmount Money
}

// CanTransfer checks if the amount is enough and within the daily amount limit.
func (l *Limits) CanTransfer(amount Money) bool {
	return l.SufficientBalance(amount) && amount <= l.MaxDailyAmount
}

// SufficientBalance checks if the amount is within the minimum and maximum limits.
func (l *Limits) SufficientBalance(amount Money) bool {
	return amount >= l.MinAmount && amount <= l.MaxAmount
}

// Sequence represents transfer sequence.
type Sequence struct {
	ID                 int
	SequenceNumber     string
	Amount             Money
	SourceAccount      string
	DestinationAccount string
	SourceName         string
	DestinationName    string
	TransactionType    string
}

func (seq *Sequence) Valid(sequenceNumber string) bool {
	return seq.SequenceNumber == sequenceNumber
}

// Remark returns the remark for the transfer sequence.
func (seq *Sequence) Remark() string {
	return fmt.Sprintf("TRF %v %v BNKYAYA %v",
		seq.SourceAccount,
		seq.DestinationAccount,
		seq.SequenceNumber,
	)
}

// Transaction represents a transfer transaction.
// It includes the transaction details, such as the transaction reference,
// the transaction amount, the transaction fee, and the transaction status.
type Transaction struct {
	ID                      int64
	UUID                    string
	UserID                  string
	WalletIDSource          int64
	Destination             string
	Amount                  Money
	TransactionType         string
	TransactionReference    string
	SequenceJournal         string
	Remarks                 string
	Note                    string
	CoreRequestPayload      string
	CoreResponsePayload     string
	EChannelRequestPayload  string
	EChannelResponsePayload string
	Status                  string
	CreatedAt               time.Time
	Fee                     string
	DestinationName         string
	InitialSourceBalance    float64
	StatusCode              string
	SequenceNumber          string
	BankCode                string
	SuccessTransactionDate  time.Time
}

// CoreStatus represents the current status of the core banking system.
// It includes the system date, overall system status, and the status
// of the stand-in processing component.
type CoreStatus struct {
	SystemDate    string
	Status        string
	StandInStatus string
}

// IsEODRunning checks if the EOD process is running and stand-in mode is not activated.
func (s *CoreStatus) IsEODRunning() bool {
	return s.Status == "STARTED" && s.StandInStatus == "N"
}

// The Account represents detailed information about a customer's bank account.
// It includes identification fields, account status, balance information,
// and customer-related metadata.
type Account struct {
	JournalSequence      string
	TransactionReference string
	AccountNumber        string
	AccountType          string
	Name                 string
	Currency             string
	Status               string
	Blocked              string
	Balance              Money
	MinBalance           Money
	AvailableBalance     Money
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

func (acc *Account) IsAccountActive() bool {
	if v, ok := accountStatus[acc.Status]; ok {
		return v
	}
	return false
}

// ABMsg is an array of strings containing transaction details from the core banking API.
type ABMsg []string

// OverbookingInput contains the required information to perform an overbooking transaction.
// This includes the source and destination accounts, the transaction amount and fee,
// as well as an optional remark for additional context.
type OverbookingInput struct {
	SourceAccount      string
	DestinationAccount string
	Amount             Money
	Fee                Money
	Remark             string
}

// OverbookingResult represents the outcome of an overbooking transaction.
// It contains details such as the journal sequence number, a transaction reference,
// and the core banking system message response.
type OverbookingResult struct {
	JournalSequence      string
	TransactionReference string
	ABMsg                ABMsg
}

// EmailData holds the information required to construct a transaction-related email notification.
// It includes sender and recipient details, transaction amounts, and metadata
// such as the transaction reference and additional notes.
type EmailData struct {
	Subject            string
	Recipient          string
	Amount             Money
	Fee                Money
	SourceName         string
	SourceAccount      string
	DestinationName    string
	DestinationAccount string
	DestinationBank    string
	TransactionRef     string
	Note               string
}

// Notification represents a transaction-related notification.
type Notification struct {
	FirebaseId  string
	Subject     string
	Amount      Money
	Destination string
	Status      string
}

// Success returns the success notification message.
func (n *Notification) String() string {
	switch n.Status {
	case SuccessStatus:
		return n.success()
	case FailedStatus:
		return n.failed()
	}
	return n.success()
}

// success returns the success notification message.
func (n *Notification) success() string {
	return fmt.Sprintf(
		"Transfer ke %s sebesar %s berhasil. Hubungi 1069 069 jika kamu tidak melakukannya.",
		n.Destination,
		n.Amount.Rupiah(),
	)
}

// failed returns the failed notification message.
func (n *Notification) failed() string {
	return fmt.Sprintf(
		"Transfer ke %s sebesar %s gagal. Hubungi 1069 069 jika kamu tidak melakukannya.",
		n.Destination,
		n.Amount.Rupiah(),
	)
}
