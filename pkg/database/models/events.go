package models

import (
	"time"
)

type Event struct {
	ID          uint             `json:"id" gorm:"primaryKey"`
	Title       string           `json:"title"`
	Description string           `json:"description"`
	Banner      string           `json:"banner"`
	StartDate   time.Time        `json:"start_date"`
	EndDate     time.Time        `json:"end_date"`
	Positions   []*EventPosition `json:"positions"`
	Signups     []*EventSignup   `json:"signups"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}
