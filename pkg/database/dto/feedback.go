package dto

import (
	dbTypes "github.com/kzdv/api/pkg/database/types"
)

type FeedbackRequest struct {
	FlightDate     string `json:"flight_date"`
	FlightCallsign string `json:"flight_callsign"`
	Controller     uint   `json:"controller"`
	Rating         string `json:"rating"`
	Position       string `json:"position"`
	Comments       string `json:"comments"`
}

type FeedbackPatchRequest struct {
	Status string `json:"status"`
}

type FeedbackResponse struct {
	ID               int    `json:"id"`
	SubmitterCID     uint   `json:"submitter"`
	SubmitterName    string `json:"submitter_name"`
	FlightDate       string `json:"flight_date"`
	FlightCallsign   string `json:"flight_callsign"`
	Rating           string `json:"rating"`
	Position         string `json:"position"`
	ControllerID     uint   `json:"controller_cid"`
	ControllerName   string `json:"controller_name"`
	ControllerRating string `json:"controller_rating"`
	Comments         string `json:"comments"`
	Status           string `json:"status"`
	CreatedAt        string `json:"created_at"`
}

func ConvertFeedbacktoResponse(feedback []dbTypes.Feedback) []FeedbackResponse {
	var ret []FeedbackResponse

	for _, f := range feedback {
		ret = append(ret, FeedbackResponse{
			ID:               f.ID,
			SubmitterCID:     f.SubmitterCID,
			SubmitterName:    f.SubmitterName,
			FlightDate:       f.FlightDate,
			FlightCallsign:   f.FlightCallsign,
			Rating:           f.Rating,
			Position:         f.Position,
			ControllerID:     f.Controller.CID,
			ControllerName:   f.Controller.FirstName + " " + f.Controller.LastName,
			ControllerRating: f.Controller.Rating.Short,
			Comments:         f.Comments,
			Status:           f.Status,
			CreatedAt:        f.CreatedAt.Format("2006-01-02"),
		})
	}

	return ret
}
