package controller

import (
	"bioskuy/api/v1/user/dto"
	"bioskuy/api/v1/user/service"
	"bioskuy/auth"
	"bioskuy/exception"
	"bioskuy/web"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) UserController {
	return &userController{userService: userService}
}

func (ctl *userController) LoginWithGoogle(c *gin.Context) {

	url := ctl.userService.GoogleLoginHandler()
	

	c.JSON(http.StatusCreated, url)
}

func (ctl *userController) CallbackFromGoogle(c *gin.Context) {

	ctx := c.Request.Context()

	code := c.Query("code")
	_, userinfo, err := auth.GetGoogleUser(code, c)
	if err != nil {
		c.Error(exception.ValidationError{Message:"Failed to get user info"}).SetType(gin.ErrorTypePublic)
		return 
	}

	user := dto.CreateUserRequest{
		Name: userinfo.Name,
		Email: userinfo.Email,
		Role: "user",
	}

	result , err := ctl.userService.Login(ctx, user, c)
	if err != nil {
		c.Error(exception.InternalServerError{Message:err.Error()}).SetType(gin.ErrorTypePublic)
		return 
	}

	response := web.FormatResponse{
		ResponseCode: http.StatusCreated,
		Data: result,
	}

	c.JSON(http.StatusCreated, response)
}

func (ctl *userController) GetUserByID(c *gin.Context) {
	response := web.FormatResponse{}
	ctx := c.Request.Context()
	id := c.Param("userId")

	result, err := ctl.userService.FindByID(ctx, id, c)
	if err != nil {
		c.Error(exception.NotFoundError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return
	}

	response.ResponseCode = http.StatusOK
	response.Data = result

	c.JSON(http.StatusOK, response)
}

func (ctl *userController) GetAllUsers(c *gin.Context) {
	ctx := c.Request.Context()

	result, err := ctl.userService.FindAll(ctx, c)
	if err != nil {
		c.Error(exception.NotFoundError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return
	} 

	response := web.FormatResponse{
		ResponseCode: http.StatusOK,
		Data: result,
	}

		c.JSON(http.StatusOK, response)
}

func (ctl *userController) UpdateUser(c *gin.Context) {
	ctx := c.Request.Context()
	user := dto.UpdateUserRequest{}

	id := c.Param("userId")

	err := c.ShouldBind(&user)
	if err != nil {
		c.Error(exception.ForbiddenError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return
	}

	user.ID = id

	result, err := ctl.userService.Update(ctx, user, c)
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return
	}

	response := web.FormatResponse{
		ResponseCode: http.StatusOK,
		Data: result,
	}

	c.JSON(http.StatusOK, response)
}

func (ctl *userController) DeleteUser(c *gin.Context) {
	ctx := c.Request.Context()
	response := web.FormatResponse{}
	id := c.Param("userId")

	err := ctl.userService.Delete(ctx, id, c)
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return

	} else{
		response.ResponseCode = http.StatusOK
		response.Data = "OK"
	
		c.JSON(http.StatusOK, response)
	}
}
