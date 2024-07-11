package controller

import (
	"bioskuy/api/v1/movies/dto"
	"bioskuy/api/v1/movies/entity"
	"bioskuy/api/v1/movies/mock/servicemock"
	"bioskuy/web"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MovieControllerTestSuite struct {
	suite.Suite
	mockService *servicemock.MockMovieService
	router      *gin.Engine
}

var mockingMovie = entity.Movie{
	ID:          "1",
	Title:       "The Matrix",
	Description: "A sci-fi action film",
	Price:       12,
	Duration:    150,
	Status:      "Active",
}

func (suite *MovieControllerTestSuite) SetupTest() {
	suite.mockService = new(servicemock.MockMovieService)
	gin.SetMode(gin.TestMode)
	suite.router = gin.Default()
	ctrl := NewMovieController(suite.mockService)
	v1 := suite.router.Group("/api/v1")
	{
		movies := v1.Group("/movies")
		{
			movies.GET("", ctrl.GetAllMovies)
			movies.POST("", ctrl.CreateMovie)
			movies.GET("/:id", ctrl.GetMovie)
			movies.PUT("/:id", ctrl.UpdateMovie)
			movies.DELETE("/:id", ctrl.DeleteMovie)
		}
	}
}

func TestMovieControllerTestSuite(t *testing.T) {
	suite.Run(t, new(MovieControllerTestSuite))
}
func (suite *MovieControllerTestSuite) TestGetAllMovies_Success() {
	movies := []entity.Movie{{Title: "Movie 1"}, {Title: "Movie 2"}}
	paging := web.Paging{Page: 1, TotalData: 2}
	suite.mockService.On("GetAllMovies", 1, 10).Return(movies, paging, nil)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/movies?page=1&size=10", nil)
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusOK, resp.Code)
	var response web.FormatResponsePaging
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, response.ResponseCode)
	assert.Equal(suite.T(), movies, response.Data)
	assert.Equal(suite.T(), paging.Page, response.Paging.Page)
	assert.Equal(suite.T(), paging.TotalData, response.Paging.TotalData)
}

func (suite *MovieControllerTestSuite) TestGetAllMovies_Error() {
	suite.mockService.On("GetAllMovies", 1, 10).Return(nil, nil, errors.New("error fetching movies"))

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/movies?page=1&size=10", nil)
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, resp.Code)
	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "error fetching movies", response["error"])
}

func (suite *MovieControllerTestSuite) TestCreateMovie_Success() {
	suite.mockService.On("CreateMovie", mock.AnythingOfType("entity.Movie")).Return(mockingMovie, nil)

	createDTO := dto.CreateMovieDTO{
		Title:       "The Matrix",
		Description: "A sci-fi action film",
		Price:       12,
		Duration:    150,
		Status:      "Active",
	}
	body, _ := json.Marshal(createDTO)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/movies", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusCreated, resp.Code)
	var response dto.MovieResponseDTO
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockingMovie.ID, response.ID)
	assert.Equal(suite.T(), mockingMovie.Title, response.Title)
}

func (suite *MovieControllerTestSuite) TestCreateMovie_BadRequest() {
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/movies", nil)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusBadRequest, resp.Code)
}

func (suite *MovieControllerTestSuite) TestCreateMovie_Error() {
	suite.mockService.On("CreateMovie", mock.AnythingOfType("entity.Movie")).Return(entity.Movie{}, errors.New("error creating movie"))

	createDTO := dto.CreateMovieDTO{
		Title:       "The Matrix",
		Description: "A sci-fi action film",
		Price:       12,
		Duration:    150,
		Status:      "Active",
	}
	body, _ := json.Marshal(createDTO)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/movies", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, resp.Code)
}

func (suite *MovieControllerTestSuite) TestGetMovie_Success() {
	suite.mockService.On("GetMovieByID", mockingMovie.ID).Return(mockingMovie, nil)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/movies/"+mockingMovie.ID, nil)
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusOK, resp.Code)
	var response dto.MovieResponseDTO
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockingMovie.ID, response.ID)
	assert.Equal(suite.T(), mockingMovie.Title, response.Title)
}

func (suite *MovieControllerTestSuite) TestGetMovie_InvalidID() {
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/movies/invalid-id", nil)
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusBadRequest, resp.Code)
}

func (suite *MovieControllerTestSuite) TestGetMovie_Error() {
	suite.mockService.On("GetMovieByID", mockingMovie.ID).Return(entity.Movie{}, errors.New("error fetching movie"))

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/movies/"+mockingMovie.ID, nil)
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, resp.Code)
}

func (suite *MovieControllerTestSuite) TestUpdateMovie_Success() {
	updatedMovie := entity.Movie{
		ID:          mockingMovie.ID,
		Title:       "Updated Matrix",
		Description: "Updated sci-fi action film",
		Price:       15,
		Duration:    160,
		Status:      "Inactive",
	}
	updateDTO := dto.UpdateMovieDTO{
		Title:       "Updated Matrix",
		Description: "Updated sci-fi action film",
		Price:       15,
		Duration:    160,
		Status:      "Inactive",
	}

	suite.mockService.On("UpdateMovie", updatedMovie).Return(updatedMovie, nil)

	body, _ := json.Marshal(updateDTO)
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/movies/"+mockingMovie.ID, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusOK, resp.Code)
	var response dto.MovieResponseDTO
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), updatedMovie.ID, response.ID)
	assert.Equal(suite.T(), updatedMovie.Title, response.Title)
}

func (suite *MovieControllerTestSuite) TestUpdateMovie_InvalidID() {
	updateDTO := dto.UpdateMovieDTO{
		Title:       "Updated Matrix",
		Description: "Updated sci-fi action film",
		Price:       15,
		Duration:    160,
		Status:      "Inactive",
	}

	body, _ := json.Marshal(updateDTO)
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/movies/invalid-id", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusBadRequest, resp.Code)
	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Invalid ID format", response["error"])
}

func (suite *MovieControllerTestSuite) TestUpdateMovie_InvalidJSON() {
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/movies/"+mockingMovie.ID, bytes.NewBuffer([]byte("{invalid}")))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusBadRequest, resp.Code)
	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Invalid JSON", response["error"])
}

func (suite *MovieControllerTestSuite) TestUpdateMovie_Error() {
	updateDTO := dto.UpdateMovieDTO{
		Title:       "Updated Matrix",
		Description: "Updated sci-fi action film",
		Price:       15,
		Duration:    160,
		Status:      "Inactive",
	}

	suite.mockService.On("UpdateMovie", mock.AnythingOfType("entity.Movie")).Return(entity.Movie{}, errors.New("error updating movie"))

	body, _ := json.Marshal(updateDTO)
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/movies/"+mockingMovie.ID, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, resp.Code)
	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "error updating movie", response["error"])
}

func (suite *MovieControllerTestSuite) TestDeleteMovie_Success() {
	suite.mockService.On("DeleteMovie", mockingMovie.ID).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/movies/"+mockingMovie.ID, nil)
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusOK, resp.Code)
}

func (suite *MovieControllerTestSuite) TestDeleteMovie_InvalidID() {
	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/movies/invalid-id", nil)
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusBadRequest, resp.Code)
	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Invalid ID format", response["error"])
}

func (suite *MovieControllerTestSuite) TestDeleteMovie_Error() {
	suite.mockService.On("DeleteMovie", mockingMovie.ID).Return(errors.New("error deleting movie"))

	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/movies/"+mockingMovie.ID, nil)
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, resp.Code)
	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "error deleting movie", response["error"])
}
