package repomock

import (
	eS "bioskuy/api/v1/seat/entity"
	"bioskuy/api/v1/studio/entity"
	"context"
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

type MockStudioRepository struct {
	mock.Mock
}

func (m *MockStudioRepository) Save(ctx context.Context, tx *sql.Tx, user entity.Studio, c *gin.Context) (entity.Studio, error) {
	args := m.Called(ctx, tx, user, c)
	return args.Get(0).(entity.Studio), args.Error(1)
}

func (m *MockStudioRepository) FindByID(ctx context.Context, tx *sql.Tx, id string, c *gin.Context) (entity.Studio, error) {
	args := m.Called(ctx, tx, id, c)
	return args.Get(0).(entity.Studio), args.Error(1)
}

func (m *MockStudioRepository) Update(ctx context.Context, tx *sql.Tx, studio entity.Studio, c *gin.Context) (entity.Studio, error) {
	args := m.Called(ctx, tx, studio, c)
	return args.Get(0).(entity.Studio), args.Error(1)
}

func (m *MockStudioRepository) FindAll(ctx context.Context, tx *sql.Tx, c *gin.Context) ([]entity.Studio, error) {
	args := m.Called(ctx, tx, c)
	return args.Get(0).([]entity.Studio), args.Error(1)
}

func (m *MockStudioRepository) Delete(ctx context.Context, tx *sql.Tx, id string, c *gin.Context) error {
	args := m.Called(ctx, tx, id, c)
	return args.Error(0)
}

type MockSeatRepository struct {
	mock.Mock
}

func (m *MockSeatRepository) Save(ctx context.Context, tx *sql.Tx, seat eS.Seat, c *gin.Context) (eS.Seat, error) {
	args := m.Called(ctx, tx, seat, c)
	return args.Get(0).(eS.Seat), args.Error(1)
}

func (m *MockSeatRepository) FindByID(ctx context.Context, tx *sql.Tx, id string, c *gin.Context) (eS.Seat, error) {
	args := m.Called(ctx, tx, id, c)
	return args.Get(0).(eS.Seat), args.Error(1)
}

func (m *MockSeatRepository) FindByIDWithNotAvailable(ctx context.Context, tx *sql.Tx, id string, c *gin.Context) (eS.Seat, error) {
	args := m.Called(ctx, tx, id, c)
	return args.Get(0).(eS.Seat), args.Error(1)
}

func (m *MockSeatRepository) Update(ctx context.Context, tx *sql.Tx, seat eS.Seat, c *gin.Context) (eS.Seat, error) {
	args := m.Called(ctx, tx, c)
	return args.Get(0).(eS.Seat), args.Error(1)
}

func (m *MockSeatRepository) FindAll(ctx context.Context, id string, tx *sql.Tx, c *gin.Context) ([]eS.Seat, error) {
	args := m.Called(ctx, id, tx, c)
	return args.Get(0).([]eS.Seat), args.Error(1)
}

func (m *MockSeatRepository) Delete(ctx context.Context, tx *sql.Tx, id string, c *gin.Context) error {
	args := m.Called(ctx, tx, id, c)
	return args.Error(0)
}
