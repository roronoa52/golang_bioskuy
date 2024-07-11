package dto

type CreateUserRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Role  string `json:"role" validate:"required"`
}

type UpdateUserRequest struct {
	ID   string `json:"id"`
	Role string `json:"role" validate:"required"`
}

type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type UserResponseLoginAndRegister struct {
	Token string `json:"token"`
}