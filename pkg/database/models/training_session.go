package models

import (
	"time"

	"github.com/adh-partnership/api/pkg/config"
	"github.com/adh-partnership/api/pkg/database/models/constants"
)

type TrainingSessionRequest struct {
	UUIDBase
	UserID uint  `json:"user_id"`
	User   *User `json:"user"`
	// Should be one of: Simulation, Live, OTS, Other -- not enforced
	TrainingType string `json:"training_type"`
	// Must be one of the fields specified in config facility.training.positions array
	TrainingFor string `json:"training_for"`
	Notes       string `json:"notes"`
	// Must be one of: none, open, accepted, completed, cancelled
	Status           string     `json:"status"`
	ScheduledSession *time.Time `json:"scheduled_session"`
	InstructorID     *uint      `json:"instructor_id"`
	Instructor       *User      `json:"instructor"`
	InstructorNotes  string     `json:"instructor_notes"`
}

func IsValidPosition(pos string) bool {
	for _, v := range config.Cfg.Facility.TrainingRequests.Positions {
		if v == pos {
			return true
		}
	}

	return false
}

func IsValidTrainingType(t string) bool {
	return t == constants.TrainingSessionStatusOpen || t == constants.TrainingSessionStatusAccepted ||
		t == constants.TrainingSessionStatusCompleted || t == constants.TrainingSessionStatusCancelled
}
