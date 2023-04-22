package models

import (
	"time"

	"github.com/google/uuid"
)

type TrainingRequestSlot struct {
	UUIDBase
	TrainingRequest   *TrainingRequest `json:"training_request" gorm:"foreignKey:TrainingRequestID;references:ID"`
	TrainingRequestID uuid.UUID        `json:"-" gorm:"not null"`
	StartTime         *time.Time       `json:"start_time" gorm:"not null"`
	EndTime           *time.Time       `json:"end_time" gorm:"not null"`
}
