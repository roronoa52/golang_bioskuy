package servicemock

import (
	"bioskuy/api/v1/genre/dto"
	"bioskuy/api/v1/movies/entity"

	"github.com/stretchr/testify/mock"
)

type MockMovieService struct {
	mock.Mock
}

func (m *MockMovieService) GetAllMovies(page int, size int) ([]entity.Movie, dto.Paging, error) {
	args := m.Called(page, size)
	return args.Get(0).([]entity.Movie), args.Get(1).(dto.Paging), args.Error(2)
}

func (m *MockMovieService) CreateMovie(movie entity.Movie) (entity.Movie, error) {
	args := m.Called(movie)
	return args.Get(0).(entity.Movie), args.Error(1)
}

func (m *MockMovieService) GetMovieByID(id string) (entity.Movie, error) {
	args := m.Called(id)
	return args.Get(0).(entity.Movie), args.Error(1)
}

func (m *MockMovieService) UpdateMovie(movie entity.Movie) (entity.Movie, error) {
	args := m.Called(movie)
	return args.Get(0).(entity.Movie), args.Error(1)
}

func (m *MockMovieService) DeleteMovie(id string) (entity.Movie, error) {
	args := m.Called(id)
	return args.Get(0).(entity.Movie), args.Error(1)
}
