package v1

import (
	"testing-api-gateway/api/handlers/tokens"
	"testing-api-gateway/config"
	"testing-api-gateway/pkg/logger"

	"testing-api-gateway/queue/kafka/producer"
	rmq "testing-api-gateway/queue/rabbitmq/producer"
	"testing-api-gateway/services"

	"github.com/casbin/casbin/v2"
)

type handlerV1 struct {
	log            logger.Logger
	serviceManager services.IServiceManager
	cfg            config.Config
	jwtHandler     tokens.JWTHandler
	enforcer       casbin.Enforcer
	kafka_writer   producer.NewProducerKafkaInit
}

// HandlerV1Config ...
type HandlerV1Config struct {
	Logger         logger.Logger
	ServiceManager services.IServiceManager
	Cfg            config.Config
	JWTHandler     tokens.JWTHandler
	KafkaWriter    producer.NewProducerKafkaInit
	Writer         rmq.RabbitMQProducer
}

// New ...
func New(c *HandlerV1Config) *handlerV1 {
	return &handlerV1{
		log:            c.Logger,
		serviceManager: c.ServiceManager,
		cfg:            c.Cfg,
		jwtHandler:     c.JWTHandler,
		kafka_writer:   c.Writer,
	}
}
