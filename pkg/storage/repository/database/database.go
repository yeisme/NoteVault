package database

import (
	"fmt"
	"sync"

	"github.com/yeisme/notevault/internal/config"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

// 声明为全局变量，包含互斥锁保证并发安全
var (
	DB     *gorm.DB
	dbOnce sync.Once
	dbLock sync.RWMutex

	// 驱动注册表
	drivers = make(map[string]DBDriver)
)

// GetDB 安全地获取数据库连接
func GetDB() *gorm.DB {
	dbLock.RLock()
	defer dbLock.RUnlock()
	return DB
}

// registerDriver 注册数据库驱动
func registerDriver(name string, driver DBDriver) {
	drivers[name] = driver
}


// InitDatabase 初始化数据库连接
// config 中应该包含数据库的配置参数
func InitDatabase(config config.Config) error {
	var err error
	dbOnce.Do(func() {
		driver, findErr := NewDatabaseDriver(config.Database.Driver)
		if findErr != nil {
			err = findErr
			return
		}

		db, connectErr := driver.Connect(config)
		if connectErr != nil {
			err = connectErr
			return
		}

		dbLock.Lock()
		DB = db
		dbLock.Unlock()
	})
	return err
}

// NewDatabaseDriver 创建数据库驱动实例
func NewDatabaseDriver(driverName string) (DBDriver, error) {
	driver, ok := drivers[driverName]
	if !ok {
		logx.Infof("Available drivers: %v", drivers)
		return nil, fmt.Errorf("unsupported database driver: %s", driverName)
	}
	return driver, nil
}
