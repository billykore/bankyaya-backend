package cron

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"go.bankyaya.org/app/backend/pkg/util/datetime"
)

type Frequency string

const (
	DailyFrequency   Frequency = "daily"
	WeeklyFrequency  Frequency = "weekly"
	MonthlyFrequency Frequency = "monthly"
)

// ParseScheduleExpr parses day or date to cron schedule expression.
// The minute and hour of the cron expression is set equal to 09:00.
func ParseScheduleExpr(frequency Frequency, day, date int) string {
	switch frequency {
	case DailyFrequency:
		// every day at 09:00
		return "0 9 * * *"
	case WeeklyFrequency:
		// every [day] at 09:00
		return fmt.Sprintf("0 9 * * %d", day)
	case MonthlyFrequency:
		// every [date] at 09:00-09:30
		return fmt.Sprintf("0-30 9 %d * *", date)
	}
	return ""
}

// DayFromExpr gets day from cron schedule expression.
//
//	Example: DayFromExpr("0 9 * * 0") returns 0.
func DayFromExpr(expr string) int {
	if expr == "" {
		return -1
	}
	s := strings.Split(expr, " ") // ["0", "9", "*", "*", "[day]"]
	if len(s) < 5 {
		return -1
	}
	day := s[4]
	intDay, err := strconv.Atoi(day)
	if err != nil {
		return -1
	}
	return intDay
}

// DateFromExpr gets date from cron schedule expression.
//
//	Example: DateFromExpr("0 9 31 * *") returns 31.
func DateFromExpr(expr string) int {
	if expr == "" {
		return 0
	}
	s := strings.Split(expr, " ") // ["0-30", "9", "[date]", "*", "*"]
	if len(s) < 5 {
		return 0
	}
	d := s[2]
	intDate, err := strconv.Atoi(d)
	if err != nil {
		return 0
	}
	return intDate
}

const (
	twentyNinthExpr = "0-30 9 29 * *"
	thirtiethExpr   = "0-30 9 30 * *"
	thirtyFirstExpr = "0-30 9 31 * *"
)

// LatestDatesCronExpr returns array of cron expression for latest date of month,
// that is 29, 30 and 31 by check the last date of current month.
func LatestDatesCronExpr() []string {
	now := time.Now()
	lastDayOfMonth := datetime.GetLastDayOfMonth(now)

	if lastDayOfMonth.Day() == 28 {
		return []string{twentyNinthExpr, thirtiethExpr, thirtyFirstExpr}
	}
	if lastDayOfMonth.Day() == 29 {
		return []string{thirtiethExpr, thirtyFirstExpr}
	}
	if lastDayOfMonth.Day() == 30 {
		return []string{thirtyFirstExpr}
	}

	return []string{}
}
