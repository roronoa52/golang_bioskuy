package service

import (
	"bioskuy/api/v1/seat/dto"
	"bioskuy/api/v1/seat/repository"
	"bioskuy/exception"
	"bioskuy/helper"
	"context"
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type seatService struct {
	Repo     repository.SeatRepository
	Validate *validator.Validate
	DB *sql.DB
}


func NewSeatervice(repo repository.SeatRepository, validate *validator.Validate, DB *sql.DB) SeatService {
	return &seatService{Repo: repo, Validate: validate, DB: DB}
}

func (s *seatService) FindByID(ctx context.Context, id string, c *gin.Context) (dto.SeatResponse, error){

	SeatResponse := dto.SeatResponse{}

	tx, err := s.DB.Begin()
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return  SeatResponse, err
	}
	defer helper.CommitAndRollback(tx, c)

	result, err := s.Repo.FindByID(ctx, tx, id, c)
	if err != nil {
		c.Error(exception.NotFoundError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return  SeatResponse, err
	}

	SeatResponse.ID = result.ID
	SeatResponse.Name = result.Name
	SeatResponse.IsAvailable = result.IsAvailable
	SeatResponse.StudioID = result.StudioID

	return SeatResponse, nil
}

func (s *seatService) FindAll(ctx context.Context, id string, c *gin.Context) ([]dto.SeatResponse, error){
	StudioResponses := []dto.SeatResponse{}

	tx, err := s.DB.Begin()
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return  StudioResponses, err
	}
	defer helper.CommitAndRollback(tx, c)

	result, err := s.Repo.FindAll(ctx, id, tx, c)
	if err != nil {
		c.Error(exception.NotFoundError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return  StudioResponses, err
	}

	for _, seat := range result {
		SeatResponse := dto.SeatResponse{}

		SeatResponse.ID = seat.ID
		SeatResponse.Name = seat.Name
		SeatResponse.IsAvailable = seat.IsAvailable
		SeatResponse.StudioID = seat.StudioID

		StudioResponses = append(StudioResponses, SeatResponse)
		
	}

	return StudioResponses, nil
}
