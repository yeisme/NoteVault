//go:build sqlite3
// +build sqlite3

package database

import (
	"github.com/yeisme/notevault/internal/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func init() {
    registerDriver("sqlite3", &SQLiteDriver{})
}

type SQLiteDriver struct{}

func (d *SQLiteDriver) Connect(config config.Config) (*gorm.DB, error) {
    dbPath := config.Database.DSN
    return gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
}
