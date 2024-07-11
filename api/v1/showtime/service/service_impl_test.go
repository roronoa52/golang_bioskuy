package service

import (
	"bioskuy/api/v1/movies/entity"
	entityMovie "bioskuy/api/v1/movies/entity"
	movieMock "bioskuy/api/v1/movies/mock/repomock"
	"bioskuy/api/v1/showtime/dto"
	ShowtimeEntity "bioskuy/api/v1/showtime/entity"
	showTimeMock "bioskuy/api/v1/showtime/mock/repomock"
	entityStudio "bioskuy/api/v1/studio/entity"
	"context"
	"database/sql"
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ShowtimeServiceTestSuite struct {
	suite.Suite
	service        *showtimesServiceImpl
	mockRepo       *showTimeMock.MockShowtimeRepository
	mockRepoMovie  *movieMock.MockMovieRepository
	mockRepoStudio *showTimeMock.MockStudioRepository
	sqlMock        sqlmock.Sqlmock
	validator      *validator.Validate
	db             *sql.DB
}

func (suite *ShowtimeServiceTestSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	assert.NoError(suite.T(), err)

	suite.db = db
	suite.sqlMock = mock
	suite.mockRepo = &showTimeMock.MockShowtimeRepository{}
	suite.mockRepoMovie = &movieMock.MockMovieRepository{}
	suite.mockRepoStudio = &showTimeMock.MockStudioRepository{}
	suite.validator = validator.New()

	suite.service = &showtimesServiceImpl{
		suite.mockRepo,
		suite.mockRepoMovie,
		suite.mockRepoStudio,
		suite.validator,
		suite.db,
	}
}

func (suite *ShowtimeServiceTestSuite) TearDownTest() {
	suite.db.Close()
}

func (suite *ShowtimeServiceTestSuite) TestCreate_Success() {
	ctx := context.Background()
	ginCtx, _ := gin.CreateTestContext(httptest.NewRecorder())
	request := dto.ShowtimeRequest{
		MovieID:   "1",
		StudioID:  "1",
		ShowStart: "2024-07-09T10:00:00Z",
	}

	suite.sqlMock.ExpectBegin()

	movie := entityMovie.Movie{
		ID:       "1",
		Duration: 2, // 2 hours
	}
	suite.mockRepoMovie.On("GetByID", request.MovieID).Return(movie, nil).Once()

	studio := entityStudio.Studio{
		ID: "1",
	}
	suite.mockRepoStudio.On("FindByID", ctx, mock.Anything, request.StudioID, ginCtx).Return(studio, nil).Once()

	showtime := ShowtimeEntity.Showtime{
		MovieID:   request.MovieID,
		StudioID:  request.StudioID,
		ShowStart: time.Date(2024, 7, 9, 10, 0, 0, 0, time.UTC),
		ShowEnd:   time.Date(2024, 7, 9, 12, 0, 0, 0, time.UTC),
	}
	suite.mockRepo.On("FindConflictingShowtimes", ctx, mock.Anything, studio, mock.Anything, ginCtx).Return(nil).Once()
	suite.mockRepo.On("Save", ctx, mock.Anything, showtime, ginCtx).Return(showtime, nil).Once()

	result, err := suite.service.Create(ctx, request, ginCtx)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), request.MovieID, result.MovieID)
	assert.Equal(suite.T(), request.StudioID, result.StudioID)

	suite.mockRepo.AssertExpectations(suite.T())
	suite.mockRepoMovie.AssertExpectations(suite.T())
	suite.mockRepoStudio.AssertExpectations(suite.T())
	suite.sqlMock.ExpectationsWereMet()
}

func (suite *ShowtimeServiceTestSuite) TestFindByID_Success() {
	ctx := context.Background()
	ginCtx, _ := gin.CreateTestContext(httptest.NewRecorder())
	id := "1"

	suite.sqlMock.ExpectBegin()
	showtime := ShowtimeEntity.Showtime{
		ID:       "1",
		StudioID: "1",
		MovieID:  "1",
	}
	suite.mockRepo.On("FindByID", ctx, mock.Anything, id, ginCtx).Return(showtime, nil).Once()
	suite.sqlMock.ExpectCommit()

	result, err := suite.service.FindByID(ctx, id, ginCtx)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), id, result.ID)
	suite.mockRepo.AssertExpectations(suite.T())
	suite.sqlMock.ExpectationsWereMet()
}

func (suite *ShowtimeServiceTestSuite) TestFindAll_Success() {
	ctx := context.Background()
	ginCtx, _ := gin.CreateTestContext(httptest.NewRecorder())

	suite.sqlMock.ExpectBegin()
	showtimes := []ShowtimeEntity.Showtime{
		{ID: "1", StudioID: "1", MovieID: "1"},
	}
	suite.mockRepo.On("FindAll", ctx, mock.Anything, ginCtx).Return(showtimes, nil).Once()
	suite.sqlMock.ExpectCommit()

	result, err := suite.service.FindAll(ctx, ginCtx)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), result)
	assert.Equal(suite.T(), showtimes[0].ID, result[0].ID)
	suite.mockRepo.AssertExpectations(suite.T())
	suite.sqlMock.ExpectationsWereMet()
}

