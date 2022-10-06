package models

import "time"

type EventSignup struct {
	ID        int   `json:"id" gorm:"primaryKey"`
	EventID   int   `json:"-"`
	Event     Event `json:"-"`
	UserID    uint
	User      User      `json:"user"`
	Choice1   string    `json:"choice1" gorm:"type:varchar(25)"`
	Choice2   string    `json:"choice2" gorm:"type:varchar(25)"`
	Choice3   string    `json:"choice3" gorm:"type:varchar(25)"`
	Notes     string    `json:"notes" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
