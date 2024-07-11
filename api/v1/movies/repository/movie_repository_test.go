package repository

import (
	"bioskuy/api/v1/genre/dto"
	"bioskuy/api/v1/movies/entity"
	"database/sql"
	"errors"
	"math"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MovieRepositoryTestSuite struct {
	suite.Suite
	mockDb  *sql.DB
	mockSql sqlmock.Sqlmock
	repo    MovieRepository
}

var mockingMovie = entity.Movie{
	ID:          "1",
	Title:       "Inception",
	Description: "A mind-bending thriller",
	Price:       100,
	Duration:    120,
	Status:      "AVAILABLE",
}

func (suite *MovieRepositoryTestSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	assert.NoError(suite.T(), err)
	suite.mockDb = db
	suite.mockSql = mock
	suite.repo = NewMovieRepository(suite.mockDb)
}

func TestMovieRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(MovieRepositoryTestSuite))
}

func (suite *MovieRepositoryTestSuite) TestGetAll_Success() {
	page, size := 1, 10
	totalRows := 15

	suite.mockSql.ExpectQuery(`SELECT id, title, description, price, duration, status FROM movies LIMIT \$1 OFFSET \$2`).WithArgs(size, (page-1)*size).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "price", "duration", "status"}).
			AddRow(mockingMovie.ID, mockingMovie.Title, mockingMovie.Description, mockingMovie.Price, mockingMovie.Duration, mockingMovie.Status).
			AddRow("1", "Avatar", "A sci-fi adventure", 150, 180, "AVAILABLE"))

	suite.mockSql.ExpectQuery(`SELECT COUNT\(\*\) FROM movies`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(totalRows))

	movies, paging, err := suite.repo.GetAll(page, size)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), movies, 2)
	assert.Equal(suite.T(), dto.Paging{
		Page:       page,
		Size:       size,
		TotalRows:  totalRows,
		TotalPages: int(math.Ceil(float64(totalRows) / float64(size))),
	}, paging)
}

func (suite *MovieRepositoryTestSuite) TestGetAll_ErrorOnQueryMovies() {
	page, size := 1, 10

	suite.mockSql.ExpectQuery(`SELECT id, title, description, price, duration, status FROM movies LIMIT \$1 OFFSET \$2`).WithArgs(size, (page-1)*size).
		WillReturnError(errors.New("Query Error"))

	_, _, err := suite.repo.GetAll(page, size)
	assert.Error(suite.T(), err)
}

func (suite *MovieRepositoryTestSuite) TestGetAll_ErrorOnCountMovies() {
	page, size := 1, 10

	suite.mockSql.ExpectQuery(`SELECT id, title, description, price, duration, status FROM movies LIMIT \$1 OFFSET \$2`).WithArgs(size, (page-1)*size).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "price", "duration", "status"}).
			AddRow(mockingMovie.ID, mockingMovie.Title, mockingMovie.Description, mockingMovie.Price, mockingMovie.Duration, mockingMovie.Status))

	suite.mockSql.ExpectQuery(`SELECT COUNT\(\*\) FROM movies`).WillReturnError(errors.New("Count Error"))

	_, _, err := suite.repo.GetAll(page, size)
	assert.Error(suite.T(), err)
}

func (suite *MovieRepositoryTestSuite) TestGetAll_ErrorOnScan() {
	page, size := 1, 10

	suite.mockSql.ExpectQuery(`SELECT id, title, description, price, duration, status FROM movies LIMIT \$1 OFFSET \$2`).WithArgs(size, (page-1)*size).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "price", "duration", "status"}).
			AddRow(mockingMovie.ID, nil, mockingMovie.Description, mockingMovie.Price, mockingMovie.Duration, mockingMovie.Status)) // nil akan menyebabkan scan error

	suite.mockSql.ExpectQuery(`SELECT COUNT\(\*\) FROM movies`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	_, _, err := suite.repo.GetAll(page, size)
	assert.Error(suite.T(), err)
}

