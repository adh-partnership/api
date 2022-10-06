package dto

import (
	"time"

	"github.com/adh-partnership/api/pkg/database/models"
)

type EventRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Banner      string     `json:"banner"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
}

type EventPositionRequest struct {
	Position string `json:"position"`
	UserID   uint   `json:"cid"`
}

type EventSignupRequest struct {
	Choice1 string `json:"choice1"`
	Choice2 string `json:"choice2"`
	Choice3 string `json:"choice3"`
	Notes   string `json:"notes"`
}

func PatchEventRequest(base *models.Event, patch EventRequest) *models.Event {
	if patch.Title != "" {
		base.Title = patch.Title
	}
	if patch.Description != "" {
		base.Description = patch.Description
	}
	if patch.Banner != "" {
		base.Banner = patch.Banner
	}
	if patch.StartDate != nil {
		base.StartDate = *patch.StartDate
	}
	if patch.EndDate != nil {
		base.EndDate = *patch.EndDate
	}
	return base
}
