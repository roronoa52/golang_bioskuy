package repository

import (
	"bioskuy/api/v1/seatbooking/entity"
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type SeatBookingRepositoryTestSuite struct {
	suite.Suite
	mockSql sqlmock.Sqlmock
	db      *sql.DB
	repo    SeatBookingRepository
}

func (suite *SeatBookingRepositoryTestSuite) SetupTest() {
	var err error
	suite.db, suite.mockSql, err = sqlmock.New()
	suite.NoError(err)
	suite.repo = NewSeatBookingRepository()
}

func (suite *SeatBookingRepositoryTestSuite) TearDownTest() {
	suite.db.Close()
}

// Test Save Method
func (suite *SeatBookingRepositoryTestSuite) TestSave_Success() {
	queryForSeatBooking := regexp.QuoteMeta("INSERT INTO seat_bookings (user_id, showtime_id) VALUES ($1, $2) RETURNING id")
	queryForSeatDetailForBooking := regexp.QuoteMeta("INSERT INTO seat_detail_for_bookings (seat_id, seatBooking_id) VALUES ($1, $2) RETURNING id")

	seatBooking := entity.SeatBooking{UserID: "user1", ShowtimeID: "showtime1", SeatID: "seat1"}

	suite.mockSql.ExpectBegin()
	suite.mockSql.ExpectQuery(queryForSeatBooking).WithArgs(seatBooking.UserID, seatBooking.ShowtimeID).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))
	suite.mockSql.ExpectQuery(queryForSeatDetailForBooking).WithArgs(seatBooking.SeatID, "1").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))

	tx, err := suite.db.Begin()
	suite.NoError(err)

	ginContext, _ := gin.CreateTestContext(nil)
	seatBooking, err = suite.repo.Save(context.Background(), tx, seatBooking, ginContext)
	suite.NoError(err)
	suite.Equal("1", seatBooking.ID)
	suite.Equal("1", seatBooking.SeatDetailForBookingID)

	suite.mockSql.ExpectCommit()
	err = tx.Commit()
	suite.NoError(err)

	err = suite.mockSql.ExpectationsWereMet()
	suite.NoError(err)
}

func (suite *SeatBookingRepositoryTestSuite) TestSave_Error() {
	queryForSeatBooking := regexp.QuoteMeta("INSERT INTO seat_bookings (user_id, showtime_id) VALUES ($1, $2) RETURNING id")

	seatBooking := entity.SeatBooking{UserID: "user1", ShowtimeID: "showtime1", SeatID: "seat1"}

	suite.mockSql.ExpectBegin()
	suite.mockSql.ExpectQuery(queryForSeatBooking).WithArgs(seatBooking.UserID, seatBooking.ShowtimeID).WillReturnError(errors.New("insert error"))

	tx, err := suite.db.Begin()
	suite.NoError(err)

	ginContext, _ := gin.CreateTestContext(nil)
	_, err = suite.repo.Save(context.Background(), tx, seatBooking, ginContext)
	suite.Error(err)
	suite.EqualError(err, "insert error")

	suite.mockSql.ExpectRollback()
	err = tx.Rollback()
	suite.NoError(err)

	err = suite.mockSql.ExpectationsWereMet()
	suite.NoError(err)
}

