package repository

import (
	"bioskuy/api/v1/payment/entity"
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PaymentRepositoryTestSuite struct {
	suite.Suite
	repo       PaymentRepository
	mockDb     *sql.DB
	mockSql    sqlmock.Sqlmock
	ctx        context.Context
	ginContext *gin.Context
}

func (suite *PaymentRepositoryTestSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	if err != nil {
		suite.T().Fatal(err)
	}

	suite.mockDb = db
	suite.mockSql = mock
	suite.repo = NewPaymentRepository()
	suite.ctx = context.TODO()
	suite.ginContext = &gin.Context{}
}

func (suite *PaymentRepositoryTestSuite) TearDownTest() {
	suite.mockDb.Close()
}

func TestPaymentRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(PaymentRepositoryTestSuite))
}

func (suite *PaymentRepositoryTestSuite) TestSave_Success() {
	payment := entity.Payment{
		UserID:                 "user-id",
		SeatDetailForBookingID: "seat-detail-id",
		TotalSeat:              5,
		TotalPrice:             100,
	}

	suite.mockSql.ExpectBegin()
	tx, err := suite.mockDb.Begin()
	assert.NoError(suite.T(), err)

	suite.mockSql.ExpectQuery("INSERT INTO payments").
		WithArgs(payment.UserID, payment.SeatDetailForBookingID, payment.TotalSeat, payment.TotalPrice).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("payment-id"))

	savedPayment, err := suite.repo.Save(suite.ctx, tx, payment, suite.ginContext)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "payment-id", savedPayment.ID)
}

