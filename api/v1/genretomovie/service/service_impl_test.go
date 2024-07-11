package service

import (
	"bioskuy/api/v1/genretomovie/dto"
	"bioskuy/api/v1/genretomovie/entity"
	"bioskuy/api/v1/genretomovie/mock/repomock"
	"bioskuy/exception"
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type GenreToMovieServiceTestSuite struct {
	suite.Suite
	mockDb     *sql.DB
	mockSql    sqlmock.Sqlmock
	mockRepo   *repomock.MockGenreToMovieRepository
	validate   *validator.Validate
	service    GenreToMovieService
	ctx        context.Context
	ginContext *gin.Context
}

func (suite *GenreToMovieServiceTestSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	assert.NoError(suite.T(), err)
	suite.mockDb = db
	suite.mockSql = mock
	suite.mockRepo = &repomock.MockGenreToMovieRepository{}
	suite.validate = validator.New()
	suite.service = NewGenreToMovieService(suite.mockRepo, suite.validate, suite.mockDb)
	suite.ctx = context.Background()
	suite.ginContext = &gin.Context{}
}

func TestGenreToMovieServiceTestSuite(t *testing.T) {
	suite.Run(t, new(GenreToMovieServiceTestSuite))
}

// Create
func (suite *GenreToMovieServiceTestSuite) TestCreate_Success() {
	request := dto.CreateGenreToMovieRequest{
		GenreID: "genre-id",
		MovieID: "movie-id",
	}
	expectedResult := entity.GenreToMovie{
		ID:      "new-id",
		GenreID: request.GenreID,
		MovieID: request.MovieID,
	}

	suite.mockSql.ExpectBegin()
	suite.mockRepo.On("Save", suite.ctx, mock.Anything, mock.AnythingOfType("entity.GenreToMovie"), suite.ginContext).Return(expectedResult, nil)
	suite.mockSql.ExpectCommit()

	response, err := suite.service.Create(suite.ctx, request, suite.ginContext)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedResult.ID, response.ID)
	assert.Equal(suite.T(), request.GenreID, response.GenreID)
	assert.Equal(suite.T(), request.MovieID, response.MovieID)
}

func (suite *GenreToMovieServiceTestSuite) TestCreate_ValidationError() {
	request := dto.CreateGenreToMovieRequest{
		GenreID: "",
		MovieID: "movie-id",
	}

	response, err := suite.service.Create(suite.ctx, request, suite.ginContext)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "", response.ID)
	assert.IsType(suite.T(), exception.ValidationError{}, err)
}

func (suite *GenreToMovieServiceTestSuite) TestCreate_SaveError() {
	request := dto.CreateGenreToMovieRequest{
		GenreID: "genre-id",
		MovieID: "movie-id",
	}
	saveError := errors.New("Save Error")

	suite.mockSql.ExpectBegin()
	suite.mockRepo.On("Save", suite.ctx, mock.Anything, mock.AnythingOfType("entity.GenreToMovie"), suite.ginContext).Return(entity.GenreToMovie{}, saveError)
	suite.mockSql.ExpectRollback()

	response, err := suite.service.Create(suite.ctx, request, suite.ginContext)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "", response.ID)
	assert.Equal(suite.T(), saveError, err)
}

// FindByID
func (suite *GenreToMovieServiceTestSuite) TestFindByID_Success() {
	expectedResult := entity.GenreToMovie{
		ID:      "some-id",
		GenreID: "genre-id",
		MovieID: "movie-id",
	}

	suite.mockSql.ExpectBegin()
	suite.mockRepo.On("FindByID", suite.ctx, mock.Anything, "some-id", suite.ginContext).Return(expectedResult, nil)
	suite.mockSql.ExpectCommit()

	response, err := suite.service.FindByID(suite.ctx, "some-id", suite.ginContext)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedResult.ID, response.ID)
	assert.Equal(suite.T(), expectedResult.GenreID, response.GenreID)
	assert.Equal(suite.T(), expectedResult.MovieID, response.MovieID)
}

