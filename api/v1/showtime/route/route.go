package movieroute

import (
	movieRepo "bioskuy/api/v1/movies/repository"
	"bioskuy/api/v1/showtime/controller"
	showtimeRepo "bioskuy/api/v1/showtime/repository"
	"bioskuy/api/v1/showtime/service"
	studioRepo "bioskuy/api/v1/studio/repository"
	"bioskuy/auth"
	"bioskuy/helper"
	"bioskuy/middleware"
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ShowtimeRoute(router *gin.Engine, validate *validator.Validate, db *sql.DB, config *helper.Config) {
	
	authService := auth.NewService(config)

	showtimeRepo := showtimeRepo.NewShowtimeRepository()
	studioRepo := studioRepo.NewStudioRepository()
	movieRepo := movieRepo.NewMovieRepository(db)
	showService := service.NewGenreToMovieService(showtimeRepo, movieRepo, studioRepo, validate, db)
	showController := controller.NewMovieController(showService)
	v1 := router.Group("/api/v1")
	{
		showtimeRoutes := v1.Group("/showtimes")
		{
			showtimeRoutes.POST("/", middleware.AuthMiddleware(authService, "admin"), showController.Create)
			showtimeRoutes.GET("/", showController.FindAll)
			showtimeRoutes.GET("/:showtimeId", showController.FindById)
			showtimeRoutes.DELETE("/:showtimeId", middleware.AuthMiddleware(authService, "admin"), showController.Delete)
		}
	}
}
