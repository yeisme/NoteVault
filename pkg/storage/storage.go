package storage

import (
	"github.com/yeisme/notevault/internal/config"
	"github.com/yeisme/notevault/pkg/storage/oss"
	"github.com/yeisme/notevault/pkg/storage/repository/database"
	"github.com/zeromicro/go-zero/core/logx"
)

// InitStorage initializes the storage system
// TODO: Add some kv storage for caching
func InitStorage(storageConfig config.StorageConfig, logConfig logx.LogConf) error {
	// Initialize database connection
	if err := database.InitDatabase(storageConfig.Database, logConfig); err != nil {
		logx.Errorf("InitDatabase err: %v", err)
		return err
	}

	if err := oss.InitOss(storageConfig.Oss); err != nil {
		logx.Errorf("InitOss err: %v", err)
		return err
	}

	return nil
}
