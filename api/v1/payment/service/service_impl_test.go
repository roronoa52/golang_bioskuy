package service

import (
	"bioskuy/api/v1/payment/dto"
	"bioskuy/api/v1/payment/entity"
	"bioskuy/api/v1/payment/mock/repomock"
	entitySeatBooking "bioskuy/api/v1/seatbooking/entity"
	"bioskuy/exception"
	"bioskuy/helper"
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type PaymentServiceTestSuite struct {
	suite.Suite
	mockDb              *sql.DB
	mockSql             sqlmock.Sqlmock
	mockRepo            *repomock.MockPaymentRepository
	mockRepoSeat        *repomock.MockSeatRepository
	mockRepoSeatBooking *repomock.MockSeatBookingRepository
	validate            *validator.Validate
	service             PaymentService
	ctx                 context.Context
	ginContext          *gin.Context
	env                 *helper.Config
}

func (suite *PaymentServiceTestSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	assert.NoError(suite.T(), err)
	suite.mockDb = db
	suite.mockSql = mock
	suite.mockRepo = &repomock.MockPaymentRepository{}
	suite.mockRepoSeat = &repomock.MockSeatRepository{}
	suite.mockRepoSeatBooking = &repomock.MockSeatBookingRepository{}
	suite.validate = validator.New()
	suite.env = &helper.Config{MIDTRANS_SERVER_KEY: "dummy-key"}
	suite.service = NewPaymentService(suite.mockRepo, suite.mockRepoSeat, suite.mockRepoSeatBooking, suite.validate, suite.mockDb, suite.env)
	suite.ctx = context.Background()
	suite.ginContext = &gin.Context{}
}

func TestPaymentServiceTestSuite(t *testing.T) {
	suite.Run(t, new(PaymentServiceTestSuite))
}

// Create
func (suite *PaymentServiceTestSuite) TestCreate_Success() {
	request := dto.PaymentRequest{
		SeatDetailForBookingID: "seat-id",
	}
	expectedResult := entity.Payment{
		ID:                     "new-id",
		UserID:                 "user-id",
		SeatDetailForBookingID: "seat-id",
		TotalSeat:              1,
		TotalPrice:             10000,
	}

	seatBooking := []entitySeatBooking.SeatBooking{
		{MoviePrice: 10000},
	}

	suite.mockSql.ExpectBegin()
	suite.mockRepoSeatBooking.On("FindAllPendingByUserID", suite.ctx, mock.Anything, "user-id", suite.ginContext).Return(seatBooking, nil)
	suite.mockRepo.On("Save", suite.ctx, mock.Anything, mock.AnythingOfType("entity.Payment"), suite.ginContext).Return(expectedResult, nil)
	suite.mockSql.ExpectCommit()

	response, err := suite.service.Create(suite.ctx, request, "user-id", suite.ginContext)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedResult.ID, response.ID)
	assert.Equal(suite.T(), request.SeatDetailForBookingID, response.SeatDetailForBookingID)
	assert.Equal(suite.T(), expectedResult.TotalSeat, response.TotalSeat)
	assert.Equal(suite.T(), expectedResult.TotalPrice, response.TotalPrice)
}

func (suite *PaymentServiceTestSuite) TestCreate_ValidationError() {
	request := dto.PaymentRequest{
		SeatDetailForBookingID: "",
	}

	response, err := suite.service.Create(suite.ctx, request, "user-id", suite.ginContext)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "", response.ID)
	assert.IsType(suite.T(), exception.ValidationError{}, err)
}

func (suite *PaymentServiceTestSuite) TestCreate_SaveError() {
	request := dto.PaymentRequest{
		SeatDetailForBookingID: "seat-id",
	}
	saveError := errors.New("Save Error")

	seatBooking := []entitySeatBooking.SeatBooking{
		{MoviePrice: 10000},
	}

	suite.mockSql.ExpectBegin()
	suite.mockRepoSeatBooking.On("FindAllPendingByUserID", suite.ctx, mock.Anything, "user-id", suite.ginContext).Return(seatBooking, nil)
	suite.mockRepo.On("Save", suite.ctx, mock.Anything, mock.AnythingOfType("entity.Payment"), suite.ginContext).Return(entity.Payment{}, saveError)
	suite.mockSql.ExpectRollback()

	response, err := suite.service.Create(suite.ctx, request, "user-id", suite.ginContext)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "", response.ID)
	assert.Equal(suite.T(), saveError, err)
}

// FindByID
func (suite *PaymentServiceTestSuite) TestFindByID_Success() {
	expectedResult := entity.Payment{
		ID:                     "some-id",
		UserID:                 "user-id",
		SeatDetailForBookingID: "seat-id",
		TotalSeat:              1,
		TotalPrice:             10000,
	}

	suite.mockSql.ExpectBegin()
	suite.mockRepo.On("FindByID", suite.ctx, mock.Anything, "some-id", suite.ginContext).Return(expectedResult, nil)
	suite.mockSql.ExpectCommit()

	response, err := suite.service.FindByID(suite.ctx, "some-id", suite.ginContext)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedResult.ID, response.ID)
	assert.Equal(suite.T(), expectedResult.UserID, response.UserID)
	assert.Equal(suite.T(), expectedResult.SeatDetailForBookingID, response.SeatDetailForBookingID)
	assert.Equal(suite.T(), expectedResult.TotalSeat, response.TotalSeat)
	assert.Equal(suite.T(), expectedResult.TotalPrice, response.TotalPrice)
}

