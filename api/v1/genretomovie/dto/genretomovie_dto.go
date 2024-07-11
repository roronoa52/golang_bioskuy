package dto

type CreateGenreToMovieRequest struct {
	GenreID string `json:"genre_id" validate:"required"`
	MovieID string `json:"movie_id" validate:"required"`
}

type GenreToMovieCreateResponse struct {
	ID      string `json:"id"`
	GenreID string `json:"genre_id"`
	MovieID string `json:"movie_id"`
}

type GenreToMovieResponse struct {
	ID               string `json:"id"`
	GenreID          string `json:"genre_id"`
	MovieID          string `json:"movie_id"`
	GenreName        string `json:"genre_name"`
	MovieTitle       string `json:"movie_title"`
	MovieDescription string `json:"movie_description"`
	MoviePrice       int    `json:"movie_price"`
	MovieDuration    int    `json:"movie_duration"`
	MovieStatus      string `json:"movie_status"`
}
