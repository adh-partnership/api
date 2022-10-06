package models

import "time"

type OAuthRefresh struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Token     string    `json:"token" gorm:"type:varchar(128);index"`
	CID       uint      `json:"cid" gorm:"type:int;index"`
	ClientID  string    `json:"client_id" gorm:"type:varchar(128);index"`
	ExpiresAt time.Time `json:"expires_at" gorm:"type:datetime"`
	CreatedAt time.Time `json:"created_at" gorm:"type:datetime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:datetime"`
}
