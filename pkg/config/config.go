// Package config contains all the services configuration values.
// The configuration is from .env file.
package config

import (
	"sync"

	"github.com/kelseyhightower/envconfig"
	"go.bankyaya.org/app/backend/pkg/config/internal"
)

// Config holds the application configuration.
type Config struct {
	HTTP        internal.HTTP
	Token       internal.Token
	Postgres    internal.Postgres
	SQLite      internal.SQLite
	CoreBanking internal.CoreBanking
	Email       internal.Email
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
