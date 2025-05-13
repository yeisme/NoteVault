package storage

import (
	"github.com/yeisme/notevault/internal/config"
	"github.com/yeisme/notevault/pkg/storage/repository/database"
	"github.com/zeromicro/go-zero/core/logx"
)

// InitStorage 初始化所有存储相关组件
func InitStorage(storageConfig config.StorageConfig) error {
	// 初始化数据库连接
	if err := database.InitDatabase(storageConfig.Database); err != nil {
		logx.Errorf("InitDatabase err: %v", err)
		return err
	}

	// 这里可以添加其他存储初始化，如对象存储、缓存等

	return nil
}
