package service

import (
	"bioskuy/api/v1/genre/dto"
	"bioskuy/api/v1/genre/entity"
	"bioskuy/api/v1/genre/mock/repomock"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type GenreServiceTestSuite struct {
	suite.Suite
	mockRepo *repomock.MockGenreRepository
	service  GenreService
}

var mockingGenre = entity.Genre{
	ID:   uuid.New(),
	Name: "Action",
}

func (suite *GenreServiceTestSuite) SetupTest() {
	suite.mockRepo = new(repomock.MockGenreRepository)
	suite.service = NewGenreService(suite.mockRepo)
}

func TestGenreServiceTestSuite(t *testing.T) {
	suite.Run(t, new(GenreServiceTestSuite))
}

func (suite *GenreServiceTestSuite) TestCreateGenre_Success() {
	suite.mockRepo.On("Create", mockingGenre).Return(mockingGenre, nil)

	createdGenre, err := suite.service.CreateGenre(mockingGenre)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockingGenre, createdGenre)
}

func (suite *GenreServiceTestSuite) TestCreateGenre_Error() {
	suite.mockRepo.On("Create", mockingGenre).Return(entity.Genre{}, errors.New("error creating genre"))

	_, err := suite.service.CreateGenre(mockingGenre)
	assert.Error(suite.T(), err)
}

func (suite *GenreServiceTestSuite) TestGetGenreByID_Success() {
	suite.mockRepo.On("GetByID", mockingGenre.ID).Return(mockingGenre, nil)

	foundGenre, err := suite.service.GetGenreByID(mockingGenre.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockingGenre, foundGenre)
}

func (suite *GenreServiceTestSuite) TestGetGenreByID_Error() {
	suite.mockRepo.On("GetByID", mockingGenre.ID).Return(entity.Genre{}, errors.New("genre not found"))

	_, err := suite.service.GetGenreByID(mockingGenre.ID)
	assert.Error(suite.T(), err)
}

func (suite *GenreServiceTestSuite) TestGetAll_Success() {
	genres := []entity.Genre{mockingGenre}
	paging := dto.Paging{Page: 1, Size: 10, TotalRows: 1, TotalPages: 1}

	suite.mockRepo.On("GetAll", 1, 10).Return(genres, paging, nil)

	foundGenres, foundPaging, err := suite.service.GetAll(1, 10)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), genres, foundGenres)
	assert.Equal(suite.T(), paging, foundPaging)
}

func (suite *GenreServiceTestSuite) TestGetAll_Error() {
	suite.mockRepo.On("GetAll", 1, 10).Return([]entity.Genre{}, dto.Paging{}, errors.New("error fetching genres"))

	_, _, err := suite.service.GetAll(1, 10)
	assert.Error(suite.T(), err)
}

func (suite *GenreServiceTestSuite) TestUpdateGenre_Success() {
	suite.mockRepo.On("Update", mockingGenre).Return(mockingGenre, nil)

	updatedGenre, err := suite.service.UpdateGenre(mockingGenre)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockingGenre, updatedGenre)
}

func (suite *GenreServiceTestSuite) TestUpdateGenre_Error() {
	suite.mockRepo.On("Update", mockingGenre).Return(entity.Genre{}, errors.New("error updating genre"))

	_, err := suite.service.UpdateGenre(mockingGenre)
	assert.Error(suite.T(), err)
}

func (suite *GenreServiceTestSuite) TestDeleteGenre_Success() {
	suite.mockRepo.On("Delete", mockingGenre.ID).Return(mockingGenre, nil)

	deletedGenre, err := suite.service.DeleteGenre(mockingGenre.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockingGenre, deletedGenre)
}

func (suite *GenreServiceTestSuite) TestDeleteGenre_Error() {
	suite.mockRepo.On("Delete", mockingGenre.ID).Return(entity.Genre{}, errors.New("error deleting genre"))

	_, err := suite.service.DeleteGenre(mockingGenre.ID)
	assert.Error(suite.T(), err)
}
