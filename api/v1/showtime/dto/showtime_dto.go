package dto

import (
	"time"
)

type ShowtimeRequest struct {
	MovieID string `json:"movie_id"`
	StudioID string `json:"studio_id"`
	ShowStart string `json:"show_start"`
	ShowEnd string `json:"show_end"`
}

type CreateShowtimeDTO struct {
	MovieID string `json:"movie_id" validate:"required"`
	StudioID string `json:"studio_id" validate:"required"`
	ShowStart time.Time `json:"show_start" validate:"required"`
	ShowEnd time.Time `json:"show_end" validate:"required"`
}

type CreateShowtimesResponseDTO struct {
	MovieID string `json:"movie_id" validate:"required"`
	StudioID string `json:"studio_id" validate:"required"`
	ShowStart time.Time `json:"show_start" validate:"required"`
	ShowEnd time.Time `json:"show_end" validate:"required"`
}

type ShowtimesResponse struct {
	ID               	string `json:"id"`
	StudioID 			string `json:"studio_id"`
	MovieID          	string `json:"movie_id"`
	StudioName        	string `json:"studio_name"`
	MovieTitle       	string `json:"movie_title"`
	MovieDescription 	string `json:"movie_description"`
	MoviePrice       	int    `json:"movie_price"`
	MovieDuration    	int    `json:"movie_duration"`
	MovieStatus      	string `json:"movie_status"`
}

