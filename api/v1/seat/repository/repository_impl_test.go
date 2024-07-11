package repository

import (
	"bioskuy/api/v1/seat/entity"
	"context"
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type SeatRepositoryTestSuite struct {
	suite.Suite
	mockSql sqlmock.Sqlmock
	db      *sql.DB
	repo    SeatRepository
}

func (suite *SeatRepositoryTestSuite) SetupTest() {
	var err error
	suite.db, suite.mockSql, err = sqlmock.New()
	suite.NoError(err)
	suite.repo = NewSeatRepository()
}

func (suite *SeatRepositoryTestSuite) TearDownTest() {
	suite.db.Close()
}

func (suite *SeatRepositoryTestSuite) TestSave_Success() {
	seat := entity.Seat{Name: "Test Seat", IsAvailable: true, StudioID: "studio1"}
	query := regexp.QuoteMeta("INSERT INTO seats (seat_name, isAvailable, studio_id) VALUES ($1, $2, $3) RETURNING id")

	suite.mockSql.ExpectBegin()
	tx, err := suite.db.Begin()
	suite.NoError(err)

	suite.mockSql.ExpectQuery(query).
		WithArgs(seat.Name, seat.IsAvailable, seat.StudioID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))

	ginContext, _ := gin.CreateTestContext(nil)
	savedSeat, err := suite.repo.Save(context.Background(), tx, seat, ginContext)
	suite.NoError(err)
	suite.Equal("1", savedSeat.ID)

	suite.mockSql.ExpectCommit()
	err = tx.Commit()
	suite.NoError(err)

	err = suite.mockSql.ExpectationsWereMet()
	suite.NoError(err)
}

func (suite *SeatRepositoryTestSuite) TestSave_Error() {
	seat := entity.Seat{Name: "Test Seat", IsAvailable: true, StudioID: "studio1"}
	query := regexp.QuoteMeta("INSERT INTO seats (seat_name, isAvailable, studio_id) VALUES ($1, $2, $3) RETURNING id")

	suite.mockSql.ExpectBegin()
	tx, err := suite.db.Begin()
	suite.NoError(err)

	suite.mockSql.ExpectQuery(query).
		WithArgs(seat.Name, seat.IsAvailable, seat.StudioID).
		WillReturnError(sql.ErrConnDone)

	ginContext, _ := gin.CreateTestContext(nil)
	_, err = suite.repo.Save(context.Background(), tx, seat, ginContext)
	suite.Error(err)

	suite.mockSql.ExpectRollback()
	err = tx.Rollback()
	suite.NoError(err)

	err = suite.mockSql.ExpectationsWereMet()
	suite.NoError(err)
}

func (suite *SeatRepositoryTestSuite) TestFindByID_Success() {
	seatID := "1"
	query := regexp.QuoteMeta(`SELECT id, seat_name, isAvailable, studio_id FROM seats WHERE id = $1`)
	rows := sqlmock.NewRows([]string{"id", "seat_name", "isAvailable", "studio_id"}).
		AddRow(seatID, "Test Seat", true, "studio1")

	suite.mockSql.ExpectBegin()
	tx, err := suite.db.Begin()
	suite.NoError(err)

	suite.mockSql.ExpectQuery(query).WithArgs(seatID).WillReturnRows(rows)

	ginContext, _ := gin.CreateTestContext(nil)
	seat, err := suite.repo.FindByID(context.Background(), tx, seatID, ginContext)
	suite.NoError(err)
	suite.Equal(seatID, seat.ID)
	suite.Equal("Test Seat", seat.Name)

	suite.mockSql.ExpectCommit()
	err = tx.Commit()
	suite.NoError(err)

	err = suite.mockSql.ExpectationsWereMet()
	suite.NoError(err)
}

func (suite *SeatRepositoryTestSuite) TestFindByID_Error() {
	seatID := "1"
	query := regexp.QuoteMeta(`SELECT id, seat_name, isAvailable, studio_id FROM seats WHERE id = $1`)

	suite.mockSql.ExpectBegin()
	tx, err := suite.db.Begin()
	suite.NoError(err)

	suite.mockSql.ExpectQuery(query).WithArgs(seatID).WillReturnError(sql.ErrNoRows)

	ginContext, _ := gin.CreateTestContext(nil)
	_, err = suite.repo.FindByID(context.Background(), tx, seatID, ginContext)
	suite.Error(err)

	suite.mockSql.ExpectRollback()
	err = tx.Rollback()
	suite.NoError(err)

	err = suite.mockSql.ExpectationsWereMet()
	suite.NoError(err)
}