func (suite *SeatBookingRepositoryTestSuite) TestFindByID_Success() {
	seatBookingID := "1"
	query := regexp.QuoteMeta(`
		SELECT
			sb.id, sb.status, sb.user_id,
			s.id as showtime_id, s.studio_id, s.movie_id, s.show_start, s.show_end,
			st.name as studio_name,
			m.title as movie_title, m.description as movie_description, m.price as movie_price, m.duration as movie_duration, m.status as movie_status,
			sdfb.id as seat_detail_for_booking_id, sdfb.seat_id,
			se.seat_name, se.isAvailable
		FROM
			seat_bookings sb
		JOIN
			showtimes s ON sb.showtime_id = s.id
		JOIN
			studios st ON s.studio_id = st.id
		JOIN
			movies m ON s.movie_id = m.id
		JOIN
			seat_detail_for_bookings sdfb ON sb.id = sdfb.seatBooking_id
		JOIN
			seats se ON sdfb.seat_id = se.id
		WHERE
			sb.id = $1
	`)
	rows := sqlmock.NewRows([]string{
		"id", "status", "user_id", "showtime_id", "studio_id", "movie_id", "show_start", "show_end", "studio_name",
		"movie_title", "movie_description", "movie_price", "movie_duration", "movie_status",
		"seat_detail_for_booking_id", "seat_id", "seat_name", "isAvailable",
	}).AddRow(seatBookingID, "booked", "user1", "showtime1", "studio1", "movie1", "2023-07-10 10:00:00", "2023-07-10 12:00:00", "Studio 1",
		"Movie 1", "Description", 100, 120, "active", "2", "seat1", "Seat 1", true)

	suite.mockSql.ExpectBegin()
	tx, err := suite.db.Begin()
	suite.NoError(err)

	suite.mockSql.ExpectQuery(query).WithArgs(seatBookingID).WillReturnRows(rows)

	ginContext, _ := gin.CreateTestContext(nil)
	seatBooking, err := suite.repo.FindByID(context.Background(), tx, seatBookingID, ginContext)
	suite.NoError(err)
	suite.Equal(seatBookingID, seatBooking.ID)
	suite.Equal("booked", seatBooking.SeatBookingStatus)

	suite.mockSql.ExpectCommit()
	err = tx.Commit()
	suite.NoError(err)

	err = suite.mockSql.ExpectationsWereMet()
	suite.NoError(err)
}

func (suite *SeatBookingRepositoryTestSuite) TestFindByID_Error() {
	seatBookingID := "1"
	query := regexp.QuoteMeta(`
		SELECT
			sb.id, sb.status, sb.user_id,
			s.id as showtime_id, s.studio_id, s.movie_id, s.show_start, s.show_end,
			st.name as studio_name,
			m.title as movie_title, m.description as movie_description, m.price as movie_price, m.duration as movie_duration, m.status as movie_status,
			sdfb.id as seat_detail_for_booking_id, sdfb.seat_id,
			se.seat_name, se.isAvailable
		FROM
			seat_bookings sb
		JOIN
			showtimes s ON sb.showtime_id = s.id
		JOIN
			studios st ON s.studio_id = st.id
		JOIN
			movies m ON s.movie_id = m.id
		JOIN
			seat_detail_for_bookings sdfb ON sb.id = sdfb.seatBooking_id
		JOIN
			seats se ON sdfb.seat_id = se.id
		WHERE
			sb.id = $1
	`)

	suite.mockSql.ExpectBegin()
	tx, err := suite.db.Begin()
	suite.NoError(err)

	suite.mockSql.ExpectQuery(query).WithArgs(seatBookingID).WillReturnError(sql.ErrNoRows)

	ginContext, _ := gin.CreateTestContext(nil)
	_, err = suite.repo.FindByID(context.Background(), tx, seatBookingID, ginContext)
	suite.Error(err)

	suite.mockSql.ExpectRollback()
	err = tx.Rollback()
	suite.NoError(err)

	err = suite.mockSql.ExpectationsWereMet()
	suite.NoError(err)
}

