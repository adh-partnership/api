package models

import "time"

type Role struct {
	ID        uint      `json:"id" gorm:"primaryKey" example:"1"`
	Name      string    `json:"name" gorm:"type:varchar(128);index" example:"wm"`
	Users     []*User   `json:"users" gorm:"many2many:user_roles"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
