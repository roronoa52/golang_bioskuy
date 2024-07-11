package dto

import "github.com/google/uuid"

type CreateGenreDTO struct {
	Name string `json:"name" binding:"required"`
}

type UpdateGenreDTO struct {
	Name string `json:"name" binding:"required"`
}

type GenreResponseDTO struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type Paging struct {
	Page       int `json:"page"`
	Size       int `json:"size"`
	TotalRows  int `json:"total_rows"`
	TotalPages int `json:"total_pages"`
}
