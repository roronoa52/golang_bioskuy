package repository

import (
	"bioskuy/api/v1/payment/entity"
	"bioskuy/exception"
	"context"
	"database/sql"
	"errors"

	"github.com/gin-gonic/gin"
)

    type paymentRepository struct {
    }

    func NewPaymentRepository() PaymentRepository {
        return &paymentRepository{}
    }

    func (r *paymentRepository) Save(ctx context.Context, tx *sql.Tx, payment entity.Payment, c *gin.Context) (entity.Payment, error){

        query := "INSERT INTO payments (user_id, seatdetailforbooking_id, total_seat, total_price) VALUES ($1, $2, $3, $4) RETURNING id"

        err := tx.QueryRowContext(ctx, query, payment.UserID, payment.SeatDetailForBookingID, payment.TotalSeat, payment.TotalPrice).Scan(&payment.ID)
        if err != nil {
            c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
            return payment, err
        }

        return payment, nil
    }

    func (r *paymentRepository) FindByID(ctx context.Context, tx *sql.Tx, id string, c *gin.Context) (entity.Payment, error) {
        query := `
        SELECT 
            p.id, p.user_id, p.seatdetailforbooking_id, p.total_seat, p.total_price, p.status,
            sb.id AS seat_booking_id, sb.status AS seat_booking_status,
            s.id AS seat_id, s.seat_name, s.isAvailable AS seat_isAvailable,
            sh.id AS showtime_id, sh.show_start, sh.show_end,
            m.id AS movie_id, m.title AS movie_title, m.description AS movie_description, 
            m.price AS movie_price, m.duration AS movie_duration, m.status AS movie_status,
            st.id AS studio_id, st.name AS studio_name
        FROM payments p
        JOIN seat_detail_for_bookings sdfb ON p.seatdetailforbooking_id = sdfb.id
        JOIN seats s ON sdfb.seat_id = s.id
        JOIN seat_bookings sb ON sdfb.seatBooking_id = sb.id
        JOIN showtimes sh ON sb.showtime_id = sh.id
        JOIN movies m ON sh.movie_id = m.id
        JOIN studios st ON sh.studio_id = st.id
        WHERE p.id = $1`
        
        paymentResponse := entity.Payment{}
        rows, err := tx.QueryContext(ctx, query, id)
        if err != nil {
            c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
            return paymentResponse, err
        }
        defer rows.Close()

        if rows.Next() {
            err := rows.Scan(
                &paymentResponse.ID, &paymentResponse.UserID, &paymentResponse.SeatDetailForBookingID, &paymentResponse.TotalSeat, &paymentResponse.TotalPrice, &paymentResponse.Status,
                &paymentResponse.SeatBookingID, &paymentResponse.SeatBookingStatus,
                &paymentResponse.SeatID, &paymentResponse.SeatName, &paymentResponse.SeatIsAvailable,
                &paymentResponse.ShowtimeID, &paymentResponse.ShowStart, &paymentResponse.ShowEnd,
                &paymentResponse.MovieID, &paymentResponse.MovieTitle, &paymentResponse.MovieDescription, &paymentResponse.MoviePrice, &paymentResponse.MovieDuration, &paymentResponse.MovieStatus,
                &paymentResponse.StudioID, &paymentResponse.StudioName,
            )
            if err != nil {
                c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
                return paymentResponse, err
            }

            return paymentResponse, nil
        } else {
            return paymentResponse, errors.New("payment not found")
        }
    }

    func (r *paymentRepository) FindAll(ctx context.Context, tx *sql.Tx, c *gin.Context) ([]entity.Payment, error) {
        query := `
        SELECT 
            p.id, p.user_id, p.seatdetailforbooking_id, p.total_seat, p.total_price, p.status,
            sb.id AS seat_booking_id, sb.status AS seat_booking_status,
            s.id AS seat_id, s.seat_name, s.isAvailable AS seat_isAvailable,
            sh.id AS showtime_id, sh.show_start, sh.show_end,
            m.id AS movie_id, m.title AS movie_title, m.description AS movie_description, 
            m.price AS movie_price, m.duration AS movie_duration, m.status AS movie_status,
            st.id AS studio_id, st.name AS studio_name
        FROM payments p
        JOIN seat_detail_for_bookings sdfb ON p.seatdetailforbooking_id = sdfb.id
        JOIN seats s ON sdfb.seat_id = s.id
        JOIN seat_bookings sb ON sdfb.seatBooking_id = sb.id
        JOIN showtimes sh ON sb.showtime_id = sh.id
        JOIN movies m ON sh.movie_id = m.id
        JOIN studios st ON sh.studio_id = st.id`
        
        rows, err := tx.QueryContext(ctx, query)
        if err != nil {
            c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
            return nil, err
        }
        defer rows.Close()

        var payments []entity.Payment

        for rows.Next() {
            var paymentResponse entity.Payment
            err := rows.Scan(
                &paymentResponse.ID, &paymentResponse.UserID, &paymentResponse.SeatDetailForBookingID, &paymentResponse.TotalSeat, &paymentResponse.TotalPrice, &paymentResponse.Status,
                &paymentResponse.SeatBookingID, &paymentResponse.SeatBookingStatus,
                &paymentResponse.SeatID, &paymentResponse.SeatName, &paymentResponse.SeatIsAvailable,
                &paymentResponse.ShowtimeID, &paymentResponse.ShowStart, &paymentResponse.ShowEnd,
                &paymentResponse.MovieID, &paymentResponse.MovieTitle, &paymentResponse.MovieDescription, &paymentResponse.MoviePrice, &paymentResponse.MovieDuration, &paymentResponse.MovieStatus,
                &paymentResponse.StudioID, &paymentResponse.StudioName,
            )
            if err != nil {
                c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
                return nil, err
            }

            payments = append(payments, paymentResponse)
        }

        if err = rows.Err(); err != nil {
            c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
            return nil, err
        }

        return payments, nil
    }

    func (r *paymentRepository) Update(ctx context.Context, tx *sql.Tx, payment entity.Payment, c *gin.Context) (entity.Payment, error){

        query := `UPDATE payments SET status = $1 WHERE id = $2`
    
        _, err := tx.ExecContext(ctx, query, payment.Status, payment.ID)
    
        if err != nil {
            c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
            return  payment, err
        }
    
        return payment, nil
    }

    func (r *paymentRepository) Delete(ctx context.Context, tx *sql.Tx, paymentId string, c *gin.Context) error {

        deleteSeatDetailQuery := `DELETE FROM payments WHERE id = $1 `
        _, err := tx.ExecContext(ctx, deleteSeatDetailQuery, paymentId)
        if err != nil {
            c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
            return err
        }
    
        return nil
    }
