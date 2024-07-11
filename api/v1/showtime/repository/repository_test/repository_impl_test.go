package repository_test

import (
	"bioskuy/api/v1/showtime/entity"
	entityStudio "bioskuy/api/v1/studio/entity"

	"bioskuy/api/v1/showtime/repository"
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type ShowtimeRepositoryTestSuite struct {
	suite.Suite
	mockSql sqlmock.Sqlmock
	db      *sql.DB
	repo    repository.ShowtimeRepository
}

func (suite *ShowtimeRepositoryTestSuite) SetupTest() {
	var err error
	suite.db, suite.mockSql, err = sqlmock.New()
	suite.NoError(err)
	suite.repo = repository.NewShowtimeRepository()
}

func (suite *ShowtimeRepositoryTestSuite) TearDownTest() {
	suite.db.Close()
}

func (suite *ShowtimeRepositoryTestSuite) TestSave_Success() {
	showtime := entity.Showtime{
		MovieID:   "1",
		StudioID:  "1",
		ShowStart: time.Now(),
		ShowEnd:   time.Now().Add(2 * time.Hour),
	}

	query := regexp.QuoteMeta("INSERT INTO showtimes (movie_id, studio_id, show_start, show_end) VALUES ($1, $2, $3, $4) RETURNING id")

	suite.mockSql.ExpectBegin()
	tx, err := suite.db.Begin()
	suite.NoError(err)

	suite.mockSql.ExpectQuery(query).
		WithArgs(showtime.MovieID, showtime.StudioID, showtime.ShowStart, showtime.ShowEnd).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))

	ginContext, _ := gin.CreateTestContext(nil)
	savedShowtime, err := suite.repo.Save(context.Background(), tx, showtime, ginContext)
	suite.NoError(err)
	suite.Equal("1", savedShowtime.ID)

	suite.mockSql.ExpectCommit()
	err = tx.Commit()
	suite.NoError(err)

	err = suite.mockSql.ExpectationsWereMet()
	suite.NoError(err)
}

func (suite *ShowtimeRepositoryTestSuite) TestSave_Error() {
	showtime := entity.Showtime{
		MovieID:   "1",
		StudioID:  "1",
		ShowStart: time.Now(),
		ShowEnd:   time.Now().Add(2 * time.Hour),
	}

	query := regexp.QuoteMeta("INSERT INTO showtimes (movie_id, studio_id, show_start, show_end) VALUES ($1, $2, $3, $4) RETURNING id")

	suite.mockSql.ExpectBegin()
	tx, err := suite.db.Begin()
	suite.NoError(err)

	suite.mockSql.ExpectQuery(query).
		WithArgs(showtime.MovieID, showtime.StudioID, showtime.ShowStart, showtime.ShowEnd).
		WillReturnError(errors.New("query error"))

	ginContext, _ := gin.CreateTestContext(nil)
	_, err = suite.repo.Save(context.Background(), tx, showtime, ginContext)
	suite.Error(err)

	suite.mockSql.ExpectRollback()
	err = tx.Rollback()
	suite.NoError(err)

	err = suite.mockSql.ExpectationsWereMet()
	suite.NoError(err)
}

func (suite *ShowtimeRepositoryTestSuite) TestFindByID_Success() {
	showtimeID := "1"
	expectedShowtime := entity.Showtime{
		ID:               "1",
		MovieID:          "1",
		StudioID:         "1",
		ShowStart:        time.Now(),
		ShowEnd:          time.Now().Add(2 * time.Hour),
		StudioName:       "Studio 1",
		MovieTitle:       "Movie 1",
		MovieDescription: "Description 1",
		MoviePrice:       10000,
		MovieDuration:    120,
		MovieStatus:      "Active",
	}

	query := regexp.QuoteMeta(`SELECT s.id, s.studio_id, s.movie_id, s.show_start, s.show_end, st.name as studio_name, m.title as movie_title, m.description as movie_description, m.price as movie_price, m.duration as movie_duration, m.status as movie_status FROM showtimes s JOIN studios st ON s.studio_id = st.id JOIN movies m ON s.movie_id = m.id WHERE s.id = $1`)

	suite.mockSql.ExpectBegin()
	tx, err := suite.db.Begin()
	suite.NoError(err)

	suite.mockSql.ExpectQuery(query).WithArgs(showtimeID).WillReturnRows(sqlmock.NewRows([]string{"id", "studio_id", "movie_id", "show_start", "show_end", "studio_name", "movie_title", "movie_description", "movie_price", "movie_duration", "movie_status"}).
		AddRow(expectedShowtime.ID, expectedShowtime.StudioID, expectedShowtime.MovieID, expectedShowtime.ShowStart, expectedShowtime.ShowEnd, expectedShowtime.StudioName, expectedShowtime.MovieTitle, expectedShowtime.MovieDescription, expectedShowtime.MoviePrice, expectedShowtime.MovieDuration, expectedShowtime.MovieStatus))

	ginContext, _ := gin.CreateTestContext(nil)
	foundShowtime, err := suite.repo.FindByID(context.Background(), tx, showtimeID, ginContext)
	suite.NoError(err)
	suite.Equal(expectedShowtime, foundShowtime)

	suite.mockSql.ExpectCommit()
	err = tx.Commit()
	suite.NoError(err)

	err = suite.mockSql.ExpectationsWereMet()
	suite.NoError(err)
}

