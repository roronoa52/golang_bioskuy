package dto

type CreateSeatRequest struct {
	Name        string `json:"name" validate:"required"`
	IsAvailable bool   `json:"is-available"`
	StudioID    string `json:"studio-id" validate:"required"`
}

type SeatResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	IsAvailable bool   `json:"is-available"`
	StudioID    string `json:"studio-id"`
}