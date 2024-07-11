package entity

import "time"

type Showtime struct {
	ID               string    `json:"id"`
	StudioID         string    `json:"studio_id"`
	MovieID          string    `json:"movie_id"`
	ShowStart        time.Time `json:"show_start"`
	ShowEnd          time.Time `json:"show_end"`
	StudioName       string    `json:"genre_name"`
	MovieTitle       string    `json:"movie_title"`
	MovieDescription string    `json:"movie_description"`
	MoviePrice       int       `json:"movie_price"`
	MovieDuration    int       `json:"movie_duration"`
	MovieStatus      string    `json:"movie_status"`
}