package models

import "time"

type ControllerStat struct {
	ID        uint      `json:"id" example:"1"`
	UserID    uint      `json:"cid" example:"1" gorm:"index"`
	User      *User     `json:"user"`
	Position  string    `json:"position" example:"ANC_TWR" gorm:"index"`
	LogonTime time.Time `json:"logon_time" example:"2020-01-01T00:00:00Z" gorm:"index"`
	Duration  int       `json:"duration" example:"3600"`
	CreatedAt time.Time `json:"created_at" example:"2020-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2020-01-01T00:00:00Z"`
}
