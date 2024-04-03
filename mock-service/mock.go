package mockservice

import (
	pbu "testing-api-gateway/genproto/user"
	"context"
)

type MockUserServiceClient interface {
	Create(ctx context.Context, in *pbu.User) (*pbu.User, error)
	Delete(ctx context.Context, in *pbu.UserRequest) (bool, error)
	Get(ctx context.Context, in *pbu.UserRequest) (*pbu.User, error)
	GetAll(ctx context.Context, in *pbu.GetAllUsersRequest) (*pbu.GetAllUsersResponse, error)
	CheckUniques(ctx context.Context, in *pbu.CheckUniquesRequest) (bool, error)
}

type mockUserServiceClient struct {
}

func NewUserMockServiceClient() MockUserServiceClient {
	return &mockUserServiceClient{}
}

// mock methods
func (u *mockUserServiceClient) Create(ctx context.Context, in *pbu.User) (*pbu.User, error) {
	return in, nil
}

func (u *mockUserServiceClient) Get(ctx context.Context, in *pbu.UserRequest) (*pbu.User, error) {
	return &pbu.User{
		Id:       in.UserId,
		Name:     "mock_Name",
		LastName: "mock_last_name",
		Username: "mock_user_name",
		Email:    "mockemail@gmail.com",
	}, nil
}

func (u *mockUserServiceClient) GetAll(ctx context.Context, in *pbu.GetAllUsersRequest) (*pbu.GetAllUsersResponse, error) {
	allUsers := pbu.GetAllUsersResponse{
		AllUsers: []*pbu.User{
			{
				Id:       "ae7ed51e-39ae-4664-a8d6-d254e98e6882",
				Name:     "mock_Name",
				LastName: "mock_last_name",
				Username: "mock_user_name",
				Email:    "mockemail@gmail.com",
			},
			{
				Id:       "ed4e2192-a915-40fc-b0a9-e355623659fb",
				Name:     "mock_Name2",
				LastName: "mock_last_name2",
				Username: "mock_user_name2",
				Email:    "mockemail2@gmail.com",
			},
			{
				Id:       "a3670b5a-0c6e-4784-9e51-45612cddfa77",
				Name:     "mock_Name3",
				LastName: "mock_last_name3",
				Username: "mock_user_name3",
				Email:    "mockemail3@gmail.com",
			},
			{
				Id:       "ce1f9f53-902c-4876-98aa-fa0f68a33d5e",
				Name:     "mock_Name4",
				LastName: "mock_last_name4",
				Username: "mock_user_name4",
				Email:    "mockemail@gmail.com",
			},
		},

	}

	return &allUsers, nil
}

func (u *mockUserServiceClient) Delete(ctx context.Context, in *pbu.UserRequest) (bool, error) {
	return true, nil
}

func (u *mockUserServiceClient) CheckUniques(ctx context.Context, in *pbu.CheckUniquesRequest) (bool, error) {
	return false, nil
}