package controller

import (
	"bioskuy/api/v1/payment/dto"
	"bioskuy/api/v1/payment/mock/servicemock"
	"bioskuy/exception"
	"bioskuy/web"
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

type PaymentControllerTestSuite struct {
	suite.Suite
	mockService *servicemock.MockPaymentService
	controller  PaymentController
	router      *gin.Engine
	ctx         context.Context
}

func (suite *PaymentControllerTestSuite) SetupTest() {
	suite.mockService = new(servicemock.MockPaymentService)
	suite.controller = NewPaymentController(suite.mockService)
	suite.router = gin.Default()
	suite.ctx = context.Background()

	suite.router.POST("/payments", suite.controller.Create)
	suite.router.GET("/payments/:paymentId", suite.controller.FindById)
	suite.router.GET("/payments", suite.controller.FindAll)
	suite.router.POST("/payments/notification", suite.controller.Notification)
}

func TestPaymentControllerTestSuite(t *testing.T) {
	suite.Run(t, new(PaymentControllerTestSuite))
}

// Create
func (suite *PaymentControllerTestSuite) TestCreate_Success() {
	paymentRequest := dto.PaymentRequest{
		SeatDetailForBookingID: "seat-id",
	}
	paymentResponse := dto.PaymentResponse{
		ID:                     "new-id",
		SeatDetailForBookingID: "seat-id",
		TotalSeat:              1,
		TotalPrice:             10000,
	}

	suite.mockService.On("Create", mock.Anything, mock.AnythingOfType("dto.PaymentRequest"), "user-id", mock.Anything).Return(paymentResponse, nil)

	payload, _ := json.Marshal(paymentRequest)
	req := httptest.NewRequest(http.MethodPost, "/payments", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), "user_id", "user-id"))
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)
	var response web.FormatResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), paymentResponse, response.Data)
}

func (suite *PaymentControllerTestSuite) TestCreate_BindError() {
	payload := `{"invalid json"}`
	req := httptest.NewRequest(http.MethodPost, "/payments", bytes.NewBuffer([]byte(payload)))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), "user_id", "user-id"))
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusForbidden, w.Code)
}

func (suite *PaymentControllerTestSuite) TestCreate_ServiceError() {
	paymentRequest := dto.PaymentRequest{
		SeatDetailForBookingID: "seat-id",
	}
	serviceError := exception.ForbiddenError{Message: "service error"}

	suite.mockService.On("Create", mock.Anything, mock.AnythingOfType("dto.PaymentRequest"), "user-id", mock.Anything).Return(dto.PaymentResponse{}, serviceError)

	payload, _ := json.Marshal(paymentRequest)
	req := httptest.NewRequest(http.MethodPost, "/payments", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), "user_id", "user-id"))
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusForbidden, w.Code)
}

// FindById
func (suite *PaymentControllerTestSuite) TestFindById_Success() {
	paymentResponse := dto.PaymentResponse{
		ID:                     "some-id",
		SeatDetailForBookingID: "seat-id",
		TotalSeat:              1,
		TotalPrice:             10000,
	}

	suite.mockService.On("FindByID", mock.Anything, "some-id", mock.Anything).Return(paymentResponse, nil)

	req := httptest.NewRequest(http.MethodGet, "/payments/some-id", nil)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	var response web.FormatResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), paymentResponse, response.Data)
}

func (suite *PaymentControllerTestSuite) TestFindById_NotFoundError() {
	serviceError := exception.NotFoundError{Message: "not found"}

	suite.mockService.On("FindByID", mock.Anything, "some-id", mock.Anything).Return(dto.PaymentResponse{}, serviceError)

	req := httptest.NewRequest(http.MethodGet, "/payments/some-id", nil)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}

// FindAll
func (suite *PaymentControllerTestSuite) TestFindAll_Success() {
	paymentResponses := []dto.PaymentResponse{
		{ID: "id1", SeatDetailForBookingID: "seat1", TotalSeat: 1, TotalPrice: 10000},
		{ID: "id2", SeatDetailForBookingID: "seat2", TotalSeat: 2, TotalPrice: 20000},
	}

	suite.mockService.On("FindAll", mock.Anything, mock.Anything).Return(paymentResponses, nil)

	req := httptest.NewRequest(http.MethodGet, "/payments", nil)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	var response web.FormatResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), paymentResponses, response.Data)
}

func (suite *PaymentControllerTestSuite) TestFindAll_ServiceError() {
	serviceError := exception.InternalServerError{Message: "internal error"}

	suite.mockService.On("FindAll", mock.Anything, mock.Anything).Return(nil, serviceError)

	req := httptest.NewRequest(http.MethodGet, "/payments", nil)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}

// Notification
func (suite *PaymentControllerTestSuite) TestNotification_Success() {
	notificationPayload := map[string]interface{}{
		"order_id": "some-id",
	}

	suite.mockService.On("Update", mock.Anything, notificationPayload, mock.Anything).Return()

	payload, _ := json.Marshal(notificationPayload)
	req := httptest.NewRequest(http.MethodPost, "/payments/notification", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *PaymentControllerTestSuite) TestNotification_BindError() {
	payload := `{"invalid json"}`
	req := httptest.NewRequest(http.MethodPost, "/payments/notification", bytes.NewBuffer([]byte(payload)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusForbidden, w.Code)
}

func (suite *PaymentControllerTestSuite) TestNotification_ValidationError() {
	notificationPayload := map[string]interface{}{
		"invalid_key": "value",
	}

	payload, _ := json.Marshal(notificationPayload)
	req := httptest.NewRequest(http.MethodPost, "/payments/notification", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusForbidden, w.Code)
}
