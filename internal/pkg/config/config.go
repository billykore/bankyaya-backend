// Package config contains all the services configuration values.
// The configuration is from .env file.
package config

import (
	"sync"

	"github.com/spf13/viper"
	"go.bankyaya.org/app/backend/internal/pkg/config/internal"
)

var _once sync.Once

// Configs hold the application configurations.
type Configs struct {
	App         internal.App
	Postgres    internal.Postgres
	CoreBanking internal.CoreBanking
	Email       internal.Email
	Clients     internal.Clients
	Token       internal.Token
}

type Config struct {
	Name    string
	Version string
	Configs Configs
}

// Load loads application configuration from a YAML file using Viper.
func Load() *Configs {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	var cfg Config
	_once.Do(func() {
		err := viper.Unmarshal(&cfg)
		if err != nil {
			panic(err)
		}
	})

	return &cfg.Configs
}
