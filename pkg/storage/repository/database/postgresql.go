//go:build postgres || all

package database

import (
	"github.com/yeisme/notevault/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type postgresDriver struct{}

func newPostgresDriver() DBDriver {
	return &postgresDriver{}
}

func (d *postgresDriver) Connect(cfg config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.Database.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
