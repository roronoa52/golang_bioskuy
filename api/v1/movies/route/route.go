package movieroute

import (
	"bioskuy/api/v1/movies/controller"
	"bioskuy/api/v1/movies/repository"
	"bioskuy/api/v1/movies/service"
	"bioskuy/auth"
	"bioskuy/helper"
	"bioskuy/middleware"
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func MovieRoute(router *gin.Engine, validate *validator.Validate, db *sql.DB, config *helper.Config) {
	
	authService := auth.NewService(config)

	movieRepo := repository.NewMovieRepository(db)
	movieService := service.NewMovieService(movieRepo)
	movieController := controller.NewMovieController(movieService)
	v1 := router.Group("/api/v1")
	{
		movieRoutes := v1.Group("/movies")
		{
			movieRoutes.POST("/", middleware.AuthMiddleware(authService, "admin"), movieController.CreateMovie)
			movieRoutes.GET("/", movieController.GetAllMovies)
			movieRoutes.GET("/:id", movieController.GetMovie)
			movieRoutes.PUT("/:id", middleware.AuthMiddleware(authService, "admin"), movieController.UpdateMovie)
			movieRoutes.DELETE("/:id", middleware.AuthMiddleware(authService, "admin"), movieController.DeleteMovie)
		}
	}
}
