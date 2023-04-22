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

type TrainingRequestCreateRequest struct {
	Position string                 `json:"position"`
	Notes    string                 `json:"notes"`
	Slots    []*TrainingRequestSlot `json:"slots"`
}

type TrainingRequestEditRequest struct {
	Position        string `json:"position"`
	Notes           string `json:"notes"`
	Status          string `json:"status"`
	InstructorNotes string `json:"instructor_notes"`
	Instructor      uint   `json:"instructor"`
	Start           string `json:"start"`
	End             string `json:"end"`
}

type TrainingRequest struct {
	ID              string                 `json:"id"`
	Student         *UserResponse          `json:"student"`
	Instructor      *UserResponse          `json:"instructor"`
	Position        string                 `json:"position"`
	Status          string                 `json:"status"`
	Notes           string                 `json:"notes"`
	InstructorNotes string                 `json:"instructor_notes"`
	Start           *time.Time             `json:"start"`
	End             *time.Time             `json:"end"`
	Slots           []*TrainingRequestSlot `json:"slots"`
	CreatedAt       *time.Time             `json:"created_at"`
	UpdatedAt       *time.Time             `json:"updated_at"`
}

type TrainingRequestSlot struct {
	ID        string     `json:"id"`
	StartTime *time.Time `json:"start_time"`
	EndTime   *time.Time `json:"end_time"`
}

func ConvertTrainingRequestToDTO(t *models.TrainingRequest) *TrainingRequest {
	ret := &TrainingRequest{
		ID:        t.ID.String(),
		Student:   ConvUserToUserResponse(t.Student),
		Position:  t.Position,
		Status:    t.Status,
		Notes:     t.Notes,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}

	if t.Instructor != nil {
		ret.Instructor = ConvUserToUserResponse(t.Instructor)
		ret.InstructorNotes = t.InstructorNotes
		ret.Start = t.Start
		ret.End = t.End
	}

	for _, v := range t.Slots {
		ret.Slots = append(ret.Slots, &TrainingRequestSlot{
			ID:        v.ID.String(),
			StartTime: v.StartTime,
			EndTime:   v.EndTime,
		})
	}

	return ret
}

func ConvertTrainingRequestsToDTO(t []*models.TrainingRequest) []*TrainingRequest {
	var res []*TrainingRequest
	for _, v := range t {
		res = append(res, ConvertTrainingRequestToDTO(v))
	}
	return res
}
