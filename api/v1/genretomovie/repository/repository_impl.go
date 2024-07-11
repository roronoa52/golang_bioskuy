package repository

import (
	"bioskuy/api/v1/genretomovie/entity"
	"bioskuy/exception"
	"context"
	"database/sql"
	"errors"

	"github.com/gin-gonic/gin"
)

type genretomovieRepository struct {
}

func NewGenreToMovieRepository() GenreToMovieRepository {
	return &genretomovieRepository{}
}

func (r *genretomovieRepository) Save(ctx context.Context, tx *sql.Tx, GenreToMovie entity.GenreToMovie, c *gin.Context) (entity.GenreToMovie, error) {
	query := "INSERT INTO genre_to_movies (genre_id, movie_id) VALUES ($1, $2) RETURNING id"

	err := tx.QueryRowContext(ctx, query, GenreToMovie.GenreID, GenreToMovie.MovieID).Scan(&GenreToMovie.ID)
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return GenreToMovie, err
	}

	return GenreToMovie, nil
}

func (r *genretomovieRepository) FindByID(ctx context.Context, tx *sql.Tx, id string, c *gin.Context) (entity.GenreToMovie, error){
	query := `SELECT gtm.id, gtm.genre_id, gtm.movie_id, g.name as genre_name, m.title as movie_title, m.description as movie_description, 
                  m.price as movie_price, m.duration as movie_duration, m.status as movie_status
	          FROM genre_to_movies gtm 
	          JOIN genres g ON gtm.genre_id = g.id 
	          JOIN movies m ON gtm.movie_id = m.id 
	          WHERE gtm.id = $1`

	genretomovie := entity.GenreToMovie{}
	rows, err := tx.QueryContext(ctx, query, id)
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return genretomovie, err
	}
	defer rows.Close()

	if rows.Next(){
		err := rows.Scan(&genretomovie.ID, &genretomovie.GenreID, &genretomovie.MovieID, &genretomovie.GenreName, &genretomovie.MovieTitle,
		                &genretomovie.MovieDescription, &genretomovie.MoviePrice, &genretomovie.MovieDuration, &genretomovie.MovieStatus)
		if err != nil {
			c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
			return genretomovie, err
		}

		return genretomovie, nil
	} else {
		return genretomovie, errors.New("genre to movie not found")
	}
}

func (r *genretomovieRepository) FindAll(ctx context.Context, tx *sql.Tx, c *gin.Context) ([]entity.GenreToMovie, error){
	query := `SELECT gtm.id, gtm.genre_id, gtm.movie_id, g.name as genre_name, m.title as movie_title, m.description as movie_description, 
                  m.price as movie_price, m.duration as movie_duration, m.status as movie_status
	          FROM genre_to_movies gtm 
	          JOIN genres g ON gtm.genre_id = g.id 
	          JOIN movies m ON gtm.movie_id = m.id`

	genretomovies := []entity.GenreToMovie{}
	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return genretomovies, err
	}
	defer rows.Close()

	for rows.Next() {
		genretomovie := entity.GenreToMovie{}
		if err := rows.Scan(&genretomovie.ID, &genretomovie.GenreID, &genretomovie.MovieID, &genretomovie.GenreName, 
			&genretomovie.MovieTitle, &genretomovie.MovieDescription, &genretomovie.MoviePrice, &genretomovie.MovieDuration, &genretomovie.MovieStatus); err != nil {
			return nil, err
		}
		genretomovies = append(genretomovies, genretomovie)
	}
	return genretomovies, nil
}

func (r *genretomovieRepository) Delete(ctx context.Context, tx *sql.Tx, id string, c *gin.Context) error {
	query := `DELETE FROM genre_to_movies WHERE id = $1`

	_, err := tx.ExecContext(ctx, query, id)

	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return err
	}

	return nil
}
