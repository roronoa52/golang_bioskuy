package service_test

import (
	eS "bioskuy/api/v1/seat/entity"
	"bioskuy/api/v1/studio/dto"
	"bioskuy/api/v1/studio/entity"
	"bioskuy/api/v1/studio/mock/repomock"
	"bioskuy/api/v1/studio/service"
	"bioskuy/exception"
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type StudioServiceTestSuite struct {
	suite.Suite
	mockStudioRepo *repomock.MockStudioRepository
	mockSeatRepo   *repomock.MockSeatRepository
	mockValidator  *validator.Validate
	mockDB         *sql.DB
	sqlMock        sqlmock.Sqlmock
	studioService  service.StudioService
	ctx            context.Context
}

func (suite *StudioServiceTestSuite) SetupTest() {
	var err error
	suite.mockStudioRepo = new(repomock.MockStudioRepository)
	suite.mockSeatRepo = new(repomock.MockSeatRepository)
	suite.mockValidator = validator.New()
	suite.mockDB, suite.sqlMock, err = sqlmock.New()
	assert.NoError(suite.T(), err)
	suite.studioService = service.NewStudioService(suite.mockStudioRepo, suite.mockValidator, suite.mockDB, suite.mockSeatRepo)
	suite.ctx = context.Background()
}

func (suite *StudioServiceTestSuite) TearDownTest() {
	suite.mockDB.Close()
}

// Create
func (suite *StudioServiceTestSuite) TestCreate_Success() {
	ginCtx, _ := gin.CreateTestContext(nil)
	request := dto.CreateStudioRequest{
		Name:       "Studio 1",
		Capacity:   50,
		MaxRowSeat: 10,
	}

	studioEntity := entity.Studio{
		ID:       "new-id",
		Name:     "Studio 1",
		Capacity: 50,
	}

	suite.mockStudioRepo.On("Save", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(studioEntity, nil).Once()

	for row := 0; row < 5; row++ {
		for seatNum := 1; seatNum <= 10; seatNum++ {
			seatName := string(rune('A'+row)) + "-" + fmt.Sprintf("%d", seatNum)
			seat := eS.Seat{
				Name:        seatName,
				IsAvailable: true,
				StudioID:    "new-id",
			}
			suite.mockSeatRepo.On("Save", mock.Anything, mock.Anything, seat, mock.Anything).Return(seat, nil).Once()
		}
	}

	suite.sqlMock.ExpectBegin()
	suite.sqlMock.ExpectCommit()

	response, err := suite.studioService.Create(suite.ctx, request, ginCtx)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), dto.StudioResponse{ID: "new-id", Name: "Studio 1", Capacity: 50}, response)
}

func (suite *StudioServiceTestSuite) TestCreate_ValidationError() {
	ginCtx, _ := gin.CreateTestContext(nil)
	request := dto.CreateStudioRequest{
		Name:       "",
		Capacity:   50,
		MaxRowSeat: 10,
	}

	response, err := suite.studioService.Create(suite.ctx, request, ginCtx)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), dto.StudioResponse{}, response)
}

func (suite *StudioServiceTestSuite) TestCreate_RepoError() {
	ginCtx, _ := gin.CreateTestContext(nil)
	request := dto.CreateStudioRequest{
		Name:       "Studio 1",
		Capacity:   50,
		MaxRowSeat: 10,
	}

	suite.mockStudioRepo.On("Save", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(entity.Studio{}, exception.InternalServerError{Message: "error"}).Once()

	suite.sqlMock.ExpectBegin()
	suite.sqlMock.ExpectRollback()

	response, err := suite.studioService.Create(suite.ctx, request, ginCtx)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), dto.StudioResponse{}, response)
}