func (suite *GenreToMovieServiceTestSuite) TestFindByID_NotFoundError() {
	notFoundError := errors.New("Not Found Error")

	suite.mockSql.ExpectBegin()
	suite.mockRepo.On("FindByID", suite.ctx, mock.Anything, "some-id", suite.ginContext).Return(entity.GenreToMovie{}, notFoundError)
	suite.mockSql.ExpectRollback()

	response, err := suite.service.FindByID(suite.ctx, "some-id", suite.ginContext)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "", response.ID)
	assert.Equal(suite.T(), notFoundError, err)
}

// FindAll
func (suite *GenreToMovieServiceTestSuite) TestFindAll_Success() {
	expectedResult := []entity.GenreToMovie{
		{ID: "id1", GenreID: "genre1", MovieID: "movie1"},
		{ID: "id2", GenreID: "genre2", MovieID: "movie2"},
	}

	suite.mockSql.ExpectBegin()
	suite.mockRepo.On("FindAll", suite.ctx, mock.Anything, suite.ginContext).Return(expectedResult, nil)
	suite.mockSql.ExpectCommit()

	response, err := suite.service.FindAll(suite.ctx, suite.ginContext)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), response, 2)
	assert.Equal(suite.T(), expectedResult[0].ID, response[0].ID)
	assert.Equal(suite.T(), expectedResult[0].GenreID, response[0].GenreID)
	assert.Equal(suite.T(), expectedResult[0].MovieID, response[0].MovieID)
}

func (suite *GenreToMovieServiceTestSuite) TestFindAll_Error() {
	findAllError := errors.New("Find All Error")

	suite.mockSql.ExpectBegin()
	suite.mockRepo.On("FindAll", suite.ctx, mock.Anything, suite.ginContext).Return([]entity.GenreToMovie{}, findAllError)
	suite.mockSql.ExpectRollback()

	response, err := suite.service.FindAll(suite.ctx, suite.ginContext)

	assert.Error(suite.T(), err)
	assert.Len(suite.T(), response, 0)
	assert.Equal(suite.T(), findAllError, err)
}

// Delete
func (suite *GenreToMovieServiceTestSuite) TestDelete_Success() {
	expectedResult := entity.GenreToMovie{
		ID:      "some-id",
		GenreID: "genre-id",
		MovieID: "movie-id",
	}

	suite.mockSql.ExpectBegin()
	suite.mockRepo.On("FindByID", suite.ctx, mock.Anything, "some-id", suite.ginContext).Return(expectedResult, nil)
	suite.mockRepo.On("Delete", suite.ctx, mock.Anything, "some-id", suite.ginContext).Return(nil)
	suite.mockSql.ExpectCommit()

	err := suite.service.Delete(suite.ctx, "some-id", suite.ginContext)

	assert.NoError(suite.T(), err)
}

func (suite *GenreToMovieServiceTestSuite) TestDelete_NotFoundError() {
	notFoundError := errors.New("Not Found Error")

	suite.mockSql.ExpectBegin()
	suite.mockRepo.On("FindByID", suite.ctx, mock.Anything, "some-id", suite.ginContext).Return(entity.GenreToMovie{}, notFoundError)
	suite.mockSql.ExpectRollback()

	err := suite.service.Delete(suite.ctx, "some-id", suite.ginContext)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), notFoundError, err)
}

func (suite *GenreToMovieServiceTestSuite) TestDelete_DeleteError() {
	expectedResult := entity.GenreToMovie{
		ID:      "some-id",
		GenreID: "genre-id",
		MovieID: "movie-id",
	}
	deleteError := errors.New("Delete Error")

	suite.mockSql.ExpectBegin()
	suite.mockRepo.On("FindByID", suite.ctx, mock.Anything, "some-id", suite.ginContext).Return(expectedResult, nil)
	suite.mockRepo.On("Delete", suite.ctx, mock.Anything, "some-id", suite.ginContext).Return(deleteError)
	suite.mockSql.ExpectRollback()

	err := suite.service.Delete(suite.ctx, "some-id", suite.ginContext)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), deleteError, err)
}
