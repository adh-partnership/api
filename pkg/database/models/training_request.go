package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/adh-partnership/api/pkg/config"
	"github.com/adh-partnership/api/pkg/database/models/constants"
)

type TrainingRequest struct {
	ID              uuid.UUID              `json:"id" gorm:"primary_key;type:char(36);not null;unique_index"`
	Student         *User                  `json:"student" gorm:"foreignKey:StudentID"`
	StudentID       uint                   `json:"-" gorm:"not null"`
	Position        string                 `json:"position" gorm:"not null"`
	Status          string                 `json:"status" gorm:"not null"`
	Notes           string                 `json:"notes"`
	Instructor      *User                  `json:"instructor" gorm:"foreignKey:InstructorID"`
	InstructorID    *uint                  `json:"-" gorm:"default:null"`
	InstructorNotes string                 `json:"instructor_notes"`
	Slots           []*TrainingRequestSlot `json:"slots" gorm:"foreignKey:TrainingRequestID"`
	Start           *time.Time             `json:"start" gorm:"default:null"`
	End             *time.Time             `json:"end" gorm:"default:null"`
	CreatedAt       *time.Time             `json:"created_at"`
	UpdatedAt       *time.Time             `json:"updated_at"`
}

func (t *TrainingRequest) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.New()
	return nil
}

func IsValidPosition(pos string) bool {
	for _, v := range config.Cfg.Facility.TrainingRequests.Positions {
		if v == pos {
			return true
		}
	}

	return false
}

func IsValidTrainingStatus(s string) bool {
	return s == constants.TrainingSessionStatusOpen || s == constants.TrainingSessionStatusAccepted ||
		s == constants.TrainingSessionStatusCompleted || s == constants.TrainingSessionStatusCancelled
}
