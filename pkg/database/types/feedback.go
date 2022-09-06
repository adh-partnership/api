package types

import "time"

type Feedback struct {
	ID             int       `json:"id" gorm:"primaryKey"`
	SubmitterCID   uint      `json:"submitter"`
	SubmitterName  string    `json:"submitter_name"`
	SubmitterEmail string    `json:"submitter_email"`
	FlightDate     string    `json:"flight_date"`
	FlightCallsign string    `json:"flight_callsign"`
	Rating         string    `json:"rating"`
	Position       string    `json:"position"`
	ControllerID   uint      `json:"-"`
	Controller     User      `json:"controller"`
	Comments       string    `json:"comments"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

var FeedbackStatus = map[string]string{
	"pending":  "pending",
	"approved": "approved",
	"rejected": "rejected",
}

var FeedbackRatings = map[string]string{
	"excellent": "excellent",
	"good":      "good",
	"fair":      "fair",
	"poor":      "poor",
}

func IsValidFeedbackRating(rating string) bool {
	_, ok := FeedbackRatings[rating]
	return ok
}

func IsValidFeedbackStatus(status string) bool {
	_, ok := FeedbackStatus[status]
	return ok
}
