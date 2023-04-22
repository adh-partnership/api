package models

import (
	"time"

	"github.com/adh-partnership/api/pkg/config"
	"github.com/adh-partnership/api/pkg/database/models/constants"
)

type TrainingRequest struct {
	UUIDBase
	Student         *User                  `json:"student" gorm:"foreignKey:StudentID;references:ID"`
	StudentID       uint                   `json:"-" gorm:"not null"`
	Position        string                 `json:"position" gorm:"not null"`
	Status          string                 `json:"status" gorm:"not null"`
	Notes           string                 `json:"notes"`
	Instructor      *User                  `json:"instructor" gorm:"foreignKey:InstructorID;references:ID"`
	InstructorID    *uint                  `json:"-" gorm:"default:null"`
	InstructorNotes string                 `json:"instructor_notes"`
	Slots           []*TrainingRequestSlot `json:"slots" gorm:"foreignKey:TrainingRequestID;references:ID"`
	Start           *time.Time             `json:"start" gorm:"default:null"`
	End             *time.Time             `json:"end" gorm:"default:null"`
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
