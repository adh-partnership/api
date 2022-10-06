package models

type APIKeys struct {
	ID   int    `json:"id" example:"1"`
	Key  string `json:"key" example:"1234567890"`
	Role string `json:"role" example:"admin"`
}
