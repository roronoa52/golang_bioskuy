package repomock

import (
	"bioskuy/api/v1/seat/entity"
	"context"
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

type SeatRepository struct {
	mock.Mock
}

func (m *SeatRepository) Save(ctx context.Context, tx *sql.Tx, seat entity.Seat, c *gin.Context) (entity.Seat, error) {
	args := m.Called(ctx, tx, seat, c)
	return args.Get(0).(entity.Seat), args.Error(1)
}

func (m *SeatRepository) FindByID(ctx context.Context, tx *sql.Tx, id string, c *gin.Context) (entity.Seat, error) {
	args := m.Called(ctx, tx, id, c)
	return args.Get(0).(entity.Seat), args.Error(1)
}

func (m *SeatRepository) FindByIDWithNotAvailable(ctx context.Context, tx *sql.Tx, id string, c *gin.Context) (entity.Seat, error) {
	args := m.Called(ctx, tx, id, c)
	return args.Get(0).(entity.Seat), args.Error(1)
}

func (m *SeatRepository) Update(ctx context.Context, tx *sql.Tx, seat entity.Seat, c *gin.Context) (entity.Seat, error) {
	args := m.Called(ctx, tx, c)
	return args.Get(0).(entity.Seat), args.Error(1)
}

func (m *SeatRepository) FindAll(ctx context.Context, id string, tx *sql.Tx, c *gin.Context) ([]entity.Seat, error) {
	args := m.Called(ctx, id, tx, c)
	return args.Get(0).([]entity.Seat), args.Error(1)
}

func (m *SeatRepository) Delete(ctx context.Context, tx *sql.Tx, id string, c *gin.Context) error {
	args := m.Called(ctx, tx, id, c)
	return args.Error(0)
}
