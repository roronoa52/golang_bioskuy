package servicemock

import (
	"bioskuy/api/v1/payment/dto"
	"context"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

// Mocking the GenreToMovieService
type MockPaymentService struct {
	mock.Mock
}

func (m *MockPaymentService) Create(ctx context.Context, request dto.PaymentRequest, userid string, c *gin.Context) (dto.CreatePaymentResponse, error) {
	args := m.Called(ctx, request, userid, c)
	return args.Get(0).(dto.CreatePaymentResponse), args.Error(1)
}

func (m *MockPaymentService) Update(ctx context.Context, notificationPayload map[string]interface{}, c *gin.Context) {
	m.Called(ctx, notificationPayload, c)
	return
}

func (m *MockPaymentService) FindByID(ctx context.Context, id string, c *gin.Context) (dto.PaymentResponse, error) {
	args := m.Called(ctx, id, c)
	return args.Get(0).(dto.PaymentResponse), args.Error(1)
}

func (m *MockPaymentService) FindAll(ctx context.Context, c *gin.Context) ([]dto.PaymentResponse, error) {
	args := m.Called(ctx, c)
	return args.Get(0).([]dto.PaymentResponse), args.Error(1)
}
