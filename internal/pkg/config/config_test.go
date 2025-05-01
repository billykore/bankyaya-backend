package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	cfg := Load()
	assert.NotEmpty(t, cfg)

	t.Log(cfg)
}
