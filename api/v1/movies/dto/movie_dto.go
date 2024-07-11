package dto

type CreateMovieDTO struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Price       int    `json:"price" binding:"required"`
	Duration    int    `json:"duration" binding:"required"`
	Status      string `json:"status" binding:"required"`
}

type UpdateMovieDTO struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Price       int    `json:"price" binding:"required"`
	Duration    int    `json:"duration" binding:"required"`
	Status      string `json:"status" binding:"required"`
}

type MovieResponseDTO struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	Duration    int    `json:"duration"`
	Status      string `json:"status"`
}
