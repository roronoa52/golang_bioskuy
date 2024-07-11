package controller

import (
	"bioskuy/api/v1/movies/dto"
	"bioskuy/api/v1/movies/entity"
	"bioskuy/api/v1/movies/service"
	"bioskuy/web"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type movieControllerImpl struct {
	Service service.MovieService
}

func NewMovieController(service service.MovieService) MovieController {
	return &movieControllerImpl{Service: service}
}

func (ctrl *movieControllerImpl) GetAllMovies(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	size, err := strconv.Atoi(c.DefaultQuery("size", "10"))
	if err != nil || size < 1 {
		size = 10
	}

	movies, paging, err := ctrl.Service.GetAllMovies(page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := web.FormatResponsePaging{
		ResponseCode: http.StatusOK,
		Data:         movies,
		Paging: web.Paging{
			Page:      paging.Page,
			TotalData: paging.TotalRows,
		},
	}
	c.JSON(http.StatusOK, response)
}

func (ctrl *movieControllerImpl) CreateMovie(c *gin.Context) {
	var createDTO dto.CreateMovieDTO
	if err := c.ShouldBindJSON(&createDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	movie := entity.Movie{
		Title:       createDTO.Title,
		Description: createDTO.Description,
		Price:       createDTO.Price,
		Duration:    createDTO.Duration,
		Status:      createDTO.Status,
	}
	createdMovie, err := ctrl.Service.CreateMovie(movie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, dto.MovieResponseDTO{ID: createdMovie.ID, Title: createdMovie.Title, Description: createdMovie.Description, Price: createdMovie.Price, Duration: createdMovie.Duration, Status: createdMovie.Status})
}

func (ctrl *movieControllerImpl) GetMovie(c *gin.Context) {
	id := c.Param("id")

	movie, err := ctrl.Service.GetMovieByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.MovieResponseDTO{ID: movie.ID, Title: movie.Title, Description: movie.Description, Price: movie.Price, Duration: movie.Duration, Status: movie.Status})
}

func (ctrl *movieControllerImpl) UpdateMovie(c *gin.Context) {

	id := c.Param("id")

	var updateDTO dto.UpdateMovieDTO
	if err := c.ShouldBindJSON(&updateDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	movie := entity.Movie{
		ID:          id,
		Title:       updateDTO.Title,
		Description: updateDTO.Description,
		Price:       updateDTO.Price,
		Duration:    updateDTO.Duration,
		Status:      updateDTO.Status,
	}
	updatedMovie, err := ctrl.Service.UpdateMovie(movie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.MovieResponseDTO{ID: updatedMovie.ID, Title: updatedMovie.Title, Description: updatedMovie.Description, Price: updatedMovie.Price, Duration: updatedMovie.Duration, Status: updatedMovie.Status})
}

func (ctrl *movieControllerImpl) DeleteMovie(c *gin.Context) {
	id := c.Param("id")

	deletedMovie, err := ctrl.Service.DeleteMovie(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.MovieResponseDTO{ID: deletedMovie.ID, Title: deletedMovie.Title, Description: deletedMovie.Description, Price: deletedMovie.Price, Duration: deletedMovie.Duration, Status: deletedMovie.Status})
}
