package service

import (
	"bioskuy/api/v1/genretomovie/dto"
	"context"

	"github.com/gin-gonic/gin"
)

type GenreToMovieService interface {
	Create(ctx context.Context, request dto.CreateGenreToMovieRequest, c *gin.Context) (dto.GenreToMovieCreateResponse, error) 
	FindByID(ctx context.Context, id string, c *gin.Context) (dto.GenreToMovieResponse, error)
	FindAll(ctx context.Context, c *gin.Context) ([]dto.GenreToMovieResponse, error)
	Delete(ctx context.Context, id string, c *gin.Context) error
}
