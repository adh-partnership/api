package dto

import (
	"time"

	"github.com/adh-partnership/api/pkg/database/models"
)

type TrainingNoteRequest struct {
	Position    string    `json:"position"`
	Type        string    `json:"type"`
	Comments    string    `json:"comments"`
	Duration    string    `json:"duration"`
	SessionDate time.Time `json:"session_date"`
}

type TrainingSessionRequestCreateRequest struct {
	TrainingType string `json:"training_type"`
	TrainingFor  string `json:"training_for"`
	Notes        string `json:"notes"`
}

type TrainingSessionRequest struct {
	ID           string        `json:"id"`
	User         *UserResponse `json:"user"`
	TrainingType string        `json:"training_type"`
	TrainingFor  string        `json:"training_for"`
	Notes        string        `json:"notes"`
	Status       string        `json:"status"`
	Instructor   *UserResponse `json:"instructor"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

func ConvertTrainingRequestToDTO(t *models.TrainingSessionRequest) *TrainingSessionRequest {
	return &TrainingSessionRequest{
		ID:           t.ID.String(),
		User:         ConvUserToUserResponse(t.User),
		TrainingType: t.TrainingType,
		TrainingFor:  t.TrainingFor,
		Notes:        t.Notes,
		Status:       t.Status,
		Instructor:   ConvUserToUserResponse(t.Instructor),
		CreatedAt:    t.CreatedAt,
		UpdatedAt:    t.UpdatedAt,
	}
}

func ConvertTrainingRequestsToDTO(t []*models.TrainingSessionRequest) []*TrainingSessionRequest {
	var res []*TrainingSessionRequest
	for _, v := range t {
		res = append(res, ConvertTrainingRequestToDTO(v))
	}
	return res
}