// FindByID
func (suite *StudioServiceTestSuite) TestFindByID_Success() {
	ginCtx, _ := gin.CreateTestContext(nil)
	studioEntity := entity.Studio{
		ID:       "some-id",
		Name:     "Studio 1",
		Capacity: 50,
	}

	suite.mockStudioRepo.On("FindByID", mock.Anything, mock.Anything, "some-id", mock.Anything).Return(studioEntity, nil).Once()

	suite.sqlMock.ExpectBegin()
	suite.sqlMock.ExpectCommit()

	response, err := suite.studioService.FindByID(suite.ctx, "some-id", ginCtx)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), dto.StudioResponse{ID: "some-id", Name: "Studio 1", Capacity: 50}, response)
}

func (suite *StudioServiceTestSuite) TestFindByID_NotFoundError() {
	ginCtx, _ := gin.CreateTestContext(nil)

	suite.mockStudioRepo.On("FindByID", mock.Anything, mock.Anything, "some-id", mock.Anything).Return(entity.Studio{}, exception.NotFoundError{Message: "not found"}).Once()

	suite.sqlMock.ExpectBegin()
	suite.sqlMock.ExpectRollback()

	response, err := suite.studioService.FindByID(suite.ctx, "some-id", ginCtx)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), dto.StudioResponse{}, response)
}

// FindAll
func (suite *StudioServiceTestSuite) TestFindAll_Success() {
	ginCtx, _ := gin.CreateTestContext(nil)
	studios := []entity.Studio{
		{ID: "id1", Name: "Studio 1", Capacity: 50},
		{ID: "id2", Name: "Studio 2", Capacity: 60},
	}

	suite.mockStudioRepo.On("FindAll", mock.Anything, mock.Anything, mock.Anything).Return(studios, nil).Once()

	suite.sqlMock.ExpectBegin()
	suite.sqlMock.ExpectCommit()

	response, err := suite.studioService.FindAll(suite.ctx, ginCtx)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), response, 2)
}

func (suite *StudioServiceTestSuite) TestFindAll_ServiceError() {
	ginCtx, _ := gin.CreateTestContext(nil)

	suite.mockStudioRepo.On("FindAll", mock.Anything, mock.Anything, mock.Anything).Return(nil, exception.InternalServerError{Message: "error"}).Once()

	suite.sqlMock.ExpectBegin()
	suite.sqlMock.ExpectRollback()

	response, err := suite.studioService.FindAll(suite.ctx, ginCtx)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), response)
}

// Update
func (suite *StudioServiceTestSuite) TestUpdate_Success() {
	ginCtx, _ := gin.CreateTestContext(nil)
	request := dto.UpdateStudioRequest{
		ID:   "some-id",
		Name: "Updated Studio",
	}

	studioEntity := entity.Studio{
		ID:       "some-id",
		Name:     "Updated Studio",
		Capacity: 50,
	}

	suite.mockStudioRepo.On("FindByID", mock.Anything, mock.Anything, "some-id", mock.Anything).Return(studioEntity, nil).Once()
	suite.mockStudioRepo.On("Update", mock.Anything, mock.Anything, studioEntity, mock.Anything).Return(studioEntity, nil).Once()

	suite.sqlMock.ExpectBegin()
	suite.sqlMock.ExpectCommit()

	response, err := suite.studioService.Update(suite.ctx, request, ginCtx)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), dto.StudioResponse{ID: "some-id", Name: "Updated Studio", Capacity: 50}, response)
}

func (suite *StudioServiceTestSuite) TestUpdate_ValidationError() {
	ginCtx, _ := gin.CreateTestContext(nil)
	request := dto.UpdateStudioRequest{
		ID:   "some-id",
		Name: "",
	}

	response, err := suite.studioService.Update(suite.ctx, request, ginCtx)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), dto.StudioResponse{}, response)
}

