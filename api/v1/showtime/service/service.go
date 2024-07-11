package service

import (
	"bioskuy/api/v1/showtime/dto"
	"context"

	"github.com/gin-gonic/gin"
)

type ShowtimeService interface {
	Create(ctx context.Context, request dto.ShowtimeRequest, c *gin.Context) (dto.CreateShowtimesResponseDTO, error) 
	FindByID(ctx context.Context, id string, c *gin.Context) (dto.ShowtimesResponse, error)
	FindAll(ctx context.Context, c *gin.Context) ([]dto.ShowtimesResponse, error)
	Delete(ctx context.Context, id string, c *gin.Context) error
}
