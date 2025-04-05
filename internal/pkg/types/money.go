package types

import "strconv"

// Money represents a monetary value stored as an unsigned 64-bit integer.
type Money uint64

func (m Money) String() string {
	return strconv.Itoa(int(m))
}

func ParseMoney(s string) (Money, error) {
	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return Money(i), nil
}
