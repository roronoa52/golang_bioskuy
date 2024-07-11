package controller

import (
	"bioskuy/api/v1/genretomovie/dto"
	"bioskuy/api/v1/genretomovie/mock/servicemock"
	"bioskuy/exception"
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// GenreControllerTestSuite struct
type GenreControllerTestSuite struct {
	suite.Suite
	mockService *servicemock.MockGenreToMovieService
	router      *gin.Engine
}

func (suite *GenreControllerTestSuite) SetupTest() {
	suite.mockService = new(servicemock.MockGenreToMovieService)
	controller := NewGenreToMovieController(suite.mockService)
	suite.router = gin.Default()
	suite.router.POST("/genretomovie", controller.Create)
	suite.router.GET("/genretomovie/:genretomovieId", controller.FindById)
	suite.router.GET("/genretomovie", controller.FindAll)
	suite.router.DELETE("/genretomovie/:genretomovieId", controller.Delete)
}

// TestCreate_Success
func (suite *GenreControllerTestSuite) TestCreate_Success() {
	mockResponse := dto.GenreToMovieResponse{ID: "1", GenreID: "1", MovieID: "1"}
	suite.mockService.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(mockResponse, nil)

	reqBody := `{"id":"1", "genre_id":"1", "movie_id":"1"}`
	req, _ := http.NewRequest(http.MethodPost, "/genretomovie", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusCreated, resp.Code)
	suite.mockService.AssertExpectations(suite.T())
}

// TestCreate_BindError
func (suite *GenreControllerTestSuite) TestCreate_BindError() {
	suite.mockService.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(dto.GenreToMovieResponse{}, exception.ForbiddenError{Message: "bind error"})

	reqBody := `{"id":1, "genre_id":1, "movie_id":1}`
	req, _ := http.NewRequest(http.MethodPost, "/genretomovie", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusBadRequest, resp.Code)
	suite.mockService.AssertExpectations(suite.T())
}

// TestFindById_Success
func (suite *GenreControllerTestSuite) TestFindById_Success() {
	mockResponse := dto.GenreToMovieResponse{ID: "1", GenreID: "1", MovieID: "1"}
	suite.mockService.On("FindByID", mock.Anything, "1", mock.Anything).Return(mockResponse, nil)

	req, _ := http.NewRequest(http.MethodGet, "/genretomovie/1", nil)
	resp := httptest.NewRecorder()

	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusOK, resp.Code)
	suite.mockService.AssertExpectations(suite.T())
}

// TestFindById_NotFound
func (suite *GenreControllerTestSuite) TestFindById_NotFound() {
	suite.mockService.On("FindByID", mock.Anything, "2", mock.Anything).Return(dto.GenreToMovieResponse{}, exception.NotFoundError{Message: "not found"})

	req, _ := http.NewRequest(http.MethodGet, "/genretomovie/2", nil)
	resp := httptest.NewRecorder()

	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusNotFound, resp.Code)
	suite.mockService.AssertExpectations(suite.T())
}

// TestFindAll_Success
func (suite *GenreControllerTestSuite) TestFindAll_Success() {
	mockResponse := []dto.GenreToMovieResponse{{ID: "1", GenreID: "1", MovieID: "1"}}
	suite.mockService.On("FindAll", mock.Anything, mock.Anything).Return(mockResponse, nil)

	req, _ := http.NewRequest(http.MethodGet, "/genretomovie", nil)
	resp := httptest.NewRecorder()

	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusOK, resp.Code)
	suite.mockService.AssertExpectations(suite.T())
}

// TestFindAll_InternalServerError
func (suite *GenreControllerTestSuite) TestFindAll_InternalServerError() {
	suite.mockService.On("FindAll", mock.Anything, mock.Anything).Return(nil, errors.New("internal server error"))

	req, _ := http.NewRequest(http.MethodGet, "/genretomovie", nil)
	resp := httptest.NewRecorder()

	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, resp.Code)
	suite.mockService.AssertExpectations(suite.T())
}

// TestDelete_Success
func (suite *GenreControllerTestSuite) TestDelete_Success() {
	suite.mockService.On("Delete", mock.Anything, "1", mock.Anything).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/genretomovie/1", nil)
	resp := httptest.NewRecorder()

	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusOK, resp.Code)
	suite.mockService.AssertExpectations(suite.T())
}

// TestDelete_InternalServerError
func (suite *GenreControllerTestSuite) TestDelete_InternalServerError() {
	suite.mockService.On("Delete", mock.Anything, "2", mock.Anything).Return(errors.New("internal server error"))

	req, _ := http.NewRequest(http.MethodDelete, "/genretomovie/2", nil)
	resp := httptest.NewRecorder()

	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, resp.Code)
	suite.mockService.AssertExpectations(suite.T())
}

// Run all the tests
func TestGenreControllerTestSuite(t *testing.T) {
	suite.Run(t, new(GenreControllerTestSuite))
}