func (suite *PaymentRepositoryTestSuite) TestFindAll_Success() {
	query := regexp.QuoteMeta(`
		SELECT 
			p.id, p.user_id, p.seatdetailforbooking_id, p.total_seat, p.total_price, p.status,
			sb.id AS seat_booking_id, sb.status AS seat_booking_status,
			s.id AS seat_id, s.seat_name, s.isAvailable AS seat_isAvailable,
			sh.id AS showtime_id, sh.show_start, sh.show_end,
			m.id AS movie_id, m.title AS movie_title, m.description AS movie_description, 
			m.price AS movie_price, m.duration AS movie_duration, m.status AS movie_status,
			st.id AS studio_id, st.name AS studio_name
		FROM payments p
		JOIN seat_detail_for_bookings sdfb ON p.seatdetailforbooking_id = sdfb.id
		JOIN seats s ON sdfb.seat_id = s.id
		JOIN seat_bookings sb ON sdfb.seatBooking_id = sb.id
		JOIN showtimes sh ON sb.showtime_id = sh.id
		JOIN movies m ON sh.movie_id = m.id
		JOIN studios st ON sh.studio_id = st.id
	`)
	now := time.Now()
	payments := []entity.Payment{
		{
			ID:                     "1",
			UserID:                 "user1",
			SeatDetailForBookingID: "seatbooking1",
			TotalSeat:              2,
			TotalPrice:             200,
			Status:                 "paid",
			SeatBookingID:          "seatbooking1",
			SeatBookingStatus:      "confirmed",
			SeatID:                 "seat1",
			SeatName:               "A1",
			SeatIsAvailable:        "isAvailable",
			ShowtimeID:             "showtime1",
			ShowStart:              now.Format(time.RFC3339),
			ShowEnd:                now.Add(time.Hour).Format(time.RFC3339),
			MovieID:                "movie1",
			MovieTitle:             "Movie 1",
			MovieDescription:       "Description 1",
			MoviePrice:             100,
			MovieDuration:          120,
			MovieStatus:            "active",
			StudioID:               "studio1",
			StudioName:             "Studio 1",
		},
		{
			ID:                     "2",
			UserID:                 "user2",
			SeatDetailForBookingID: "seatbooking2",
			TotalSeat:              3,
			TotalPrice:             300,
			Status:                 "pending",
			SeatBookingID:          "seatbooking2",
			SeatBookingStatus:      "pending",
			SeatID:                 "seat2",
			SeatName:               "B1",
			SeatIsAvailable:        "isAvailable",
			ShowtimeID:             "showtime2",
			ShowStart:              now.Format(time.RFC3339),
			ShowEnd:                now.Add(time.Hour).Format(time.RFC3339),
			MovieID:                "movie2",
			MovieTitle:             "Movie 2",
			MovieDescription:       "Description 2",
			MoviePrice:             120,
			MovieDuration:          110,
			MovieStatus:            "inactive",
			StudioID:               "studio2",
			StudioName:             "Studio 2",
		},
	}

	rows := sqlmock.NewRows([]string{
		"id", "user_id", "seatdetailforbooking_id", "total_seat", "total_price", "status",
		"seat_booking_id", "seat_booking_status", "seat_id", "seat_name", "seat_isAvailable",
		"showtime_id", "show_start", "show_end", "movie_id", "movie_title", "movie_description",
		"movie_price", "movie_duration", "movie_status", "studio_id", "studio_name",
	}).AddRow(
		payments[0].ID, payments[0].UserID, payments[0].SeatDetailForBookingID, payments[0].TotalSeat, payments[0].TotalPrice, payments[0].Status,
		payments[0].SeatBookingID, payments[0].SeatBookingStatus, payments[0].SeatID, payments[0].SeatName, payments[0].SeatIsAvailable,
		payments[0].ShowtimeID, payments[0].ShowStart, payments[0].ShowEnd, payments[0].MovieID, payments[0].MovieTitle, payments[0].MovieDescription,
		payments[0].MoviePrice, payments[0].MovieDuration, payments[0].MovieStatus, payments[0].StudioID, payments[0].StudioName,
	).AddRow(
		payments[1].ID, payments[1].UserID, payments[1].SeatDetailForBookingID, payments[1].TotalSeat, payments[1].TotalPrice, payments[1].Status,
		payments[1].SeatBookingID, payments[1].SeatBookingStatus, payments[1].SeatID, payments[1].SeatName, payments[1].SeatIsAvailable,
		payments[1].ShowtimeID, payments[1].ShowStart, payments[1].ShowEnd, payments[1].MovieID, payments[1].MovieTitle, payments[1].MovieDescription,
		payments[1].MoviePrice, payments[1].MovieDuration, payments[1].MovieStatus, payments[1].StudioID, payments[1].StudioName,
	)

	suite.mockSql.ExpectBegin()
	suite.mockSql.ExpectQuery(query).WillReturnRows(rows)

	tx, err := suite.mockDb.Begin()
	suite.NoError(err)

	ginContext, _ := gin.CreateTestContext(nil)
	result, err := suite.repo.FindAll(context.Background(), tx, ginContext)
	suite.NoError(err)
	suite.Len(result, 2)
	suite.Equal(payments[0].ID, result[0].ID)
	suite.Equal(payments[1].UserID, result[1].UserID)

	suite.mockSql.ExpectCommit()
	err = tx.Commit()
	suite.NoError(err)

	err = suite.mockSql.ExpectationsWereMet()
	suite.NoError(err)
}

func (suite *PaymentRepositoryTestSuite) TestFindAll_Error() {
	query := regexp.QuoteMeta(`
		SELECT 
			p.id, p.user_id, p.seatdetailforbooking_id, p.total_seat, p.total_price, p.status,
			sb.id AS seat_booking_id, sb.status AS seat_booking_status,
			s.id AS seat_id, s.seat_name, s.isAvailable AS seat_isAvailable,
			sh.id AS showtime_id, sh.show_start, sh.show_end,
			m.id AS movie_id, m.title AS movie_title, m.description AS movie_description, 
			m.price AS movie_price, m.duration AS movie_duration, m.status AS movie_status,
			st.id AS studio_id, st.name AS studio_name
		FROM payments p
		JOIN seat_detail_for_bookings sdfb ON p.seatdetailforbooking_id = sdfb.id
		JOIN seats s ON sdfb.seat_id = s.id
		JOIN seat_bookings sb ON sdfb.seatBooking_id = sb.id
		JOIN showtimes sh ON sb.showtime_id = sh.id
		JOIN movies m ON sh.movie_id = m.id
		JOIN studios st ON sh.studio_id = st.id
	`)

	suite.mockSql.ExpectBegin()
	suite.mockSql.ExpectQuery(query).WillReturnError(errors.New("query error"))

	tx, err := suite.mockDb.Begin()
	suite.NoError(err)

	ginContext, _ := gin.CreateTestContext(nil)
	_, err = suite.repo.FindAll(context.Background(), tx, ginContext)
	suite.Error(err)
	suite.EqualError(err, "query error")

	suite.mockSql.ExpectRollback()
	err = tx.Rollback()
	suite.NoError(err)

	err = suite.mockSql.ExpectationsWereMet()
	suite.NoError(err)
}

