package models

import (
	"time"

	"github.com/adh-partnership/api/pkg/database/types"
)

type TrainingRequestSlot struct {
	ID                types.BinaryData `json:"id" gorm:"primary_key;default:(UUID_TO_BIN(UUID()));"`
	TrainingRequest   *TrainingRequest `json:"training_request" gorm:"foreignKey:TrainingRequestID;references:ID"`
	TrainingRequestID types.BinaryData `json:"-" gorm:"not null"`
	StartTime         *time.Time       `json:"start_time" gorm:"not null"`
	EndTime           *time.Time       `json:"end_time" gorm:"not null"`
	CreatedAt         *time.Time       `json:"created_at"`
	UpdatedAt         *time.Time       `json:"updated_at"`
}
