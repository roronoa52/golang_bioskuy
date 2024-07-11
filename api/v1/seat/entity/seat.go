package entity

type Seat struct {
	ID          string `json:"id" `
	Name        string `json:"name"`
	IsAvailable bool   `json:"is-available"`
	StudioID    string `json:"studio-id"`
}