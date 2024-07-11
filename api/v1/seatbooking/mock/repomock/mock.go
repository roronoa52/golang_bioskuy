package mock

import (
	eS "bioskuy/api/v1/seat/entity"
	"bioskuy/api/v1/seatbooking/entity"
	eSO "bioskuy/api/v1/showtime/entity"
	eSt "bioskuy/api/v1/studio/entity"
	"context"
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

type SeatBookingRepositoryMock struct {
	mock.Mock
}

func (m *SeatBookingRepositoryMock) Save(ctx context.Context, tx *sql.Tx, seatbooking entity.SeatBooking, c *gin.Context) (entity.SeatBooking, error) {
	args := m.Called(ctx, tx, seatbooking, c)
	return args.Get(0).(entity.SeatBooking), args.Error(1)
}

func (m *SeatBookingRepositoryMock) FindByID(ctx context.Context, tx *sql.Tx, id string, c *gin.Context) (entity.SeatBooking, error) {
	args := m.Called(ctx, tx, id, c)
	return args.Get(0).(entity.SeatBooking), args.Error(1)
}

func (m *SeatBookingRepositoryMock) FindAllPendingByUserID(ctx context.Context, tx *sql.Tx, userID string, c *gin.Context) ([]entity.SeatBooking, error) {
	args := m.Called(ctx, tx, userID, c)
	return args.Get(0).([]entity.SeatBooking), args.Error(1)
}

func (m *SeatBookingRepositoryMock) FindAll(ctx context.Context, tx *sql.Tx, c *gin.Context) ([]entity.SeatBooking, error) {
	args := m.Called(ctx, tx, c)
	return args.Get(0).([]entity.SeatBooking), args.Error(1)
}

func (m *SeatBookingRepositoryMock) Delete(ctx context.Context, tx *sql.Tx, id string, c *gin.Context) error {
	args := m.Called(ctx, tx, id, c)
	return args.Error(0)
}

func (m *SeatBookingRepositoryMock) Update(ctx context.Context, tx *sql.Tx, payment entity.SeatBooking, c *gin.Context) (entity.SeatBooking, error) {
	args := m.Called(ctx, tx, payment, c)
	return args.Get(0).(entity.SeatBooking), args.Error(1)
}

type MockShowtimeRepository struct {
	mock.Mock
}

func (m *MockShowtimeRepository) Save(ctx context.Context, tx *sql.Tx, showtime eSO.Showtime, c *gin.Context) (eSO.Showtime, error) {
	args := m.Called(ctx, tx, showtime, c)
	return args.Get(0).(eSO.Showtime), args.Error(1)
}

func (m *MockShowtimeRepository) FindByID(ctx context.Context, tx *sql.Tx, id string, c *gin.Context) (eSO.Showtime, error) {
	args := m.Called(ctx, tx, id, c)
	return args.Get(0).(eSO.Showtime), args.Error(1)
}

func (m *MockShowtimeRepository) FindAll(ctx context.Context, tx *sql.Tx, c *gin.Context) ([]eSO.Showtime, error) {
	args := m.Called(ctx, tx, c)
	return args.Get(0).([]eSO.Showtime), args.Error(1)
}

func (m *MockShowtimeRepository) Delete(ctx context.Context, tx *sql.Tx, id string, c *gin.Context) error {
	args := m.Called(ctx, tx, id, c)
	return args.Error(0)
}

func (m *MockShowtimeRepository) FindConflictingShowtimes(ctx context.Context, tx *sql.Tx, studio eSt.Studio, showtime eSO.Showtime, c *gin.Context) error {
	args := m.Called(ctx, tx, studio, showtime, c)
	return args.Error(0)
}

type SeatRepository struct {
	mock.Mock
}

func (m *SeatRepository) Save(ctx context.Context, tx *sql.Tx, seat eS.Seat, c *gin.Context) (eS.Seat, error) {
	args := m.Called(ctx, tx, seat, c)
	return args.Get(0).(eS.Seat), args.Error(1)
}

func (m *SeatRepository) FindByID(ctx context.Context, tx *sql.Tx, id string, c *gin.Context) (eS.Seat, error) {
	args := m.Called(ctx, tx, id, c)
	return args.Get(0).(eS.Seat), args.Error(1)
}

func (m *SeatRepository) FindByIDWithNotAvailable(ctx context.Context, tx *sql.Tx, id string, c *gin.Context) (eS.Seat, error) {
	args := m.Called(ctx, tx, id, c)
	return args.Get(0).(eS.Seat), args.Error(1)
}

func (m *SeatRepository) Update(ctx context.Context, tx *sql.Tx, seat eS.Seat, c *gin.Context) (eS.Seat, error) {
	args := m.Called(ctx, tx, c)
	return args.Get(0).(eS.Seat), args.Error(1)
}

func (m *SeatRepository) FindAll(ctx context.Context, id string, tx *sql.Tx, c *gin.Context) ([]eS.Seat, error) {
	args := m.Called(ctx, id, tx, c)
	return args.Get(0).([]eS.Seat), args.Error(1)
}

func (m *SeatRepository) Delete(ctx context.Context, tx *sql.Tx, id string, c *gin.Context) error {
	args := m.Called(ctx, tx, id, c)
	return args.Error(0)
}
