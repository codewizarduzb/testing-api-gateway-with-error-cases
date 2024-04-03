package services

import (
	"testing-api-gateway/config"
	pbu "testing-api-gateway/genproto/user"
	mockservice "testing-api-gateway/mock-service"

	"google.golang.org/grpc/resolver"
)

type IServiceManager interface {
	UserService() pbu.UserServiceClient
	MockService() mockservice.MockUserServiceClient
}

type serviceManager struct {
	userService pbu.UserServiceClient
	mockService mockservice.MockUserServiceClient
}

func (s *serviceManager) UserService() pbu.UserServiceClient {
	return s.userService
}
func (s *serviceManager) MockService() mockservice.MockUserServiceClient {
	return s.mockService
}

func NewServiceManager(conf *config.Config) (IServiceManager, error) {
	resolver.SetDefaultScheme("dns")

	serviceManager := &serviceManager{
		mockService: mockservice.NewUserMockServiceClient(),
	}

	return serviceManager, nil
}
