package main

import (
	"testing-api-gateway/api"
	"testing-api-gateway/config"
	"testing-api-gateway/pkg/logger"
	producer "testing-api-gateway/queue/rabbitmq/producer"
	"testing-api-gateway/services"

	"github.com/casbin/casbin/v2"
	defaultrolemanager "github.com/casbin/casbin/v2/rbac/default-role-manager"
	"github.com/casbin/casbin/v2/util"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.LogLevel, "api_gateway")

	serviceManager, err := services.NewServiceManager(&cfg)
	if err != nil {
		log.Error("gRPC dial error", logger.Error(err))
	}

	enforcer, err := casbin.NewEnforcer("auth.conf")
	if err != nil {
		log.Error("NewEnforcer error", logger.Error(err))
		return
	}

	enforcer.GetRoleManager().(*defaultrolemanager.RoleManagerImpl).AddMatchingFunc("keyMatch", util.KeyMatch)
	enforcer.GetRoleManager().(*defaultrolemanager.RoleManagerImpl).AddMatchingFunc("keyMatch3", util.KeyMatch3)

	writer, err := producer.NewRabbitMQProducer("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal("error creating RabbitMQ producer", logger.Error(err))
	}

	server := api.New(api.Option{
		Conf:           cfg,
		Logger:         log,
		Enforcer:       *enforcer,
		ServiceManager: serviceManager,
		Writer:         writer,
	})
	if err := server.Run(cfg.HTTPPort); err != nil {
		log.Fatal("failed to run http server", logger.Error(err))
		panic(err)
	}
}
