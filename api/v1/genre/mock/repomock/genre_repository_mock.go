package repomock

import (
	"bioskuy/api/v1/genre/dto"
	"bioskuy/api/v1/genre/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockGenreRepository struct {
	mock.Mock
}

func (m *MockGenreRepository) Create(genre entity.Genre) (entity.Genre, error) {
	args := m.Called(genre)
	return args.Get(0).(entity.Genre), args.Error(1)
}

func (m *MockGenreRepository) GetByID(id uuid.UUID) (entity.Genre, error) {
	args := m.Called(id)
	return args.Get(0).(entity.Genre), args.Error(1)
}

func (m *MockGenreRepository) GetAll(page int, size int) ([]entity.Genre, dto.Paging, error) {
	args := m.Called(page, size)
	return args.Get(0).([]entity.Genre), args.Get(1).(dto.Paging), args.Error(2)
}

func (m *MockGenreRepository) Update(genre entity.Genre) (entity.Genre, error) {
	args := m.Called(genre)
	return args.Get(0).(entity.Genre), args.Error(1)
}

func (m *MockGenreRepository) Delete(id uuid.UUID) (entity.Genre, error) {
	args := m.Called(id)
	return args.Get(0).(entity.Genre), args.Error(1)
}
