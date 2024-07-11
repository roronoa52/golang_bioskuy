package entity

type User struct {
	ID    string `json:"id" `
	Name  string `json:"name"`
	Email string `json:"email"`
	Token string `json:"token"`
	Role  string `json:"role"`
}