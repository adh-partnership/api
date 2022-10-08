package models

import "time"

type VisitorApplication struct {
	ID        uint      `json:"id" example:"1"`
	UserID    uint      `json:"-" example:"1" gorm:"index"`
	User      *User     `json:"user"`
	CreatedAt time.Time `json:"created_at" example:"2020-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2020-01-01T00:00:00Z"`
}
