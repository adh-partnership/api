package models

import "time"

type EventPosition struct {
	ID        uint      `json:"id"`
	EventID   uint      `json:"-" gorm:"index:event_position"`
	Event     Event     `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Position  string    `json:"position" gorm:"index:event_position"`
	UserID    *uint     `json:"-"`
	User      *User     `json:"user"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