func (suite *SeatBookingRepositoryTestSuite) TestFindAll_Success() {
	query := regexp.QuoteMeta(`
		SELECT
			sb.id, sb.status, sb.user_id,
			s.id as showtime_id, s.studio_id, s.movie_id, s.show_start, s.show_end,
			st.name as studio_name,
			m.title as movie_title, m.description as movie_description, m.price as movie_price, m.duration as movie_duration, m.status as movie_status,
			sdfb.id as seat_detail_for_booking_id, sdfb.seat_id,
			se.seat_name, se.isAvailable
		FROM
			seat_bookings sb
		JOIN
			showtimes s ON sb.showtime_id = s.id
		JOIN
			studios st ON s.studio_id = st.id
		JOIN
			movies m ON s.movie_id = m.id
		JOIN
			seat_detail_for_bookings sdfb ON sb.id = sdfb.seatBooking_id
		JOIN
			seats se ON sdfb.seat_id = se.id
	`)
	rows := sqlmock.NewRows([]string{
		"id", "status", "user_id", "showtime_id", "studio_id", "movie_id", "show_start", "show_end", "studio_name",
		"movie_title", "movie_description", "movie_price", "movie_duration", "movie_status",
		"seat_detail_for_booking_id", "seat_id", "seat_name", "isAvailable",
	}).
		AddRow("1", "booked", "user1", "showtime1", "studio1", "movie1", "2023-07-10 10:00:00", "2023-07-10 12:00:00", "Studio 1",
			"Movie 1", "Description 1", 100, 120, "active", "2", "seat1", "Seat 1", true).
		AddRow("2", "booked", "user2", "showtime2", "studio2", "movie2", "2023-07-11 14:00:00", "2023-07-11 16:00:00", "Studio 2",
			"Movie 2", "Description 2", 120, 140, "active", "3", "seat2", "Seat 2", true)

	suite.mockSql.ExpectBegin()
	tx, err := suite.db.Begin()
	suite.NoError(err)

	suite.mockSql.ExpectQuery(query).WillReturnRows(rows)

	ginContext, _ := gin.CreateTestContext(nil)
	seatBookings, err := suite.repo.FindAll(context.Background(), tx, ginContext)
	suite.NoError(err)
	suite.Len(seatBookings, 2)

	suite.mockSql.ExpectCommit()
	err = tx.Commit()
	suite.NoError(err)

	err = suite.mockSql.ExpectationsWereMet()
	suite.NoError(err)
}

func (suite *SeatBookingRepositoryTestSuite) TestFindAll_Error() {
	query := regexp.QuoteMeta(`
		SELECT
			sb.id, sb.status, sb.user_id,
			s.id as showtime_id, s.studio_id, s.movie_id, s.show_start, s.show_end,
			st.name as studio_name,
			m.title as movie_title, m.description as movie_description, m.price as movie_price, m.duration as movie_duration, m.status as movie_status,
			sdfb.id as seat_detail_for_booking_id, sdfb.seat_id,
			se.seat_name, se.isAvailable
		FROM
			seat_bookings sb
		JOIN
			showtimes s ON sb.showtime_id = s.id
		JOIN
			studios st ON s.studio_id = st.id
		JOIN
			movies m ON s.movie_id = m.id
		JOIN
			seat_detail_for_bookings sdfb ON sb.id = sdfb.seatBooking_id
		JOIN
			seats se ON sdfb.seat_id = se.id
	`)

	suite.mockSql.ExpectBegin()
	tx, err := suite.db.Begin()
	suite.NoError(err)

	suite.mockSql.ExpectQuery(query).WillReturnError(sql.ErrConnDone)

	ginContext, _ := gin.CreateTestContext(nil)
	_, err = suite.repo.FindAll(context.Background(), tx, ginContext)
	suite.Error(err)

	suite.mockSql.ExpectRollback()
	err = tx.Rollback()
	suite.NoError(err)

	err = suite.mockSql.ExpectationsWereMet()
	suite.NoError(err)
}

func (suite *SeatBookingRepositoryTestSuite) TestDelete_Success() {
	seatBookingID := "1"
	updateSeatQuery := regexp.QuoteMeta(`UPDATE seats SET isAvailable = true WHERE id IN (SELECT seat_id FROM seat_detail_for_bookings WHERE seatBooking_id = $1)`)
	deleteSeatDetailQuery := regexp.QuoteMeta(`DELETE FROM seat_detail_for_bookings WHERE seatBooking_id = $1`)
	updateSeatBookingQuery := regexp.QuoteMeta(`UPDATE seat_bookings SET status = 'active' WHERE id = $1`)
	deleteSeatBookingQuery := regexp.QuoteMeta(`DELETE FROM seat_bookings WHERE id = $1`)

	suite.mockSql.ExpectBegin()
	tx, err := suite.db.Begin()
	suite.NoError(err)

	suite.mockSql.ExpectExec(updateSeatQuery).WithArgs(seatBookingID).WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mockSql.ExpectExec(deleteSeatDetailQuery).WithArgs(seatBookingID).WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mockSql.ExpectExec(updateSeatBookingQuery).WithArgs(seatBookingID).WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mockSql.ExpectExec(deleteSeatBookingQuery).WithArgs(seatBookingID).WillReturnResult(sqlmock.NewResult(1, 1))

	ginContext, _ := gin.CreateTestContext(nil)
	err = suite.repo.Delete(context.Background(), tx, seatBookingID, ginContext)
	suite.NoError(err)

	suite.mockSql.ExpectCommit()
	err = tx.Commit()
	suite.NoError(err)

	err = suite.mockSql.ExpectationsWereMet()
	suite.NoError(err)
}

