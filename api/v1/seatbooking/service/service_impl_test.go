package service

import (
	"bioskuy/api/v1/seatbooking/dto"
	"bioskuy/api/v1/seatbooking/entity"
	"bioskuy/api/v1/seatbooking/mock/entitymock"
	mockSB "bioskuy/api/v1/seatbooking/mock/repomock"
	"bioskuy/helper"
	"context"
	"database/sql"
	"sync"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SeatBookingServiceTestSuite struct {
	suite.Suite
	repoSBMock *mockSB.SeatBookingRepositoryMock
	repoSTMock *mockSB.MockShowtimeRepository
	repoS      *mockSB.SeatRepository
	validate   *validator.Validate
	mockDb     *sql.DB
	mockSql    sqlmock.Sqlmock
	Mutex      sync.Mutex
	sBSB       SeatBookingService
	ctx        context.Context
	ginContext *gin.Context
}

func (suite *SeatBookingServiceTestSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	assert.NoError(suite.T(), err)

	suite.repoSBMock = new(mockSB.SeatBookingRepositoryMock)
	suite.repoSTMock = new(mockSB.MockShowtimeRepository)
	suite.repoS = new(mockSB.SeatRepository)
	suite.validate = &validator.Validate{}
	suite.validate = validator.New()
	suite.mockDb = db
	suite.mockSql = mock
	suite.Mutex = sync.Mutex{}
	suite.sBSB = NewSeatBookingService(suite.repoSBMock, suite.repoSTMock, suite.repoS, suite.validate, suite.mockDb)
	suite.ctx = context.Background()
	suite.ginContext = &gin.Context{}
}

func TestSeatBookingServiceTestSuite(t *testing.T) {
	suite.Run(t, new(SeatBookingServiceTestSuite))
}

func (suite *SeatBookingServiceTestSuite) TestCreate_Success() {
	var SeatBookingResponse dto.CreateSeatBookingResponse
	var SeatBookingRequest dto.CreateSeatBookingRequest

	SeatBookingRequest.UserID = entitymock.MockSeatBookingEntity.ID
	SeatBookingRequest.ShowtimeID = entitymock.MockSeatBookingRequest.ShowtimeID

	err := suite.validate.Struct(SeatBookingRequest)
	assert.NoError(suite.T(), err)

	suite.mockSql.ExpectBegin()
	tx, err := suite.mockDb.Begin()
	assert.NoError(suite.T(), err)
	defer helper.CommitAndRollback(tx, suite.ginContext)

	suite.repoSTMock.On("FindByID", suite.ctx, tx, entitymock.MockSeatBookingRequest.ShowtimeID, suite.ginContext).Return(entitymock.MockShowtimeEntity, nil)
	suite.repoS.On("FindByIDWithNotAvailable", suite.ctx, tx, entitymock.MockSeatBookingRequest.SeatID, suite.ginContext).Return(entitymock.MockSeatEntity, nil)
	seatbooking := entitymock.MockSeatBookingEntity
	suite.repoSBMock.On("Save", suite.ctx, tx, seatbooking, entitymock.MockSeatBookingRequest.SeatID, suite.ginContext).Return(entitymock.MockSeatBookingEntity, nil)

	updatedSeat := entitymock.MockSeatEntity
	updatedSeat.IsAvailable = false
	suite.repoS.On("Update", suite.ctx, tx, updatedSeat, suite.ginContext).Return(updatedSeat, nil)

	response, err := suite.sBSB.Create(suite.ctx, entitymock.MockSeatBookingRequest, entitymock.MockCreateSeatBookingRequest.UserID, suite.ginContext)

	// Validasi respons
	assert.Equal(suite.T(), SeatBookingResponse.ID, response.ID)
	assert.Equal(suite.T(), SeatBookingResponse.SeatID, response.SeatID)
	assert.Equal(suite.T(), SeatBookingResponse.SeatBookingID, response.SeatBookingID)

	err = suite.mockSql.ExpectationsWereMet()
	assert.NoError(suite.T(), err)
}

func (suite *SeatBookingServiceTestSuite) TestCreate_ValidationError() {
	var invalidRequest = entitymock.MockSeatBookingRequest
	invalidRequest.ShowtimeID = "" // UserID kosong untuk memicu error validasi

	_, err := suite.sBSB.Create(suite.ctx, invalidRequest, invalidRequest.SeatID, suite.ginContext)
	assert.Error(suite.T(), err)
}

func (suite *SeatBookingServiceTestSuite) TestFindByID_NotFound() {
	suite.mockSql.ExpectBegin()
	tx, err := suite.mockDb.Begin()
	assert.NoError(suite.T(), err)
	defer helper.CommitAndRollback(tx, suite.ginContext)

	suite.repoSBMock.On("FindByID", suite.ctx, tx, entitymock.MockSeatBookingEntity.ID, suite.ginContext).Return(entity.SeatBooking{}, sql.ErrNoRows)

	response, err := suite.sBSB.FindByID(suite.ctx, entitymock.MockSeatBookingEntity.ID, suite.ginContext)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), dto.SeatBookingResponse{}, response)
}

func (suite *SeatBookingServiceTestSuite) TestFindAll_NoSeatBookings() {
	suite.mockSql.ExpectBegin()
	tx, err := suite.mockDb.Begin()
	assert.NoError(suite.T(), err)
	defer helper.CommitAndRollback(tx, suite.ginContext)

	suite.repoSBMock.On("FindAll", suite.ctx, tx, suite.ginContext).Return([]entity.SeatBooking{}, sql.ErrNoRows)

	responses, err := suite.sBSB.FindAll(suite.ctx, suite.ginContext)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), []dto.SeatBookingResponse{}, responses)
}

func (suite *SeatBookingServiceTestSuite) TestDelete_Success() {
	suite.mockSql.ExpectBegin()
	tx, err := suite.mockDb.Begin()
	assert.NoError(suite.T(), err)
	defer helper.CommitAndRollback(tx, suite.ginContext)

	suite.repoSBMock.On("Delete", suite.ctx, tx, entitymock.MockSeatBookingEntity.ID, suite.ginContext).Return(sql.ErrConnDone)

	assert.NoError(suite.T(), err)
	assert.Nil(suite.T(), err)
}

func (suite *SeatBookingServiceTestSuite) TestDelete_Error() {
	suite.mockSql.ExpectBegin()
	tx, err := suite.mockDb.Begin()
	assert.NoError(suite.T(), err)
	defer helper.CommitAndRollback(tx, suite.ginContext)

	suite.repoSBMock.On("Delete", suite.ctx, tx, entitymock.MockSeatBookingEntity.ID, suite.ginContext).Return(sql.ErrConnDone)

	err = suite.sBSB.Delete(suite.ctx, entitymock.MockSeatBookingEntity.ID, suite.ginContext)
	assert.Error(suite.T(), err)
}
