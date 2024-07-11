package service

import (
	"bioskuy/api/v1/genre/dto"
	"bioskuy/api/v1/genre/entity"
	"bioskuy/api/v1/genre/repository"

	"github.com/google/uuid"
)

type genreServiceImpl struct {
	repo repository.GenreRepository
}

func NewGenreService(repo repository.GenreRepository) GenreService {
	return &genreServiceImpl{repo}
}

func (s *genreServiceImpl) GetAll(page int, size int) ([]entity.Genre, dto.Paging, error) {
	return s.repo.GetAll(page, size)
}

func (s *genreServiceImpl) CreateGenre(genre entity.Genre) (entity.Genre, error) {
	return s.repo.Create(genre)
}

func (s *genreServiceImpl) GetGenreByID(id uuid.UUID) (entity.Genre, error) {
	return s.repo.GetByID(id)
}

func (s *genreServiceImpl) UpdateGenre(genre entity.Genre) (entity.Genre, error) {
	return s.repo.Update(genre)
}
func (s *genreServiceImpl) DeleteGenre(id uuid.UUID) (entity.Genre, error) {
	return s.repo.Delete(id)
}
