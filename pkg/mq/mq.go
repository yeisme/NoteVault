package mq

import (
	"fmt"

	"github.com/yeisme/notevault/internal/config"
	"github.com/zeromicro/go-zero/core/logx"
)

var currentMQType string

// InitMQ initializes the message queue based on the provided configuration
// Note: This function will be called before using any MQ functionality
func InitMQ(mqConfig config.MQConfig) error {
	currentMQType = mqConfig.Type

	switch mqConfig.Type {
	case "nats":
		initNats(mqConfig.NATS)
	default:
		logx.Slow("No message queue specified or unsupported type, continuing without MQ")
	}

	return nil
}

// GetMQClient returns the appropriate MQ client based on what was initialized
// if the MQ type is not supported or not initialized, it returns an error and nil client
func GetMQClient() (any, error) {
	switch currentMQType {
	case "nats":
		conn, err := GetNATSConn()
		if err != nil {
			return nil, fmt.Errorf("NATS client was configured but not properly initialized: %w", err)
		}
		return conn, nil
	default:
		return nil, fmt.Errorf("no message queue client available, type configured: %s", currentMQType)
	}
}

// Close closes any active message queue connections
func Close() {
	// Close connections based on the current MQ type
	switch currentMQType {
	case "nats":
		if natsConn != nil {
			natsConn.Close()
			logx.Info("NATS connection closed")
		}
	default:
		logx.Info("No active message queue connection to close")
	}
}
