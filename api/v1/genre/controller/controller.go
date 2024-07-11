package controller

import (
	"github.com/gin-gonic/gin"
)

type GenreController interface {
	CreateGenre(c *gin.Context)
	GetGenre(c *gin.Context)
	GetAll(c *gin.Context)
	UpdateGenre(c *gin.Context)
	DeleteGenre(c *gin.Context)
}
