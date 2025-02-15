package cron

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCronExprSchedule(t *testing.T) {
	dailyCron := ParseScheduleExpr(DailyFrequency, 0, 0)
	assert.Equal(t, "0 9 * * *", dailyCron)
	weeklyCron := ParseScheduleExpr(WeeklyFrequency, 0, 0)
	assert.Equal(t, "0 9 * * 0", weeklyCron)
	monthlyCron := ParseScheduleExpr(MonthlyFrequency, 0, 25)
	assert.Equal(t, "0-30 9 25 * *", monthlyCron)
	yearlyCron := ParseScheduleExpr("yearly", 0, 0)
	assert.Equal(t, "", yearlyCron)
	undefinedCron := ParseScheduleExpr("", 0, 0)
	assert.Equal(t, "", undefinedCron)
}

func TestDayFromExpr(t *testing.T) {
	day := DayFromExpr("0 9 * * 0")
	assert.Equal(t, 0, day)
	invalid := DayFromExpr("0 9 * * *")
	assert.Equal(t, -1, invalid)
}

func TestDateFromExpr(t *testing.T) {
	date := DateFromExpr("0 9 31 * *")
	assert.Equal(t, 31, date)
	invalid := DateFromExpr("")
	assert.Equal(t, 0, invalid)
	invalidDay := DayFromExpr("0 9 31 * *")
	assert.Equal(t, -1, invalidDay)
}

func TestLatestDatesCronExpr(t *testing.T) {
	cronExpr := LatestDatesCronExpr()
	assert.Len(t, cronExpr, 3)
}