func (suite *SeatBookingRepositoryTestSuite) TestDelete_Error() {
	seatBookingID := "1"
	updateSeatQuery := regexp.QuoteMeta(`UPDATE seats SET isAvailable = true WHERE id IN (SELECT seat_id FROM seat_detail_for_bookings WHERE seatBooking_id = $1)`)
	deleteSeatDetailQuery := regexp.QuoteMeta(`DELETE FROM seat_detail_for_bookings WHERE seatBooking_id = $1`)
	updateSeatBookingQuery := regexp.QuoteMeta(`UPDATE seat_bookings SET status = 'active' WHERE id = $1`)
	deleteSeatBookingQuery := regexp.QuoteMeta(`DELETE FROM seat_bookings WHERE id = $1`)

	suite.mockSql.ExpectBegin()
	tx, err := suite.db.Begin()
	suite.NoError(err)

	suite.mockSql.ExpectExec(updateSeatQuery).WithArgs(seatBookingID).WillReturnError(sql.ErrConnDone)
	suite.mockSql.ExpectRollback()
	ginContext, _ := gin.CreateTestContext(nil)
	err = suite.repo.Delete(context.Background(), tx, seatBookingID, ginContext)
	suite.Error(err)

	err = tx.Rollback()
	suite.NoError(err)

	err = suite.mockSql.ExpectationsWereMet()
	suite.NoError(err)

	_ = deleteSeatDetailQuery
	_ = updateSeatBookingQuery
	_ = deleteSeatBookingQuery
}

func (suite *SeatBookingRepositoryTestSuite) TestDelete_DeleteSeatDetailQueryError() {
	seatBookingID := "1"
	updateSeatQuery := regexp.QuoteMeta(`UPDATE seats SET isAvailable = true WHERE id IN (SELECT seat_id FROM seat_detail_for_bookings WHERE seatBooking_id = $1)`)
	deleteSeatDetailQuery := regexp.QuoteMeta(`DELETE FROM seat_detail_for_bookings WHERE seatBooking_id = $1`)

	suite.mockSql.ExpectBegin()
	tx, err := suite.db.Begin()
	suite.NoError(err)

	suite.mockSql.ExpectExec(updateSeatQuery).WithArgs(seatBookingID).WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mockSql.ExpectExec(deleteSeatDetailQuery).WithArgs(seatBookingID).WillReturnError(sql.ErrConnDone)
	suite.mockSql.ExpectRollback()
	ginContext, _ := gin.CreateTestContext(nil)
	err = suite.repo.Delete(context.Background(), tx, seatBookingID, ginContext)
	suite.Error(err)

	err = tx.Rollback()
	suite.NoError(err)

	err = suite.mockSql.ExpectationsWereMet()
	suite.NoError(err)
}

func (suite *SeatBookingRepositoryTestSuite) TestDelete_UpdateSeatBookingQueryError() {
	seatBookingID := "1"
	updateSeatQuery := regexp.QuoteMeta(`UPDATE seats SET isAvailable = true WHERE id IN (SELECT seat_id FROM seat_detail_for_bookings WHERE seatBooking_id = $1)`)
	deleteSeatDetailQuery := regexp.QuoteMeta(`DELETE FROM seat_detail_for_bookings WHERE seatBooking_id = $1`)
	updateSeatBookingQuery := regexp.QuoteMeta(`UPDATE seat_bookings SET status = 'active' WHERE id = $1`)

	suite.mockSql.ExpectBegin()
	tx, err := suite.db.Begin()
	suite.NoError(err)

	suite.mockSql.ExpectExec(updateSeatQuery).WithArgs(seatBookingID).WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mockSql.ExpectExec(deleteSeatDetailQuery).WithArgs(seatBookingID).WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mockSql.ExpectExec(updateSeatBookingQuery).WithArgs(seatBookingID).WillReturnError(sql.ErrConnDone)
	suite.mockSql.ExpectRollback()
	ginContext, _ := gin.CreateTestContext(nil)
	err = suite.repo.Delete(context.Background(), tx, seatBookingID, ginContext)
	suite.Error(err)

	err = tx.Rollback()
	suite.NoError(err)

	err = suite.mockSql.ExpectationsWereMet()
	suite.NoError(err)
}

