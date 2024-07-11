package repository

import (
	"bioskuy/api/v1/seatbooking/entity"
	"bioskuy/exception"
	"context"
	"database/sql"
	"errors"

	"github.com/gin-gonic/gin"
)

type seatBookingRepository struct {
}

func NewSeatBookingRepository() SeatBookingRepository {
	return &seatBookingRepository{}
}

func (r *seatBookingRepository) Save(ctx context.Context, tx *sql.Tx, seatbooking entity.SeatBooking, c *gin.Context) (entity.SeatBooking, error){

	queryForSeatBooking := "INSERT INTO seat_bookings (user_id, showtime_id) VALUES ($1, $2) RETURNING id"

	err := tx.QueryRowContext(ctx, queryForSeatBooking, seatbooking.UserID, seatbooking.ShowtimeID).Scan(&seatbooking.ID)
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return seatbooking, err
	}

	queryForSeatDetailForBooking := "INSERT INTO seat_detail_for_bookings (seat_id, seatBooking_id) VALUES ($1, $2) RETURNING id"

	err = tx.QueryRowContext(ctx, queryForSeatDetailForBooking, seatbooking.SeatID, seatbooking.ID).Scan(&seatbooking.SeatDetailForBookingID)
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return seatbooking, err
	}

	return seatbooking, nil
}

func (r *seatBookingRepository) FindByID(ctx context.Context, tx *sql.Tx, id string, c *gin.Context) (entity.SeatBooking, error) {
	query := `
		SELECT
			sb.id, sb.status, sb.user_id,
			s.id as showtime_id, s.studio_id, s.movie_id, s.show_start, s.show_end,
			st.name as studio_name,
			m.title as movie_title, m.description as movie_description, m.price as movie_price, m.duration as movie_duration, m.status as movie_status,
			sdfb.id as seat_detail_for_booking_id, sdfb.seat_id,
			se.seat_name, se.isAvailable
		FROM
			seat_bookings sb
		JOIN
			showtimes s ON sb.showtime_id = s.id
		JOIN
			studios st ON s.studio_id = st.id
		JOIN
			movies m ON s.movie_id = m.id
		JOIN
			seat_detail_for_bookings sdfb ON sb.id = sdfb.seatBooking_id
		JOIN
			seats se ON sdfb.seat_id = se.id
		WHERE
			sb.id = $1
	`
	seatBookingResponse := entity.SeatBooking{}
	rows, err := tx.QueryContext(ctx, query, id)
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return seatBookingResponse, err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(
			&seatBookingResponse.ID, &seatBookingResponse.SeatBookingStatus, &seatBookingResponse.UserID,
			&seatBookingResponse.ShowtimeID, &seatBookingResponse.StudioID, &seatBookingResponse.MovieID, &seatBookingResponse.ShowStart, &seatBookingResponse.ShowEnd,
			&seatBookingResponse.StudioName,
			&seatBookingResponse.MovieTitle, &seatBookingResponse.MovieDescription, &seatBookingResponse.MoviePrice, &seatBookingResponse.MovieDuration, &seatBookingResponse.MovieStatus,
			&seatBookingResponse.SeatDetailForBookingID, &seatBookingResponse.SeatID,
			&seatBookingResponse.SeatName, &seatBookingResponse.SeatIsAvailable,
		)
		if err != nil {
			c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
			return seatBookingResponse, err
		}

		return seatBookingResponse, nil
	} else {
		return seatBookingResponse, errors.New("seatbooking not found")
	}
}

func (r *seatBookingRepository) FindAll(ctx context.Context, tx *sql.Tx, c *gin.Context) ([]entity.SeatBooking, error) {
	query := `
		SELECT
			sb.id, sb.status, sb.user_id,
			s.id as showtime_id, s.studio_id, s.movie_id, s.show_start, s.show_end,
			st.name as studio_name,
			m.title as movie_title, m.description as movie_description, m.price as movie_price, m.duration as movie_duration, m.status as movie_status,
			sdfb.id as seat_detail_for_booking_id, sdfb.seat_id,
			se.seat_name, se.isAvailable
		FROM
			seat_bookings sb
		JOIN
			showtimes s ON sb.showtime_id = s.id
		JOIN
			studios st ON s.studio_id = st.id
		JOIN
			movies m ON s.movie_id = m.id
		JOIN
			seat_detail_for_bookings sdfb ON sb.id = sdfb.seatBooking_id
		JOIN
			seats se ON sdfb.seat_id = se.id
	`
	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return nil, err
	}
	defer rows.Close()

	var seatBookings []entity.SeatBooking

	for rows.Next() {
		var seatBookingResponse entity.SeatBooking
		err := rows.Scan(
			&seatBookingResponse.ID, &seatBookingResponse.SeatBookingStatus, &seatBookingResponse.UserID,
			&seatBookingResponse.ShowtimeID, &seatBookingResponse.StudioID, &seatBookingResponse.MovieID, &seatBookingResponse.ShowStart, &seatBookingResponse.ShowEnd,
			&seatBookingResponse.StudioName,
			&seatBookingResponse.MovieTitle, &seatBookingResponse.MovieDescription, &seatBookingResponse.MoviePrice, &seatBookingResponse.MovieDuration, &seatBookingResponse.MovieStatus,
			&seatBookingResponse.SeatDetailForBookingID, &seatBookingResponse.SeatID,
			&seatBookingResponse.SeatName, &seatBookingResponse.SeatIsAvailable,
		)
		if err != nil {
			c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
			return nil, err
		}

		seatBookings = append(seatBookings, seatBookingResponse)
	}

	if err = rows.Err(); err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return nil, err
	}

	return seatBookings, nil
}

