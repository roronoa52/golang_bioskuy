package route

import (
	repoSeat "bioskuy/api/v1/seat/repository"
	"bioskuy/api/v1/studio/controller"
	"bioskuy/api/v1/studio/repository"
	"bioskuy/api/v1/studio/service"
	"bioskuy/auth"
	"bioskuy/helper"
	"bioskuy/middleware"
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func StudioRoute(router *gin.Engine, validate *validator.Validate, db *sql.DB, config *helper.Config) {

	authService := auth.NewService(config)

	seatRepo := repoSeat.NewSeatRepository()

	studioRepo := repository.NewStudioRepository()
	studioService := service.NewStudioService(studioRepo, validate, db, seatRepo)
	studioController := controller.NewStudioController(studioService)

	v1 := router.Group("/api/v1")
	{
		studios := v1.Group("/studios")
		{
			studios.POST("/", middleware.AuthMiddleware(authService, "admin"), studioController.Create)
			studios.GET("/:studioId", studioController.FindById)
			studios.GET("/", studioController.FindAll)
			studios.PUT("/:studioId", middleware.AuthMiddleware(authService, "admin"), studioController.Update)
			studios.DELETE("/:studioId", middleware.AuthMiddleware(authService, "admin"), studioController.Delete)
		}
	}
}
