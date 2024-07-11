package genreroute_test

import (
	"bioskuy/api/v1/genre/controller"
	"bioskuy/api/v1/genre/mock/servicemock"
	"bioskuy/api/v1/genre/repository"
	"bioskuy/api/v1/genre/service"
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type GenreRouteTestSuite struct {
	suite.Suite
	router      *gin.Engine
	mockService *servicemock.MockAuthService
	mockConfig  *servicemock.MockConfig
}

func (suite *GenreRouteTestSuite) SetupTest() {
	suite.mockService = new(servicemock.MockAuthService)
	suite.mockConfig = servicemock.NewMockConfig()
	db, _ := sql.Open("postgres", "your_db_url")

	gin.SetMode(gin.TestMode)
	suite.router = gin.Default()

	genreRepo := repository.NewGenreRepository(db)
	genreService := service.NewGenreService(genreRepo)
	genreController := controller.NewGenreController(genreService)

	v1 := suite.router.Group("/api/v1")
	{
		genre := v1.Group("/genres")
		{
			genre.GET("/", genreController.GetAll)
			genre.POST("/", servicemock.MockAuthMiddleware(suite.mockService, "admin"), genreController.CreateGenre)
			genre.GET("/:id", genreController.GetGenre)
			genre.PUT("/:id", servicemock.MockAuthMiddleware(suite.mockService, "admin"), genreController.UpdateGenre)
			genre.DELETE("/:id", servicemock.MockAuthMiddleware(suite.mockService, "admin"), genreController.DeleteGenre)
		}
	}
}

func (suite *GenreRouteTestSuite) TestCreateGenre_Success() {
	claims := map[string]interface{}{
		"user_id": "test-user-id",
		"name":    "test-name",
		"role":    "admin",
		"email":   "test@example.com",
	}
	suite.mockService.On("ValidateToken", "valid-token").Return(claims, nil)

	createDTO := map[string]string{"name": "Action"}
	body, _ := json.Marshal(createDTO)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/genres", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer valid-token")
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	suite.Equal(http.StatusCreated, resp.Code)
	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	suite.Equal("Action", response["name"])
}

func TestGenreRouteTestSuite(t *testing.T) {
	suite.Run(t, new(GenreRouteTestSuite))
}
