package repository

import (
	"bioskuy/api/v1/studio/entity"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"golang.org/x/net/context"
)

type StudioRepositoryTestSuite struct {
	suite.Suite
	mockDb  *sql.DB
	mockSql sqlmock.Sqlmock
	repo    StudioRepository
	ctx     context.Context
	ginCtx  *gin.Context
}

var mockingStudio = entity.Studio{
	ID:       "1231cmf1m",
	Name:     "Studio 1",
	Capacity: 100,
}

func (suite *StudioRepositoryTestSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	assert.NoError(suite.T(), err)
	suite.mockDb = db
	suite.mockSql = mock
	suite.repo = NewStudioRepository()
	suite.ctx = context.Background()
	suite.ginCtx = &gin.Context{}
}

func TestStudioRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(StudioRepositoryTestSuite))
}

func (suite *StudioRepositoryTestSuite) TestSave_Success() {
	suite.mockSql.ExpectBegin()
	suite.mockSql.ExpectQuery(`INSERT INTO studios \(name, capacity\) VALUES \(\$1, \$2\) RETURNING id`).
		WithArgs(mockingStudio.Name, mockingStudio.Capacity).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(mockingStudio.ID))
	suite.mockSql.ExpectCommit()

	tx, err := suite.mockDb.Begin()
	assert.NoError(suite.T(), err)

	result, err := suite.repo.Save(suite.ctx, tx, mockingStudio, suite.ginCtx)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockingStudio.Name, result.Name)
	assert.Equal(suite.T(), mockingStudio.Capacity, result.Capacity)
	assert.NoError(suite.T(), tx.Commit())
}

func (suite *StudioRepositoryTestSuite) TestSave_Failed() {
	suite.mockSql.ExpectBegin()
	suite.mockSql.ExpectQuery(`INSERT INTO studios \(name, capacity\) VALUES \(\$1, \$2\) RETURNING id`).
		WithArgs(mockingStudio.Name, mockingStudio.Capacity).
		WillReturnError(errors.New("Insert Studio Failed"))
	suite.mockSql.ExpectRollback()

	tx, err := suite.mockDb.Begin()
	assert.NoError(suite.T(), err)

	_, err = suite.repo.Save(suite.ctx, tx, mockingStudio, suite.ginCtx)
	assert.Error(suite.T(), err)
	assert.NoError(suite.T(), tx.Rollback())
}

func (suite *StudioRepositoryTestSuite) TestFindByID_Success() {
	suite.mockSql.ExpectBegin()
	suite.mockSql.ExpectQuery(`SELECT id, name, capacity FROM studios WHERE id = \$1`).
		WithArgs(mockingStudio.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "capacity"}).AddRow(mockingStudio.ID, mockingStudio.Name, mockingStudio.Capacity))
	suite.mockSql.ExpectCommit()

	tx, err := suite.mockDb.Begin()
	assert.NoError(suite.T(), err)

	result, err := suite.repo.FindByID(suite.ctx, tx, mockingStudio.ID, suite.ginCtx)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockingStudio.ID, result.ID)
	assert.Equal(suite.T(), mockingStudio.Name, result.Name)
	assert.Equal(suite.T(), mockingStudio.Capacity, result.Capacity)
	assert.NoError(suite.T(), tx.Commit())
}

func (suite *StudioRepositoryTestSuite) TestFindByID_NotFound() {
	suite.mockSql.ExpectBegin()
	suite.mockSql.ExpectQuery(`SELECT id, name, capacity FROM studios WHERE id = \$1`).
		WithArgs(mockingStudio.ID).
		WillReturnError(sql.ErrNoRows)
	suite.mockSql.ExpectRollback()

	tx, err := suite.mockDb.Begin()
	assert.NoError(suite.T(), err)

	_, err = suite.repo.FindByID(suite.ctx, tx, mockingStudio.ID, suite.ginCtx)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "studio not found", err.Error())
	assert.NoError(suite.T(), tx.Rollback())
}

func (suite *StudioRepositoryTestSuite) TestFindAll_Success() {
	suite.mockSql.ExpectBegin()
	suite.mockSql.ExpectQuery(`SELECT id, name, capacity FROM studios`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "capacity"}).
			AddRow(mockingStudio.ID, mockingStudio.Name, mockingStudio.Capacity).
			AddRow(uuid.New(), "Studio 2", 200))
	suite.mockSql.ExpectCommit()

	tx, err := suite.mockDb.Begin()
	assert.NoError(suite.T(), err)

	result, err := suite.repo.FindAll(suite.ctx, tx, suite.ginCtx)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 2)
	assert.NoError(suite.T(), tx.Commit())
}

func (suite *StudioRepositoryTestSuite) TestUpdate_Success() {
	suite.mockSql.ExpectBegin()
	suite.mockSql.ExpectExec(`UPDATE studios SET name = \$1, capacity = \$2 WHERE id = \$3`).
		WithArgs(mockingStudio.Name, mockingStudio.Capacity, mockingStudio.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mockSql.ExpectCommit()

	tx, err := suite.mockDb.Begin()
	assert.NoError(suite.T(), err)

	result, err := suite.repo.Update(suite.ctx, tx, mockingStudio, suite.ginCtx)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockingStudio, result)
	assert.NoError(suite.T(), tx.Commit())
}

func (suite *StudioRepositoryTestSuite) TestUpdate_Failed() {
	suite.mockSql.ExpectBegin()
	suite.mockSql.ExpectExec(`UPDATE studios SET name = \$1, capacity = \$2 WHERE id = \$3`).
		WithArgs(mockingStudio.Name, mockingStudio.Capacity, mockingStudio.ID).
		WillReturnError(errors.New("Update Studio Failed"))
	suite.mockSql.ExpectRollback()

	tx, err := suite.mockDb.Begin()
	assert.NoError(suite.T(), err)

	_, err = suite.repo.Update(suite.ctx, tx, mockingStudio, suite.ginCtx)
	assert.Error(suite.T(), err)
	assert.NoError(suite.T(), tx.Rollback())
}

func (suite *StudioRepositoryTestSuite) TestDelete_Success() {
	suite.mockSql.ExpectBegin()
	suite.mockSql.ExpectExec(`DELETE FROM studios WHERE id = \$1`).
		WithArgs(mockingStudio.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mockSql.ExpectCommit()

	tx, err := suite.mockDb.Begin()
	assert.NoError(suite.T(), err)

	err = suite.repo.Delete(suite.ctx, tx, mockingStudio.ID, suite.ginCtx)
	assert.NoError(suite.T(), err)
	assert.NoError(suite.T(), tx.Commit())
}

func (suite *StudioRepositoryTestSuite) TestDelete_Failed() {
	suite.mockSql.ExpectBegin()
	suite.mockSql.ExpectExec(`DELETE FROM studios WHERE id = \$1`).
		WithArgs(mockingStudio.ID).
		WillReturnError(errors.New("Delete Studio Failed"))
	suite.mockSql.ExpectRollback()

	tx, err := suite.mockDb.Begin()
	assert.NoError(suite.T(), err)

	err = suite.repo.Delete(suite.ctx, tx, mockingStudio.ID, suite.ginCtx)
	assert.Error(suite.T(), err)
	assert.NoError(suite.T(), tx.Rollback())
}
