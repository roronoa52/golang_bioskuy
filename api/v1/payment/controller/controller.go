package controller

import (
	"github.com/gin-gonic/gin"
)

type PaymentController interface {
	Create(c *gin.Context)
	FindById(c *gin.Context)
	FindAll(c *gin.Context)
	Notification(c *gin.Context)
}
