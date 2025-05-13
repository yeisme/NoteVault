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

func (d *MySQLDriver) Connect(databaseConfig config.DatabaseConfig) (*gorm.DB, error) {
	dsn := databaseConfig.DSN
	return gorm.Open(mysql.Open(dsn), gormConfig)
}
