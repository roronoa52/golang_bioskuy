package repository

import (
	"bioskuy/api/v1/showtime/entity"
	entityStudio "bioskuy/api/v1/studio/entity"
	"context"
	"database/sql"

	"github.com/gin-gonic/gin"
)

type ShowtimeRepository interface {
	Save(ctx context.Context, tx *sql.Tx, user entity.Showtime, c *gin.Context) (entity.Showtime, error)
	FindByID(ctx context.Context, tx *sql.Tx, id string, c *gin.Context) (entity.Showtime, error)
	FindAll(ctx context.Context, tx *sql.Tx, c *gin.Context) ([]entity.Showtime, error)
	Delete(ctx context.Context, tx *sql.Tx, id string, c *gin.Context) error
	FindConflictingShowtimes(ctx context.Context, tx *sql.Tx, studio entityStudio.Studio, showtime entity.Showtime, c *gin.Context) error 
}
