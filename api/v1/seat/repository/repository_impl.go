package repository

import (
	"bioskuy/api/v1/seat/entity"
	"bioskuy/exception"
	"context"
	"database/sql"
	"errors"

	"github.com/gin-gonic/gin"
)

type seatRepository struct {
}

func NewSeatRepository() SeatRepository {
	return &seatRepository{}
}

func (r *seatRepository) Save(ctx context.Context, tx *sql.Tx, seat entity.Seat, c *gin.Context) (entity.Seat, error) {
	query := "INSERT INTO seats (seat_name, isAvailable, studio_id) VALUES ($1, $2, $3) RETURNING id"

	err := tx.QueryRowContext(ctx, query, seat.Name, seat.IsAvailable, seat.StudioID).Scan(&seat.ID)
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return seat, err
	}

	return seat, nil
}

func (r *seatRepository) FindByID(ctx context.Context, tx *sql.Tx, id string, c *gin.Context) (entity.Seat, error) {

	query := `SELECT id, seat_name, isAvailable, studio_id FROM seats WHERE id = $1`

	seat := entity.Seat{}
	rows, err := tx.QueryContext(ctx, query, id)
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return seat, err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&seat.ID, &seat.Name, &seat.IsAvailable, &seat.StudioID)
		if err != nil {
			c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
			return seat, err
		}

		return seat, nil
	} else {
		return seat, errors.New("seat not found")
	}
}

func (r *seatRepository) FindByIDWithNotAvailable(ctx context.Context, tx *sql.Tx, id string, c *gin.Context) (entity.Seat, error){

	query := `SELECT id, seat_name, isAvailable, studio_id FROM seats WHERE id = $1 AND isAvailable = true`
	
	seat := entity.Seat{}
	rows, err := tx.QueryContext(ctx, query, id)
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return  seat, err
	}
	defer rows.Close()

	if rows.Next(){
		err := rows.Scan(&seat.ID, &seat.Name, &seat.IsAvailable, &seat.StudioID)
		if err != nil {
			c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
			return  seat, err
		}

		return seat, nil
	}else{
		return seat, errors.New("seat not found or have been taken")
	}
}

func (r *seatRepository) Update(ctx context.Context, tx *sql.Tx, seat entity.Seat, c *gin.Context) (entity.Seat, error){

	query := `UPDATE seats SET isAvailable = $1 WHERE id = $2`

	_, err := tx.ExecContext(ctx, query, seat.IsAvailable, seat.ID)

	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return  seat, err
	}

	return seat, nil
}

func (r *seatRepository) FindAll(ctx context.Context, id string, tx *sql.Tx, c *gin.Context) ([]entity.Seat, error){

	query := `SELECT  id, seat_name, isAvailable, studio_id FROM seats WHERE studio_id = $1`

	seats := []entity.Seat{}
	rows, err := tx.QueryContext(ctx, query, id)
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return seats, err
	}
	defer rows.Close()

	for rows.Next() {
		seat := entity.Seat{}
		if err := rows.Scan(&seat.ID, &seat.Name, &seat.IsAvailable, &seat.StudioID); err != nil {
			return nil, err
		}
		seats = append(seats, seat)
	}
	return seats, nil
}

func (r *seatRepository) Delete(ctx context.Context, tx *sql.Tx, id string, c *gin.Context) error {
	query := `DELETE FROM seats WHERE studio_id = $1`

	_, err := tx.ExecContext(ctx, query, id)

	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return err
	}

	return nil
}
