//go:build sqlite || all

package database

import (
	"github.com/yeisme/notevault/internal/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type sqliteDriver struct{}

func newSQLiteDriver() DBDriver {
	return &sqliteDriver{}
}

func (d *sqliteDriver) Connect(cfg config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(cfg.Database.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