func (suite *MovieRepositoryTestSuite) TestCreate_Success() {
	suite.mockSql.ExpectQuery(`INSERT INTO movies \(title, description, price, duration, status\) VALUES \(\$1, \$2, \$3, \$4, \$5\) RETURNING id`).
		WithArgs(mockingMovie.Title, mockingMovie.Description, mockingMovie.Price, mockingMovie.Duration, mockingMovie.Status).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(mockingMovie.ID))

	result, err := suite.repo.Create(mockingMovie)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), mockingMovie.Title, result.Title)
}

func (suite *MovieRepositoryTestSuite) TestCreate_Failed() {
	suite.mockSql.ExpectQuery(`INSERT INTO movies \(title, description, price, duration, status\) VALUES \(\$1, \$2, \$3, \$4, \$5\) RETURNING id`).
		WithArgs(mockingMovie.Title, mockingMovie.Description, mockingMovie.Price, mockingMovie.Duration, mockingMovie.Status).
		WillReturnError(errors.New("Insert Movie Failed"))

	_, err := suite.repo.Create(mockingMovie)
	assert.Error(suite.T(), err)
}

func (suite *MovieRepositoryTestSuite) TestGetByID_ErrorOnQuery() {
	id := "1"

	suite.mockSql.ExpectQuery(`SELECT id, title, description, price, duration, status FROM movies WHERE id = \$1`).WithArgs(id).
		WillReturnError(errors.New("Query Error"))

	_, err := suite.repo.GetByID(id)
	assert.Error(suite.T(), err)
}

func (suite *MovieRepositoryTestSuite) TestGetByID_Success() {
	id := "1"

	suite.mockSql.ExpectQuery(`SELECT id, title, description, price, duration, status FROM movies WHERE id = \$1`).WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "price", "duration", "status"}).AddRow(mockingMovie.ID, mockingMovie.Title, mockingMovie.Description, mockingMovie.Price, mockingMovie.Duration, mockingMovie.Status))

	movie, err := suite.repo.GetByID(id)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockingMovie, movie)
}

func (suite *MovieRepositoryTestSuite) TestUpdate_Success() {
	suite.mockSql.ExpectExec(`UPDATE movies SET title = \$1, description = \$2, price = \$3, duration = \$4, status = \$5 WHERE id = \$6`).
		WithArgs(mockingMovie.Title, mockingMovie.Description, mockingMovie.Price, mockingMovie.Duration, mockingMovie.Status, mockingMovie.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := suite.repo.Update(mockingMovie)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockingMovie, result)
}

func (suite *MovieRepositoryTestSuite) TestUpdate_Failed() {
	suite.mockSql.ExpectExec(`UPDATE movies SET title = \$1, description = \$2, price = \$3, duration = \$4, status = \$5 WHERE id = \$6`).
		WithArgs(mockingMovie.Title, mockingMovie.Description, mockingMovie.Price, mockingMovie.Duration, mockingMovie.Status, mockingMovie.ID).
		WillReturnError(errors.New("Update Movie Failed"))

	_, err := suite.repo.Update(mockingMovie)
	assert.Error(suite.T(), err)
}

func (suite *MovieRepositoryTestSuite) TestDelete_Success() {
	id := mockingMovie.ID

	suite.mockSql.ExpectQuery(`DELETE FROM movies WHERE id = \$1 RETURNING id, title, description, price, duration, status`).WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "price", "duration", "status"}).AddRow(mockingMovie.ID, mockingMovie.Title, mockingMovie.Description, mockingMovie.Price, mockingMovie.Duration, mockingMovie.Status))

	result, err := suite.repo.Delete(id)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockingMovie, result)
}

func (suite *MovieRepositoryTestSuite) TestDelete_Failed() {
	id := "1"

	suite.mockSql.ExpectQuery(`DELETE FROM movies WHERE id = \$1 RETURNING id, title, description, price, duration, status`).WithArgs(id).
		WillReturnError(errors.New("Delete Movie Failed"))

	_, err := suite.repo.Delete(id)
	assert.Error(suite.T(), err)
}
