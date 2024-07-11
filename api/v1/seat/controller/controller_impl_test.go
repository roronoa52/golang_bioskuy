package controller

import (
	"bioskuy/api/v1/seat/dto"
	"bioskuy/api/v1/seat/mock/servicemock"
	"bioskuy/exception"
	"bioskuy/web"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type SeatControllerTestSuite struct {
	suite.Suite
	controller  SeatController
	router      *gin.Engine
	mockService *servicemock.SeatService
}

func (suite *SeatControllerTestSuite) SetupTest() {
	suite.mockService = new(servicemock.SeatService)
	suite.controller = NewSeatController(suite.mockService)
}

func (suite *SeatControllerTestSuite) TestFindById_Success() {
	gin.SetMode(gin.TestMode)

	expectedSeat := dto.SeatResponse{ID: "1", Name: "Seat A", IsAvailable: true, StudioID: "1"}
	suite.mockService.On("FindByID", mock.Anything, "1", mock.Anything).Return(expectedSeat, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest("GET", "/seats/1", nil)
	c.Request = req
	c.Params = gin.Params{{Key: "seatId", Value: "1"}}

	suite.controller.FindById(c)

	var response web.FormatResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	// Extracting Data field from response
	responseData, err := json.Marshal(response.Data)
	assert.NoError(suite.T(), err)

	var result dto.SeatResponse
	err = json.Unmarshal(responseData, &result)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedSeat, result)

	suite.mockService.AssertExpectations(suite.T())
}

func (suite *SeatControllerTestSuite) TestFindById_NotFound() {
	gin.SetMode(gin.TestMode)

	suite.mockService.On("FindByID", mock.Anything, "1", mock.Anything).Return(dto.SeatResponse{}, exception.NotFoundError{Message: "Seat not found"})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest("GET", "/seats/1", nil)
	c.Request = req
	c.Params = gin.Params{{Key: "seatId", Value: "1"}}

	suite.controller.FindById(c)

	assert.Equal(suite.T(), http.StatusNotFound, w.Code)

	suite.mockService.AssertExpectations(suite.T())
}

func (suite *SeatControllerTestSuite) TestFindAll_Success() {
	gin.SetMode(gin.TestMode)

	expectedSeats := []dto.SeatResponse{
		{ID: "1", Name: "Seat A", IsAvailable: true, StudioID: "1"},
		{ID: "2", Name: "Seat B", IsAvailable: false, StudioID: "1"},
	}
	suite.mockService.On("FindAll", mock.Anything, "1", mock.Anything).Return(expectedSeats, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest("GET", "/studios/1/seats", nil)
	c.Request = req
	c.Params = gin.Params{{Key: "studioId", Value: "1"}}

	suite.controller.FindAll(c)

	var response web.FormatResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	// Extracting Data field from response
	responseData, err := json.Marshal(response.Data)
	assert.NoError(suite.T(), err)

	var result []dto.SeatResponse
	err = json.Unmarshal(responseData, &result)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedSeats, result)

	suite.mockService.AssertExpectations(suite.T())
}

func (suite *SeatControllerTestSuite) TestFindAll_InternalServerError() {
	gin.SetMode(gin.TestMode)

	suite.mockService.On("FindAll", mock.Anything, "1", mock.Anything).Return(nil, exception.InternalServerError{Message: "Database error"})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest("GET", "/studios/1/seats", nil)
	c.Request = req
	c.Params = gin.Params{{Key: "studioId", Value: "1"}}

	suite.controller.FindAll(c)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)

	suite.mockService.AssertExpectations(suite.T())
}

func TestSeatControllerTestSuite(t *testing.T) {
	suite.Run(t, new(SeatControllerTestSuite))
}
