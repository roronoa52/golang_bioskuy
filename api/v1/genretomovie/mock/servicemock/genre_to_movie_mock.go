package servicemock

import (
	"bioskuy/api/v1/genretomovie/dto"
	"context"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

// Mocking the GenreToMovieService
type MockGenreToMovieService struct {
	mock.Mock
}

func (m *MockGenreToMovieService) Create(ctx context.Context, request dto.CreateGenreToMovieRequest, c *gin.Context) (dto.GenreToMovieCreateResponse, error) {
	args := m.Called(ctx, request, c)
	return args.Get(0).(dto.GenreToMovieCreateResponse), args.Error(1)
}

func (m *MockGenreToMovieService) FindByID(ctx context.Context, id string, c *gin.Context) (dto.GenreToMovieResponse, error) {
	args := m.Called(ctx, id, c)
	return args.Get(0).(dto.GenreToMovieResponse), args.Error(1)
}

func (m *MockGenreToMovieService) FindAll(ctx context.Context, c *gin.Context) ([]dto.GenreToMovieResponse, error) {
	args := m.Called(ctx, c)
	return args.Get(0).([]dto.GenreToMovieResponse), args.Error(1)
}

func (m *MockGenreToMovieService) Delete(ctx context.Context, id string, c *gin.Context) error {
	args := m.Called(ctx, id, c)
	return args.Error(0)
}
