package service

import (
	"bioskuy/api/v1/genre/dto"
	"bioskuy/api/v1/genre/entity"

	"github.com/google/uuid"
)

type GenreService interface {
	CreateGenre(genre entity.Genre) (entity.Genre, error)
	GetAll(page int, size int) ([]entity.Genre, dto.Paging, error)
	GetGenreByID(id uuid.UUID) (entity.Genre, error)
	UpdateGenre(genre entity.Genre) (entity.Genre, error)
	DeleteGenre(id uuid.UUID) (entity.Genre, error)
}
