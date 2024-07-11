package service

import (
	RepoMovie "bioskuy/api/v1/movies/repository"
	"bioskuy/api/v1/showtime/dto"
	"bioskuy/api/v1/showtime/entity"
	"bioskuy/api/v1/showtime/repository"
	RepoStudio "bioskuy/api/v1/studio/repository"
	"bioskuy/exception"
	"bioskuy/helper"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type showtimesServiceImpl struct {
	Repo repository.ShowtimeRepository
	RepoMovie RepoMovie.MovieRepository
	RepoStudio RepoStudio.StudioRepository
	Validate *validator.Validate
	DB *sql.DB
}

func NewGenreToMovieService(repo repository.ShowtimeRepository, RepoMovie RepoMovie.MovieRepository, RepoStudio RepoStudio.StudioRepository, validate *validator.Validate, DB *sql.DB) ShowtimeService {
	return &showtimesServiceImpl{
		Repo: repo,
		RepoMovie: RepoMovie,
		RepoStudio: RepoStudio,
		Validate: validate,
		DB: DB,
	}
}

func (s *showtimesServiceImpl) Create(ctx context.Context, request dto.ShowtimeRequest, c *gin.Context) (dto.CreateShowtimesResponseDTO, error) {
	var ShowtimeResponse = dto.CreateShowtimesResponseDTO{}
	var ShowtimeRequest = dto.CreateShowtimeDTO{}
	
	ShowtimeRequest.MovieID = request.MovieID
	ShowtimeRequest.StudioID = request.StudioID
	ShowtimeRequest.ShowStart = helper.StringToDate(request.ShowStart, c)

	err := s.Validate.Struct(request)
	if err != nil {
		c.Error(exception.ValidationError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return ShowtimeResponse, err
	}

	tx, err := s.DB.Begin()
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return ShowtimeResponse, err
	}
	defer helper.CommitAndRollback(tx, c)

	movie , err := s.RepoMovie.GetByID(request.MovieID)
	if err != nil {
		c.Error(exception.NotFoundError{Message: "Movie Not Found"}).SetType(gin.ErrorTypePublic)
		return ShowtimeResponse, err
	}

	studio , err := s.RepoStudio.FindByID(ctx, tx, request.StudioID, c)
	if err != nil {
		c.Error(exception.NotFoundError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return ShowtimeResponse, err
	}

	duration := time.Duration(movie.Duration) * time.Hour

	ShowtimeRequest.ShowEnd = ShowtimeRequest.ShowStart.Add(duration)

	showtime := entity.Showtime{
		StudioID: ShowtimeRequest.StudioID ,
		MovieID: ShowtimeRequest.MovieID,
		ShowStart: ShowtimeRequest.ShowStart,
		ShowEnd: ShowtimeRequest.ShowEnd ,
	}

	fmt.Println(showtime)

	err = s.Repo.FindConflictingShowtimes(ctx, tx, studio, showtime, c)
	if err != nil {
		c.Error(exception.ValidationError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return ShowtimeResponse, err
	}

	result, err := s.Repo.Save(ctx, tx, showtime, c)
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return ShowtimeResponse, err
	}

	ShowtimeResponse.MovieID = result.MovieID
	ShowtimeResponse.StudioID = result.StudioID
	ShowtimeResponse.ShowStart = result.ShowStart
	ShowtimeResponse.ShowEnd = result.ShowEnd

	return ShowtimeResponse, nil
}

func (s *showtimesServiceImpl) FindByID(ctx context.Context, id string, c *gin.Context) (dto.ShowtimesResponse, error){
	ShowtimeResponse := dto.ShowtimesResponse{}

	tx, err := s.DB.Begin()
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return  ShowtimeResponse, err
	}
	defer helper.CommitAndRollback(tx, c)

	result, err := s.Repo.FindByID(ctx, tx, id, c)
	if err != nil {
		c.Error(exception.NotFoundError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return  ShowtimeResponse, err
	}

	ShowtimeResponse.ID = result.ID
	ShowtimeResponse.StudioID = result.StudioID
	ShowtimeResponse.MovieID = result.MovieID
	ShowtimeResponse.StudioName = result.StudioName
	ShowtimeResponse.MovieTitle = result.MovieTitle
	ShowtimeResponse.MovieDescription = result.MovieDescription
	ShowtimeResponse.MoviePrice = result.MoviePrice
	ShowtimeResponse.MovieDuration = result.MovieDuration
	ShowtimeResponse.MovieStatus = result.MovieStatus

	return ShowtimeResponse, nil
}

func (s *showtimesServiceImpl) FindAll(ctx context.Context, c *gin.Context) ([]dto.ShowtimesResponse, error){
	ShowtimeResponses := []dto.ShowtimesResponse{}

	tx, err := s.DB.Begin()
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return  ShowtimeResponses, err
	}
	defer helper.CommitAndRollback(tx, c)

	results, err := s.Repo.FindAll(ctx, tx, c)
	if err != nil {
		c.Error(exception.NotFoundError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return  ShowtimeResponses, err
	}

	for _, result := range results {
		ShowtimeResponse := dto.ShowtimesResponse{
			ID: result.ID,
			StudioID: result.StudioID,
			MovieID: result.MovieID,
			StudioName: result.StudioName,
			MovieTitle: result.MovieTitle,
			MovieDescription: result.MovieDescription,
			MoviePrice: result.MoviePrice,
			MovieDuration: result.MovieDuration,
			MovieStatus: result.MovieStatus,
		}
		ShowtimeResponses = append(ShowtimeResponses, ShowtimeResponse)
	}

	return ShowtimeResponses, nil
}

func (s *showtimesServiceImpl) Delete(ctx context.Context, id string, c *gin.Context) error{
	showtime := entity.Showtime{}

	tx, err := s.DB.Begin()
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return err
	}
	defer helper.CommitAndRollback(tx, c)

	resultUser, err := s.Repo.FindByID(ctx, tx, id, c)
	if err != nil {
		c.Error(exception.NotFoundError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return err
	}

	showtime.ID = resultUser.ID

	err = s.Repo.Delete(ctx, tx, id, c)
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return err
	}

	return nil
}

