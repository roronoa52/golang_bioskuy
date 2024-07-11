package dto

type PaymentRequest struct {
	SeatDetailForBookingID string `json:"seat_detail_for_booking_id" validate:"required"`
}

type CreatePaymentRequest struct {
	TotalSeat              int    `json:"total_seat"`
	TotalPrice             int    `json:"total_price"`
	SeatDetailForBookingID string `json:"seat_detail_for_booking_id" validate:"required"`
	UserID                 string `json:"user_id" validate:"required"`
}

type CreatePaymentResponse struct {
	ID                     string `json:"id"`
	TotalSeat              int    `json:"total_seat"`
	TotalPrice             int    `json:"total_price"`
	SeatDetailForBookingID string `json:"seat_detail_for_booking_id" validate:"required"`
	UserID                 string `json:"user_id" validate:"required"`
	URL                    string `json:"url"`
}

type PaymentResponse struct {
	ID string `json:"id"`

	UserID string `json:"user_id"`

	ShowtimeID string `json:"showtime_id"`
	ShowStart  string `json:"show_start"`
	ShowEnd    string `json:"show_end"`

	MovieID          string `json:"movie_id"`
	MovieTitle       string `json:"movie_title"`
	MovieDescription string `json:"movie_description"`
	MoviePrice       int    `json:"movie_price"`
	MovieDuration    int    `json:"movie_duration"`
	MovieStatus      string `json:"movie_status"`

	StudioID   string `json:"studio_id" validate:"required"`
	StudioName string `json:"studio_name"`

	SeatID          string `json:"seat_id"`
	SeatName        string `json:"seat_name"`
	SeatIsAvailable string `json:"seat_isAvailabe"`

	SeatDetailForBookingID string `json:"seat_detail_for_booking_id"`

	SeatBookingID     string `json:"seat_booking_id"`
	SeatBookingStatus string `json:"seat_booking_status"`

	TotalSeat  int    `json:"total_seat"`
	TotalPrice int    `json:"total_price"`
	Status     string `json:"status"`
}
