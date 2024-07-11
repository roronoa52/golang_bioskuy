package controller

import (
	"bioskuy/api/v1/genre/dto"
	"bioskuy/api/v1/genre/entity"
	"bioskuy/api/v1/genre/service"
	"bioskuy/web"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type genreControllerImpl struct {
	Service service.GenreService
}

func NewGenreController(service service.GenreService) GenreController {
	return &genreControllerImpl{Service: service}
}

func (gc *genreControllerImpl) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	genres, paging, err := gc.Service.GetAll(page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, web.FormatResponsePaging{
		ResponseCode: http.StatusOK,
		Data:         genres,
		Paging: web.Paging{
			Page:      paging.Page,
			TotalData: paging.TotalRows,
		},
	})
}

func (ctrl *genreControllerImpl) CreateGenre(c *gin.Context) {
	var createDTO dto.CreateGenreDTO
	if err := c.ShouldBindJSON(&createDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	genre := entity.Genre{
		Name: createDTO.Name,
	}
	createdGenre, err := ctrl.Service.CreateGenre(genre)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, dto.GenreResponseDTO{ID: createdGenre.ID, Name: createdGenre.Name})
}

func (ctrl *genreControllerImpl) GetGenre(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	genre, err := ctrl.Service.GetGenreByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.GenreResponseDTO{ID: genre.ID, Name: genre.Name})
}

func (ctrl *genreControllerImpl) UpdateGenre(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var updateDTO dto.UpdateGenreDTO
	if err := c.ShouldBindJSON(&updateDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	genre := entity.Genre{
		ID:   id,
		Name: updateDTO.Name,
	}

	updatedGenre, err := ctrl.Service.UpdateGenre(genre)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.GenreResponseDTO{ID: updatedGenre.ID, Name: updatedGenre.Name})
}

func (ctrl *genreControllerImpl) DeleteGenre(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	deletedGenre, err := ctrl.Service.DeleteGenre(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.GenreResponseDTO{ID: deletedGenre.ID, Name: deletedGenre.Name})
}
