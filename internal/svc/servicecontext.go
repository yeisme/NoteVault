package svc

import (
	"github.com/yeisme/notevault/internal/config"
	"github.com/yeisme/notevault/pkg/storage/repository/database"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		DB:     database.GetDB(),
	}
}
