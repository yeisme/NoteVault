package config

import "github.com/zeromicro/go-zero/rest"

type (
	Config struct {
		rest.RestConf
		Auth struct {
			AccessSecret string
			AccessExpire int64
		}

		Storage StorageConfig
	}

	StorageConfig struct {
		Database DatabaseConfig
		Oss      OssConfig
	}

	DatabaseConfig struct {
		Driver      string //Database driver type mysql/sqlite/postgres
		DSN         string //Database connection string
		MaxOpenConn int    //Maximum number of connections
		MaxIdleConn int    //Maximum number of idle connections
		MaxLifetime int    //Maximum survival time (seconds)
	}

	OssConfig struct {
		Endpoint        string
		AccessKeyID     string
		SecretAccessKey string
		BucketName      string
		Region          string
	}
)