func (suite *SeatRepositoryTestSuite) TestFindByIDWithNotAvailable_Success() {
	seatID := "1"
	query := regexp.QuoteMeta(`SELECT id, seat_name, isAvailable, studio_id FROM seats WHERE id = $1 AND isAvailable = true`)
	rows := sqlmock.NewRows([]string{"id", "seat_name", "isAvailable", "studio_id"}).
		AddRow(seatID, "Test Seat", true, "studio1")

	suite.mockSql.ExpectBegin()
	tx, err := suite.db.Begin()
	suite.NoError(err)

	suite.mockSql.ExpectQuery(query).WithArgs(seatID).WillReturnRows(rows)

	ginContext, _ := gin.CreateTestContext(nil)
	seat, err := suite.repo.FindByIDWithNotAvailable(context.Background(), tx, seatID, ginContext)
	suite.NoError(err)
	suite.Equal(seatID, seat.ID)
	suite.Equal("Test Seat", seat.Name)
	suite.Equal(true, seat.IsAvailable)

	suite.mockSql.ExpectCommit()
	err = tx.Commit()
	suite.NoError(err)

	err = suite.mockSql.ExpectationsWereMet()
	suite.NoError(err)
}

func (suite *SeatRepositoryTestSuite) TestFindByIDWithNotAvailable_Error() {
	seatID := "1"
	query := regexp.QuoteMeta(`SELECT id, seat_name, isAvailable, studio_id FROM seats WHERE id = $1 AND isAvailable = true`)

	suite.mockSql.ExpectBegin()
	tx, err := suite.db.Begin()
	suite.NoError(err)

	suite.mockSql.ExpectQuery(query).WithArgs(seatID).WillReturnError(sql.ErrNoRows)

	ginContext, _ := gin.CreateTestContext(nil)
	_, err = suite.repo.FindByIDWithNotAvailable(context.Background(), tx, seatID, ginContext)
	suite.Error(err)

	suite.mockSql.ExpectRollback()
	err = tx.Rollback()
	suite.NoError(err)

	err = suite.mockSql.ExpectationsWereMet()
	suite.NoError(err)
}

func (suite *SeatRepositoryTestSuite) TestUpdate_Success() {
	seat := entity.Seat{ID: "1", IsAvailable: false}
	query := regexp.QuoteMeta(`UPDATE seats SET isAvailable = $1 WHERE id = $2`)

	suite.mockSql.ExpectBegin()
	tx, err := suite.db.Begin()
	suite.NoError(err)

	suite.mockSql.ExpectExec(query).WithArgs(seat.IsAvailable, seat.ID).WillReturnResult(sqlmock.NewResult(1, 1))

	ginContext, _ := gin.CreateTestContext(nil)
	updatedSeat, err := suite.repo.Update(context.Background(), tx, seat, ginContext)
	suite.NoError(err)
	suite.Equal(seat.ID, updatedSeat.ID)
	suite.Equal(seat.IsAvailable, updatedSeat.IsAvailable)

	suite.mockSql.ExpectCommit()
	err = tx.Commit()
	suite.NoError(err)

	err = suite.mockSql.ExpectationsWereMet()
	suite.NoError(err)
}

func (suite *SeatRepositoryTestSuite) TestUpdate_Error() {
	seat := entity.Seat{ID: "1", IsAvailable: false}
	query := regexp.QuoteMeta(`UPDATE seats SET isAvailable = $1 WHERE id = $2`)

	suite.mockSql.ExpectBegin()
	tx, err := suite.db.Begin()
	suite.NoError(err)

	suite.mockSql.ExpectExec(query).WithArgs(seat.IsAvailable, seat.ID).WillReturnError(sql.ErrConnDone)

	ginContext, _ := gin.CreateTestContext(nil)
	_, err = suite.repo.Update(context.Background(), tx, seat, ginContext)
	suite.Error(err)

	suite.mockSql.ExpectRollback()
	err = tx.Rollback()
	suite.NoError(err)

	err = suite.mockSql.ExpectationsWereMet()
	suite.NoError(err)
}

