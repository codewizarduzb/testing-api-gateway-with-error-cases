package models

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}


type UserSwagger struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type AddPolicyRequest struct {
	Role     string `json:"role"`
	Endpoint string `json:"endpoint"`
	Method   string `json:"method"`
  }
  