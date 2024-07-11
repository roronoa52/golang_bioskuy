package servicemock

import (
	"bioskuy/api/v1/studio/dto"
	"context"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

type MockStudioService struct {
	mock.Mock
}

func (m *MockStudioService) Create(ctx context.Context, request dto.CreateStudioRequest, c *gin.Context) (dto.StudioResponse, error) {
	args := m.Called(ctx, request, request, c)
	return args.Get(0).(dto.StudioResponse), args.Error(1)
}

func (m *MockStudioService) FindByID(ctx context.Context, id string, c *gin.Context) (dto.StudioResponse, error) {
	args := m.Called(ctx, id, c)
	return args.Get(0).(dto.StudioResponse), args.Error(1)
}

func (m *MockStudioService) Update(ctx context.Context, request dto.UpdateStudioRequest, c *gin.Context) (dto.StudioResponse, error) {
	args := m.Called(ctx, request, c)
	return args.Get(0).(dto.StudioResponse), args.Error(1)
}

func (m *MockStudioService) FindAll(ctx context.Context, c *gin.Context) ([]dto.StudioResponse, error) {
	args := m.Called(ctx, c)
	return args.Get(0).([]dto.StudioResponse), args.Error(1)
}

func (m *MockStudioService) Delete(ctx context.Context, id string, c *gin.Context) error {
	args := m.Called(ctx, id, c)
	return args.Error(0)
}
