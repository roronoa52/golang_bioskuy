package service

import (
	"bioskuy/api/v1/studio/dto"
	"context"

	"github.com/gin-gonic/gin"
)

type StudioService interface {
	Create(ctx context.Context, request dto.CreateStudioRequest, c *gin.Context) (dto.StudioResponse, error)
	FindByID(ctx context.Context, id string, c *gin.Context) (dto.StudioResponse, error)
	FindAll(ctx context.Context, c *gin.Context) ([]dto.StudioResponse, error)
	Update(ctx context.Context, request dto.UpdateStudioRequest, c *gin.Context) (dto.StudioResponse, error)
	Delete(ctx context.Context, id string, c *gin.Context) error
}