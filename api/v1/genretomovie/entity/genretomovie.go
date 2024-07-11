package entity

type GenreToMovie struct {
	ID               string `db:"id"`
	GenreID          string `db:"genre_id"`
	MovieID          string `db:"movie_id"`
	GenreName        string `db:"genre_name"`
	MovieTitle       string `db:"movie_title"`
	MovieDescription string `db:"movie_description"`
	MoviePrice       int    `db:"movie_price"`
	MovieDuration    int    `db:"movie_duration"`
	MovieStatus      string `db:"movie_status"`
}