func (suite *PaymentRepositoryTestSuite) TestUpdate_Success() {
	payment := entity.Payment{
		ID:     "payment-id",
		Status: "new-status",
	}

	suite.mockSql.ExpectBegin()
	tx, err := suite.mockDb.Begin()
	assert.NoError(suite.T(), err)

	query := `UPDATE payments SET status = $1 WHERE id = $2`

	suite.mockSql.ExpectExec(query).WithArgs(payment.Status, payment.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	updatedPayment, err := suite.repo.Update(suite.ctx, tx, payment, suite.ginContext)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "new-status", updatedPayment.Status)
}

func (suite *PaymentRepositoryTestSuite) TestDelete_Success() {
	id := "payment-id"

	suite.mockSql.ExpectBegin()
	tx, err := suite.mockDb.Begin()
	assert.NoError(suite.T(), err)

	query := `DELETE FROM payments WHERE id = $1`

	suite.mockSql.ExpectExec(query).WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = suite.repo.Delete(suite.ctx, tx, id, suite.ginContext)
	assert.NoError(suite.T(), err)
}

func (suite *PaymentRepositoryTestSuite) TestSave_InsertError() {
	payment := entity.Payment{
		UserID:                 "user-id",
		SeatDetailForBookingID: "seat-detail-id",
		TotalSeat:              5,
		TotalPrice:             100,
	}

	suite.mockSql.ExpectBegin()
	tx, err := suite.mockDb.Begin()
	assert.NoError(suite.T(), err)

	suite.mockSql.ExpectQuery("INSERT INTO payments").
		WithArgs(payment.UserID, payment.SeatDetailForBookingID, payment.TotalSeat, payment.TotalPrice).
		WillReturnError(errors.New("insert error"))

	savedPayment, err := suite.repo.Save(suite.ctx, tx, payment, suite.ginContext)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "insert error", err.Error())
	assert.Equal(suite.T(), "", savedPayment.ID)
}

func (suite *PaymentRepositoryTestSuite) TestFindAll_QueryError() {
	suite.mockSql.ExpectBegin()
	tx, err := suite.mockDb.Begin()
	assert.NoError(suite.T(), err)

	query := `SELECT p.id, p.user_id, p.seatdetailforbooking_id, p.total_seat, p.total_price, p.status,
                     sb.id AS seat_booking_id, sb.status AS seat_booking_status,
                     s.id AS seat_id, s.seat_name, s.isAvailable AS seat_isAvailable,
                     sh.id AS showtime_id, sh.show_start, sh.show_end,
                     m.id AS movie_id, m.title AS movie_title, m.description AS movie_description, 
                     m.price AS movie_price, m.duration AS movie_duration, m.status AS movie_status,
                     st.id AS studio_id, st.name AS studio_name
              FROM payments p
              JOIN seat_detail_for_bookings sdfb ON p.seatdetailforbooking_id = sdfb.id
              JOIN seats s ON sdfb.seat_id = s.id
              JOIN seat_bookings sb ON sdfb.seatBooking_id = sb.id
              JOIN showtimes sh ON sb.showtime_id = sh.id
              JOIN movies m ON sh.movie_id = m.id
              JOIN studios st ON sh.studio_id = st.id`

	suite.mockSql.ExpectQuery(query).WillReturnError(errors.New("query error"))

	payments, err := suite.repo.FindAll(suite.ctx, tx, suite.ginContext)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "query error", err.Error())
	assert.Nil(suite.T(), payments)
}

