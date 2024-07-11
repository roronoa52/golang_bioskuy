package repository

import (
	"bioskuy/api/v1/genre/dto"
	"bioskuy/api/v1/movies/entity"

)

type MovieRepository interface {
	GetAll(page int, size int) ([]entity.Movie, dto.Paging, error)
	Create(movie entity.Movie) (entity.Movie, error)
	GetByID(id string) (entity.Movie, error)
	Update(movie entity.Movie) (entity.Movie, error)
	Delete(id string) (entity.Movie, error)
}
