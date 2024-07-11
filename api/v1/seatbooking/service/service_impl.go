package service

import (
	RepoSeat "bioskuy/api/v1/seat/repository"
	"bioskuy/api/v1/seatbooking/dto"
	"bioskuy/api/v1/seatbooking/entity"
	"bioskuy/api/v1/seatbooking/repository"
	RepoShowtime "bioskuy/api/v1/showtime/repository"
	"bioskuy/exception"
	"bioskuy/helper"
	"context"
	"database/sql"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type seatbookingServiceImpl struct {
	Repo         repository.SeatBookingRepository
	RepoShowtime RepoShowtime.ShowtimeRepository
	RepoSeat     RepoSeat.SeatRepository
	Validate     *validator.Validate
	DB           *sql.DB
	Mutex        sync.Mutex
}

func NewSeatBookingService(repo repository.SeatBookingRepository, RepoShowtime RepoShowtime.ShowtimeRepository,
	RepoSeat RepoSeat.SeatRepository, validate *validator.Validate, DB *sql.DB) SeatBookingService {
	return &seatbookingServiceImpl{
		Repo:         repo,
		RepoShowtime: RepoShowtime,
		RepoSeat:     RepoSeat,
		Validate:     validate,
		DB:           DB,
	}
}

func (s *seatbookingServiceImpl) Create(ctx context.Context, request dto.SeatBookingRequest, userid string, c *gin.Context) (dto.CreateSeatBookingResponse, error) {
	var SeatBookingResponse = dto.CreateSeatBookingResponse{}
	var SeatBookingRequest = dto.CreateSeatBookingRequest{}

	SeatBookingRequest.UserID = userid
	SeatBookingRequest.ShowtimeID = request.ShowtimeID

	err := s.Validate.Struct(request)
	if err != nil {
		c.Error(exception.ValidationError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return SeatBookingResponse, err
	}

	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	tx, err := s.DB.Begin()
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return SeatBookingResponse, err
	}
	defer helper.CommitAndRollback(tx, c)

	_, err = s.RepoShowtime.FindByID(ctx, tx, SeatBookingRequest.ShowtimeID, c)
	if err != nil {
		c.Error(exception.NotFoundError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return SeatBookingResponse, err
	}

	seat, err := s.RepoSeat.FindByIDWithNotAvailable(ctx, tx, request.SeatID, c)
	if err != nil {
		c.Error(exception.NotFoundError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return SeatBookingResponse, err
	}

	seatbooking := entity.SeatBooking{
		UserID:            userid,
		ShowtimeID:        SeatBookingRequest.ShowtimeID,
		SeatBookingStatus: SeatBookingRequest.Status,
		SeatID:            request.SeatID,
	}

	result, err := s.Repo.Save(ctx, tx, seatbooking, c)
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return SeatBookingResponse, err
	}

	seat.IsAvailable = false

	_, err = s.RepoSeat.Update(ctx, tx, seat, c)
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return SeatBookingResponse, err
	}

	SeatBookingResponse.ID = result.SeatDetailForBookingID
	SeatBookingResponse.SeatBookingID = result.ID
	SeatBookingResponse.SeatID = result.SeatID

	return SeatBookingResponse, nil
}

func (s *seatbookingServiceImpl) FindByID(ctx context.Context, id string, c *gin.Context) (dto.SeatBookingResponse, error) {
	seatBookingResponse := dto.SeatBookingResponse{}

	tx, err := s.DB.Begin()
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return seatBookingResponse, err
	}
	defer helper.CommitAndRollback(tx, c)

	result, err := s.Repo.FindByID(ctx, tx, id, c)
	if err != nil {
		c.Error(exception.NotFoundError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return seatBookingResponse, err
	}

	seatBookingResponse.ID = result.ID
	seatBookingResponse.SeatBookingStatus = result.SeatBookingStatus
	seatBookingResponse.UserID = result.UserID
	seatBookingResponse.ShowtimeID = result.ShowtimeID
	seatBookingResponse.ShowStart = result.ShowStart
	seatBookingResponse.ShowEnd = result.ShowEnd
	seatBookingResponse.StudioID = result.StudioID
	seatBookingResponse.StudioName = result.StudioName
	seatBookingResponse.MovieID = result.MovieID
	seatBookingResponse.MovieTitle = result.MovieTitle
	seatBookingResponse.MovieDescription = result.MovieDescription
	seatBookingResponse.MoviePrice = result.MoviePrice
	seatBookingResponse.MovieDuration = result.MovieDuration
	seatBookingResponse.MovieStatus = result.MovieStatus
	seatBookingResponse.SeatID = result.SeatID
	seatBookingResponse.SeatDetailForBookingID = result.SeatDetailForBookingID
	seatBookingResponse.SeatName = result.SeatName
	seatBookingResponse.SeatIsAvailable = result.SeatIsAvailable
	

	return seatBookingResponse, nil
}

func (s *seatbookingServiceImpl) FindAll(ctx context.Context, c *gin.Context) ([]dto.SeatBookingResponse, error) {
	seatBookingResponses := []dto.SeatBookingResponse{}

	tx, err := s.DB.Begin()
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return seatBookingResponses, err
	}
	defer helper.CommitAndRollback(tx, c)

	results, err := s.Repo.FindAll(ctx, tx, c)
	if err != nil {
		c.Error(exception.NotFoundError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return seatBookingResponses, err
	}

	for _, result := range results {
		seatBookingResponse := dto.SeatBookingResponse{
			ID:                     result.ID,
			SeatBookingStatus:      result.SeatBookingStatus,
			UserID:                 result.UserID,
			ShowtimeID:             result.ShowtimeID,
			ShowStart:              result.ShowStart,
			ShowEnd:                result.ShowEnd,
			StudioID:               result.StudioID,
			StudioName:             result.StudioName,
			MovieID:                result.MovieID,
			MovieTitle:             result.MovieTitle,
			MovieDescription:       result.MovieDescription,
			MoviePrice:             result.MoviePrice,
			MovieDuration:          result.MovieDuration,
			MovieStatus:            result.MovieStatus,
			SeatID:                 result.SeatID,
			SeatDetailForBookingID: result.SeatDetailForBookingID,
			SeatName:               result.SeatName,
			SeatIsAvailable:        result.SeatIsAvailable,
		}
		seatBookingResponses = append(seatBookingResponses, seatBookingResponse)
	}

	return seatBookingResponses, nil
}

func (s *seatbookingServiceImpl) Delete(ctx context.Context, id string, c *gin.Context) error {
	tx, err := s.DB.Begin()
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return err
	}
	defer helper.CommitAndRollback(tx, c)

	err = s.Repo.Delete(ctx, tx, id, c)
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return err
	}

	return nil
}
