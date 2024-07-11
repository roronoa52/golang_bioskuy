package helper

import (
	"bioskuy/exception"
	"time"

	"github.com/gin-gonic/gin"
)

func StringToDate(value string, c *gin.Context) (time.Time){
	result, err := time.Parse(time.RFC3339, value)
	if err != nil {
		c.Error(exception.ValidationError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return  result
	}

	return result
}