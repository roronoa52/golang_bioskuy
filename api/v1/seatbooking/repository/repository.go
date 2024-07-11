package repository

import (
	"bioskuy/api/v1/seatbooking/entity"
	"context"
	"database/sql"

	"github.com/gin-gonic/gin"
)

type SeatBookingRepository interface {
	Save(ctx context.Context, tx *sql.Tx, seatbooking entity.SeatBooking, c *gin.Context) (entity.SeatBooking, error)
	FindByID(ctx context.Context, tx *sql.Tx, id string, c *gin.Context) (entity.SeatBooking, error)
	FindAll(ctx context.Context, tx *sql.Tx, c *gin.Context) ([]entity.SeatBooking, error)
	FindAllPendingByUserID(ctx context.Context, tx *sql.Tx, userID string, c *gin.Context) ([]entity.SeatBooking, error) 
	Delete(ctx context.Context, tx *sql.Tx, id string, c *gin.Context) error
	Update(ctx context.Context, tx *sql.Tx, payment entity.SeatBooking, c *gin.Context) (entity.SeatBooking, error)
}
