package genreroute

import (
	"bioskuy/api/v1/genre/controller"
	"bioskuy/api/v1/genre/repository"
	"bioskuy/api/v1/genre/service"
	"bioskuy/auth"
	"bioskuy/helper"
	"bioskuy/middleware"
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func GenreRoute(router *gin.Engine, validate *validator.Validate, db *sql.DB, config *helper.Config) {

	authService := auth.NewService(config)

	genreRepo := repository.NewGenreRepository(db)
	genreService := service.NewGenreService(genreRepo)
	genreController := controller.NewGenreController(genreService)
	v1 := router.Group("/api/v1")
	{
		genre := v1.Group("/genres")
		{
			genre.POST("/", middleware.AuthMiddleware(authService, "admin"), genreController.CreateGenre)
			genre.GET("/", genreController.GetAll)
			genre.GET("/:id", genreController.GetGenre)
			genre.PUT("/:id", middleware.AuthMiddleware(authService, "admin"), genreController.UpdateGenre)
			genre.DELETE("/:id", middleware.AuthMiddleware(authService, "admin"), genreController.DeleteGenre)
		}
	}
}
