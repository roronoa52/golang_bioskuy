package route

import (
	"bioskuy/api/v1/seat/controller"
	"bioskuy/api/v1/seat/repository"
	"bioskuy/api/v1/seat/service"
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func SeatRoute(router *gin.Engine, validate *validator.Validate, db *sql.DB) {

	seatRepo := repository.NewSeatRepository()
	seatService := service.NewSeatervice(seatRepo, validate, db)
	seatController := controller.NewSeatController(seatService)

	v1 := router.Group("/api/v1")
	{
		seats := v1.Group("/seats")
		{
			seats.GET("/:seatId", seatController.FindById)
			seats.GET("/studio/:studioId", seatController.FindAll)
		}
	}
}
