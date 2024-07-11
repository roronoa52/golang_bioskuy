package servicemock

import (
	"bioskuy/api/v1/seatbooking/dto"
	"context"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

type MockSeatBookingService struct {
	mock.Mock
}

func (m *MockSeatBookingService) Create(ctx context.Context, request dto.SeatBookingRequest, userid string, c *gin.Context) (dto.CreateSeatBookingResponse, error) {
	args := m.Called(ctx, request, userid, c)
	return args.Get(0).(dto.CreateSeatBookingResponse), args.Error(1)
}

func (m *MockSeatBookingService) FindByID(ctx context.Context, id string, c *gin.Context) (dto.SeatBookingResponse, error) {
	args := m.Called(ctx, id, c)
	return args.Get(0).(dto.SeatBookingResponse), args.Error(1)
}

func (m *MockSeatBookingService) FindAll(ctx context.Context, c *gin.Context) ([]dto.SeatBookingResponse, error) {
	args := m.Called(ctx, c)
	return args.Get(0).([]dto.SeatBookingResponse), args.Error(1)
}

func (m *MockSeatBookingService) Delete(ctx context.Context, id string, c *gin.Context) error {
	args := m.Called(ctx, id, c)
	return args.Error(0)
}