func (suite *SeatBookingRepositoryTestSuite) TestDelete_DeleteSeatBookingQueryError() {
	seatBookingID := "1"
	updateSeatQuery := regexp.QuoteMeta(`UPDATE seats SET isAvailable = true WHERE id IN (SELECT seat_id FROM seat_detail_for_bookings WHERE seatBooking_id = $1)`)
	deleteSeatDetailQuery := regexp.QuoteMeta(`DELETE FROM seat_detail_for_bookings WHERE seatBooking_id = $1`)
	updateSeatBookingQuery := regexp.QuoteMeta(`UPDATE seat_bookings SET status = 'active' WHERE id = $1`)
	deleteSeatBookingQuery := regexp.QuoteMeta(`DELETE FROM seat_bookings WHERE id = $1`)

	suite.mockSql.ExpectBegin()
	tx, err := suite.db.Begin()
	suite.NoError(err)

	suite.mockSql.ExpectExec(updateSeatQuery).WithArgs(seatBookingID).WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mockSql.ExpectExec(deleteSeatDetailQuery).WithArgs(seatBookingID).WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mockSql.ExpectExec(updateSeatBookingQuery).WithArgs(seatBookingID).WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mockSql.ExpectExec(deleteSeatBookingQuery).WithArgs(seatBookingID).WillReturnError(sql.ErrConnDone)
	suite.mockSql.ExpectRollback()
	ginContext, _ := gin.CreateTestContext(nil)
	err = suite.repo.Delete(context.Background(), tx, seatBookingID, ginContext)
	suite.Error(err)

	err = tx.Rollback()
	suite.NoError(err)

	err = suite.mockSql.ExpectationsWereMet()
	suite.NoError(err)
}

func (suite *SeatBookingRepositoryTestSuite) TestFindByID_NotFoundError() {
	seatBookingID := "1"
	query := regexp.QuoteMeta(`
		SELECT
			sb.id, sb.status, sb.user_id,
			s.id as showtime_id, s.studio_id, s.movie_id, s.show_start, s.show_end,
			st.name as studio_name,
			m.title as movie_title, m.description as movie_description, m.price as movie_price, m.duration as movie_duration, m.status as movie_status,
			sdfb.id as seat_detail_for_booking_id, sdfb.seat_id,
			se.seat_name, se.isAvailable
		FROM
			seat_bookings sb
		JOIN
			showtimes s ON sb.showtime_id = s.id
		JOIN
			studios st ON s.studio_id = st.id
		JOIN
			movies m ON s.movie_id = m.id
		JOIN
			seat_detail_for_bookings sdfb ON sb.id = sdfb.seatBooking_id
		JOIN
			seats se ON sdfb.seat_id = se.id
		WHERE
			sb.id = $1
	`)
	rows := sqlmock.NewRows([]string{
		"id", "status", "user_id", "showtime_id", "studio_id", "movie_id", "show_start", "show_end", "studio_name",
		"movie_title", "movie_description", "movie_price", "movie_duration", "movie_status",
		"seat_detail_for_booking_id", "seat_id", "seat_name", "isAvailable",
	})

	suite.mockSql.ExpectBegin()
	tx, err := suite.db.Begin()
	suite.NoError(err)

	suite.mockSql.ExpectQuery(query).WithArgs(seatBookingID).WillReturnRows(rows)

	ginContext, _ := gin.CreateTestContext(nil)
	_, err = suite.repo.FindByID(context.Background(), tx, seatBookingID, ginContext)
	suite.Error(err)

	suite.mockSql.ExpectRollback()
	err = tx.Rollback()
	suite.NoError(err)

	err = suite.mockSql.ExpectationsWereMet()
	suite.NoError(err)
}

