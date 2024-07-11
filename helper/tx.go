package helper

import (
	"bioskuy/exception"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func CommitAndRollback(tx *sql.Tx, c *gin.Context) {
	if r := recover(); r != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			c.Error(exception.InternalServerError{Message: errRollback.Error()}).SetType(gin.ErrorTypePublic)
		}
		panic(r)
	} else {
		if len(c.Errors) > 0 {
			if errRollback := tx.Rollback(); errRollback != nil {
				c.Error(exception.InternalServerError{Message: errRollback.Error()}).SetType(gin.ErrorTypePublic)
			}
		} else {
			if errCommit := tx.Commit(); errCommit != nil {
				c.Error(exception.InternalServerError{Message: errCommit.Error()}).SetType(gin.ErrorTypePublic)
			}
		}
	}
}
