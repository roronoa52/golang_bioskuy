package service

import (
	"bioskuy/api/v1/payment/dto"
	"context"

	"github.com/gin-gonic/gin"
)

type PaymentService interface {
	Create(ctx context.Context, request dto.PaymentRequest, userid string, c *gin.Context) (dto.CreatePaymentResponse, error)
	Update(ctx context.Context, notificationPayload map[string]interface{}, c *gin.Context) 
	FindByID(ctx context.Context, id string, c *gin.Context) (dto.PaymentResponse, error)
	FindAll(ctx context.Context, c *gin.Context) ([]dto.PaymentResponse, error)
}
