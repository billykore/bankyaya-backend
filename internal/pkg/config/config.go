// Package config contains all the service configuration values.
// The configuration is from the config.yaml file.
package config

import (
	"sync"

	"github.com/spf13/viper"
	"go.bankyaya.org/app/backend/internal/pkg/config/internal"
	"go.bankyaya.org/app/backend/internal/pkg/logger"
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
	log := logger.New()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../../../.")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("failed to read config file: %v", err)
	}

	var cfg Config
	_once.Do(func() {
		err := viper.Unmarshal(&cfg)
		if err != nil {
			log.Fatalf("failed to unmarshal config: %v", err)
		}
	})

	return &cfg.Configs
}