func (suite *PaymentRepositoryTestSuite) TestUpdate_UpdateError() {
	payment := entity.Payment{
		ID:     "payment-id",
		Status: "paid",
	}

	suite.mockSql.ExpectBegin()
	tx, err := suite.mockDb.Begin()
	assert.NoError(suite.T(), err)

	query := `UPDATE payments SET status = $1 WHERE id = $2`

	suite.mockSql.ExpectExec(query).
		WithArgs(payment.Status, payment.ID).
		WillReturnError(errors.New("update error"))

	updatedPayment, err := suite.repo.Update(suite.ctx, tx, payment, suite.ginContext)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "update error", err.Error())
	assert.Equal(suite.T(), payment.ID, updatedPayment.ID)
	assert.Equal(suite.T(), "", updatedPayment.Status)
}

func (suite *PaymentRepositoryTestSuite) TestUpdate_Error() {
	payment := entity.Payment{
		ID:     "payment-id",
		Status: "paid",
	}

	suite.mockSql.ExpectBegin()
	tx, err := suite.mockDb.Begin()
	assert.NoError(suite.T(), err)

	query := `UPDATE payments SET status = $1 WHERE id = $2`

	suite.mockSql.ExpectExec(query).
		WithArgs(payment.Status, payment.ID).
		WillReturnError(errors.New("update error"))

	updatedPayment, err := suite.repo.Update(suite.ctx, tx, payment, suite.ginContext)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "update error", err.Error())
	assert.Equal(suite.T(), payment.ID, updatedPayment.ID)
	assert.Equal(suite.T(), "", updatedPayment.Status)
}

func (suite *PaymentRepositoryTestSuite) TestDelete_Error() {
	paymentId := "payment-id"

	suite.mockSql.ExpectBegin()
	tx, err := suite.mockDb.Begin()
	assert.NoError(suite.T(), err)

	query := `DELETE FROM payments WHERE id = $1`

	suite.mockSql.ExpectExec(query).
		WithArgs(paymentId).
		WillReturnError(errors.New("delete error"))

	err = suite.repo.Delete(suite.ctx, tx, paymentId, suite.ginContext)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "delete error", err.Error())
}

func (suite *PaymentRepositoryTestSuite) TestSave_NullOrEmptyValues() {
	payment := entity.Payment{
		UserID:                 "",
		SeatDetailForBookingID: "",
		TotalSeat:              0,
		TotalPrice:             0,
	}

	suite.mockSql.ExpectBegin()
	tx, err := suite.mockDb.Begin()
	assert.NoError(suite.T(), err)

	suite.mockSql.ExpectQuery("INSERT INTO payments").
		WithArgs(payment.UserID, payment.SeatDetailForBookingID, payment.TotalSeat, payment.TotalPrice).
		WillReturnError(errors.New("insert error"))

	savedPayment, err := suite.repo.Save(suite.ctx, tx, payment, suite.ginContext)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "insert error", err.Error())
	assert.Equal(suite.T(), "", savedPayment.ID)
}

