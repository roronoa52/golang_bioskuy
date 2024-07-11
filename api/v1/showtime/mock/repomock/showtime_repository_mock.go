package repomock

import (
	dtoGenre "bioskuy/api/v1/genre/dto"
	entityMovie "bioskuy/api/v1/movies/entity"

	"bioskuy/api/v1/showtime/entity"
	entityStudio "bioskuy/api/v1/studio/entity"
	"context"
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

type MockShowtimeRepository struct {
	mock.Mock
}

func (m *MockShowtimeRepository) Save(ctx context.Context, tx *sql.Tx, showtime entity.Showtime, c *gin.Context) (entity.Showtime, error) {
	args := m.Called(ctx, tx, showtime, c)
	return args.Get(0).(entity.Showtime), args.Error(1)
}

func (m *MockShowtimeRepository) FindByID(ctx context.Context, tx *sql.Tx, id string, c *gin.Context) (entity.Showtime, error) {
	args := m.Called(ctx, tx, id, c)
	return args.Get(0).(entity.Showtime), args.Error(1)
}

func (m *MockShowtimeRepository) FindAll(ctx context.Context, tx *sql.Tx, c *gin.Context) ([]entity.Showtime, error) {
	args := m.Called(ctx, tx, c)
	return args.Get(0).([]entity.Showtime), args.Error(1)
}

func (m *MockShowtimeRepository) Delete(ctx context.Context, tx *sql.Tx, id string, c *gin.Context) error {
	args := m.Called(ctx, tx, id, c)
	return args.Error(0)
}

func (m *MockShowtimeRepository) FindConflictingShowtimes(ctx context.Context, tx *sql.Tx, studio entityStudio.Studio, showtime entity.Showtime, c *gin.Context) error {
	args := m.Called(ctx, tx, studio, showtime, c)
	return args.Error(0)
}

type MockMovieRepository struct {
	mock.Mock
}

func (m *MockMovieRepository) GetAll(page int, size int) ([]entityMovie.Movie, dtoGenre.Paging, error) {
	args := m.Called(page, size)
	return args.Get(0).([]entityMovie.Movie), args.Get(1).(dtoGenre.Paging), args.Error(2)
}

func (m *MockMovieRepository) Create(movie entityMovie.Movie) (entityMovie.Movie, error) {
	args := m.Called(movie)
	return args.Get(0).(entityMovie.Movie), args.Error(1)
}

func (m *MockMovieRepository) GetByID(id string) (entityMovie.Movie, error) {
	args := m.Called(id)
	return args.Get(0).(entityMovie.Movie), args.Error(1)
}

func (m *MockMovieRepository) Update(movie entityMovie.Movie) (entityMovie.Movie, error) {
	args := m.Called(movie)
	return args.Get(0).(entityMovie.Movie), args.Error(1)
}

func (m *MockMovieRepository) Delete(id string) (entityMovie.Movie, error) {
	args := m.Called(id)
	return args.Get(0).(entityMovie.Movie), args.Error(1)
}

type MockStudioRepository struct {
	mock.Mock
}

func (m *MockStudioRepository) Save(ctx context.Context, tx *sql.Tx, studio entityStudio.Studio, c *gin.Context) (entityStudio.Studio, error) {
	args := m.Called(ctx, tx, studio, c)
	return args.Get(0).(entityStudio.Studio), args.Error(1)
}

func (m *MockStudioRepository) FindByID(ctx context.Context, tx *sql.Tx, id string, c *gin.Context) (entityStudio.Studio, error) {
	args := m.Called(ctx, tx, id, c)
	return args.Get(0).(entityStudio.Studio), args.Error(1)
}

func (m *MockStudioRepository) FindAll(ctx context.Context, tx *sql.Tx, c *gin.Context) ([]entityStudio.Studio, error) {
	args := m.Called(ctx, tx, c)
	return args.Get(0).([]entityStudio.Studio), args.Error(1)
}

func (m *MockStudioRepository) Update(ctx context.Context, tx *sql.Tx, studio entityStudio.Studio, c *gin.Context) (entityStudio.Studio, error) {
	args := m.Called(ctx, tx, studio, c)
	return args.Get(0).(entityStudio.Studio), args.Error(1)
}

func (m *MockStudioRepository) Delete(ctx context.Context, tx *sql.Tx, id string, c *gin.Context) error {
	args := m.Called(ctx, tx, id, c)
	return args.Error(0)
}
