// Package config contains all the services configuration values.
// The configuration is from .env file.
package config

import (
	"sync"

	"github.com/kelseyhightower/envconfig"
	internal2 "go.bankyaya.org/app/backend/internal/pkg/config/internal"
)

// Config holds the application configuration.
type Config struct {
	HTTP        internal2.HTTP
	Token       internal2.Token
	Postgres    internal2.Postgres
	SQLite      internal2.SQLite
	CoreBanking internal2.CoreBanking
	Email       internal2.Email
	QRIS        internal2.QRIS
	Rabbit      internal2.Rabbit
}

var (
	_cfg  *Config
	_once sync.Once
)

// Get initializes and returns the singleton Config instance.
func Get() *Config {
	_once.Do(func() {
		_cfg = new(Config)
		envconfig.MustProcess("", _cfg)
	})
	return _cfg
}
