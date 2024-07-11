package service

import (
	entitySeat "bioskuy/api/v1/seat/entity"
	repoSeat "bioskuy/api/v1/seat/repository"
	"bioskuy/api/v1/studio/dto"
	"bioskuy/api/v1/studio/entity"
	"bioskuy/api/v1/studio/repository"
	"bioskuy/exception"
	"bioskuy/helper"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type studioService struct {
	RepoStudio     repository.StudioRepository
	RepoSeat     repoSeat.SeatRepository
	Validate *validator.Validate
	DB *sql.DB
}


func NewStudioService(RepoStudio repository.StudioRepository, validate *validator.Validate, DB *sql.DB, RepoSeat repoSeat.SeatRepository) StudioService {
	return &studioService{RepoStudio: RepoStudio, Validate: validate, DB: DB, RepoSeat: RepoSeat}
}

func (s *studioService) Create(ctx context.Context, request dto.CreateStudioRequest, c *gin.Context) (dto.StudioResponse, error) {
    var StudioResponse = dto.StudioResponse{}

    err := s.Validate.Struct(request)
    if err != nil {
        c.Error(exception.ValidationError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
        return StudioResponse, err
    }

    if request.MaxRowSeat == 0 {
        err := errors.New("maxRowSeat tidak boleh nol")
        c.Error(exception.ValidationError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
        return StudioResponse, err
    }

    numRows := (request.Capacity + request.MaxRowSeat - 1) / request.MaxRowSeat
    if numRows > 26 {
        err := errors.New("jumlah baris melebihi batas maksimal (26)")
        c.Error(exception.ValidationError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
        return StudioResponse, err
    }

    tx, err := s.DB.Begin()
    if err != nil {
        c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
        return StudioResponse, err
    }
    defer helper.CommitAndRollback(tx, c)

    studio := entity.Studio{
        Name:     request.Name,
        Capacity: request.Capacity,
    }

    result, err := s.RepoStudio.Save(ctx, tx, studio, c)
    if err != nil {
        c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
        return StudioResponse, err
    }

    for row := 0; row < numRows; row++ {
        for seatNum := 1; seatNum <= request.MaxRowSeat; seatNum++ {
            actualSeatNum := row*request.MaxRowSeat + seatNum
            if actualSeatNum > request.Capacity {
                break
            }
            seatName := string(rune('A'+row)) + "-" + fmt.Sprintf("%d", seatNum)

            seat := entitySeat.Seat{
                Name:       seatName,
                IsAvailable: true,
                StudioID:   result.ID,
            }

            _, err := s.RepoSeat.Save(ctx, tx, seat, c)
            if err != nil {
                c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
                return StudioResponse, err
            }
        }
    }

    StudioResponse.ID = result.ID
    StudioResponse.Name = result.Name
    StudioResponse.Capacity = result.Capacity

    return StudioResponse, nil
}

func (s *studioService) FindByID(ctx context.Context, id string, c *gin.Context) (dto.StudioResponse, error){

	StudioResponse := dto.StudioResponse{}

	tx, err := s.DB.Begin()
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return  StudioResponse, err
	}
	defer helper.CommitAndRollback(tx, c)

	result, err := s.RepoStudio.FindByID(ctx, tx, id, c)
	if err != nil {
		c.Error(exception.NotFoundError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return  StudioResponse, err
	}

	StudioResponse.ID = result.ID
	StudioResponse.Name = result.Name
	StudioResponse.Capacity = result.Capacity

	return StudioResponse, nil
}

func (s *studioService) FindAll(ctx context.Context, c *gin.Context) ([]dto.StudioResponse, error){
	StudioResponses := []dto.StudioResponse{}

	tx, err := s.DB.Begin()
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return  StudioResponses, err
	}
	defer helper.CommitAndRollback(tx, c)

	result, err := s.RepoStudio.FindAll(ctx, tx, c)
	if err != nil {
		c.Error(exception.NotFoundError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return  StudioResponses, err
	}

	for _, studio := range result {
		StudioResponse := dto.StudioResponse{}

		StudioResponse.ID = studio.ID
		StudioResponse.Name = studio.Name
		StudioResponse.Capacity = studio.Capacity

		StudioResponses = append(StudioResponses, StudioResponse)
		
	}

	return StudioResponses, nil
}

func (s *studioService) Update(ctx context.Context, request dto.UpdateStudioRequest, c *gin.Context) (dto.StudioResponse, error){
	StudioResponse := dto.StudioResponse{}
    var studio entity.Studio

    err := s.Validate.Struct(request)
    if err != nil {
		c.Error(exception.ValidationError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return  StudioResponse, err
	}

    resultStudio, err := s.FindByID(ctx, request.ID, c)
    if err != nil {
		c.Error(exception.NotFoundError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return  StudioResponse, err
	}

    tx, err := s.DB.Begin()
    if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return  StudioResponse, err
	}
    defer helper.CommitAndRollback(tx, c)

    studio.ID = resultStudio.ID
	studio.Name = resultStudio.Name
	studio.Capacity = resultStudio.Capacity

	if request.Name != "" {
		studio.Name = request.Name
	}

    result, err := s.RepoStudio.Update(ctx, tx, studio, c)
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return  StudioResponse, err
	}

    StudioResponse.ID = result.ID
	StudioResponse.Name = resultStudio.Name
	StudioResponse.Capacity = resultStudio.Capacity

    return StudioResponse, nil
}

func (s *studioService) Delete(ctx context.Context, id string, c *gin.Context) error {
    tx, err := s.DB.Begin()
    if err != nil {
        c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
        return err
    }
    defer helper.CommitAndRollback(tx, c)

    studio, err := s.RepoStudio.FindByID(ctx, tx, id, c)
    if err != nil {
        c.Error(exception.NotFoundError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
        return err
    }

    err = s.RepoSeat.Delete(ctx, tx, studio.ID, c)
    if err != nil {
        c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
        return err
    }

    err = s.RepoStudio.Delete(ctx, tx, studio.ID, c)
    if err != nil {
        c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
        return err
    }

    return nil
}