func (suite *ShowtimeRepositoryTestSuite) TestFindByID_NotFound() {
	showtimeID := "1"
	query := regexp.QuoteMeta(`SELECT s.id, s.studio_id, s.movie_id, s.show_start, s.show_end, st.name as studio_name, m.title as movie_title, m.description as movie_description, m.price as movie_price, m.duration as movie_duration, m.status as movie_status FROM showtimes s JOIN studios st ON s.studio_id = st.id JOIN movies m ON s.movie_id = m.id WHERE s.id = $1`)

	suite.mockSql.ExpectBegin()
	tx, err := suite.db.Begin()
	suite.NoError(err)

	suite.mockSql.ExpectQuery(query).WithArgs(showtimeID).WillReturnError(sql.ErrNoRows)

	ginContext, _ := gin.CreateTestContext(nil)
	_, err = suite.repo.FindByID(context.Background(), tx, showtimeID, ginContext)
	suite.Error(err)

	suite.mockSql.ExpectRollback()
	err = tx.Rollback()
	suite.NoError(err)

	err = suite.mockSql.ExpectationsWereMet()
	suite.NoError(err)
}

func (suite *ShowtimeRepositoryTestSuite) TestFindByID_ScanError() {
	showtimeID := "1"
	query := regexp.QuoteMeta(`SELECT s.id, s.studio_id, s.movie_id, s.show_start, s.show_end, st.name as studio_name, m.title as movie_title, m.description as movie_description, m.price as movie_price, m.duration as movie_duration, m.status as movie_status FROM showtimes s JOIN studios st ON s.studio_id = st.id JOIN movies m ON s.movie_id = m.id WHERE s.id = $1`)

	suite.mockSql.ExpectBegin()
	tx, err := suite.db.Begin()
	suite.NoError(err)

	suite.mockSql.ExpectQuery(query).WithArgs(showtimeID).WillReturnRows(sqlmock.NewRows([]string{"id", "studio_id", "movie_id", "show_start", "show_end", "studio_name", "movie_title", "movie_description", "movie_price", "movie_duration", "movie_status"}).
		AddRow("invalid_id", "studio_id", "movie_id", "show_start", "show_end", "studio_name", "movie_title", "movie_description", "movie_price", "movie_duration", "movie_status"))

	ginContext, _ := gin.CreateTestContext(nil)
	foundShowtime, err := suite.repo.FindByID(context.Background(), tx, showtimeID, ginContext)
	suite.Error(err)
	suite.Empty(foundShowtime)

	suite.mockSql.ExpectRollback()
	err = tx.Rollback()
	suite.NoError(err)

	err = suite.mockSql.ExpectationsWereMet()
	suite.NoError(err)
}

func (suite *ShowtimeRepositoryTestSuite) TestFindByID_ShowtimeNotFound() {
	showtimeID := "1"
	query := regexp.QuoteMeta(`SELECT s.id, s.studio_id, s.movie_id, s.show_start, s.show_end, st.name as studio_name, m.title as movie_title, m.description as movie_description, m.price as movie_price, m.duration as movie_duration, m.status as movie_status FROM showtimes s JOIN studios st ON s.studio_id = st.id JOIN movies m ON s.movie_id = m.id WHERE s.id = $1`)

	suite.mockSql.ExpectBegin()
	tx, err := suite.db.Begin()
	suite.NoError(err)

	suite.mockSql.ExpectQuery(query).WithArgs(showtimeID).WillReturnRows(sqlmock.NewRows([]string{"id", "studio_id", "movie_id", "show_start", "show_end", "studio_name", "movie_title", "movie_description", "movie_price", "movie_duration", "movie_status"}))

	ginContext, _ := gin.CreateTestContext(nil)
	foundShowtime, err := suite.repo.FindByID(context.Background(), tx, showtimeID, ginContext)
	suite.Error(err)
	suite.EqualError(err, "showtime not found")
	suite.Empty(foundShowtime)

	suite.mockSql.ExpectRollback()
	err = tx.Rollback()
	suite.NoError(err)

	err = suite.mockSql.ExpectationsWereMet()
	suite.NoError(err)
}

