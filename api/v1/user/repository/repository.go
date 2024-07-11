package repository

import (
	"bioskuy/api/v1/user/entity"
	"context"
	"database/sql"

	"github.com/gin-gonic/gin"
)

type UserRepository interface {
	Save(ctx context.Context, tx *sql.Tx, user entity.User, c *gin.Context) (entity.User, error)
	FindByEmail(ctx context.Context, tx *sql.Tx, email string, c *gin.Context) (entity.User, error)
	FindByID(ctx context.Context, tx *sql.Tx, email string, c *gin.Context) (entity.User, error)
	FindAll(ctx context.Context, tx *sql.Tx, c *gin.Context) ([]entity.User, error)
	UpdateToken(ctx context.Context, tx *sql.Tx, user entity.User, c *gin.Context) (entity.User, error)
	Update(ctx context.Context, tx *sql.Tx, user entity.User, c *gin.Context) (entity.User, error)
	Delete(ctx context.Context, tx *sql.Tx, id string, c *gin.Context) error
}