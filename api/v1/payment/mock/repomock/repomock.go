package repomock

import (
	"bioskuy/api/v1/payment/entity"
	eS "bioskuy/api/v1/seat/entity"
	eSB "bioskuy/api/v1/seatbooking/entity"
	"context"
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

type MockPaymentRepository struct {
	mock.Mock
}

func (m *MockPaymentRepository) Save(ctx context.Context, tx *sql.Tx, payment entity.Payment, c *gin.Context) (entity.Payment, error) {
	args := m.Called(ctx, tx, payment, c)
	return args.Get(0).(entity.Payment), args.Error(1)
}

func (m *MockPaymentRepository) FindByID(ctx context.Context, tx *sql.Tx, id string, c *gin.Context) (entity.Payment, error) {
	args := m.Called(ctx, tx, id, c)
	return args.Get(0).(entity.Payment), args.Error(1)
}

func (m *MockPaymentRepository) Update(ctx context.Context, tx *sql.Tx, payment entity.Payment, c *gin.Context) (entity.Payment, error) {
	args := m.Called(ctx, tx, c)
	return args.Get(0).(entity.Payment), args.Error(1)
}

func (m *MockPaymentRepository) FindAll(ctx context.Context, tx *sql.Tx, c *gin.Context) ([]entity.Payment, error) {
	args := m.Called(ctx, tx, c)
	return args.Get(0).([]entity.Payment), args.Error(1)
}

func (m *MockPaymentRepository) Delete(ctx context.Context, tx *sql.Tx, id string, c *gin.Context) error {
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

type MockSeatBookingRepository struct {
	mock.Mock
}

func (m *MockSeatBookingRepository) Save(ctx context.Context, tx *sql.Tx, seatbooking eSB.SeatBooking, c *gin.Context) (eSB.SeatBooking, error) {
	args := m.Called(ctx, tx, seatbooking, c)
	return args.Get(0).(eSB.SeatBooking), args.Error(1)
}

func (m *MockSeatBookingRepository) FindByID(ctx context.Context, tx *sql.Tx, id string, c *gin.Context) (eSB.SeatBooking, error) {
	args := m.Called(ctx, tx, id, c)
	return args.Get(0).(eSB.SeatBooking), args.Error(1)
}

func (m *MockSeatBookingRepository) FindAllPendingByUserID(ctx context.Context, tx *sql.Tx, userID string, c *gin.Context) ([]eSB.SeatBooking, error) {
	args := m.Called(ctx, tx, userID, c)
	return args.Get(0).([]eSB.SeatBooking), args.Error(1)
}

func (m *MockSeatBookingRepository) FindAll(ctx context.Context, tx *sql.Tx, c *gin.Context) ([]eSB.SeatBooking, error) {
	args := m.Called(ctx, tx, c)
	return args.Get(0).([]eSB.SeatBooking), args.Error(1)
}

func (m *MockSeatBookingRepository) Delete(ctx context.Context, tx *sql.Tx, id string, c *gin.Context) error {
	args := m.Called(ctx, tx, id, c)
	return args.Error(0)
}

func (m *MockSeatBookingRepository) Update(ctx context.Context, tx *sql.Tx, payment eSB.SeatBooking, c *gin.Context) (eSB.SeatBooking, error) {
	args := m.Called(ctx, tx, payment, c)
	return args.Get(0).(eSB.SeatBooking), args.Error(1)
}
