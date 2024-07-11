package repository

import (
	"bioskuy/api/v1/studio/entity"
	"context"
	"database/sql"

	"github.com/gin-gonic/gin"
)

type StudioRepository interface {
	Save(ctx context.Context, tx *sql.Tx, user entity.Studio, c *gin.Context) (entity.Studio, error)
	FindByID(ctx context.Context, tx *sql.Tx, id string, c *gin.Context) (entity.Studio, error)
	FindAll(ctx context.Context, tx *sql.Tx, c *gin.Context) ([]entity.Studio, error)
	Update(ctx context.Context, tx *sql.Tx, studio entity.Studio, c *gin.Context) (entity.Studio, error)
	Delete(ctx context.Context, tx *sql.Tx, id string, c *gin.Context) error
}