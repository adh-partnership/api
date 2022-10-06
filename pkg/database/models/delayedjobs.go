package models

import "time"

type DelayedJob struct {
	ID        int       `json:"id" example:"1"`
	Queue     string    `json:"queue" gorm:"type:varchar(128)" example:"email"`
	Body      string    `json:"body" gorm:"type:text"`
	NotBefore time.Time `json:"not_before" example:"2020-01-01T00:00:00Z"`
	CreatedAt time.Time `json:"created_at" example:"2020-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2020-01-01T00:00:00Z"`
}
