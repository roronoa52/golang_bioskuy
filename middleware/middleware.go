package middleware

import (
	"bioskuy/auth"
	"bioskuy/exception"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authService auth.Auth, allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			c.Error(exception.ForbiddenError{Message: "Unauthorized"}).SetType(gin.ErrorTypePublic)
			c.AbortWithStatus(403)
			return
		}

		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		claims, err := authService.ValidateToken(tokenString)
		if err != nil {
			c.Error(exception.ForbiddenError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
			c.AbortWithStatus(403)
			return
		}

		role, ok := claims["role"].(string)
		if !ok || !contains(allowedRoles, role) {
			c.Error(exception.ForbiddenError{Message: "you don't have access this feature"}).SetType(gin.ErrorTypePublic)
			c.AbortWithStatus(403)
			return
		}

		c.Set("user_id", claims["user_id"])
		c.Set("name", claims["name"])
		c.Set("role", role)
		c.Set("email", claims["email"])

		c.Next()
	}
}

func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}
