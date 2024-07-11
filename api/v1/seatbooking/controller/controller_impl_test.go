package controller

import (
	"bioskuy/api/v1/seatbooking/dto"
	"bioskuy/api/v1/seatbooking/mock/servicemock"
	"bioskuy/exception"
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type SeatBookingControllerTestSuite struct {
	suite.Suite
	mockService *servicemock.MockSeatBookingService
	router      *gin.Engine
}

func (suite *SeatBookingControllerTestSuite) SetupTest() {
	suite.mockService = new(servicemock.MockSeatBookingService)
	controller := NewSeatbookingController(suite.mockService)
	suite.router = gin.Default()
	suite.router.POST("/seatbooking", controller.Create)
	suite.router.GET("/seatbooking/:seatbookingId", controller.FindById)
	suite.router.GET("/seatbooking", controller.FindAll)
	suite.router.DELETE("/seatbooking/:seatbookingId", controller.Delete)
}

func (suite *SeatBookingControllerTestSuite) TestCreate_Success() {
	mockResponse := dto.CreateSeatBookingResponse{ID: "1"}
	suite.mockService.On("Create", mock.Anything, mock.Anything, "test_user", mock.Anything).Return(mockResponse, nil)

	reqBody := `{"show_id":"1", "seat_id":"A1", "customer":"John Doe"}`
	req, _ := http.NewRequest(http.MethodPost, "/seatbooking", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), "user_id", "test_user"))
	resp := httptest.NewRecorder()

	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusCreated, resp.Code)
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *SeatBookingControllerTestSuite) TestCreate_BindError() {
	reqBody := `{"show_id":1, "seat_id":"A1", "customer":"John Doe"}`
	req, _ := http.NewRequest(http.MethodPost, "/seatbooking", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	suite.router.ServeHTTP(resp, req)

	suite.mockService.AssertExpectations(suite.T())
}

func (suite *SeatBookingControllerTestSuite) TestFindById_Success() {
	mockResponse := dto.SeatBookingResponse{ID: "1", ShowtimeID: "1", SeatID: "A1", UserID: "John Doe"}
	suite.mockService.On("FindByID", mock.Anything, "1", mock.Anything).Return(mockResponse, nil)

	req, _ := http.NewRequest(http.MethodGet, "/seatbooking/1", nil)
	resp := httptest.NewRecorder()

	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusOK, resp.Code)
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *SeatBookingControllerTestSuite) TestFindById_NotFound() {
	suite.mockService.On("FindByID", mock.Anything, "2", mock.Anything).Return(dto.SeatBookingResponse{}, exception.NotFoundError{Message: "not found"})

	req, _ := http.NewRequest(http.MethodGet, "/seatbooking/2", nil)
	resp := httptest.NewRecorder()

	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusNotFound, resp.Code)
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *SeatBookingControllerTestSuite) TestFindAll_Success() {
	mockResponse := []dto.SeatBookingResponse{{ID: "1", ShowtimeID: "1", SeatID: "A1", UserID: "John Doe"}}
	suite.mockService.On("FindAll", mock.Anything, mock.Anything).Return(mockResponse, nil)

	req, _ := http.NewRequest(http.MethodGet, "/seatbooking", nil)
	resp := httptest.NewRecorder()

	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusOK, resp.Code)
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *SeatBookingControllerTestSuite) TestFindAll_InternalServerError() {
	suite.mockService.On("FindAll", mock.Anything, mock.Anything).Return(nil, errors.New("internal server error"))

	req, _ := http.NewRequest(http.MethodGet, "/seatbooking", nil)
	resp := httptest.NewRecorder()

	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, resp.Code)
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *SeatBookingControllerTestSuite) TestDelete_Success() {
	suite.mockService.On("Delete", mock.Anything, "1", mock.Anything).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/seatbooking/1", nil)
	resp := httptest.NewRecorder()

	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusOK, resp.Code)
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *SeatBookingControllerTestSuite) TestDelete_InternalServerError() {
	suite.mockService.On("Delete", mock.Anything, "2", mock.Anything).Return(errors.New("internal server error"))

	req, _ := http.NewRequest(http.MethodDelete, "/seatbooking/2", nil)
	resp := httptest.NewRecorder()

	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, resp.Code)
	suite.mockService.AssertExpectations(suite.T())
}

func TestSeatBookingControllerTestSuite(t *testing.T) {
	suite.Run(t, new(SeatBookingControllerTestSuite))
}
