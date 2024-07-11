package service

import (
	"bioskuy/api/v1/seatbooking/dto"
	"context"

	"github.com/gin-gonic/gin"
)

type SeatBookingService interface {
	Create(ctx context.Context, request dto.SeatBookingRequest, userid string, c *gin.Context) (dto.CreateSeatBookingResponse, error)
	FindByID(ctx context.Context, id string, c *gin.Context) (dto.SeatBookingResponse, error)
	FindAll(ctx context.Context, c *gin.Context) ([]dto.SeatBookingResponse, error)
	Delete(ctx context.Context, id string, c *gin.Context) error
}
