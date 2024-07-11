package controller

import (
	"github.com/gin-gonic/gin"
)

type MovieController interface {
	CreateMovie(c *gin.Context)
	GetMovie(c *gin.Context)
	GetAllMovies(c *gin.Context)
	UpdateMovie(c *gin.Context)
	DeleteMovie(c *gin.Context)
}
