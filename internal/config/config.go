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
		Type string     `json:",default=nats"` // Message queue type: nats, kafka, rabbitmq
		NATS NATSConfig `json:",optional"`     // NATS configuration
	}

	NATSConfig struct {
		URL            string   `json:",default=nats://localhost:4222"` // NATS server URL
		ClientID       string   `json:",optional"`                      // Client ID for connection identification
		QueueGroup     string   `json:",default=notevault"`             // Queue group name for load balancing
		ConnectTimeout int      `json:",default=10"`                    // Connection timeout (seconds)
		MaxReconnects  int      `json:",default=60"`                    // Maximum reconnection attempts
		ReconnectWait  int      `json:",default=2"`                     // Reconnection wait time (seconds)
		Servers        []string `json:",optional"`                      // Backup server list

		UseCredentials  bool   `json:",optional"` // Whether to use credentials file authentication
		CredentialsFile string `json:",optional"` // Credentials file path

		UseToken bool   `json:",optional"` // Whether to use token authentication
		Token    string `json:",optional"` // Authentication token

		UseUserCredentials bool   `json:",optional"` // Whether to use username and password authentication
		User               string `json:",optional"` // Username
		Password           string `json:",optional"` // Password

		EnableTLS bool   `json:",optional"` // Enable TLS
		TLSCert   string `json:",optional"` // TLS certificate path
		TLSKey    string `json:",optional"` // TLS key path
		TLSCaCert string `json:",optional"` // TLS CA certificate path
		// JetStream support
		EnableJetStream bool   `json:",optional"`                    // Whether to enable JetStream
		JetStreamDomain string `json:",optional"`                    // JetStream domain
		JetStreamPrefix string `json:",optional,default=NOTEVAULT_"` // JetStream prefix
	}
)
