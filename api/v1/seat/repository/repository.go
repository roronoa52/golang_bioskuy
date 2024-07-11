package repository

import (
	"bioskuy/api/v1/seat/entity"
	"context"
	"database/sql"

	"github.com/gin-gonic/gin"
)

type SeatRepository interface {
	Save(ctx context.Context, tx *sql.Tx, user entity.Seat, c *gin.Context) (entity.Seat, error)
	FindByID(ctx context.Context, tx *sql.Tx, id string, c *gin.Context) (entity.Seat, error)
	FindByIDWithNotAvailable(ctx context.Context, tx *sql.Tx, id string, c *gin.Context) (entity.Seat, error)
	FindAll(ctx context.Context, id string, tx *sql.Tx, c *gin.Context) ([]entity.Seat, error)
	Update(ctx context.Context, tx *sql.Tx, seat entity.Seat, c *gin.Context) (entity.Seat, error)
	Delete(ctx context.Context, tx *sql.Tx, id string, c *gin.Context) error
}