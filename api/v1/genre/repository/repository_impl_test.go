package repository

import (
	"bioskuy/api/v1/genre/dto"
	"bioskuy/api/v1/genre/entity"
	"database/sql"
	"errors"
	"math"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type GenreRepositoryTestSuite struct {
	suite.Suite
	mockDb  *sql.DB
	mockSql sqlmock.Sqlmock
	repo    GenreRepository
}

var mockingGenre = entity.Genre{
	ID:   uuid.New(),
	Name: "Action",
}

func (suite *GenreRepositoryTestSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	assert.NoError(suite.T(), err)
	suite.mockDb = db
	suite.mockSql = mock
	suite.repo = NewGenreRepository(suite.mockDb)
}

func TestGenreRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(GenreRepositoryTestSuite))
}

func (suite *GenreRepositoryTestSuite) TestGetAll_Success() {
	page, size := 1, 10
	totalRows := 15

	suite.mockSql.ExpectQuery(`SELECT id, name FROM genres LIMIT \$1 OFFSET \$2`).WithArgs(size, (page-1)*size).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(mockingGenre.ID, mockingGenre.Name).
			AddRow(uuid.New(), "Drama"))

	suite.mockSql.ExpectQuery(`SELECT COUNT\(\*\) FROM genres`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(totalRows))

	genres, paging, err := suite.repo.GetAll(page, size)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), genres, 2)
	assert.Equal(suite.T(), dto.Paging{
		Page:       page,
		Size:       size,
		TotalRows:  totalRows,
		TotalPages: int(math.Ceil(float64(totalRows) / float64(size))),
	}, paging)
}

func (suite *GenreRepositoryTestSuite) TestGetAll_ErrorOnQueryGenres() {
	page, size := 1, 10

	suite.mockSql.ExpectQuery(`SELECT id, name FROM genres LIMIT \$1 OFFSET \$2`).WithArgs(size, (page-1)*size).
		WillReturnError(errors.New("Query Error"))

	_, _, err := suite.repo.GetAll(page, size)
	assert.Error(suite.T(), err)
}

func (suite *GenreRepositoryTestSuite) TestGetAll_ErrorOnCountGenres() {
	page, size := 1, 10

	suite.mockSql.ExpectQuery(`SELECT id, name FROM genres LIMIT \$1 OFFSET \$2`).WithArgs(size, (page-1)*size).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(mockingGenre.ID, mockingGenre.Name))

	suite.mockSql.ExpectQuery(`SELECT COUNT\(\*\) FROM genres`).WillReturnError(errors.New("Count Error"))

	_, _, err := suite.repo.GetAll(page, size)
	assert.Error(suite.T(), err)
}

func (suite *GenreRepositoryTestSuite) TestGetAll_ErrorOnScan() {
	page, size := 1, 10

	suite.mockSql.ExpectQuery(`SELECT id, name FROM genres LIMIT \$1 OFFSET \$2`).WithArgs(size, (page-1)*size).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(mockingGenre.ID, nil)) // nil akan menyebabkan scan error

	suite.mockSql.ExpectQuery(`SELECT COUNT\(\*\) FROM genres`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	_, _, err := suite.repo.GetAll(page, size)
	assert.Error(suite.T(), err)
}

func (suite *GenreRepositoryTestSuite) TestCreate_Success() {
	suite.mockSql.ExpectQuery(`INSERT INTO genres \(id, name\) VALUES \(\$1, \$2\) RETURNING id, name`).
		WithArgs(sqlmock.AnyArg(), mockingGenre.Name).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(mockingGenre.ID, mockingGenre.Name))

	result, err := suite.repo.Create(mockingGenre)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), mockingGenre.Name, result.Name)
}

func (suite *GenreRepositoryTestSuite) TestCreate_Failed() {
	suite.mockSql.ExpectBegin().
		WillReturnError(errors.New("Insert Genre Failed"))

	_, err := suite.repo.Create(mockingGenre)
	assert.Error(suite.T(), err)
}

func (suite *GenreRepositoryTestSuite) TestGetByID_ErrorOnQuery() {
	id := uuid.New()

	suite.mockSql.ExpectQuery(`SELECT id, name FROM genres WHERE id = \$1`).WithArgs(id).
		WillReturnError(errors.New("Query Error"))

	_, err := suite.repo.GetByID(id)
	assert.Error(suite.T(), err)
}

func (suite *GenreRepositoryTestSuite) TestGetByID_ErrorNoRows() {
	id := entity.Genre{
		ID: uuid.New(),
	}

	suite.mockSql.ExpectQuery(`SELECT id, name FROM genres WHERE id = \$1`).WithArgs(id).
		WillReturnError(sql.ErrNoRows)

	genre, err := suite.repo.GetByID(id.ID)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), entity.Genre{}, genre)
}

func (suite *GenreRepositoryTestSuite) TestGetByID_Success() {
	id := uuid.New()

	suite.mockSql.ExpectQuery(`SELECT id, name FROM genres WHERE id = \$1`).WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(mockingGenre.ID, mockingGenre.Name))

	genre, err := suite.repo.GetByID(id)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockingGenre, genre)
}

func (suite *GenreRepositoryTestSuite) TestUpdate_Success() {
	suite.mockSql.ExpectQuery(`UPDATE genres SET name = \$1 WHERE id = \$2 RETURNING id, name`).
		WithArgs(mockingGenre.Name, mockingGenre.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(mockingGenre.ID, mockingGenre.Name))

	result, err := suite.repo.Update(mockingGenre)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockingGenre, result)
}

func (suite *GenreRepositoryTestSuite) TestUpdate_Failed() {
	suite.mockSql.ExpectQuery(`UPDATE genres SET name = \$1 WHERE id = \$2 RETURNING id, name`).
		WithArgs(mockingGenre.Name, mockingGenre.ID).
		WillReturnError(errors.New("Update Genre Failed"))

	_, err := suite.repo.Update(mockingGenre)
	assert.Error(suite.T(), err)
}

func (suite *GenreRepositoryTestSuite) TestDelete_Success() {
	id := mockingGenre.ID

	suite.mockSql.ExpectQuery(`DELETE FROM genres WHERE id = \$1 RETURNING id, name`).WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(mockingGenre.ID, mockingGenre.Name))

	result, _ := suite.repo.Delete(id)
	assert.Equal(suite.T(), mockingGenre, result)
}

func (suite *GenreRepositoryTestSuite) TestDelete_Failed() {
	id := entity.Genre{
		ID: uuid.New(),
	}

	suite.mockSql.ExpectQuery(`DELETE FROM genres WHERE id = \$1 RETURNING id, name`).WithArgs(id).
		WillReturnError(sql.ErrConnDone)

	_, err := suite.repo.Delete(id.ID)
	assert.Error(suite.T(), err)
}
