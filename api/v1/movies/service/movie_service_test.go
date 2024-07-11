package service

import (
	"bioskuy/api/v1/genre/dto"
	"bioskuy/api/v1/movies/entity"
	"bioskuy/api/v1/movies/mock/repomock"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MovieServiceTestSuite struct {
	suite.Suite
	repoMock *repomock.MockMovieRepository
	service  MovieService
}

var mockMovie = entity.Movie{
	ID:          "1",
	Title:       "Inception",
	Description: "A mind-bending thriller",
	Price:       10.0,
	Duration:    148,
	Status:      "available",
}

func (suite *MovieServiceTestSuite) SetupTest() {
	suite.repoMock = new(repomock.MockMovieRepository)
	suite.service = NewMovieService(suite.repoMock)
}

func TestMovieServiceTestSuite(t *testing.T) {
	suite.Run(t, new(MovieServiceTestSuite))
}

func (suite *MovieServiceTestSuite) TestGetAllMovies_Success() {
	page, size := 1, 10
	totalRows := 1

	suite.repoMock.On("GetAll", page, size).Return([]entity.Movie{mockMovie}, dto.Paging{
		Page:       page,
		Size:       size,
		TotalRows:  totalRows,
		TotalPages: 1,
	}, nil)

	movies, paging, err := suite.service.GetAllMovies(page, size)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), movies, 1)
	assert.Equal(suite.T(), mockMovie.Title, movies[0].Title)
	assert.Equal(suite.T(), dto.Paging{
		Page:       page,
		Size:       size,
		TotalRows:  totalRows,
		TotalPages: 1,
	}, paging)

	suite.repoMock.AssertExpectations(suite.T())
}

func (suite *MovieServiceTestSuite) TestCreateMovie_Success() {
	suite.repoMock.On("Create", mockMovie).Return(mockMovie, nil)

	result, err := suite.service.CreateMovie(mockMovie)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockMovie.Title, result.Title)

	suite.repoMock.AssertExpectations(suite.T())
}

func (suite *MovieServiceTestSuite) TestCreateMovie_Failed() {
	suite.repoMock.On("Create", mockMovie).Return(entity.Movie{}, errors.New("error"))

	_, err := suite.service.CreateMovie(mockMovie)
	assert.Error(suite.T(), err)

	suite.repoMock.AssertExpectations(suite.T())
}

func (suite *MovieServiceTestSuite) TestGetMovieByID_Success() {
	id := mockMovie.ID

	suite.repoMock.On("GetByID", id).Return(mockMovie, nil)

	result, err := suite.service.GetMovieByID(id)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockMovie.Title, result.Title)

	suite.repoMock.AssertExpectations(suite.T())
}

func (suite *MovieServiceTestSuite) TestGetMovieByID_Failed() {
	id := mockMovie.ID

	suite.repoMock.On("GetByID", id).Return(entity.Movie{}, errors.New("error"))

	_, err := suite.service.GetMovieByID(id)
	assert.Error(suite.T(), err)

	suite.repoMock.AssertExpectations(suite.T())
}

func (suite *MovieServiceTestSuite) TestUpdateMovie_Success() {
	suite.repoMock.On("Update", mockMovie).Return(mockMovie, nil)

	result, err := suite.service.UpdateMovie(mockMovie)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockMovie.Title, result.Title)

	suite.repoMock.AssertExpectations(suite.T())
}

func (suite *MovieServiceTestSuite) TestUpdateMovie_Failed() {
	suite.repoMock.On("Update", mockMovie).Return(entity.Movie{}, errors.New("error"))

	_, err := suite.service.UpdateMovie(mockMovie)
	assert.Error(suite.T(), err)

	suite.repoMock.AssertExpectations(suite.T())
}

func (suite *MovieServiceTestSuite) TestDeleteMovie_Success() {
	id := mockMovie.ID

	suite.repoMock.On("Delete", id).Return(mockMovie, nil)

	result, err := suite.service.DeleteMovie(id)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockMovie.Title, result.Title)

	suite.repoMock.AssertExpectations(suite.T())
}

func (suite *MovieServiceTestSuite) TestDeleteMovie_Failed() {
	id := mockMovie.ID

	suite.repoMock.On("Delete", id).Return(entity.Movie{}, errors.New("error"))

	_, err := suite.service.DeleteMovie(id)
	assert.Error(suite.T(), err)

	suite.repoMock.AssertExpectations(suite.T())
}
