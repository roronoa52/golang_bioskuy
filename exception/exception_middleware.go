package exception

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler(c *gin.Context) {
	c.Next()

	for _, ginErr := range c.Errors {
		if ginErr != nil {
			err := ginErr.Err
			if validationError(c, err) {
				return
			}

			if forbiddenError(c, err) {
				return
			}

			if notFoundError(c, err) {
				return
			}

			if internalServerError(c, err) {
				return
			}
		}
	}
}

func forbiddenError(c *gin.Context, err error) bool {
	if e, ok := err.(ForbiddenError); ok {
		c.JSON(http.StatusForbidden, gin.H{"error": e.Error()})
		return true
	}
	return false
}

func validationError(c *gin.Context, err error) bool {
	if e, ok := err.(ValidationError); ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": e.Error()})
		return true
	}
	return false
}

func notFoundError(c *gin.Context, err error) bool {
	if e, ok := err.(NotFoundError); ok {
		c.JSON(http.StatusNotFound, gin.H{"error": e.Message})
		return true
	}
	return false
}

func internalServerError(c *gin.Context, err error) bool {
	if e, ok := err.(InternalServerError); ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": e.Message})
		panic(e)
	}
	return false
}