func (r *seatBookingRepository) Delete(ctx context.Context, tx *sql.Tx, seatBookingID string, c *gin.Context) error {

	updateSeatQuery := `UPDATE seats SET isAvailable = true WHERE id IN (SELECT seat_id FROM seat_detail_for_bookings WHERE seatBooking_id = $1)
	`
	_, err := tx.ExecContext(ctx, updateSeatQuery, seatBookingID)
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return err
	}

	deleteSeatDetailQuery := `DELETE FROM seat_detail_for_bookings  WHERE seatBooking_id = $1 `
	_, err = tx.ExecContext(ctx, deleteSeatDetailQuery, seatBookingID)
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return err
	}

	updateSeatBookingQuery := `UPDATE seat_bookings SET status = 'active' WHERE id = $1`
	_, err = tx.ExecContext(ctx, updateSeatBookingQuery, seatBookingID)
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return err
	}

	deleteSeatBookingQuery := `DELETE FROM seat_bookings WHERE id = $1`
	_, err = tx.ExecContext(ctx, deleteSeatBookingQuery, seatBookingID)
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return err
	}

	return nil
}

func (r *seatBookingRepository) FindAllPendingByUserID(ctx context.Context, tx *sql.Tx, userID string, c *gin.Context) ([]entity.SeatBooking, error) {
	query := `
		SELECT
			sb.id, sb.status, sb.user_id,
			s.id as showtime_id, s.studio_id, s.movie_id, s.show_start, s.show_end,
			st.name as studio_name,
			m.title as movie_title, m.description as movie_description, m.price as movie_price, m.duration as movie_duration, m.status as movie_status,
			sdfb.id as seat_detail_for_booking_id, sdfb.seat_id,
			se.seat_name, se.isAvailable
		FROM
			seat_bookings sb
		JOIN
			showtimes s ON sb.showtime_id = s.id
		JOIN
			studios st ON s.studio_id = st.id
		JOIN
			movies m ON s.movie_id = m.id
		JOIN
			seat_detail_for_bookings sdfb ON sb.id = sdfb.seatBooking_id
		JOIN
			seats se ON sdfb.seat_id = se.id
		WHERE
			sb.user_id = $1 AND sb.status = 'pending'
	`
	rows, err := tx.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seatBookings []entity.SeatBooking

	for rows.Next() {
		var seatBookingResponse entity.SeatBooking
		err := rows.Scan(
			&seatBookingResponse.ID, &seatBookingResponse.SeatBookingStatus, &seatBookingResponse.UserID,
			&seatBookingResponse.ShowtimeID, &seatBookingResponse.StudioID, &seatBookingResponse.MovieID, &seatBookingResponse.ShowStart, &seatBookingResponse.ShowEnd,
			&seatBookingResponse.StudioName,
			&seatBookingResponse.MovieTitle, &seatBookingResponse.MovieDescription, &seatBookingResponse.MoviePrice, &seatBookingResponse.MovieDuration, &seatBookingResponse.MovieStatus,
			&seatBookingResponse.SeatDetailForBookingID, &seatBookingResponse.SeatID,
			&seatBookingResponse.SeatName, &seatBookingResponse.SeatIsAvailable,
		)
		if err != nil {
			return nil, err
		}

		seatBookings = append(seatBookings, seatBookingResponse)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return seatBookings, nil
}

func (r *seatBookingRepository) Update(ctx context.Context, tx *sql.Tx, seatbooking entity.SeatBooking, c *gin.Context) (entity.SeatBooking, error){

	query := `UPDATE seat_bookings SET status = $1 WHERE id = $2`

	_, err := tx.ExecContext(ctx, query, seatbooking.SeatBookingStatus, seatbooking.ID)

	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return  seatbooking, err
	}

	return seatbooking, nil
}


