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
	jsCtx    nats.JetStreamContext
)

func GetNATSConn() (*nats.Conn, error) {
	if natsConn == nil {
		return nil, fmt.Errorf("NATS client was configured but not properly initialized")
	}
	return natsConn, nil
}

// DefaultClientName is the default client name for NATS connections
const DefaultClientName = "NotevaultService"

func initNats(natsConfig config.NATSConfig) {
	logx.Infof("Initializing NATS with URL: %s", natsConfig.URL)

	// Set connection options
	opts := []nats.Option{
		nats.Name(DefaultClientName), // Client name
		nats.ReconnectWait(time.Duration(natsConfig.ReconnectWait) * time.Second),
		nats.MaxReconnects(natsConfig.MaxReconnects),
		nats.Timeout(time.Duration(natsConfig.ConnectTimeout) * time.Second),
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
		serverURL = natsConfig.URL
	}

	// Connect to NATS server
	var err error
	natsConn, err = nats.Connect(serverURL, opts...)
	if err != nil {
		logx.Errorf("Failed to connect to NATS: %v", err)
		logx.Info("Continuing without NATS messaging system")
		return
	}

	// Initialize JetStream if enabled
	if natsConfig.EnableJetStream && natsConn != nil {
		logx.Info("Initializing JetStream")

		jsOpts := []nats.JSOpt{}

		if natsConfig.JetStreamDomain != "" {
			jsOpts = append(jsOpts, nats.Domain(natsConfig.JetStreamDomain))
			logx.Infof("Using JetStream domain: %s", natsConfig.JetStreamDomain)
		}

		if natsConfig.JetStreamPrefix != "" {
			jsOpts = append(jsOpts, nats.APIPrefix(natsConfig.JetStreamPrefix))
			logx.Infof("Using JetStream API prefix: %s", natsConfig.JetStreamPrefix)
		}

		// 修改: 存储JetStreamContext对象
		jsCtx, err = natsConn.JetStream(jsOpts...)
		if err != nil {
			logx.Errorf("Failed to initialize JetStream: %v", err)
		} else {
			logx.Info("JetStream initialized successfully")
		}
	}
}

func GetJetStreamContext() (nats.JetStreamContext, error) {
	if jsCtx == nil {
		return nil, fmt.Errorf("JetStream not initialized or not enabled in configuration")
	}
	return jsCtx, nil
}
