package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewVersion(t *testing.T) {
	version := "1.0.0"
	i := NewVersion(version)
	assert.Equal(t, Version(100), i)
}

func TestEqual(t *testing.T) {
	v1 := NewVersion("1.0.0")
	v2 := NewVersion("1.0.0")
	v3 := NewVersion("1.0.1")

	assert.True(t, v1.Equal(v1))
	assert.True(t, v1.Equal(v2))
	assert.True(t, v2.Equal(v1))
	assert.True(t, v2.Equal(v2))
	assert.False(t, v1.Equal(v3))
	assert.False(t, v2.Equal(v3))
	assert.False(t, v3.Equal(v1))
	assert.False(t, v3.Equal(v2))
}

func TestLessThan(t *testing.T) {
	v100 := NewVersion("1.0.0")
	v101 := NewVersion("1.0.1")

	assert.True(t, v100.LessThanOrEqual(v100))
	assert.True(t, v100.LessThanOrEqual(v101))
	assert.False(t, v101.LessThanOrEqual(v100))
}

func TestGreaterThan(t *testing.T) {
	v100 := NewVersion("1.0.0")
	v101 := NewVersion("1.0.1")

	assert.True(t, v100.GreaterThanOrEqual(v100))
	assert.True(t, v101.GreaterThanOrEqual(v100))
	assert.False(t, v100.GreaterThanOrEqual(v101))
}

func TestBetween(t *testing.T) {
	v100 := NewVersion("1.0.0")
	v101 := NewVersion("1.0.1")
	v102 := NewVersion("1.0.2")

	assert.True(t, v101.Between(v100, v102))
	assert.True(t, v101.Between(v100, v101))
	assert.True(t, v101.Between(v101, v101))
	assert.True(t, v101.Between(v101, v102))
	assert.False(t, v101.Between(v100, v100))
	assert.False(t, v101.Between(v102, v102))
	assert.False(t, v101.Between(v102, v100))
	assert.False(t, v101.Between(v102, v101))
	assert.False(t, v101.Between(v102, v102))
}
