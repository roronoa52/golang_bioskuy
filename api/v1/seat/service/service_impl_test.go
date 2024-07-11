package service

import (
	"bioskuy/api/v1/seat/dto"
	"bioskuy/api/v1/seat/entity"
	"bioskuy/api/v1/seat/mock/repomock"
	"bioskuy/exception"
	"context"
	"database/sql"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type SeatServiceTestSuite struct {
	suite.Suite
	service  *seatService
	mockRepo *repomock.SeatRepository
	sqlMock  sqlmock.Sqlmock
	db       *sql.DB
}

func (suite *SeatServiceTestSuite) SetupTest() {
	suite.mockRepo = new(repomock.SeatRepository)
	validate := validator.New()
	db, sqlMock, _ := sqlmock.New()
	suite.db = db
	suite.sqlMock = sqlMock
	suite.service = &seatService{Repo: suite.mockRepo, Validate: validate, DB: db}
}

func (suite *SeatServiceTestSuite) TearDownTest() {
	suite.db.Close()
}

func (suite *SeatServiceTestSuite) TestFindByID_Success() {
	ctx := context.Background()
	ginCtx, _ := gin.CreateTestContext(httptest.NewRecorder())
	id := "1"

	mockSeat := entity.Seat{
		ID:          "1",
		Name:        "A1",
		IsAvailable: true,
		StudioID:    "1",
	}

	suite.sqlMock.ExpectBegin()
	suite.sqlMock.ExpectCommit()

	suite.mockRepo.On("FindByID", ctx, mock.Anything, id, ginCtx).Return(mockSeat, nil).Once()

	result, err := suite.service.FindByID(ctx, id, ginCtx)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), dto.SeatResponse{
		ID:          "1",
		Name:        "A1",
		IsAvailable: true,
		StudioID:    "1",
	}, result)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *SeatServiceTestSuite) TestFindByID_NotFound() {
	ctx := context.Background()
	ginCtx, _ := gin.CreateTestContext(httptest.NewRecorder())
	id := "1"

	suite.sqlMock.ExpectBegin()
	suite.sqlMock.ExpectRollback()

	suite.mockRepo.On("FindByID", ctx, mock.Anything, id, ginCtx).Return(entity.Seat{}, exception.NotFoundError{Message: "seat not found"}).Once()

	result, err := suite.service.FindByID(ctx, id, ginCtx)
	assert.Error(suite.T(), err)
	assert.IsType(suite.T(), exception.NotFoundError{}, err)
	assert.Equal(suite.T(), dto.SeatResponse{}, result)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *SeatServiceTestSuite) TestFindAll_Success() {
	ctx := context.Background()
	ginCtx, _ := gin.CreateTestContext(httptest.NewRecorder())
	id := "1"

	mockSeats := []entity.Seat{
		{
			ID:          "1",
			Name:        "A1",
			IsAvailable: true,
			StudioID:    "1",
		},
		{
			ID:          "2",
			Name:        "A2",
			IsAvailable: false,
			StudioID:    "1",
		},
	}

	suite.sqlMock.ExpectBegin()
	suite.sqlMock.ExpectCommit()

	suite.mockRepo.On("FindAll", ctx, id, mock.Anything, ginCtx).Return(mockSeats, nil).Once()

	result, err := suite.service.FindAll(ctx, id, ginCtx)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), []dto.SeatResponse{
		{
			ID:          "1",
			Name:        "A1",
			IsAvailable: true,
			StudioID:    "1",
		},
		{
			ID:          "2",
			Name:        "A2",
			IsAvailable: false,
			StudioID:    "1",
		},
	}, result)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *SeatServiceTestSuite) TestFindAll_NotFound() {
	ctx := context.Background()
	ginCtx, _ := gin.CreateTestContext(httptest.NewRecorder())
	id := "1"

	suite.sqlMock.ExpectBegin()
	suite.sqlMock.ExpectRollback()

	suite.mockRepo.On("FindAll", ctx, id, mock.Anything, ginCtx).Return([]entity.Seat{}, exception.NotFoundError{Message: "seats not found"}).Once()

	result, err := suite.service.FindAll(ctx, id, ginCtx)
	assert.Error(suite.T(), err)
	assert.IsType(suite.T(), exception.NotFoundError{}, err)
	assert.Equal(suite.T(), []dto.SeatResponse{}, result)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *SeatServiceTestSuite) TestFindByID_DBBeginError() {
	ctx := context.Background()
	ginCtx, _ := gin.CreateTestContext(httptest.NewRecorder())
	id := "1"

	// Simulate error when beginning a transaction
	suite.sqlMock.ExpectBegin().WillReturnError(errors.New("begin error"))

	result, err := suite.service.FindByID(ctx, id, ginCtx)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "begin error", err.Error())
	assert.Equal(suite.T(), dto.SeatResponse{}, result)
	suite.mockRepo.AssertNotCalled(suite.T(), "FindByID", ctx, mock.Anything, id, ginCtx)
}

func (suite *SeatServiceTestSuite) TestFindAll_DBBeginError() {
	ctx := context.Background()
	ginCtx, _ := gin.CreateTestContext(httptest.NewRecorder())
	id := "1"

	// Simulate error when beginning a transaction
	suite.sqlMock.ExpectBegin().WillReturnError(errors.New("begin error"))

	result, err := suite.service.FindAll(ctx, id, ginCtx)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "begin error", err.Error())
	assert.Equal(suite.T(), []dto.SeatResponse{}, result)
	suite.mockRepo.AssertNotCalled(suite.T(), "FindAll", ctx, id, mock.Anything, ginCtx)
}

func TestSeatServiceTestSuite(t *testing.T) {
	suite.Run(t, new(SeatServiceTestSuite))
}