func (suite *ShowtimeRepositoryTestSuite) TestFindAll_Success() {
	query := regexp.QuoteMeta(`
        SELECT s.id, s.studio_id, s.movie_id, s.show_start, s.show_end, st.name as studio_name, m.title as movie_title, m.description as movie_description, m.price as movie_price, m.duration as movie_duration, m.status as movie_status
        FROM showtimes s
        JOIN studios st ON s.studio_id = st.id
        JOIN movies m ON s.movie_id = m.id
    `)

	rows := sqlmock.NewRows([]string{"id", "studio_id", "movie_id", "show_start", "show_end", "studio_name", "movie_title", "movie_description", "movie_price", "movie_duration", "movie_status"}).
		AddRow("1", "1", "1", time.Now(), time.Now().Add(2*time.Hour), "Studio 1", "Movie 1", "Description 1", 10000, 120, "Active").
		AddRow("2", "2", "2", time.Now(), time.Now().Add(3*time.Hour), "Studio 2", "Movie 2", "Description 2", 15000, 150, "Active")

	suite.mockSql.ExpectBegin()
	tx, err := suite.db.Begin()
	suite.NoError(err)

	suite.mockSql.ExpectQuery(query).WillReturnRows(rows)

	ginContext, _ := gin.CreateTestContext(nil)
	showtimes, err := suite.repo.FindAll(context.Background(), tx, ginContext)
	suite.NoError(err)
	suite.Len(showtimes, 2)

	suite.mockSql.ExpectCommit()
	err = tx.Commit()
	suite.NoError(err)

	err = suite.mockSql.ExpectationsWereMet()
	suite.NoError(err)
}

func (suite *ShowtimeRepositoryTestSuite) TestFindAll_Error() {
	query := regexp.QuoteMeta(`
        SELECT s.id, s.studio_id, s.movie_id, s.show_start, s.show_end, st.name as studio_name, m.title as movie_title, m.description as movie_description, m.price as movie_price, m.duration as movie_duration, m.status as movie_status
        FROM showtimes s
        JOIN studios st ON s.studio_id = st.id
        JOIN movies m ON s.movie_id = m.id
    `)

	suite.mockSql.ExpectBegin()
	tx, err := suite.db.Begin()
	suite.NoError(err)

	suite.mockSql.ExpectQuery(query).WillReturnError(errors.New("query error"))

	ginContext, _ := gin.CreateTestContext(nil)
	_, err = suite.repo.FindAll(context.Background(), tx, ginContext)
	suite.Error(err)

	suite.mockSql.ExpectRollback()
	err = tx.Rollback()
	suite.NoError(err)

	err = suite.mockSql.ExpectationsWereMet()
	suite.NoError(err)
}

func (suite *ShowtimeRepositoryTestSuite) TestDelete_Success() {
	showtimeID := "1"
	query := regexp.QuoteMeta(`DELETE FROM showtimes WHERE id = $1`)

	suite.mockSql.ExpectBegin()
	tx, err := suite.db.Begin()
	suite.NoError(err)

	suite.mockSql.ExpectExec(query).WithArgs(showtimeID).WillReturnResult(sqlmock.NewResult(1, 1))

	ginContext, _ := gin.CreateTestContext(nil)
	err = suite.repo.Delete(context.Background(), tx, showtimeID, ginContext)
	suite.NoError(err)

	suite.mockSql.ExpectCommit()
	err = tx.Commit()
	suite.NoError(err)

	err = suite.mockSql.ExpectationsWereMet()
	suite.NoError(err)
}

func (suite *ShowtimeRepositoryTestSuite) TestDelete_Error() {
	showtimeID := "1"
	query := regexp.QuoteMeta(`DELETE FROM showtimes WHERE id = $1`)

	suite.mockSql.ExpectBegin()
	tx, err := suite.db.Begin()
	suite.NoError(err)

	suite.mockSql.ExpectExec(query).WithArgs(showtimeID).WillReturnError(errors.New("delete error"))

	ginContext, _ := gin.CreateTestContext(nil)
	err = suite.repo.Delete(context.Background(), tx, showtimeID, ginContext)
	suite.Error(err)

	suite.mockSql.ExpectRollback()
	err = tx.Rollback()
	suite.NoError(err)

	err = suite.mockSql.ExpectationsWereMet()
	suite.NoError(err)
}

func (suite *ShowtimeRepositoryTestSuite) TestFindConflictingShowtimes_Success() {
	studio := entityStudio.Studio{ID: "1"}
	showtime := entity.Showtime{
		ShowStart: time.Now(),
		ShowEnd:   time.Now().Add(2 * time.Hour),
	}

	query := `SELECT s.id
              FROM showtimes s
              WHERE s.studio_id = \$1
              AND (
                  (s.show_start >= \$2 AND s.show_start <= \$3) OR
                  (s.show_end >= \$2 AND s.show_end <= \$3) OR
                  (s.show_start <= \$2 AND s.show_end >= \$2) OR
                  (s.show_start <= \$3 AND s.show_end >= \$3)
              )`

	suite.mockSql.ExpectBegin()
	tx, err := suite.db.Begin()
	suite.NoError(err)

	suite.mockSql.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(studio.ID, showtime.ShowStart, showtime.ShowEnd).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	ginContext, _ := gin.CreateTestContext(nil)
	err = suite.repo.FindConflictingShowtimes(context.Background(), tx, studio, showtime, ginContext)
	suite.NoError(err)

	suite.mockSql.ExpectCommit()
	err = tx.Commit()
	suite.NoError(err)

	err = suite.mockSql.ExpectationsWereMet()
	suite.NoError(err)
}

func TestShowtimeRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(ShowtimeRepositoryTestSuite))
}
