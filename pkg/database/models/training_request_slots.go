package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TrainingRequestSlot struct {
	ID                uuid.UUID        `json:"id" gorm:"primary_key;type:char(36);not null;unique_index"`
	TrainingRequest   *TrainingRequest `json:"training_request" gorm:"foreignKey:TrainingRequestID;references:ID"`
	TrainingRequestID uuid.UUID        `json:"-" gorm:"not null"`
	Start             *time.Time       `json:"start" gorm:"not null"`
	End               *time.Time       `json:"end" gorm:"not null"`
	CreatedAt         *time.Time       `json:"created_at"`
	UpdatedAt         *time.Time       `json:"updated_at"`
}

func (t *TrainingRequestSlot) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.New()
	return nil
}
