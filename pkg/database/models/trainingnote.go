package models

import "time"

type TrainingNote struct {
	ID           uint       `json:"id" gorm:"primaryKey"`
	ControllerID uint       `json:"controller_id"`
	Controller   *User      `json:"controller"`
	InstructorID uint       `json:"instructor_id"`
	Instructor   *User      `json:"instructor"`
	Position     string     `json:"position"`
	Type         string     `json:"type"`
	Comments     string     `json:"comments"`
	SessionDate  *time.Time `json:"session_date"`
	Duration     string     `json:"duration"`
	VATUSAID     uint       `json:"vatusa_id"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

var TrainingNoteTypes = map[string]string{
	"classroom":      "classroom",
	"live":           "live",
	"simulation":     "simulation",
	"live-ots":       "live-ots",
	"simulation-ots": "simulation-ots",
	"no-show":        "no-show",
	"other":          "other",
}

var TrainingNoteTypesToVATUSA = map[string]string{
	"classroom":      "0",
	"live":           "1",
	"simulation":     "2",
	"live-ots":       "1",
	"simulation-ots": "2",
	"no-show":        "0",
	"other":          "0",
}
