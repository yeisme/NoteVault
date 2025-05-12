package database

import (
	"github.com/yeisme/notevault/internal/config"
	"gorm.io/gorm"
)

// DBDriver 定义数据库驱动接口
type DBDriver interface {
	Connect(config.DatabaseConfig) (*gorm.DB, error)
}
