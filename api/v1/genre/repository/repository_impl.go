package repository

import (
	"bioskuy/api/v1/genre/dto"
	"bioskuy/api/v1/genre/entity"
	"bioskuy/exception"
	"database/sql"
	"math"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type genreRepositoryImpl struct {
	DB *sql.DB
}

func NewGenreRepository(DB *sql.DB) GenreRepository {
	return &genreRepositoryImpl{
		DB: DB,
	}
}

func (g *genreRepositoryImpl) GetAll(page int, size int) ([]entity.Genre, dto.Paging, error) {
	var listData []entity.Genre

	skip := (page - 1) * size

	rows, err := g.DB.Query(`SELECT id, name FROM genres LIMIT $1 OFFSET $2`, size, skip)
	if err != nil {
		return nil, dto.Paging{}, err
	}
	defer rows.Close()

	totalRows := 0
	err = g.DB.QueryRow(`SELECT COUNT(*) FROM genres`).Scan(&totalRows)
	if err != nil {
		return nil, dto.Paging{}, err
	}

	for rows.Next() {
		var genre entity.Genre

		err := rows.Scan(&genre.ID, &genre.Name)
		if err != nil {
			return nil, dto.Paging{}, err
		}
		listData = append(listData, genre)
	}

	paging := dto.Paging{
		Page:       page,
		Size:       size,
		TotalRows:  totalRows,
		TotalPages: int(math.Ceil(float64(totalRows) / float64(size))),
	}
	return listData, paging, nil
}

func (r *genreRepositoryImpl) Create(genre entity.Genre) (entity.Genre, error) {

	var C *gin.Context
	genre.ID = uuid.New()
	err := r.DB.QueryRow(
		"INSERT INTO genres (id, name) VALUES ($1, $2) RETURNING id, name",
		genre.ID, genre.Name,
	).Scan(&genre.ID, &genre.Name)
	if err != nil {
		C.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return entity.Genre{}, err
	}
	return genre, nil
}

func (r *genreRepositoryImpl) GetByID(id uuid.UUID) (entity.Genre, error) {

	var C *gin.Context
	var genre entity.Genre
	err := r.DB.QueryRow(
		"SELECT id, name FROM genres WHERE id = $1",
		id,
	).Scan(&genre.ID, &genre.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			C.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
			return entity.Genre{}, nil
		}
		return entity.Genre{}, err
	}
	return genre, nil
}

func (r *genreRepositoryImpl) Update(genre entity.Genre) (entity.Genre, error) {
	err := r.DB.QueryRow(
		"UPDATE genres SET name = $1 WHERE id = $2 RETURNING id, name",
		genre.Name, genre.ID,
	).Scan(&genre.ID, &genre.Name)
	if err != nil {
		return entity.Genre{}, err
	}
	return genre, nil
}

func (r *genreRepositoryImpl) Delete(id uuid.UUID) (entity.Genre, error) {

	var C *gin.Context
	var genre entity.Genre
	err := r.DB.QueryRow(
		"DELETE FROM genres WHERE id = $1 RETURNING id, name",
		id,
	).Scan(&genre.ID, &genre.Name)
	if err != nil {
		C.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return entity.Genre{}, err
	}
	return genre, nil
}
