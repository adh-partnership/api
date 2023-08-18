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
	"fmt"
	"time"

	"github.com/adh-partnership/api/pkg/database"
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
	Comments string `json:"comments"`
	Status   string `json:"status"`
}

type FeedbackResponse struct {
	ID           int           `json:"id"`
	Submitter    *UserResponse `json:"submitter"`
	Controller   *UserResponse `json:"controller"`
	Rating       string        `json:"rating"`
	Status       string        `json:"status"`
	Position     string        `json:"position"`
	Callsign     string        `json:"callsign"`
	Comments     string        `json:"comments"`
	ContactEmail string        `json:"contact_email"`
	CreatedAt    *time.Time    `json:"created_at"`
}

func ConvertFeedbacktoResponse(feedback []*models.Feedback, includeEmail bool) []FeedbackResponse {
	ret := []FeedbackResponse{}

	for _, f := range feedback {
		controller, _ := database.FindUserByCID(fmt.Sprint(f.Controller.CID))
		submitter, _ := database.FindUserByCID(fmt.Sprint(f.Submitter.CID))
		fdbk := &FeedbackResponse{
			ID:         f.ID,
			Submitter:  ConvUserToUserResponse(submitter),
			Controller: ConvUserToUserResponse(controller),
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

func ConvertSingleFeedbacktoResponse(feedback *models.Feedback, includeEmail bool) *FeedbackResponse {
	controller, _ := database.FindUserByCID(fmt.Sprint(feedback.Controller.CID))
	submitter, _ := database.FindUserByCID(fmt.Sprint(feedback.Submitter.CID))
	fdbk := &FeedbackResponse{
		ID:         feedback.ID,
		Submitter:  ConvUserToUserResponse(submitter),
		Controller: ConvUserToUserResponse(controller),
		Rating:     feedback.Rating,
		Position:   feedback.Position,
		Callsign:   feedback.Callsign,
		Comments:   feedback.Comments,
		Status:     feedback.Status,
		CreatedAt:  feedback.CreatedAt,
	}

	if includeEmail {
		fdbk.ContactEmail = feedback.ContactEmail
	}

	return fdbk
}