func (suite *PaymentRepositoryTestSuite) TestSave_InvalidDataTypes() {
	payment := entity.Payment{
		UserID:                 "user-id",
		SeatDetailForBookingID: "seat-detail-id",
		TotalSeat:              5, // Invalid data type
		TotalPrice:             100,
	}

	suite.mockSql.ExpectBegin()
	tx, err := suite.mockDb.Begin()
	assert.NoError(suite.T(), err)

	suite.mockSql.ExpectQuery("INSERT INTO payments").
		WithArgs(payment.UserID, payment.SeatDetailForBookingID, payment.TotalSeat, payment.TotalPrice).
		WillReturnError(errors.New("invalid data type"))

	savedPayment, err := suite.repo.Save(suite.ctx, tx, payment, suite.ginContext)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "invalid data type", err.Error())
	assert.Equal(suite.T(), "", savedPayment.ID)
}
func (suite *PaymentRepositoryTestSuite) TestFindByID_Success() {
	query := regexp.QuoteMeta(`
        SELECT 
            p.id, p.user_id, p.seatdetailforbooking_id, p.total_seat, p.total_price, p.status,
            sb.id AS seat_booking_id, sb.status AS seat_booking_status,
            s.id AS seat_id, s.seat_name, s.isAvailable AS seat_isAvailable,
            sh.id AS showtime_id, sh.show_start, sh.show_end,
            m.id AS movie_id, m.title AS movie_title, m.description AS movie_description, 
            m.price AS movie_price, m.duration AS movie_duration, m.status AS movie_status,
            st.id AS studio_id, st.name AS studio_name
        FROM payments p
        JOIN seat_detail_for_bookings sdfb ON p.seatdetailforbooking_id = sdfb.id
        JOIN seats s ON sdfb.seat_id = s.id
        JOIN seat_bookings sb ON sdfb.seatBooking_id = sb.id
        JOIN showtimes sh ON sb.showtime_id = sh.id
        JOIN movies m ON sh.movie_id = m.id
        JOIN studios st ON sh.studio_id = st.id
        WHERE p.id = $1`)
	now := time.Now()
	expectedPayment := entity.Payment{
		ID:                     "1",
		UserID:                 "1",
		SeatDetailForBookingID: "1",
		TotalSeat:              2,
		TotalPrice:             50000,
		Status:                 "PAID",
		SeatBookingID:          "1",
		SeatBookingStatus:      "CONFIRMED",
		SeatID:                 "1",
		SeatName:               "A1",
		SeatIsAvailable:        "isAvailable",
		ShowtimeID:             "showtime1",
		ShowStart:              now.Format(time.RFC3339),
		ShowEnd:                now.Add(time.Hour).Format(time.RFC3339),
		MovieID:                "1",
		MovieTitle:             "Avengers",
		MovieDescription:       "Superhero movie",
		MoviePrice:             25000,
		MovieDuration:          120,
		MovieStatus:            "AVAILABLE",
		StudioID:               "1",
		StudioName:             "Studio 1",
	}

	rows := sqlmock.NewRows([]string{
		"id", "user_id", "seatdetailforbooking_id", "total_seat", "total_price", "status",
		"seat_booking_id", "seat_booking_status",
		"seat_id", "seat_name", "seat_isAvailable",
		"showtime_id", "show_start", "show_end",
		"movie_id", "movie_title", "movie_description", "movie_price", "movie_duration", "movie_status",
		"studio_id", "studio_name",
	}).AddRow(
		expectedPayment.ID, expectedPayment.UserID, expectedPayment.SeatDetailForBookingID, expectedPayment.TotalSeat, expectedPayment.TotalPrice, expectedPayment.Status,
		expectedPayment.SeatBookingID, expectedPayment.SeatBookingStatus,
		expectedPayment.SeatID, expectedPayment.SeatName, expectedPayment.SeatIsAvailable,
		expectedPayment.ShowtimeID, expectedPayment.ShowStart, expectedPayment.ShowEnd,
		expectedPayment.MovieID, expectedPayment.MovieTitle, expectedPayment.MovieDescription, expectedPayment.MoviePrice, expectedPayment.MovieDuration, expectedPayment.MovieStatus,
		expectedPayment.StudioID, expectedPayment.StudioName,
	)

	tx, err := suite.mockDb.Begin()
	assert.NoError(suite.T(), err)
	suite.mockSql.ExpectBegin()
	suite.mockSql.ExpectQuery(query).WithArgs(expectedPayment.ID).WillReturnRows(rows)

	result, err := suite.repo.FindByID(context.Background(), tx, expectedPayment.ID, suite.ginContext)
	suite.NoError(err)
	suite.Equal(expectedPayment, result)
}

func (suite *PaymentRepositoryTestSuite) TestFindByID_NotFound() {
	query := regexp.QuoteMeta(`
        SELECT 
            p.id, p.user_id, p.seatdetailforbooking_id, p.total_seat, p.total_price, p.status,
            sb.id AS seat_booking_id, sb.status AS seat_booking_status,
            s.id AS seat_id, s.seat_name, s.isAvailable AS seat_isAvailable,
            sh.id AS showtime_id, sh.show_start, sh.show_end,
            m.id AS movie_id, m.title AS movie_title, m.description AS movie_description, 
            m.price AS movie_price, m.duration AS movie_duration, m.status AS movie_status,
            st.id AS studio_id, st.name AS studio_name
        FROM payments p
        JOIN seat_detail_for_bookings sdfb ON p.seatdetailforbooking_id = sdfb.id
        JOIN seats s ON sdfb.seat_id = s.id
        JOIN seat_bookings sb ON sdfb.seatBooking_id = sb.id
        JOIN showtimes sh ON sb.showtime_id = sh.id
        JOIN movies m ON sh.movie_id = m.id
        JOIN studios st ON sh.studio_id = st.id
        WHERE p.id = $1`)

	tx, err := suite.mockDb.Begin()
	assert.NoError(suite.T(), err)
	suite.mockSql.ExpectBegin()

	suite.mockSql.ExpectQuery(query).WithArgs("invalid_id").WillReturnError(sql.ErrNoRows)

	_, err = suite.repo.FindByID(context.Background(), tx, "invalid_id", suite.ginContext)
	suite.Error(err)
	suite.EqualError(err, "payment not found")
}

