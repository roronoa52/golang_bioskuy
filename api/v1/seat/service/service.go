package service

import (
	"bioskuy/api/v1/seat/dto"
	"context"

	"github.com/gin-gonic/gin"
)

type SeatService interface {
	FindByID(ctx context.Context, id string, c *gin.Context) (dto.SeatResponse, error)
	FindAll(ctx context.Context, id string, c *gin.Context) ([]dto.SeatResponse, error)
}