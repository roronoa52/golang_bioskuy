package servicemock

import (
	"bioskuy/api/v1/genre/dto"
	"bioskuy/api/v1/genre/entity"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockGenreService struct {
	mock.Mock
}

func (m *MockGenreService) CreateGenre(genre entity.Genre) (entity.Genre, error) {
	args := m.Called(genre)
	return args.Get(0).(entity.Genre), args.Error(1)
}

func (m *MockGenreService) GetGenreByID(id uuid.UUID) (entity.Genre, error) {
	args := m.Called(id)
	return args.Get(0).(entity.Genre), args.Error(1)
}

func (m *MockGenreService) GetAll(page int, size int) ([]entity.Genre, dto.Paging, error) {
	args := m.Called(page, size)
	return args.Get(0).([]entity.Genre), args.Get(1).(dto.Paging), args.Error(2)
}

func (m *MockGenreService) UpdateGenre(genre entity.Genre) (entity.Genre, error) {
	args := m.Called(genre)
	return args.Get(0).(entity.Genre), args.Error(1)
}

func (m *MockGenreService) DeleteGenre(id uuid.UUID) (entity.Genre, error) {
	args := m.Called(id)
	return args.Get(0).(entity.Genre), args.Error(1)
}

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) ValidateToken(token string) (map[string]interface{}, error) {
	args := m.Called(token)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func MockAuthMiddleware(authService *MockAuthService, allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := authService.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}

		role, ok := claims["role"].(string)
		if !ok || !contains(allowedRoles, role) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "you don't have access to this feature"})
			return
		}

		c.Set("user_id", claims["user_id"])
		c.Set("name", claims["name"])
		c.Set("role", role)
		c.Set("email", claims["email"])

		c.Next()
	}
}

func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

type MockConfig struct {
	DBUser             string
	DBPassword         string
	DBHost             string
	DBName             string
	DBPort             string
	DriverName         string
	SecretKey          string
	DurationJWT        string
	GoogleClientID     string
	GoogleClientSecret string
}

func NewMockConfig() *MockConfig {
	return &MockConfig{
		DBUser:             "test_user",
		DBPassword:         "test_password",
		DBHost:             "localhost",
		DBName:             "test_db",
		DBPort:             "5432",
		DriverName:         "postgres",
		SecretKey:          "test_secret_key",
		DurationJWT:        "24h",
		GoogleClientID:     "test_google_client_id",
		GoogleClientSecret: "test_google_client_secret",
	}
}
