package controller

import (
	"github.com/gin-gonic/gin"
)

type SeatController interface {
	FindById(c *gin.Context)
	FindAll(c *gin.Context)
}
