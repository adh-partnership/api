/*
 * Copyright ADH Partnership
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

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

type EventsResponse struct {
	ID          uint                     `json:"id"`
	Title       string                   `json:"title"`
	Description string                   `json:"description"`
	Banner      string                   `json:"banner"`
	StartDate   time.Time                `json:"start_date"`
	EndDate     time.Time                `json:"end_date"`
	Positions   []*EventPositionResponse `json:"positions"`
	Signups     []*EventSignupResponse   `json:"signups"`
	CreatedAt   time.Time                `json:"created_at"`
	UpdatedAt   time.Time                `json:"updated_at"`
}

type EventPositionResponse struct {
	ID       uint          `json:"id"`
	Position string        `json:"position"`
	UserID   *uint         `json:"cid"`
	User     *UserResponse `json:"user"`
}

type EventSignupResponse struct {
	ID      uint          `json:"id"`
	Choice1 string        `json:"choice1"`
	Choice2 string        `json:"choice2"`
	Choice3 string        `json:"choice3"`
	Notes   string        `json:"notes"`
	UserID  *uint         `json:"cid"`
	User    *UserResponse `json:"user"`
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

func ConvEventsToEventsResponse(events []*models.Event) []*EventsResponse {
	res := []*EventsResponse{}

	for _, event := range events {
		res = append(res, ConvEventToEventsResponse(event))
	}
	return res
}

func ConvEventToEventsResponse(event *models.Event) *EventsResponse {
	return &EventsResponse{
		ID:          event.ID,
		Title:       event.Title,
		Description: event.Description,
		Banner:      event.Banner,
		StartDate:   event.StartDate,
		EndDate:     event.EndDate,
		Positions:   ConvEventPositionsToEventPositionResponse(event.Positions),
		Signups:     ConvEventSignupsToEventSignupResponse(event.Signups),
		CreatedAt:   event.CreatedAt,
		UpdatedAt:   event.UpdatedAt,
	}
}

func ConvEventPositionsToEventPositionResponse(positions []*models.EventPosition) []*EventPositionResponse {
	res := []*EventPositionResponse{}
	for _, position := range positions {
		res = append(res, ConvEventPositionToEventPositionResponse(position))
	}
	return res
}

func ConvEventPositionToEventPositionResponse(position *models.EventPosition) *EventPositionResponse {
	pos := &EventPositionResponse{
		ID:       position.ID,
		Position: position.Position,
		UserID:   position.UserID,
	}
	if position.User != nil {
		pos.User = ConvUserToUserResponse(position.User)
	}

	return pos
}

func ConvEventSignupsToEventSignupResponse(signups []*models.EventSignup) []*EventSignupResponse {
	res := []*EventSignupResponse{}
	for _, signup := range signups {
		res = append(res, ConvEventSignupToEventSignupResponse(signup))
	}
	return res
}

func ConvEventSignupToEventSignupResponse(signup *models.EventSignup) *EventSignupResponse {
	sup := &EventSignupResponse{
		ID:      signup.ID,
		Choice1: signup.Choice1,
		Choice2: signup.Choice2,
		Choice3: signup.Choice3,
		Notes:   signup.Notes,
		UserID:  signup.UserID,
	}

	if signup.User != nil {
		sup.User = ConvUserToUserResponse(signup.User)
	}

	return sup
}
