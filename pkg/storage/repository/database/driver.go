package database

import (
	"github.com/yeisme/notevault/internal/config"
	"gorm.io/gorm"
)

// 声明为全局变量，而不是类型
var DB *gorm.DB

type DBDriver interface {
	Connect(config config.Config) (*gorm.DB, error)
}

func InitDatabase(config config.Config) error {
	driver, err := NewDatabaseDriver(config.Database.Driver)
	if err != nil {
		return err
	}

	db, err := driver.Connect(config)
	if err != nil {
		return err
	}

	// 将数据库实例保存到全局变量中
	DB = db
	return nil
}
