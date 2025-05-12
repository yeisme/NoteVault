//go:build mysql || all

package database

import (
	"github.com/yeisme/notevault/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type mysqlDriver struct{}

func newMySQLDriver() DBDriver {
	return &mysqlDriver{}
}

func (d *mysqlDriver) Connect(cfg config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(cfg.Database.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
