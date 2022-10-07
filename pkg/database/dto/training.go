package dto

import "time"

type TrainingNoteRequest struct {
	Position    string    `json:"position"`
	Type        string    `json:"type"`
	Comments    string    `json:"comments"`
	Duration    string    `json:"duration"`
	SessionDate time.Time `json:"session_date"`
}
