package entitymock

import (
	entityMovie "bioskuy/api/v1/movies/entity"

	"bioskuy/api/v1/showtime/entity"
	entityStudio "bioskuy/api/v1/studio/entity"
	"time"
)

var MockMovie = entityMovie.Movie{
	ID:          "1",
	Title:       "Mock Movie",
	Description: "A mock movie for testing",
	Price:       10000,
	Duration:    2,
	Status:      "Active",
}

var MockStudio = entityStudio.Studio{
	ID:   "1",
	Name: "Mock Studio",
}

var MockShowtime = entity.Showtime{
	ID:               "1",
	MovieID:          "1",
	StudioID:         "1",
	ShowStart:        time.Now(),
	ShowEnd:          time.Now().Add(2 * time.Hour),
	StudioName:       "Mock Studio",
	MovieTitle:       "Mock Movie",
	MovieDescription: "A mock movie for testing",
	MoviePrice:       10000,
	MovieDuration:    2,
	MovieStatus:      "Active",
}
