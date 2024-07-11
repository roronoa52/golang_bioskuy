package controller

import (
	"github.com/gin-gonic/gin"
)

type UserController interface {
	LoginWithGoogle(c *gin.Context)
	CallbackFromGoogle(c *gin.Context)
	GetUserByID(c *gin.Context)
	GetAllUsers(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
}