func (suite *SeatBookingRepositoryTestSuite) TestFindAllPendingByUserID_Success() {
	query := regexp.QuoteMeta(`
		SELECT
			sb.id, sb.status, sb.user_id,
			s.id as showtime_id, s.studio_id, s.movie_id, s.show_start, s.show_end,
			st.name as studio_name,
			m.title as movie_title, m.description as movie_description, m.price as movie_price, m.duration as movie_duration, m.status as movie_status,
			sdfb.id as seat_detail_for_booking_id, sdfb.seat_id,
			se.seat_name, se.isAvailable
		FROM
			seat_bookings sb
		JOIN
			showtimes s ON sb.showtime_id = s.id
		JOIN
			studios st ON s.studio_id = st.id
		JOIN
			movies m ON s.movie_id = m.id
		JOIN
			seat_detail_for_bookings sdfb ON sb.id = sdfb.seatBooking_id
		JOIN
			seats se ON sdfb.seat_id = se.id
		WHERE
			sb.user_id = $1 AND sb.status = 'pending'
	`)

	userID := "user1"
	rows := sqlmock.NewRows([]string{
		"id", "status", "user_id",
		"showtime_id", "studio_id", "movie_id", "show_start", "show_end",
		"studio_name",
		"movie_title", "movie_description", "movie_price", "movie_duration", "movie_status",
		"seat_detail_for_booking_id", "seat_id",
		"seat_name", "isAvailable",
	}).AddRow(
		"1", "pending", userID,
		"showtime1", "studio1", "movie1", "2022-07-20 14:00:00", "2022-07-20 16:00:00",
		"Studio 1",
		"Movie 1", "Description", 100, 120, "active",
		"1", "seat1",
		"Seat 1", true,
	)

	suite.mockSql.ExpectBegin()
	suite.mockSql.ExpectQuery(query).WithArgs(userID).WillReturnRows(rows)

	tx, err := suite.db.Begin()
	suite.NoError(err)

	ginContext, _ := gin.CreateTestContext(nil)
	seatBookings, err := suite.repo.FindAllPendingByUserID(context.Background(), tx, userID, ginContext)
	suite.NoError(err)
	suite.Len(seatBookings, 1)
	suite.Equal("1", seatBookings[0].ID)
	suite.Equal("Seat 1", seatBookings[0].SeatName)

	suite.mockSql.ExpectCommit()
	err = tx.Commit()
	suite.NoError(err)

	err = suite.mockSql.ExpectationsWereMet()
	suite.NoError(err)
}

func (suite *SeatBookingRepositoryTestSuite) TestFindAllPendingByUserID_Error() {
	query := regexp.QuoteMeta(`
		SELECT
			sb.id, sb.status, sb.user_id,
			s.id as showtime_id, s.studio_id, s.movie_id, s.show_start, s.show_end,
			st.name as studio_name,
			m.title as movie_title, m.description as movie_description, m.price as movie_price, m.duration as movie_duration, m.status as movie_status,
			sdfb.id as seat_detail_for_booking_id, sdfb.seat_id,
			se.seat_name, se.isAvailable
		FROM
			seat_bookings sb
		JOIN
			showtimes s ON sb.showtime_id = s.id
		JOIN
			studios st ON s.studio_id = st.id
		JOIN
			movies m ON s.movie_id = m.id
		JOIN
			seat_detail_for_bookings sdfb ON sb.id = sdfb.seatBooking_id
		JOIN
			seats se ON sdfb.seat_id = se.id
		WHERE
			sb.user_id = $1 AND sb.status = 'pending'
	`)

	userID := "user1"
	suite.mockSql.ExpectBegin()
	suite.mockSql.ExpectQuery(query).WithArgs(userID).WillReturnError(errors.New("query error"))

	tx, err := suite.db.Begin()
	suite.NoError(err)

	ginContext, _ := gin.CreateTestContext(nil)
	_, err = suite.repo.FindAllPendingByUserID(context.Background(), tx, userID, ginContext)
	suite.Error(err)
	suite.EqualError(err, "query error")

	suite.mockSql.ExpectRollback()
	err = tx.Rollback()
	suite.NoError(err)

	err = suite.mockSql.ExpectationsWereMet()
	suite.NoError(err)
}

