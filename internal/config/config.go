package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}

	Database struct {
		Driver      string // 数据库驱动类型 mysql/sqlite/postgres
		DSN         string // 数据库连接字符串
		MaxOpenConn int    // 最大连接数
		MaxIdleConn int    // 最大空闲连接数
		MaxLifetime int    // 连接最大生存时间(秒)
	}
}