func (suite *ShowtimeServiceTestSuite) TestDelete_Success() {
	ctx := context.Background()
	ginCtx, _ := gin.CreateTestContext(httptest.NewRecorder())
	id := "1"

	suite.sqlMock.ExpectBegin()
	showtime := ShowtimeEntity.Showtime{ID: id}
	suite.mockRepo.On("FindByID", ctx, mock.Anything, id, ginCtx).Return(showtime, nil).Once()
	suite.mockRepo.On("Delete", ctx, mock.Anything, id, ginCtx).Return(nil).Once()
	suite.sqlMock.ExpectCommit()

	err := suite.service.Delete(ctx, id, ginCtx)
	assert.NoError(suite.T(), err)
	suite.mockRepo.AssertExpectations(suite.T())
	suite.sqlMock.ExpectationsWereMet()
}

func (suite *ShowtimeServiceTestSuite) TestCreate_ValidationError() {
	ctx := context.Background()
	ginCtx, _ := gin.CreateTestContext(httptest.NewRecorder())
	request := dto.ShowtimeRequest{
		MovieID:   "",
		StudioID:  "",
		ShowStart: "invalid-date",
	}

	result, err := suite.service.Create(ctx, request, ginCtx)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), dto.CreateShowtimesResponseDTO{}, result)
}

func (suite *ShowtimeServiceTestSuite) TestCreate_DBError() {
	ctx := context.Background()
	ginCtx, _ := gin.CreateTestContext(httptest.NewRecorder())
	request := dto.ShowtimeRequest{
		MovieID:   "1",
		StudioID:  "1",
		ShowStart: "2024-07-09T10:00:00Z",
	}

	suite.sqlMock.ExpectBegin().WillReturnError(errors.New("db error"))

	result, err := suite.service.Create(ctx, request, ginCtx)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "db error", err.Error())
	assert.Equal(suite.T(), dto.CreateShowtimesResponseDTO{}, result)
	suite.sqlMock.ExpectationsWereMet()
}

func (suite *ShowtimeServiceTestSuite) TestCreate_MovieNotFound() {
	ctx := context.Background()
	ginCtx, _ := gin.CreateTestContext(httptest.NewRecorder())
	request := dto.ShowtimeRequest{
		MovieID:   "1",
		StudioID:  "1",
		ShowStart: "2024-07-09T10:00:00Z",
	}

	suite.sqlMock.ExpectBegin()
	suite.mockRepoMovie.On("GetByID", request.MovieID).Return(entity.Movie{}, errors.New("movie not found")).Once()
	suite.sqlMock.ExpectRollback()

	result, err := suite.service.Create(ctx, request, ginCtx)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "movie not found", err.Error())
	assert.Equal(suite.T(), dto.CreateShowtimesResponseDTO{}, result)
	suite.mockRepoMovie.AssertExpectations(suite.T())
	suite.sqlMock.ExpectationsWereMet()
}

func (suite *ShowtimeServiceTestSuite) TestFindByID_NotFoundError() {
	ctx := context.Background()
	ginCtx, _ := gin.CreateTestContext(httptest.NewRecorder())
	id := "1"

	suite.sqlMock.ExpectBegin()
	suite.mockRepo.On("FindByID", ctx, mock.Anything, id, ginCtx).Return(ShowtimeEntity.Showtime{}, errors.New("not found")).Once()
	suite.sqlMock.ExpectRollback()

	result, err := suite.service.FindByID(ctx, id, ginCtx)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "not found", err.Error())
	assert.Equal(suite.T(), dto.ShowtimesResponse{}, result)
	suite.mockRepo.AssertExpectations(suite.T())
	suite.sqlMock.ExpectationsWereMet()
}

func (suite *ShowtimeServiceTestSuite) TestFindAll_DBError() {
	ctx := context.Background()
	ginCtx, _ := gin.CreateTestContext(httptest.NewRecorder())

	suite.sqlMock.ExpectBegin().WillReturnError(errors.New("db error"))

	result, err := suite.service.FindAll(ctx, ginCtx)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "db error", err.Error())
	assert.Nil(suite.T(), result)
	suite.sqlMock.ExpectationsWereMet()
}

func (suite *ShowtimeServiceTestSuite) TestDelete_NotFoundError() {
	ctx := context.Background()
	ginCtx, _ := gin.CreateTestContext(httptest.NewRecorder())
	id := "1"

	suite.sqlMock.ExpectBegin()
	suite.mockRepo.On("FindByID", ctx, mock.Anything, id, ginCtx).Return(ShowtimeEntity.Showtime{}, errors.New("not found")).Once()
	suite.sqlMock.ExpectRollback()

	err := suite.service.Delete(ctx, id, ginCtx)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "not found", err.Error())
	suite.mockRepo.AssertExpectations(suite.T())
	suite.sqlMock.ExpectationsWereMet()
}

func (suite *ShowtimeServiceTestSuite) TestDelete_DBError() {
	ctx := context.Background()
	ginCtx, _ := gin.CreateTestContext(httptest.NewRecorder())
	id := "1"

	suite.sqlMock.ExpectBegin()
	suite.mockRepo.On("FindByID", ctx, mock.Anything, id, ginCtx).Return(ShowtimeEntity.Showtime{ID: id}, nil).Once()
	suite.mockRepo.On("Delete", ctx, mock.Anything, id, ginCtx).Return(errors.New("db error")).Once()
	suite.sqlMock.ExpectRollback()

	err := suite.service.Delete(ctx, id, ginCtx)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "db error", err.Error())
	suite.mockRepo.AssertExpectations(suite.T())
	suite.sqlMock.ExpectationsWereMet()
}

func TestShowtimeServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ShowtimeServiceTestSuite))
}