func (suite *PaymentRepositoryTestSuite) TestFindByID_ErrorScan() {
	query := regexp.QuoteMeta(`
        SELECT 
            p.id, p.user_id, p.seatdetailforbooking_id, p.total_seat, p.total_price, p.status,
            sb.id AS seat_booking_id, sb.status AS seat_booking_status,
            s.id AS seat_id, s.seat_name, s.isAvailable AS seat_isAvailable,
            sh.id AS showtime_id, sh.show_start, sh.show_end,
            m.id AS movie_id, m.title AS movie_title, m.description AS movie_description, 
            m.price AS movie_price, m.duration AS movie_duration, m.status AS movie_status,
            st.id AS studio_id, st.name AS studio_name
        FROM payments p
        JOIN seat_detail_for_bookings sdfb ON p.seatdetailforbooking_id = sdfb.id
        JOIN seats s ON sdfb.seat_id = s.id
        JOIN seat_bookings sb ON sdfb.seatBooking_id = sb.id
        JOIN showtimes sh ON sb.showtime_id = sh.id
        JOIN movies m ON sh.movie_id = m.id
        JOIN studios st ON sh.studio_id = st.id
        WHERE p.id = $1`)

	expectedErr := errors.New("scan error")
	tx, err := suite.mockDb.Begin()
	assert.NoError(suite.T(), err)
	suite.mockSql.ExpectBegin()

	rows := sqlmock.NewRows([]string{
		"id", "user_id", "seatdetailforbooking_id", "total_seat", "total_price", "status",
		"seat_booking_id", "seat_booking_status",
		"seat_id", "seat_name", "seat_isAvailable",
		"showtime_id", "show_start", "show_end",
		"movie_id", "movie_title", "movie_description", "movie_price", "movie_duration", "movie_status",
		"studio_id", "studio_name",
	}).AddRow("1", "1", "1", 2, 50000, "PAID", "1", "CONFIRMED", "1", "A1", true,
		"1", time.Now(), time.Now().Add(2*time.Hour), "1", "Avengers", "Superhero movie", 25000, 120, "AVAILABLE",
		"1", "Studio 1")

	suite.mockSql.ExpectQuery(query).WithArgs("1").WillReturnRows(rows)

	// Cause a scan error by expecting more columns than returned
	rows = sqlmock.NewRows([]string{
		"id", "user_id", "seatdetailforbooking_id", "total_seat", "total_price", "status",
		"seat_booking_id", "seat_booking_status",
		"seat_id", "seat_name", "seat_isAvailable",
		"showtime_id", "show_start", "show_end",
		"movie_id", "movie_title", "movie_description", "movie_price", "movie_duration", "movie_status",
		"studio_id", "studio_name",
	}).AddRow("1", "1", "1", 2, 50000, "PAID", "1", "CONFIRMED", "1", "A1", true,
		"1", time.Now(), time.Now().Add(2*time.Hour), "1", "Avengers", "Superhero movie", 25000, 120, "AVAILABLE",
		"1", "Studio 1", "extra_column")

	suite.mockSql.ExpectQuery(query).WithArgs("1").WillReturnRows(rows)

	_, err = suite.repo.FindByID(context.Background(), tx, "1", suite.ginContext)
	suite.Error(err)
	suite.EqualError(err, expectedErr.Error())
}

func (suite *PaymentRepositoryTestSuite) TestDatabaseConnectionError() {
	suite.mockSql.ExpectBegin().WillReturnError(errors.New("db connection error"))

	_, err := suite.mockDb.Begin()
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "db connection error", err.Error())
}
