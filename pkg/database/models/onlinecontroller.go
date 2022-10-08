package models

import "time"

type OnlineController struct {
	ID        int       `json:"id" example:"1"`
	UserID    uint      `json:"cid" example:"1"`
	User      *User     `json:"user"`
	Position  string    `json:"position" example:"ANC_TWR" gorm:"index"`
	Frequency string    `json:"frequency" example:"118.000"`
	LogonTime time.Time `json:"logon_time" example:"2020-01-01T00:00:00Z"`
	UpdateID  string    `json:"update_id" example:"1"`
	CreatedAt time.Time `json:"created_at" example:"2020-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2020-01-01T00:00:00Z"`
}
