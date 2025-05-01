package version

import (
	"strconv"
	"strings"
)

type Version int64

func NewVersion(version string) Version {
	v := strings.ReplaceAll(version, ".", "")
	i, err := strconv.Atoi(v)
	if err != nil {
		return 0
	}
	return Version(i)
}

func (v Version) Equal(version Version) bool {
	return v == version
}

func (v Version) LessThanOrEqual(version Version) bool {
	return v <= version
}

func (v Version) GreaterThanOrEqual(version Version) bool {
	return v >= version
}

func (v Version) Between(low, high Version) bool {
	return v.GreaterThanOrEqual(low) && v.LessThanOrEqual(high)
}
