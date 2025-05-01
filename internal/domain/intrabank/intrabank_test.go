package intrabank

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseMoney(t *testing.T) {
	s := "100000"
	m, err := ParseMoney(s)
	assert.NoError(t, err)
	assert.Equal(t, Money(100000), m)
}

func TestParseMoneyFailed(t *testing.T) {
	s := "100000.00"
	m, err := ParseMoney(s)
	assert.Error(t, err)
	assert.Equal(t, ErrFailedParseMoney, err)
	assert.Equal(t, Money(0), m)
}
