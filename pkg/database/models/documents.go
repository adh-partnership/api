package models

import "time"

type Document struct {
	ID          int       `json:"id" example:"1"`
	Name        string    `json:"name" gorm:"type:varchar(100)" example:"document name"`
	Description string    `json:"description" gorm:"type:varchar(255)" example:"Description of document"`
	Category    string    `json:"category" gorm:"type:varchar(100)" example:"sops"`
	URL         string    `json:"url" gorm:"type:varchar(255)" example:"https://www.example.com/document.pdf"`
	CreatedAt   time.Time `json:"created_at" example:"2020-01-01T00:00:00Z"`
	CreatedByID uint      `json:"created_by_id" example:"1"`
	CreatedBy   User      `json:"created_by"`
	UpdatedAt   time.Time `json:"updated_at" example:"2020-01-01T00:00:00Z"`
	UpdatedByID uint      `json:"updated_by_id" example:"1"`
	UpdatedBy   User      `json:"updated_by"`
}
