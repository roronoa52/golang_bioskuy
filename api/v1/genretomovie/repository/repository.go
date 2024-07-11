package repository

import (
	"bioskuy/api/v1/genretomovie/entity"
	"context"
	"database/sql"

	"github.com/gin-gonic/gin"
)

type GenreToMovieRepository interface {
	Save(ctx context.Context, tx *sql.Tx, user entity.GenreToMovie, c *gin.Context) (entity.GenreToMovie, error)
	FindByID(ctx context.Context, tx *sql.Tx, id string, c *gin.Context) (entity.GenreToMovie, error)
	FindAll(ctx context.Context, tx *sql.Tx, c *gin.Context) ([]entity.GenreToMovie, error)
	Delete(ctx context.Context, tx *sql.Tx, id string, c *gin.Context) error
}
