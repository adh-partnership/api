package models

import "time"

type EventSignup struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	EventID   uint      `json:"-"`
	Event     Event     `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID    *uint     `json:"-"`
	User      *User     `json:"user"`
	Choice1   string    `json:"choice1" gorm:"type:varchar(25)"`
	Choice2   string    `json:"choice2" gorm:"type:varchar(25)"`
	Choice3   string    `json:"choice3" gorm:"type:varchar(25)"`
	Notes     string    `json:"notes" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
