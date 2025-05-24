package mq

import (
	"github.com/yeisme/notevault/internal/config"
	"github.com/zeromicro/go-zero/core/logx"
)

var currentMQType string

// MQClientInterface 定义了所有消息队列客户端应实现的接口
type MQClientInterface interface {
	Type() string
	IsAvailable() bool
}

// MQClient 统一的客户端结构体
type MQClient struct {
	Client MQClientInterface
}

type NullMQClient struct{}

func (n *NullMQClient) Type() string      { return "null" }
func (n *NullMQClient) IsAvailable() bool { return false }

// InitMQ initializes the message queue based on the provided configuration
// Note: This function will be called before using any MQ functionality
func InitMQ(mqConfig config.MQConfig) error {
	currentMQType = mqConfig.Type

	switch mqConfig.Type {
	case "nats":
		initNats(mqConfig.NATS)
		logx.Infof("NATS messaging system configuration applied")
	case "kafka":
		logx.Slow("Kafka messaging system support is planned but not yet implemented")
		//TODO: initKafka(mqConfig.Kafka) // 未来实现
	case "rabbitmq":
		logx.Slow("RabbitMQ messaging system support is planned but not yet implemented")
		//TODO: initRabbitMQ(mqConfig.RabbitMQ) // 未来实现
	case "":
		logx.Info("No message queue type specified, skipping MQ initialization")
	default:
		logx.Slowf("Unsupported message queue type: %s, continuing without MQ", mqConfig.Type)
	}

	return nil
}

// GetMQClient returns the appropriate MQ client
func GetMQClient() *MQClient {
	switch currentMQType {
	case "nats":
		conn, err := getNATSConn()
		if err != nil {
			logx.Debugf("Failed to get NATS connection: %v", err)
			return &MQClient{Client: &NullMQClient{}}
		}

		js, err := getJetStreamContext()
		if err != nil {
			logx.Debugf("No JetStream context available: %v", err)
		}
		natsClient := &NATSClient{
			Conn: conn,
			JS:   js,
		}
		return &MQClient{Client: natsClient}
	default:
		return &MQClient{Client: &NullMQClient{}}
	}
}
