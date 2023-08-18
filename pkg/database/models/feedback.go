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

package models

import (
	"time"

	"github.com/adh-partnership/api/pkg/database/models/constants"
)

type Feedback struct {
	ID           int        `json:"id" gorm:"primaryKey"`
	SubmitterID  uint       `json:"-"`
	Submitter    *User      `json:"submitter"`
	ControllerID uint       `json:"-"`
	Controller   *User      `json:"controller"`
	Rating       string     `json:"rating" gorm:"type:varchar(20);not null"`
	Status       string     `json:"status" gorm:"type:varchar(20);not null"`
	Position     string     `json:"position" gorm:"type:varchar(20);not null"`
	Callsign     string     `json:"callsign" gorm:"type:varchar(20);not null"`
	Comments     string     `json:"comments" gorm:"type:text"`
	ContactEmail string     `json:"contact_email" gorm:"type:varchar(255)"`
	CreatedAt    *time.Time `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
}

func IsValidFeedbackRating(rating string) bool {
	switch rating {
	case constants.FeedbackRatingExcellent,
		constants.FeedbackRatingGood,
		constants.FeedbackRatingFair,
		constants.FeedbackRatingPoor:
		return true
	default:
		return false
	}
}

func IsValidFeedbackStatus(status string) bool {
	switch status {
	case constants.FeedbackStatusPending,
		constants.FeedbackStatusApproved,
		constants.FeedbackStatusRejected:
		return true
	default:
		return false
	}
}
