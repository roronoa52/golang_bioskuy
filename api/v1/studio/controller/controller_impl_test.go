package controller

import (
	"bioskuy/api/v1/studio/dto"
	"bioskuy/api/v1/studio/mock/servicemock"
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

type StudioControllerTestSuite struct {
	suite.Suite
	mockService *servicemock.MockStudioService
	controller  StudioController
	router      *gin.Engine
	ctx         context.Context
}

func (suite *StudioControllerTestSuite) SetupTest() {
	suite.mockService = new(servicemock.MockStudioService)
	suite.controller = NewStudioController(suite.mockService)
	suite.router = gin.Default()
	suite.ctx = context.Background()

	suite.router.POST("/studios", suite.controller.Create)
	suite.router.GET("/studios/:studioId", suite.controller.FindById)
	suite.router.GET("/studios", suite.controller.FindAll)
	suite.router.PUT("/studios/:studioId", suite.controller.Update)
	suite.router.DELETE("/studios/:studioId", suite.controller.Delete)
}

func TestStudioControllerTestSuite(t *testing.T) {
	suite.Run(t, new(StudioControllerTestSuite))
}

// Create
func (suite *StudioControllerTestSuite) TestCreate_Success() {
	studioRequest := dto.CreateStudioRequest{
		Name: "Studio 1",
	}
	studioResponse := dto.StudioResponse{
		ID:   "new-id",
		Name: "Studio 1",
	}

	suite.mockService.On("Create", mock.Anything, mock.AnythingOfType("dto.CreateStudioRequest"), mock.Anything).Return(studioResponse, nil)

	payload, _ := json.Marshal(studioRequest)
	req := httptest.NewRequest(http.MethodPost, "/studios", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)
	var response web.FormatResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), studioResponse, response.Data)
}

func (suite *StudioControllerTestSuite) TestCreate_BindError() {
	payload := `{"invalid json"}`
	req := httptest.NewRequest(http.MethodPost, "/studios", bytes.NewBuffer([]byte(payload)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusForbidden, w.Code)
}

func (suite *StudioControllerTestSuite) TestCreate_ServiceError() {
	studioRequest := dto.CreateStudioRequest{
		Name: "Studio 1",
	}
	serviceError := exception.ForbiddenError{Message: "service error"}

	suite.mockService.On("Create", mock.Anything, mock.AnythingOfType("dto.CreateStudioRequest"), mock.Anything).Return(dto.StudioResponse{}, serviceError)

	payload, _ := json.Marshal(studioRequest)
	req := httptest.NewRequest(http.MethodPost, "/studios", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusForbidden, w.Code)
}

// FindById
func (suite *StudioControllerTestSuite) TestFindById_Success() {
	studioResponse := dto.StudioResponse{
		ID:   "some-id",
		Name: "Studio 1",
	}

	suite.mockService.On("FindByID", mock.Anything, "some-id", mock.Anything).Return(studioResponse, nil)

	req := httptest.NewRequest(http.MethodGet, "/studios/some-id", nil)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	var response web.FormatResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), studioResponse, response.Data)
}

func (suite *StudioControllerTestSuite) TestFindById_NotFoundError() {
	serviceError := exception.NotFoundError{Message: "not found"}

	suite.mockService.On("FindByID", mock.Anything, "some-id", mock.Anything).Return(dto.StudioResponse{}, serviceError)

	req := httptest.NewRequest(http.MethodGet, "/studios/some-id", nil)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}

// FindAll
func (suite *StudioControllerTestSuite) TestFindAll_Success() {
	studioResponses := []dto.StudioResponse{
		{ID: "id1", Name: "Studio 1"},
		{ID: "id2", Name: "Studio 2"},
	}

	suite.mockService.On("FindAll", mock.Anything, mock.Anything).Return(studioResponses, nil)

	req := httptest.NewRequest(http.MethodGet, "/studios", nil)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	var response web.FormatResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), studioResponses, response.Data)
}

func (suite *StudioControllerTestSuite) TestFindAll_ServiceError() {
	serviceError := exception.InternalServerError{Message: "internal error"}

	suite.mockService.On("FindAll", mock.Anything, mock.Anything).Return(nil, serviceError)

	req := httptest.NewRequest(http.MethodGet, "/studios", nil)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}

// Update
func (suite *StudioControllerTestSuite) TestUpdate_Success() {
	studioRequest := dto.UpdateStudioRequest{
		Name: "Updated Studio",
	}
	studioResponse := dto.StudioResponse{
		ID:   "some-id",
		Name: "Updated Studio",
	}

	suite.mockService.On("Update", mock.Anything, mock.AnythingOfType("dto.UpdateStudioRequest"), mock.Anything).Return(studioResponse, nil)

	payload, _ := json.Marshal(studioRequest)
	req := httptest.NewRequest(http.MethodPut, "/studios/some-id", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	var response web.FormatResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), studioResponse, response.Data)
}

func (suite *StudioControllerTestSuite) TestUpdate_BindError() {
	payload := `{"invalid json"}`
	req := httptest.NewRequest(http.MethodPut, "/studios/some-id", bytes.NewBuffer([]byte(payload)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusForbidden, w.Code)
}

func (suite *StudioControllerTestSuite) TestUpdate_ServiceError() {
	studioRequest := dto.UpdateStudioRequest{
		Name: "Updated Studio",
	}
	serviceError := exception.InternalServerError{Message: "internal error"}

	suite.mockService.On("Update", mock.Anything, mock.AnythingOfType("dto.UpdateStudioRequest"), mock.Anything).Return(dto.StudioResponse{}, serviceError)

	payload, _ := json.Marshal(studioRequest)
	req := httptest.NewRequest(http.MethodPut, "/studios/some-id", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}

// Delete
func (suite *StudioControllerTestSuite) TestDelete_Success() {
	suite.mockService.On("Delete", mock.Anything, "some-id", mock.Anything).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/studios/some-id", nil)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	var response web.FormatResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "OK", response.Data)
}

func (suite *StudioControllerTestSuite) TestDelete_ServiceError() {
	serviceError := exception.InternalServerError{Message: "internal error"}

	suite.mockService.On("Delete", mock.Anything, "some-id", mock.Anything).Return(serviceError)

	req := httptest.NewRequest(http.MethodDelete, "/studios/some-id", nil)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}
