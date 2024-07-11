package entity

type Movie struct {
	ID          string `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       int       `json:"price"`
	Duration    int       `json:"duration"`
	Status      string    `json:"status"`
}
