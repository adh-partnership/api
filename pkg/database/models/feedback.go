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
