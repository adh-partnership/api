package models

import (
	"time"

	"github.com/adh-partnership/api/pkg/config"
	"github.com/adh-partnership/api/pkg/database/models/constants"
	"github.com/adh-partnership/api/pkg/database/types"
)

type TrainingRequest struct {
	ID           types.BinaryData `json:"id" gorm:"primary_key;default:(UUID_TO_BIN(UUID()));"`
	Student      *User            `json:"student" gorm:"foreignKey:StudentID;references:ID"`
	StudentID    uint             `json:"-" gorm:"not null"`
	Position     string           `json:"position" gorm:"not null"`
	Status       string           `json:"status" gorm:"not null"`
	Notes        string           `json:"notes"`
	Instructor   *User            `json:"instructor" gorm:"foreignKey:InstructorID;references:ID"`
	InstructorID *uint            `json:"-" gorm:"default:null"`
	Availability string           `json:"availability"`
	CreatedAt    *time.Time       `json:"created_at"`
	UpdatedAt    *time.Time       `json:"updated_at"`
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
