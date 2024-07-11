package controller

import (
	"bioskuy/api/v1/showtime/dto"
	"bioskuy/exception"
	"bioskuy/helper"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockShowtimeService struct {
	mock.Mock
}

func (m *MockShowtimeService) Create(ctx context.Context, request dto.ShowtimeRequest, c *gin.Context) (dto.CreateShowtimesResponseDTO, error) {
	args := m.Called(ctx, request, c)
	return args.Get(0).(dto.CreateShowtimesResponseDTO), args.Error(1)
}

func (m *MockShowtimeService) FindByID(ctx context.Context, id string, c *gin.Context) (dto.ShowtimesResponse, error) {
	args := m.Called(ctx, id, c)
	return args.Get(0).(dto.ShowtimesResponse), args.Error(1)
}

func (m *MockShowtimeService) FindAll(ctx context.Context, c *gin.Context) ([]dto.ShowtimesResponse, error) {
	args := m.Called(ctx, c)
	return args.Get(0).([]dto.ShowtimesResponse), args.Error(1)
}

func (m *MockShowtimeService) Delete(ctx context.Context, id string, c *gin.Context) error {
	args := m.Called(ctx, id, c)
	return args.Error(0)
}

type ShowtimeControllerTestSuite struct {
	suite.Suite
	mockService *MockShowtimeService
	controller  ShowtimeController
}

func (suite *ShowtimeControllerTestSuite) SetupTest() {
	suite.mockService = new(MockShowtimeService)
	suite.controller = NewMovieController(suite.mockService)
}

func TestShowtimeControllerTestSuite(t *testing.T) {
	suite.Run(t, new(ShowtimeControllerTestSuite))
}

func (suite *ShowtimeControllerTestSuite) TestCreate_Success() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/showtimes", suite.controller.Create)

	showtimeRequest := dto.ShowtimeRequest{
		MovieID:   "1",
		StudioID:  "1",
		ShowStart: "2024-07-09T10:00:00Z",
	}

	showtimeResponse := dto.CreateShowtimesResponseDTO{
		MovieID:   "1",
		StudioID:  "1",
		ShowStart: helper.StringToDate("2024-07-09T10:00:00Z", nil),
		ShowEnd:   helper.StringToDate("2024-07-09T12:00:00Z", nil),
	}

	suite.mockService.On("Create", mock.Anything, showtimeRequest, mock.Anything).Return(showtimeResponse, nil).Once()

	body, _ := json.Marshal(showtimeRequest)
	req, _ := http.NewRequest(http.MethodPost, "/showtimes", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *ShowtimeControllerTestSuite) TestFindById_Success() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/showtimes/:showtimeId", suite.controller.FindById)

	showtimeResponse := dto.ShowtimesResponse{
		ID:               "1",
		StudioID:         "1",
		MovieID:          "1",
		StudioName:       "Studio 1",
		MovieTitle:       "Movie 1",
		MovieDescription: "Description 1",
		MoviePrice:       100,
		MovieDuration:    2,
		MovieStatus:      "Available",
	}

	suite.mockService.On("FindByID", mock.Anything, "1", mock.Anything).Return(showtimeResponse, nil).Once()

	req, _ := http.NewRequest(http.MethodGet, "/showtimes/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *ShowtimeControllerTestSuite) TestFindAll_Success() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/showtimes", suite.controller.FindAll)

	showtimeResponses := []dto.ShowtimesResponse{
		{
			ID:               "1",
			StudioID:         "1",
			MovieID:          "1",
			StudioName:       "Studio 1",
			MovieTitle:       "Movie 1",
			MovieDescription: "Description 1",
			MoviePrice:       100,
			MovieDuration:    2,
			MovieStatus:      "Available",
		},
	}

	suite.mockService.On("FindAll", mock.Anything, mock.Anything).Return(showtimeResponses, nil).Once()

	req, _ := http.NewRequest(http.MethodGet, "/showtimes", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *ShowtimeControllerTestSuite) TestDelete_Success() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.DELETE("/showtimes/:showtimeId", suite.controller.Delete)

	suite.mockService.On("Delete", mock.Anything, "1", mock.Anything).Return(nil).Once()

	req, _ := http.NewRequest(http.MethodDelete, "/showtimes/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *ShowtimeControllerTestSuite) TestCreate_ValidationError() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/showtimes", suite.controller.Create)

	showtimeRequest := dto.ShowtimeRequest{}

	suite.mockService.On("Create", mock.Anything, showtimeRequest, mock.Anything).Return(dto.CreateShowtimesResponseDTO{}, exception.ValidationError{Message: "Validation Error"}).Once()

	body, _ := json.Marshal(showtimeRequest)
	req, _ := http.NewRequest(http.MethodPost, "/showtimes", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusForbidden, w.Code)
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *ShowtimeControllerTestSuite) TestCreate_InternalServerError() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/showtimes", suite.controller.Create)

	showtimeRequest := dto.ShowtimeRequest{
		MovieID:   "1",
		StudioID:  "1",
		ShowStart: "2024-07-09T10:00:00Z",
	}

	suite.mockService.On("Create", mock.Anything, showtimeRequest, mock.Anything).Return(dto.CreateShowtimesResponseDTO{}, exception.InternalServerError{Message: "Internal Server Error"}).Once()

	body, _ := json.Marshal(showtimeRequest)
	req, _ := http.NewRequest(http.MethodPost, "/showtimes", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *ShowtimeControllerTestSuite) TestFindById_NotFoundError() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/showtimes/:showtimeId", suite.controller.FindById)

	suite.mockService.On("FindByID", mock.Anything, "1", mock.Anything).Return(dto.ShowtimesResponse{}, exception.NotFoundError{Message: "Showtime Not Found"}).Once()

	req, _ := http.NewRequest(http.MethodGet, "/showtimes/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *ShowtimeControllerTestSuite) TestFindById_InternalServerError() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/showtimes/:showtimeId", suite.controller.FindById)

	suite.mockService.On("FindByID", mock.Anything, "1", mock.Anything).Return(dto.ShowtimesResponse{}, exception.InternalServerError{Message: "Internal Server Error"}).Once()

	req, _ := http.NewRequest(http.MethodGet, "/showtimes/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *ShowtimeControllerTestSuite) TestFindAll_InternalServerError() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/showtimes", suite.controller.FindAll)

	suite.mockService.On("FindAll", mock.Anything, mock.Anything).Return(nil, exception.InternalServerError{Message: "Internal Server Error"}).Once()

	req, _ := http.NewRequest(http.MethodGet, "/showtimes", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *ShowtimeControllerTestSuite) TestDelete_InternalServerError() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.DELETE("/showtimes/:showtimeId", suite.controller.Delete)

	suite.mockService.On("Delete", mock.Anything, "1", mock.Anything).Return(exception.InternalServerError{Message: "Internal Server Error"}).Once()

	req, _ := http.NewRequest(http.MethodDelete, "/showtimes/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *ShowtimeControllerTestSuite) TestDelete_NotFoundError() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.DELETE("/showtimes/:showtimeId", suite.controller.Delete)

	suite.mockService.On("Delete", mock.Anything, "1", mock.Anything).Return(exception.NotFoundError{Message: "Showtime Not Found"}).Once()

	req, _ := http.NewRequest(http.MethodDelete, "/showtimes/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
	suite.mockService.AssertExpectations(suite.T())
}
