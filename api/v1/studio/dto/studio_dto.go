package dto

type CreateStudioRequest struct {
	Name       string `json:"name" validate:"required"`
	Capacity   int    `json:"capacity" validate:"required"`
	MaxRowSeat int    `json:"max-row-seat" validate:"required"`
}

type UpdateStudioRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type StudioResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Capacity int    `json:"capacity"`
}