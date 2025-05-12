// 默认使用 PostgreSQL 数据库
// 需要安装 PostgreSQL 驱动
// 如果需要 mysql sqlite 及其他驱动，请在编译时指定

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
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
