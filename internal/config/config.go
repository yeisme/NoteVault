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
		MQ      MQConfig
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

	MQConfig struct {
		Type string
		NATS NATSConfig
	}

	NATSConfig struct {
		URL                string   // NATS服务器URL，例如: nats://localhost:4222
		Cluster            *string  `json:"cluster,optional"`
		ClientID           string   // 客户端ID，可选
		QueueGroup         string   // 队列组名称，用于负载均衡
		ConnectTimeout     int      // 连接超时时间(秒)
		MaxReconnects      int      // 最大重连次数
		ReconnectWait      int      // 重连等待时间(秒)
		Servers            []string // 备用服务器列表
		UseCredentials     bool     // 是否使用凭证文件认证
		CredentialsFile    string   // 凭证文件路径
		UseToken           bool     // 是否使用Token认证
		Token              string   // 认证Token
		UseUserCredentials bool     // 是否使用用户名密码认证
		User               string   // 用户名
		Password           string   // 密码
		EnableTLS          bool     // 启用TLS
		TLSCert            string   // TLS证书路径
		TLSKey             string   // TLS密钥路径
		TLSCaCert          string   // TLS CA证书路径
	}
)
