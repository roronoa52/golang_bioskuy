package service

import (
	"bioskuy/api/v1/user/dto"
	"context"

	"github.com/gin-gonic/gin"
)

type UserService interface {
	GoogleLoginHandler() string
	Login(ctx context.Context, request dto.CreateUserRequest, c *gin.Context) (dto.UserResponseLoginAndRegister, error)
	FindByEmail(ctx context.Context, email string, c *gin.Context) (dto.UserResponse, error)
	FindByID(ctx context.Context, id string, c *gin.Context) (dto.UserResponse, error)
	FindAll(ctx context.Context, c *gin.Context) ([]dto.UserResponse, error)
	Update(ctx context.Context, request dto.UpdateUserRequest, c *gin.Context) (dto.UserResponse, error)
	Delete(ctx context.Context, id string, c *gin.Context) error
}