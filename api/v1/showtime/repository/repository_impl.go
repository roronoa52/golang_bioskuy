package repository

import (
	"bioskuy/api/v1/showtime/entity"
	entityStudio "bioskuy/api/v1/studio/entity"
	"bioskuy/exception"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

type showtimeRepository struct {
}

func NewShowtimeRepository() ShowtimeRepository {
	return &showtimeRepository{}
}

func (r *showtimeRepository) Save(ctx context.Context, tx *sql.Tx, showtime entity.Showtime, c *gin.Context) (entity.Showtime, error){

	fmt.Println(showtime.ShowStart)
	fmt.Println(showtime.ShowEnd)

	query := "INSERT INTO showtimes (movie_id, studio_id, show_start, show_end) VALUES ($1, $2, $3, $4) RETURNING id"

	err := tx.QueryRowContext(ctx, query, showtime.MovieID, showtime.StudioID, showtime.ShowStart, showtime.ShowEnd).Scan(&showtime.ID)
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return showtime, err
	}

	return showtime, nil
}

func (r *showtimeRepository) FindByID(ctx context.Context, tx *sql.Tx, id string, c *gin.Context) (entity.Showtime, error){

	query := `SELECT s.id, s.studio_id, s.movie_id, s.show_start, s.show_end, st.name as studio_name, m.title as movie_title, m.description as movie_description, m.price as movie_price, m.duration as movie_duration, m.status as movie_status
    FROM showtimes s
    JOIN studios st ON s.studio_id = st.id
    JOIN movies m ON s.movie_id = m.id
    WHERE s.id = $1
    `
	showtime := entity.Showtime{}
	rows, err := tx.QueryContext(ctx, query, id)
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return  showtime, err
	}
	defer rows.Close()

	if rows.Next(){
		err := rows.Scan(
			&showtime.ID, &showtime.StudioID, &showtime.MovieID, &showtime.ShowStart, &showtime.ShowEnd,
			&showtime.StudioName, &showtime.MovieTitle, &showtime.MovieDescription, &showtime.MoviePrice,
			&showtime.MovieDuration, &showtime.MovieStatus,
		)
		if err != nil {
			c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
			return  showtime, err
		}

		return showtime, nil
	}else{
		return showtime, errors.New("showtime not found")
	}
}

func (r *showtimeRepository) FindAll(ctx context.Context, tx *sql.Tx, c *gin.Context) ([]entity.Showtime, error){

	query := `
    SELECT s.id, s.studio_id, s.movie_id, s.show_start, s.show_end, st.name as studio_name, m.title as movie_title, m.description as movie_description, m.price as movie_price, m.duration as movie_duration, m.status as movie_status
    FROM showtimes s
    JOIN studios st ON s.studio_id = st.id
    JOIN movies m ON s.movie_id = m.id
    `

	showtimes := []entity.Showtime{}
	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return  showtimes, err
	}
	defer rows.Close()

	for rows.Next() {
		showtime := entity.Showtime{}
		if err := rows.Scan(
            &showtime.ID, &showtime.StudioID, &showtime.MovieID, &showtime.ShowStart, &showtime.ShowEnd,
            &showtime.StudioName, &showtime.MovieTitle, &showtime.MovieDescription, &showtime.MoviePrice,
            &showtime.MovieDuration, &showtime.MovieStatus,
        ); err != nil {
			return nil, err
		}
		showtimes = append(showtimes, showtime)
	}
	return showtimes, nil
}

func (r *showtimeRepository) Delete(ctx context.Context, tx *sql.Tx, id string, c *gin.Context) error{
	query := `DELETE FROM showtimes WHERE id = $1`

	_, err := tx.ExecContext(ctx, query, id)

	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return  err
	}

	return nil
}

func (r *showtimeRepository) FindConflictingShowtimes(ctx context.Context, tx *sql.Tx, studio entityStudio.Studio, showtime entity.Showtime, c *gin.Context) error {
    query := `SELECT s.id
              FROM showtimes s
              WHERE s.studio_id = $1
              AND (
                  (s.show_start >= $2 AND s.show_start <= $3) OR
                  (s.show_end >= $2 AND s.show_end <= $3) OR
                  (s.show_start <= $2 AND s.show_end >= $2) OR
                  (s.show_start <= $3 AND s.show_end >= $3)
              )`
    
    rows, err := tx.QueryContext(ctx, query, studio.ID, showtime.ShowStart, showtime.ShowEnd)
    if err != nil {
        c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
        return err
    }
    defer rows.Close()

    if rows.Next() {
        return errors.New("conflicting showtimes found")
    }
    
    return nil
}


