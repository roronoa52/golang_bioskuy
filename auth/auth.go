package auth

import (
	"bioskuy/api/v1/user/entity"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Auth interface {
	GenerateToken(user entity.User, c *gin.Context) (string, error)
	ValidateToken(encodedToken string) (jwt.MapClaims, error)
}
