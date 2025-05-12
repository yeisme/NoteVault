//go:build mysql
// +build mysql

package database

import (
	"github.com/yeisme/notevault/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	registerDriver("mysql", &MySQLDriver{})
}

type MySQLDriver struct{}

func (d *MySQLDriver) Connect(config config.Config) (*gorm.DB, error) {
	dsn := config.Database.DSN
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
