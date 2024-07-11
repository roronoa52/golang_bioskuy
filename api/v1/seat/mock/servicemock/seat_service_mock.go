package servicemock

import (
	"bioskuy/api/v1/seat/dto"
	"context"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

type SeatService struct {
	mock.Mock
}

func (m *SeatService) FindByID(ctx context.Context, id string, c *gin.Context) (dto.SeatResponse, error) {
	args := m.Called(ctx, id, c)
	return args.Get(0).(dto.SeatResponse), args.Error(1)
}

func (m *SeatService) FindAll(ctx context.Context, id string, c *gin.Context) ([]dto.SeatResponse, error) {
	args := m.Called(ctx, id, c)
	return args.Get(0).([]dto.SeatResponse), args.Error(1)
}