func (suite *PaymentServiceTestSuite) TestFindByID_NotFoundError() {
	notFoundError := errors.New("Not Found Error")

	suite.mockSql.ExpectBegin()
	suite.mockRepo.On("FindByID", suite.ctx, mock.Anything, "some-id", suite.ginContext).Return(entity.Payment{}, notFoundError)
	suite.mockSql.ExpectRollback()

	response, err := suite.service.FindByID(suite.ctx, "some-id", suite.ginContext)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "", response.ID)
	assert.Equal(suite.T(), notFoundError, err)
}

// FindAll
func (suite *PaymentServiceTestSuite) TestFindAll_Success() {
	expectedResults := []entity.Payment{
		{ID: "id1", UserID: "user1", SeatDetailForBookingID: "seat1", TotalSeat: 1, TotalPrice: 10000},
		{ID: "id2", UserID: "user2", SeatDetailForBookingID: "seat2", TotalSeat: 2, TotalPrice: 20000},
	}

	suite.mockSql.ExpectBegin()
	suite.mockRepo.On("FindAll", suite.ctx, mock.Anything, suite.ginContext).Return(expectedResults, nil)
	suite.mockSql.ExpectCommit()

	response, err := suite.service.FindAll(suite.ctx, suite.ginContext)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), response, 2)
	assert.Equal(suite.T(), expectedResults[0].ID, response[0].ID)
	assert.Equal(suite.T(), expectedResults[0].UserID, response[0].UserID)
	assert.Equal(suite.T(), expectedResults[0].SeatDetailForBookingID, response[0].SeatDetailForBookingID)
}

func (suite *PaymentServiceTestSuite) TestFindAll_Error() {
	findAllError := errors.New("Find All Error")

	suite.mockSql.ExpectBegin()
	suite.mockRepo.On("FindAll", suite.ctx, mock.Anything, suite.ginContext).Return([]entity.Payment{}, findAllError)
	suite.mockSql.ExpectRollback()

	response, err := suite.service.FindAll(suite.ctx, suite.ginContext)

	assert.Error(suite.T(), err)
	assert.Len(suite.T(), response, 0)
	assert.Equal(suite.T(), findAllError, err)
}

func (suite *PaymentServiceTestSuite) TestUpdate_Success() {
	notificationPayload := map[string]interface{}{
		"order_id":           "some-id",
		"transaction_status": "settlement",
	}
	expectedResult := entity.Payment{
		ID:                     "some-id",
		UserID:                 "user-id",
		SeatDetailForBookingID: "seat-id",
		Status:                 "paid",
	}

	suite.mockSql.ExpectBegin()
	suite.mockRepo.On("FindByID", suite.ctx, mock.Anything, "some-id", suite.ginContext).Return(expectedResult, nil)
	suite.mockRepo.On("Update", suite.ctx, mock.Anything, mock.AnythingOfType("entity.Payment"), suite.ginContext).Return(expectedResult, nil)
	suite.mockRepoSeatBooking.On("Update", suite.ctx, mock.Anything, mock.AnythingOfType("entitySeatBooking.SeatBooking"), suite.ginContext).Return(nil, nil)
	suite.mockSql.ExpectCommit()

	suite.service.Update(suite.ctx, notificationPayload, suite.ginContext)
}

func (suite *PaymentServiceTestSuite) TestUpdate_DenyStatus() {
	notificationPayload := map[string]interface{}{
		"order_id":           "some-id",
		"transaction_status": "deny",
	}
	expectedResult := entity.Payment{
		ID:                     "some-id",
		UserID:                 "user-id",
		SeatDetailForBookingID: "seat-id",
		Status:                 "cancelled",
	}

	suite.mockSql.ExpectBegin()
	suite.mockRepo.On("FindByID", suite.ctx, mock.Anything, "some-id", suite.ginContext).Return(expectedResult, nil)
	suite.mockRepo.On("Update", suite.ctx, mock.Anything, mock.AnythingOfType("entity.Payment"), suite.ginContext).Return(expectedResult, nil)
	suite.mockRepoSeatBooking.On("Update", suite.ctx, mock.Anything, mock.AnythingOfType("entitySeatBooking.SeatBooking"), suite.ginContext).Return(nil, nil)
	suite.mockSql.ExpectCommit()

	suite.service.Update(suite.ctx, notificationPayload, suite.ginContext)
}

func (suite *PaymentServiceTestSuite) TestUpdate_PendingStatus() {
	notificationPayload := map[string]interface{}{
		"order_id":           "some-id",
		"transaction_status": "pending",
	}
	expectedResult := entity.Payment{
		ID:                     "some-id",
		UserID:                 "user-id",
		SeatDetailForBookingID: "seat-id",
		Status:                 "pending",
	}

	suite.mockSql.ExpectBegin()
	suite.mockRepo.On("FindByID", suite.ctx, mock.Anything, "some-id", suite.ginContext).Return(expectedResult, nil)
	suite.mockRepo.On("Update", suite.ctx, mock.Anything, mock.AnythingOfType("entity.Payment"), suite.ginContext).Return(expectedResult, nil)
	suite.mockRepoSeatBooking.On("Update", suite.ctx, mock.Anything, mock.AnythingOfType("entitySeatBooking.SeatBooking"), suite.ginContext).Return(nil, nil)
	suite.mockSql.ExpectCommit()

	suite.service.Update(suite.ctx, notificationPayload, suite.ginContext)
}

func (suite *PaymentServiceTestSuite) TestUpdate_NotFound() {
	notificationPayload := map[string]interface{}{
		"order_id":           "some-id",
		"transaction_status": "settlement",
	}
	notFoundError := errors.New("Not Found Error")

	suite.mockSql.ExpectBegin()
	suite.mockRepo.On("FindByID", suite.ctx, mock.Anything, "some-id", suite.ginContext).Return(entity.Payment{}, notFoundError)
	suite.mockSql.ExpectRollback()

	suite.service.Update(suite.ctx, notificationPayload, suite.ginContext)
}
