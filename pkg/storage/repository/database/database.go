package database

import (
	"fmt"
	"sync"
	"time"

	"github.com/yeisme/notevault/internal/config"
	"github.com/yeisme/notevault/pkg/storage/repository/model"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

var (
	db         *gorm.DB
	gormDbOnce sync.Once
	mu         sync.RWMutex
	drivers    = make(map[string]DBDriver)

	gormConfig = &gorm.Config{
		Logger: NewLogxAdapter(),
	}
)

// registerDriver 注册数据库驱动
func registerDriver(name string, driver DBDriver) {
	mu.Lock()
	defer mu.Unlock()
	drivers[name] = driver
}

// InitDatabase 初始化数据库连接
func InitDatabase(dbConfig config.DatabaseConfig, logConfig logx.LogConf) error {
	var err error
	gormDbOnce.Do(func() {
		logx.Infof("Initializing database connection with driver: %s", dbConfig.Driver)

		mu.RLock()
		driver, ok := drivers[dbConfig.Driver]
		mu.RUnlock()

		if !ok {
			err = fmt.Errorf("unsupported database driver: %s", dbConfig.Driver)
			return
		}

		// 连接数据库
		db, err = driver.Connect(dbConfig)
		if err != nil {
			logx.Errorf("Failed to connect to database: %v", err)
			return
		}

		// 配置连接池
		sqlDB, sqlErr := db.DB()
		if sqlErr != nil {
			err = fmt.Errorf("failed to get sql.DB: %w", sqlErr)
			return
		}

		if dbConfig.MaxOpenConn > 0 {
			sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConn)
		}
		if dbConfig.MaxIdleConn > 0 {
			sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConn)
		}
		if dbConfig.MaxLifetime > 0 {
			sqlDB.SetConnMaxLifetime(time.Duration(dbConfig.MaxLifetime) * time.Second)
		}

		logx.Info("Database connection initialized successfully")
	})

	return err
}

// GetDB 获取数据库连接
func GetDB() *gorm.DB {
	if db == nil {
		logx.Error("Database connection is not initialized")
		return nil
	}
	return db
}

// AutoMigrate 自动迁移数据库
func AutoMigrate() error {
	if db == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	mu.RLock()
	defer mu.RUnlock()

	return db.AutoMigrate(
		&model.File{},
		&model.FileTag{},
		&model.FileVersion{},
		&model.Tag{},
	)
}
