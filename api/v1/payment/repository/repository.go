package repository

import (
	"bioskuy/api/v1/payment/entity"
	"context"
	"database/sql"

	"github.com/gin-gonic/gin"
)

type PaymentRepository interface {
	Save(ctx context.Context, tx *sql.Tx, payment entity.Payment, c *gin.Context) (entity.Payment, error)
	FindByID(ctx context.Context, tx *sql.Tx, id string, c *gin.Context) (entity.Payment, error)
	Update(ctx context.Context, tx *sql.Tx, payment entity.Payment, c *gin.Context) (entity.Payment, error)
	FindAll(ctx context.Context, tx *sql.Tx, c *gin.Context) ([]entity.Payment, error)
	Delete(ctx context.Context, tx *sql.Tx, id string, c *gin.Context) error
}
