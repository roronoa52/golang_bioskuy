package repomock

import (
	"bioskuy/api/v1/genretomovie/entity"
	"context"
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

type MockGenreToMovieRepository struct {
	mock.Mock
}

func (m *MockGenreToMovieRepository) Save(ctx context.Context, tx *sql.Tx, user entity.GenreToMovie, c *gin.Context) (entity.GenreToMovie, error) {
	args := m.Called(ctx, tx, user, c)
	return args.Get(0).(entity.GenreToMovie), args.Error(1)
}

func (m *MockGenreToMovieRepository) FindByID(ctx context.Context, tx *sql.Tx, id string, c *gin.Context) (entity.GenreToMovie, error) {
	args := m.Called(ctx, tx, id, c)
	return args.Get(0).(entity.GenreToMovie), args.Error(1)
}

func (m *MockGenreToMovieRepository) FindAll(ctx context.Context, tx *sql.Tx, c *gin.Context) ([]entity.GenreToMovie, error) {
	args := m.Called(ctx, tx, c)
	return args.Get(0).([]entity.GenreToMovie), args.Error(1)
}

func (m *MockGenreToMovieRepository) Delete(ctx context.Context, tx *sql.Tx, id string, c *gin.Context) error {
	args := m.Called(ctx, tx, id, c)
	return args.Error(0)
}
