// default use the postgres driver, even if the driver is not specified, it will use the postgres driver

package database

import (
	"github.com/yeisme/notevault/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	registerDriver("postgres", &PostgresDriver{})
}

type PostgresDriver struct{}

func (d *PostgresDriver) Connect(databaseConfig config.DatabaseConfig) (*gorm.DB, error) {
	dsn := databaseConfig.DSN
	return gorm.Open(postgres.Open(dsn), gormConfig)
}
