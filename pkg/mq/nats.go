package mq

import (
	"fmt"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/yeisme/notevault/internal/config"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
	natsConn *nats.Conn
)

// Default NATS configuration values
const (
	DefaultNATSURL           = "nats://localhost:4222"
	DefaultConnectTimeout    = 10 // seconds
	DefaultMaxReconnects     = 60
	DefaultReconnectWait     = 2 // seconds
	DefaultQueueGroup        = "notevault"
	DefaultClientName        = "NotevaultService"
)

func initNats(natsConfig config.NATSConfig) {
	logx.Infof("Initializing NATS with URL: %s", natsConfig.URL)

	// Use default values if not provided
	url := natsConfig.URL
	if url == "" {
		url = DefaultNATSURL
		logx.Infof("Using default NATS URL: %s", url)
	}

	connectTimeout := natsConfig.ConnectTimeout
	if connectTimeout <= 0 {
		connectTimeout = DefaultConnectTimeout
		logx.Infof("Using default connect timeout: %d seconds", connectTimeout)
	}

	maxReconnects := natsConfig.MaxReconnects
	if maxReconnects <= 0 {
		maxReconnects = DefaultMaxReconnects
		logx.Infof("Using default max reconnects: %d", maxReconnects)
	}

	reconnectWait := natsConfig.ReconnectWait
	if reconnectWait <= 0 {
		reconnectWait = DefaultReconnectWait
		logx.Infof("Using default reconnect wait: %d seconds", reconnectWait)
	}

	queueGroup := natsConfig.QueueGroup
	if queueGroup == "" {
		queueGroup = DefaultQueueGroup
		logx.Infof("Using default queue group: %s", queueGroup)
	}

	// Set connection options
	opts := []nats.Option{
		nats.Name(DefaultClientName), // Client name
		nats.ReconnectWait(time.Duration(reconnectWait) * time.Second),
		nats.MaxReconnects(maxReconnects),
		nats.Timeout(time.Duration(connectTimeout) * time.Second),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			logx.Errorf("NATS connection disconnected: %v", err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			logx.Infof("NATS reconnected to %s", nc.ConnectedUrl())
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			logx.Errorf("NATS connection closed: %v", nc.LastError())
		}),
	}

	// Add authentication options
	if natsConfig.UseUserCredentials {
		if natsConfig.User != "" && natsConfig.Password != "" {
			logx.Info("Using username/password authentication")
			opts = append(opts, nats.UserInfo(natsConfig.User, natsConfig.Password))
		} else {
			logx.Slow("UseUserCredentials is true but username or password is empty")
		}
	} else if natsConfig.UseToken {
		if natsConfig.Token != "" {
			logx.Info("Using token authentication")
			opts = append(opts, nats.Token(natsConfig.Token))
		} else {
			logx.Slow("UseToken is true but token is empty")
		}
	} else if natsConfig.UseCredentials {
		if natsConfig.CredentialsFile != "" {
			logx.Info("Using credentials file authentication")
			opts = append(opts, nats.UserCredentials(natsConfig.CredentialsFile))
		} else {
			logx.Slow("UseCredentials is true but credentials file path is empty")
		}
	} else {
		logx.Slow("No authentication method specified, connecting without authentication")
	}

	// Add TLS options
	if natsConfig.EnableTLS {
		logx.Info("TLS is enabled for NATS connection")
		if natsConfig.TLSCert != "" && natsConfig.TLSKey != "" {
			opts = append(opts, nats.ClientCert(natsConfig.TLSCert, natsConfig.TLSKey))
			logx.Info("Using client certificate for TLS")
		}
		if natsConfig.TLSCaCert != "" {
			opts = append(opts, nats.RootCAs(natsConfig.TLSCaCert))
			logx.Info("Using CA certificate for TLS")
		}
	}

	// Use multiple servers if provided
	var serverURL string
	if len(natsConfig.Servers) > 0 {
		serverURL = strings.Join(natsConfig.Servers, ",")
		logx.Infof("Using multiple NATS servers: %s", serverURL)
	} else {
		serverURL = url
	}

	// Connect to NATS server
	var err error
	natsConn, err = nats.Connect(serverURL, opts...)
	if err != nil {
		logx.Errorf("Failed to connect to NATS: %v", err)
		logx.Info("Continuing without NATS messaging system")
		return
	}

	logx.Infof("Successfully connected to NATS server: %s", natsConn.ConnectedUrl())
}

// GetNATSConn returns the NATS connection for direct use
func GetNATSConn() (*nats.Conn, error) {
	if natsConn == nil {
		return nil, fmt.Errorf("NATS connection not initialized")
	}
	return natsConn, nil
}
