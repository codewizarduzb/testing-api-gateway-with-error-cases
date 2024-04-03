package api

import (
	_ "testing-api-gateway/api/docs" // swag
	v1 "testing-api-gateway/api/handlers/v1"
	casb "testing-api-gateway/api/middlware/casbin"
	"testing-api-gateway/config"
	"testing-api-gateway/pkg/logger"

	rmq "testing-api-gateway/queue/rabbitmq/producer"
	"testing-api-gateway/services"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/gin-contrib/cors"
)

// Option ...
type Option struct {
	Conf           config.Config
	Logger         logger.Logger
	ServiceManager services.IServiceManager
	Enforcer       casbin.Enforcer
	// Writer         producer.NewProducerKafkaInit
	Writer rmq.RabbitMQProducer
}

// New ...

// @Title Welcome to this page
// @Version 1.0
// @Description you can use this as social network
// @Host localhost:8080
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func New(option Option) *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{"*"}
	corsConfig.AllowBrowserExtensions = true
	corsConfig.AllowMethods = []string{"*"}
	router.Use(cors.New(corsConfig))

	handlerV1 := v1.New(&v1.HandlerV1Config{
		Logger:         option.Logger,
		ServiceManager: option.ServiceManager,
		Cfg:            option.Conf,
		Writer:         option.Writer,
	})
	router.Use(casb.NewAuthorizer())

	api := router.Group("/v1")

	api.POST("/createuser", handlerV1.CreateUser)
	api.GET("/getuser", handlerV1.GetUser)
	api.GET("/users", handlerV1.ListUsers)
	api.DELETE("/deleteuser", handlerV1.DeleteUser)

	url := ginSwagger.URL("swaggerdoc.json")
	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	return router
}
