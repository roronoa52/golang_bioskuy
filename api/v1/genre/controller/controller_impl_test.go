package controller

import (
	"bioskuy/api/v1/genre/dto"
	"bioskuy/api/v1/genre/entity"
	"bioskuy/api/v1/genre/mock/servicemock"
	"bioskuy/web"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type GenreControllerTestSuite struct {
	suite.Suite
	mockService *servicemock.MockGenreService
	router      *gin.Engine
}

var mockingGenre = entity.Genre{
	ID:   uuid.New(),
	Name: "Action",
}
var mockingGenres = []entity.Genre{mockingGenre}

func (suite *GenreControllerTestSuite) SetupTest() {
	suite.mockService = new(servicemock.MockGenreService)
	gin.SetMode(gin.TestMode)
	suite.router = gin.Default()
	ctrl := NewGenreController(suite.mockService)
	v1 := suite.router.Group("/api/v1")
	{
		genres := v1.Group("/genres")
		{
			genres.GET("", ctrl.GetAll)
			genres.POST("", ctrl.CreateGenre)
			genres.GET("/:id", ctrl.GetGenre)
			genres.PUT("/:id", ctrl.UpdateGenre)
			genres.DELETE("/:id", ctrl.DeleteGenre)
		}
	}
}

func TestGenreControllerTestSuite(t *testing.T) {
	suite.Run(t, new(GenreControllerTestSuite))
}

func (suite *GenreControllerTestSuite) TestGetAll_Success() {
	paging := dto.Paging{Page: 1, Size: 10, TotalRows: 1, TotalPages: 1}

	suite.mockService.On("GetAll", 1, 10).Return(mockingGenres, paging, nil)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/genres?page=1&size=10", nil)
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusOK, resp.Code)
	var response web.FormatResponsePaging
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), paging.Page, response.Paging.Page)
	assert.Equal(suite.T(), paging.TotalRows, response.Paging.TotalData)
}

func (suite *GenreControllerTestSuite) TestCreateGenre_Success() {
	suite.mockService.On("CreateGenre", mock.AnythingOfType("entity.Genre")).Return(mockingGenre, nil)

	createDTO := dto.CreateGenreDTO{Name: "Action"}
	body, _ := json.Marshal(createDTO)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/genres", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusCreated, resp.Code)
	var response dto.GenreResponseDTO
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockingGenre.ID, response.ID)
	assert.Equal(suite.T(), mockingGenre.Name, response.Name)
}

func (suite *GenreControllerTestSuite) TestCreateGenre_BadRequest() {
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/genres", nil)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusBadRequest, resp.Code)
}

func (suite *GenreControllerTestSuite) TestCreateGenre_Error() {
	suite.mockService.On("CreateGenre", mock.AnythingOfType("entity.Genre")).Return(entity.Genre{}, errors.New("error creating genre"))

	createDTO := dto.CreateGenreDTO{Name: "Action"}
	body, _ := json.Marshal(createDTO)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/genres", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, resp.Code)
}

func (suite *GenreControllerTestSuite) TestGetGenre_Success() {
	suite.mockService.On("GetGenreByID", mockingGenre.ID).Return(mockingGenre, nil)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/genres/"+mockingGenre.ID.String(), nil)
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusOK, resp.Code)
	var response dto.GenreResponseDTO
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockingGenre.ID, response.ID)
	assert.Equal(suite.T(), mockingGenre.Name, response.Name)
}

func (suite *GenreControllerTestSuite) TestGetGenre_InvalidID() {
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/genres/invalid-id", nil)
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusBadRequest, resp.Code)
}

func (suite *GenreControllerTestSuite) TestGetGenre_Error() {
	suite.mockService.On("GetGenreByID", mockingGenre.ID).Return(entity.Genre{}, errors.New("error fetching genre"))

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/genres/"+mockingGenre.ID.String(), nil)
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, resp.Code)
}

func (suite *GenreControllerTestSuite) TestUpdateGenre_Success() {
	updatedGenre := entity.Genre{ID: mockingGenre.ID, Name: "Updated Action"}
	updateDTO := dto.UpdateGenreDTO{Name: "Updated Action"}

	suite.mockService.On("UpdateGenre", updatedGenre).Return(updatedGenre, nil)

	body, _ := json.Marshal(updateDTO)
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/genres/"+mockingGenre.ID.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusOK, resp.Code)
	var response dto.GenreResponseDTO
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), updatedGenre.ID, response.ID)
	assert.Equal(suite.T(), updatedGenre.Name, response.Name)
}

func (suite *GenreControllerTestSuite) TestUpdateGenre_InvalidID() {
	updateDTO := dto.UpdateGenreDTO{Name: "Updated Action"}

	body, _ := json.Marshal(updateDTO)
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/genres/invalid-id", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusBadRequest, resp.Code)
	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Invalid ID format", response["error"])
}

func (suite *GenreControllerTestSuite) TestUpdateGenre_InvalidJSON() {
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/genres/"+mockingGenre.ID.String(), bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusBadRequest, resp.Code)
}

func (suite *GenreControllerTestSuite) TestUpdateGenre_ServiceError() {
	updatedGenre := entity.Genre{ID: mockingGenre.ID, Name: "Updated Action"}
	updateDTO := dto.UpdateGenreDTO{Name: "Updated Action"}

	suite.mockService.On("UpdateGenre", updatedGenre).Return(entity.Genre{}, errors.New("service error"))

	body, _ := json.Marshal(updateDTO)
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/genres/"+mockingGenre.ID.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, resp.Code)
	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "service error", response["error"])
}

func (suite *GenreControllerTestSuite) TestDeleteGenre_Success() {
	suite.mockService.On("DeleteGenre", mockingGenre.ID).Return(mockingGenre, nil)

	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/genres/"+mockingGenre.ID.String(), nil)
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusOK, resp.Code)
	var response dto.GenreResponseDTO
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockingGenre.ID, response.ID)
	assert.Equal(suite.T(), mockingGenre.Name, response.Name)
}

func (suite *GenreControllerTestSuite) TestDeleteGenre_InvalidID() {
	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/genres/invalid-id", nil)
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusBadRequest, resp.Code)
}

func (suite *GenreControllerTestSuite) TestDeleteGenre_Error() {
	suite.mockService.On("DeleteGenre", mockingGenre.ID).Return(entity.Genre{}, errors.New("error deleting genre"))

	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/genres/"+mockingGenre.ID.String(), nil)
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, resp.Code)
}