// Update
func (suite *StudioServiceTestSuite) TestUpdate_NotFoundError() {
	ginCtx, _ := gin.CreateTestContext(nil)
	request := dto.UpdateStudioRequest{
		ID:   "some-id",
		Name: "Updated Studio",
	}

	suite.mockStudioRepo.On("FindByID", mock.Anything, mock.Anything, "some-id", mock.Anything).Return(entity.Studio{}, exception.NotFoundError{Message: "not found"}).Once()

	suite.sqlMock.ExpectBegin()
	suite.sqlMock.ExpectRollback()

	response, err := suite.studioService.Update(suite.ctx, request, ginCtx)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), dto.StudioResponse{}, response)
}

// Delete
func (suite *StudioServiceTestSuite) TestDelete_Success() {
	ginCtx, _ := gin.CreateTestContext(nil)
	studioEntity := entity.Studio{
		ID:       "some-id",
		Name:     "Studio 1",
		Capacity: 50,
	}

	suite.mockStudioRepo.On("FindByID", mock.Anything, mock.Anything, "some-id", mock.Anything).Return(studioEntity, nil).Once()
	suite.mockSeatRepo.On("Delete", mock.Anything, mock.Anything, "some-id", mock.Anything).Return(nil).Once()
	suite.mockStudioRepo.On("Delete", mock.Anything, mock.Anything, "some-id", mock.Anything).Return(nil).Once()

	suite.sqlMock.ExpectBegin()
	suite.sqlMock.ExpectCommit()

	err := suite.studioService.Delete(suite.ctx, "some-id", ginCtx)

	assert.NoError(suite.T(), err)
}

func (suite *StudioServiceTestSuite) TestDelete_StudioNotFound() {
	ginCtx, _ := gin.CreateTestContext(nil)

	suite.mockStudioRepo.On("FindByID", mock.Anything, mock.Anything, "some-id", mock.Anything).Return(entity.Studio{}, exception.NotFoundError{Message: "not found"}).Once()

	suite.sqlMock.ExpectBegin()
	suite.sqlMock.ExpectRollback()

	err := suite.studioService.Delete(suite.ctx, "some-id", ginCtx)

	assert.Error(suite.T(), err)
}

func (suite *StudioServiceTestSuite) TestDelete_SeatDeleteError() {
	ginCtx, _ := gin.CreateTestContext(nil)
	studioEntity := entity.Studio{
		ID:       "some-id",
		Name:     "Studio 1",
		Capacity: 50,
	}

	suite.mockStudioRepo.On("FindByID", mock.Anything, mock.Anything, "some-id", mock.Anything).Return(studioEntity, nil).Once()
	suite.mockSeatRepo.On("Delete", mock.Anything, mock.Anything, "some-id", mock.Anything).Return(exception.InternalServerError{Message: "error"}).Once()

	suite.sqlMock.ExpectBegin()
	suite.sqlMock.ExpectRollback()

	err := suite.studioService.Delete(suite.ctx, "some-id", ginCtx)

	assert.Error(suite.T(), err)
}

func (suite *StudioServiceTestSuite) TestDelete_StudioDeleteError() {
	ginCtx, _ := gin.CreateTestContext(nil)
	studioEntity := entity.Studio{
		ID:       "some-id",
		Name:     "Studio 1",
		Capacity: 50,
	}

	suite.mockStudioRepo.On("FindByID", mock.Anything, mock.Anything, "some-id", mock.Anything).Return(studioEntity, nil).Once()
	suite.mockSeatRepo.On("Delete", mock.Anything, mock.Anything, "some-id", mock.Anything).Return(nil).Once()
	suite.mockStudioRepo.On("Delete", mock.Anything, mock.Anything, "some-id", mock.Anything).Return(exception.InternalServerError{Message: "error"}).Once()

	suite.sqlMock.ExpectBegin()
	suite.sqlMock.ExpectRollback()

	err := suite.studioService.Delete(suite.ctx, "some-id", ginCtx)

	assert.Error(suite.T(), err)
}

func TestStudioServiceTestSuite(t *testing.T) {
	suite.Run(t, new(StudioServiceTestSuite))
}
