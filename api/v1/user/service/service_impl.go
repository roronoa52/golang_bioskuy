package service

import (
	"bioskuy/api/v1/user/dto"
	"bioskuy/api/v1/user/entity"
	"bioskuy/api/v1/user/repository"
	"bioskuy/auth"
	"bioskuy/exception"
	"bioskuy/helper"
	"context"
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type userService struct {
	Repo     repository.UserRepository
	Validate *validator.Validate
	DB       *sql.DB
	Jwt      auth.Auth
}

func NewUserService(repo repository.UserRepository, validate *validator.Validate, DB *sql.DB, jwt auth.Auth) UserService {
	return &userService{Repo: repo, Validate: validate, DB: DB, Jwt: jwt}
}

func (s *userService) GoogleLoginHandler() string {
	url := auth.GetGoogleLoginURL("state")

	return url
}

func (s *userService) Login(ctx context.Context, request dto.CreateUserRequest, c *gin.Context) (dto.UserResponseLoginAndRegister, error) {
	var UserResponse = dto.UserResponseLoginAndRegister{}

	err := s.Validate.Struct(request)
	if err != nil {
		c.Error(exception.ValidationError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return UserResponse, err
	}

	tx, err := s.DB.Begin()
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return UserResponse, err
	}
	defer helper.CommitAndRollback(tx, c)

	user := entity.User{
		Name:  request.Name,
		Email: request.Email,
	}

	var result entity.User

	result, err = s.Repo.FindByEmail(ctx, tx, request.Email, c)

	if err != nil {
		result, err = s.Repo.Save(ctx, tx, user, c)
		if err != nil {
			c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
			return UserResponse, err
		}

	}

	fmt.Println(result)

	Token, err := s.Jwt.GenerateToken(result, c)
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return UserResponse, err
	}

	user.ID = result.ID
	user.Token = Token

	user, err = s.Repo.UpdateToken(ctx, tx, user, c)
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return UserResponse, err
	}

	UserResponse.Token = Token

	return UserResponse, nil
}

func (s *userService) FindByEmail(ctx context.Context, email string, c *gin.Context) (dto.UserResponse, error) {

	UserResponse := dto.UserResponse{}

	tx, err := s.DB.Begin()
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return UserResponse, err
	}
	defer helper.CommitAndRollback(tx, c)

	result, err := s.Repo.FindByEmail(ctx, tx, email, c)
	if err != nil {
		c.Error(exception.NotFoundError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return UserResponse, err
	}

	UserResponse.ID = result.ID
	UserResponse.Name = result.Name
	UserResponse.Email = result.Email
	UserResponse.Role = result.Role

	return UserResponse, nil
}

func (s *userService) FindByID(ctx context.Context, id string, c *gin.Context) (dto.UserResponse, error) {

	UserResponse := dto.UserResponse{}

	tx, err := s.DB.Begin()
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return UserResponse, err
	}
	defer helper.CommitAndRollback(tx, c)

	result, err := s.Repo.FindByID(ctx, tx, id, c)
	if err != nil {
		c.Error(exception.NotFoundError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return UserResponse, err
	}

	UserResponse.ID = result.ID
	UserResponse.Name = result.Name
	UserResponse.Email = result.Email
	UserResponse.Role = result.Role

	return UserResponse, nil
}

func (s *userService) FindAll(ctx context.Context, c *gin.Context) ([]dto.UserResponse, error) {
	UserResponses := []dto.UserResponse{}

	tx, err := s.DB.Begin()
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return UserResponses, err
	}
	defer helper.CommitAndRollback(tx, c)

	result, err := s.Repo.FindAll(ctx, tx, c)
	if err != nil {
		c.Error(exception.NotFoundError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return UserResponses, err
	}

	for _, category := range result {
		UserResponse := dto.UserResponse{}

		UserResponse.ID = category.ID
		UserResponse.Name = category.Name
		UserResponse.Email = category.Email
		UserResponse.Role = category.Role

		UserResponses = append(UserResponses, UserResponse)

	}

	return UserResponses, nil
}

func (s *userService) Update(ctx context.Context, request dto.UpdateUserRequest, c *gin.Context) (dto.UserResponse, error) {
	UserResponse := dto.UserResponse{}
	var user entity.User

	err := s.Validate.Struct(request)
	if err != nil {
		c.Error(exception.ValidationError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return UserResponse, err
	}

	resultCustomer, err := s.FindByID(ctx, request.ID, c)
	if err != nil {
		c.Error(exception.NotFoundError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return UserResponse, err
	}

	tx, err := s.DB.Begin()
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return UserResponse, err
	}
	defer helper.CommitAndRollback(tx, c)

	user.ID = resultCustomer.ID
	user.Role = resultCustomer.Role

	if request.Role != "" {
		user.Role = request.Role
	}

	result, err := s.Repo.Update(ctx, tx, user, c)
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return UserResponse, err
	}

	UserResponse.ID = result.ID
	UserResponse.Name = resultCustomer.Name
	UserResponse.Email = resultCustomer.Email
	UserResponse.Role = result.Role

	return UserResponse, nil
}

func (s *userService) Delete(ctx context.Context, id string, c *gin.Context) error {
	customer := entity.User{}

	tx, err := s.DB.Begin()
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return err
	}
	defer helper.CommitAndRollback(tx, c)

	resultUser, err := s.Repo.FindByID(ctx, tx, id, c)
	if err != nil {
		c.Error(exception.NotFoundError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return err
	}

	customer.ID = resultUser.ID

	err = s.Repo.Delete(ctx, tx, id, c)
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return err
	}

	return nil
}