func (suite *SeatBookingRepositoryTestSuite) TestFindAllPendingByUserID_NoRows() {
	query := regexp.QuoteMeta(`
		SELECT
			sb.id, sb.status, sb.user_id,
			s.id as showtime_id, s.studio_id, s.movie_id, s.show_start, s.show_end,
			st.name as studio_name,
			m.title as movie_title, m.description as movie_description, m.price as movie_price, m.duration as movie_duration, m.status as movie_status,
			sdfb.id as seat_detail_for_booking_id, sdfb.seat_id,
			se.seat_name, se.isAvailable
		FROM
			seat_bookings sb
		JOIN
			showtimes s ON sb.showtime_id = s.id
		JOIN
			studios st ON s.studio_id = st.id
		JOIN
			movies m ON s.movie_id = m.id
		JOIN
			seat_detail_for_bookings sdfb ON sb.id = sdfb.seatBooking_id
		JOIN
			seats se ON sdfb.seat_id = se.id
		WHERE
			sb.user_id = $1 AND sb.status = 'pending'
	`)

	userID := "user1"
	rows := sqlmock.NewRows([]string{
		"id", "status", "user_id",
		"showtime_id", "studio_id", "movie_id", "show_start", "show_end",
		"studio_name",
		"movie_title", "movie_description", "movie_price", "movie_duration", "movie_status",
		"seat_detail_for_booking_id", "seat_id",
		"seat_name", "isAvailable",
	})

	suite.mockSql.ExpectBegin()
	suite.mockSql.ExpectQuery(query).WithArgs(userID).WillReturnRows(rows)

	tx, err := suite.db.Begin()
	suite.NoError(err)

	ginContext, _ := gin.CreateTestContext(nil)
	seatBookings, err := suite.repo.FindAllPendingByUserID(context.Background(), tx, userID, ginContext)
	suite.NoError(err)
	suite.Len(seatBookings, 0)

	suite.mockSql.ExpectCommit()
	err = tx.Commit()
	suite.NoError(err)

	err = suite.mockSql.ExpectationsWereMet()
	suite.NoError(err)
}

func (suite *SeatBookingRepositoryTestSuite) TestUpdate_Success() {
	query := regexp.QuoteMeta("UPDATE seat_bookings SET status = $1 WHERE id = $2")

	seatBooking := entity.SeatBooking{ID: "1", SeatBookingStatus: "confirmed"}

	suite.mockSql.ExpectBegin()
	suite.mockSql.ExpectExec(query).WithArgs(seatBooking.SeatBookingStatus, seatBooking.ID).WillReturnResult(sqlmock.NewResult(1, 1))

	tx, err := suite.db.Begin()
	suite.NoError(err)

	ginContext, _ := gin.CreateTestContext(nil)
	seatBooking, err = suite.repo.Update(context.Background(), tx, seatBooking, ginContext)
	suite.NoError(err)
	suite.Equal("confirmed", seatBooking.SeatBookingStatus)

	suite.mockSql.ExpectCommit()
	err = tx.Commit()
	suite.NoError(err)

	err = suite.mockSql.ExpectationsWereMet()
	suite.NoError(err)
}

func (suite *SeatBookingRepositoryTestSuite) TestUpdate_Error() {
	query := regexp.QuoteMeta("UPDATE seat_bookings SET status = $1 WHERE id = $2")

	seatBooking := entity.SeatBooking{ID: "1", SeatBookingStatus: "confirmed"}

	suite.mockSql.ExpectBegin()
	suite.mockSql.ExpectExec(query).WithArgs(seatBooking.SeatBookingStatus, seatBooking.ID).WillReturnError(errors.New("update error"))

	tx, err := suite.db.Begin()
	suite.NoError(err)

	ginContext, _ := gin.CreateTestContext(nil)
	_, err = suite.repo.Update(context.Background(), tx, seatBooking, ginContext)
	suite.Error(err)
	suite.EqualError(err, "update error")

	suite.mockSql.ExpectRollback()
	err = tx.Rollback()
	suite.NoError(err)

	err = suite.mockSql.ExpectationsWereMet()
	suite.NoError(err)
}

func TestSeatBookingRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(SeatBookingRepositoryTestSuite))
}
