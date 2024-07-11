package entitymock

import (
	eS "bioskuy/api/v1/seat/entity"
	"bioskuy/api/v1/seatbooking/dto"
	"bioskuy/api/v1/seatbooking/entity"
	eST "bioskuy/api/v1/showtime/entity"
	"time"
)

var MockSeatBookingEntity = entity.SeatBooking{
	ID: "booking123",

	ShowtimeID: "showtime456",
	ShowStart:  "2024-07-10T15:00:00Z",
	ShowEnd:    "2024-07-10T17:00:00Z",

	MovieID:          "movie789",
	MovieTitle:       "Avengers: Endgame",
	MovieDescription: "The epic conclusion to the Infinity Saga.",
	MoviePrice:       15000,
	MovieDuration:    180,
	MovieStatus:      "released",

	StudioID:   "studio456",
	StudioName: "Studio 1",

	SeatID:          "seat123",
	SeatName:        "A1",
	SeatIsAvailable: "false",

	UserID:                 "user123",
	SeatDetailForBookingID: "seatdetail789",
	SeatBookingStatus:      "confirmed",
}

var MockSeatBookingRequest = dto.SeatBookingRequest{
	SeatID:     "seat123",
	ShowtimeID: "showtime456",
}

var MockCreateSeatBookingRequest = dto.CreateSeatBookingRequest{
	ID:         "booking789",
	UserID:     "user123",
	ShowtimeID: "showtime456",
	Status:     "confirmed",
}

var MockCreateSeatBookingResponse = dto.CreateSeatBookingResponse{
	ID:            "booking789",
	SeatID:        "seat123",
	SeatBookingID: "seatbooking789",
}

var MockSeatBookingResponse = dto.SeatBookingResponse{
	ID: "booking789",

	ShowtimeID: "showtime456",
	ShowStart:  "2024-07-10T15:00:00Z",
	ShowEnd:    "2024-07-10T17:00:00Z",

	MovieID:          "movie789",
	MovieTitle:       "Avengers: Endgame",
	MovieDescription: "The epic conclusion to the Infinity Saga.",
	MoviePrice:       15000,
	MovieDuration:    180,
	MovieStatus:      "released",

	StudioID:   "studio456",
	StudioName: "Studio 1",

	SeatID:          "seat123",
	SeatName:        "A1",
	SeatIsAvailable: "false",

	UserID:                 "user123",
	SeatDetailForBookingID: "seatdetail789",
	SeatBookingStatus:      "confirmed",
}

var MockShowtimeEntity = eST.Showtime{
	ID:               "showtime456",
	StudioID:         "studio456",
	MovieID:          "movie789",
	ShowStart:        time.Date(2024, 7, 10, 15, 0, 0, 0, time.UTC),
	ShowEnd:          time.Date(2024, 7, 10, 17, 0, 0, 0, time.UTC),
	StudioName:       "Studio 1",
	MovieTitle:       "Avengers: Endgame",
	MovieDescription: "The epic conclusion to the Infinity Saga.",
	MoviePrice:       15000,
	MovieDuration:    180,
	MovieStatus:      "released",
}

var MockSeatEntity = eS.Seat{
	ID:          "seat123",
	Name:        "A1",
	IsAvailable: false,
	StudioID:    "studio456",
}
