package datetime

import "time"

// DefaultTimeLayout is the standard project time layout.
const DefaultTimeLayout = "2006-01-02T15:04:05"

var IndonesianDayNames = []string{
	"Minggu",
	"Senin",
	"Selasa",
	"Rabu",
	"Kamis",
	"Jumat",
	"Sabtu",
}

// IndonesianWeekdayValue gets Indonesian day
// and return integer representation of the day.
// The integer value are similar to Go time.Weekday.
func IndonesianWeekdayValue(day string) int {
	for i, v := range IndonesianDayNames {
		if v == day {
			return i
		}
	}
	return 0
}

// IndonesianWeekdayName gets an integer
// and return Indonesian name of the day that represented by the integer.
// The integer value are similar to Go time.Weekday.
func IndonesianWeekdayName(val int) string {
	if int(time.Sunday) <= val && val <= int(time.Saturday) {
		return IndonesianDayNames[val]
	}
	return ""
}

// GetLastDayOfMonth returns the last day of the month for a given time.Time
func GetLastDayOfMonth(t time.Time) time.Time {
	// First day of the next month
	nextMonth := t.AddDate(0, 1, -t.Day()+1)
	// Subtract one day to get the last day of the current month
	lastDay := nextMonth.AddDate(0, 0, -1)
	return lastDay
}

// IsLastDayOfMonth checks if today is the last day of the month.
func IsLastDayOfMonth() bool {
	now := time.Now()
	lastDay := GetLastDayOfMonth(now)
	return now.Day() == lastDay.Day()
}

// IsBeforeToday checks if the provided time is before today's date.
func IsBeforeToday(t time.Time) bool {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, t.Location())
	return t.Before(today)
}
