package dto

import (
	"time"

	"github.com/adh-partnership/api/pkg/database/models"
)

type FeedbackRequest struct {
	Controller string `json:"controller" binding:"required"`
	Callsign   string `json:"callsign" binding:"required"`
	Position   string `json:"position" binding:"required"`
	Rating     string `json:"rating" binding:"required"`
	Comments   string `json:"comments" binding:"required"`
}

type FeedbackPatchRequest struct {
	Status string `json:"status"`
}

type FeedbackResponse struct {
	ID           int          `json:"id"`
	Submitter    *models.User `json:"submitter"`
	Controller   *models.User `json:"controller"`
	Rating       string       `json:"rating"`
	Status       string       `json:"status"`
	Position     string       `json:"position"`
	Callsign     string       `json:"callsign"`
	Comments     string       `json:"comments"`
	ContactEmail string       `json:"contact_email"`
	CreatedAt    *time.Time   `json:"created_at"`
}

func ConvertFeedbacktoResponse(feedback []*models.Feedback, includeEmail bool) []FeedbackResponse {
	var ret []FeedbackResponse

	for _, f := range feedback {
		fdbk := &FeedbackResponse{
			ID:         f.ID,
			Submitter:  f.Submitter,
			Controller: f.Controller,
			Rating:     f.Rating,
			Position:   f.Position,
			Callsign:   f.Callsign,
			Comments:   f.Comments,
			Status:     f.Status,
			CreatedAt:  f.CreatedAt,
		}

		if includeEmail {
			fdbk.ContactEmail = f.ContactEmail
		}

		ret = append(ret, *fdbk)
	}

	return ret
}
