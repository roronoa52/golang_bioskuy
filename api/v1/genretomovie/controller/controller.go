package controller

import (
	"github.com/gin-gonic/gin"
)

type GenretomovieController interface {
	Create(c *gin.Context)
	FindById(c *gin.Context)
	FindAll(c *gin.Context)
	Delete(c *gin.Context)
}
