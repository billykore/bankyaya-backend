package postgres

import (
	"go.bankyaya.org/app/backend/internal/pkg/config"
	"go.bankyaya.org/app/backend/internal/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// New returns new postgres db connection.
func New(cfg *config.Configs) *gorm.DB {
	log := logger.New()
	dsn := cfg.Postgres.DSN
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	return db
}
