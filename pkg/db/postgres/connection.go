package postgres

import (
	"context"
	"time"

	"go.bankyaya.org/app/backend/pkg/config"
	"go.bankyaya.org/app/backend/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// New returns new postgres db connection.
func New(cfg *config.Config) *gorm.DB {
	log := logger.New()
	dbCh := make(chan *gorm.DB, 1)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	go func() {
		dsn := cfg.Postgres.DSN
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("failed to connect database: %v", err)
			return
		}
		dbCh <- db
		close(dbCh)
	}()

	select {
	case <-ctx.Done():
		log.Fatalf("failed to connect database: %v", ctx.Err())
		return nil
	case db := <-dbCh:
		return db
	}
}
