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