func (suite *SeatRepositoryTestSuite) TestFindAll_Success() {
	studioID := "studio1"
	query := regexp.QuoteMeta(`SELECT id, seat_name, isAvailable, studio_id FROM seats WHERE studio_id = $1`)
	rows := sqlmock.NewRows([]string{"id", "seat_name", "isAvailable", "studio_id"}).
		AddRow("1", "Seat 1", true, studioID).
		AddRow("2", "Seat 2", true, studioID)

	suite.mockSql.ExpectBegin()
	tx, err := suite.db.Begin()
	suite.NoError(err)

	suite.mockSql.ExpectQuery(query).WithArgs(studioID).WillReturnRows(rows)

	ginContext, _ := gin.CreateTestContext(nil)
	seats, err := suite.repo.FindAll(context.Background(), studioID, tx, ginContext)
	suite.NoError(err)
	suite.Len(seats, 2)

	suite.mockSql.ExpectCommit()
	err = tx.Commit()
	suite.NoError(err)

	err = suite.mockSql.ExpectationsWereMet()
	suite.NoError(err)
}

func (suite *SeatRepositoryTestSuite) TestFindAll_Error() {
	studioID := "studio1"
	query := regexp.QuoteMeta(`SELECT id, seat_name, isAvailable, studio_id FROM seats WHERE studio_id = $1`)

	suite.mockSql.ExpectBegin()
	tx, err := suite.db.Begin()
	suite.NoError(err)

	suite.mockSql.ExpectQuery(query).WithArgs(studioID).WillReturnError(sql.ErrConnDone)

	ginContext, _ := gin.CreateTestContext(nil)
	_, err = suite.repo.FindAll(context.Background(), studioID, tx, ginContext)
	suite.Error(err)

	suite.mockSql.ExpectRollback()
	err = tx.Rollback()
	suite.NoError(err)

	err = suite.mockSql.ExpectationsWereMet()
	suite.NoError(err)
}

func (suite *SeatRepositoryTestSuite) TestDelete_Success() {
	seatID := "1"
	query := regexp.QuoteMeta(`DELETE FROM seats WHERE studio_id = $1`)

	suite.mockSql.ExpectBegin()
	tx, err := suite.db.Begin()
	suite.NoError(err)

	suite.mockSql.ExpectExec(query).WithArgs(seatID).WillReturnResult(sqlmock.NewResult(1, 1))

	ginContext, _ := gin.CreateTestContext(nil)
	err = suite.repo.Delete(context.Background(), tx, seatID, ginContext)
	suite.NoError(err)

	suite.mockSql.ExpectCommit()
	err = tx.Commit()
	suite.NoError(err)

	err = suite.mockSql.ExpectationsWereMet()
	suite.NoError(err)
}

func (suite *SeatRepositoryTestSuite) TestDelete_Error() {
	seatID := "1"
	query := regexp.QuoteMeta(`DELETE FROM seats WHERE studio_id = $1`)

	suite.mockSql.ExpectBegin()
	tx, err := suite.db.Begin()
	suite.NoError(err)

	suite.mockSql.ExpectExec(query).WithArgs(seatID).WillReturnError(sql.ErrConnDone)

	ginContext, _ := gin.CreateTestContext(nil)
	err = suite.repo.Delete(context.Background(), tx, seatID, ginContext)
	suite.Error(err)

	suite.mockSql.ExpectRollback()
	err = tx.Rollback()
	suite.NoError(err)

	err = suite.mockSql.ExpectationsWereMet()
	suite.NoError(err)
}

func (suite *SeatRepositoryTestSuite) TestFindByID_NotFoundError() {
	seatID := "1"
	query := regexp.QuoteMeta(`SELECT id, seat_name, isAvailable, studio_id FROM seats WHERE id = $1`)
	rows := sqlmock.NewRows([]string{"id", "seat_name", "isAvailable", "studio_id"})

	suite.mockSql.ExpectBegin()
	tx, err := suite.db.Begin()
	suite.NoError(err)

	suite.mockSql.ExpectQuery(query).WithArgs(seatID).WillReturnRows(rows)

	ginContext, _ := gin.CreateTestContext(nil)
	_, err = suite.repo.FindByID(context.Background(), tx, seatID, ginContext)
	suite.Error(err)

	suite.mockSql.ExpectRollback()
	err = tx.Rollback()
	suite.NoError(err)

	err = suite.mockSql.ExpectationsWereMet()
	suite.NoError(err)
}

func TestSeatRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(SeatRepositoryTestSuite))
}
