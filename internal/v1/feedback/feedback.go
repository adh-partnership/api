package feedback

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/adh-partnership/api/pkg/database"
	"github.com/adh-partnership/api/pkg/database/dto"
	dbTypes "github.com/adh-partnership/api/pkg/database/types"
	"github.com/adh-partnership/api/pkg/gin/response"
)

// Get Pilot Feedback
// @Summary Get Pilot Feedback
// @Description Get feedback for a pilot
// @Tags Feedback
// @Param cid query string false "Controller CID filter"
// @Param status query string false "Status filter"
// @Success 200 {object} dto.FeedbackResponse[]
// @Failure 400 {object} response.R
// @Failure 403 {object} response.R
// @Failure 500 {object} response.R
func getFeedback(c *gin.Context) {
	var feedback []dbTypes.Feedback

	res := database.DB
	if c.Query("cid") != "" {
		res = res.Where("controller_id = ?", c.Query("cid"))
	}
	if c.Query("status") != "" {
		res = res.Where("status = ?", c.Query("status"))
	}
	if err := res.Find(&feedback).Error; err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	response.Respond(c, http.StatusOK, dto.ConvertFeedbacktoResponse(feedback))
}

// Submit Pilot Feedback
// @Summary Submit Pilot Feedback
// @Description Submit feedback for a pilot
// @Tags Feedback
// @Param data body dto.FeedbackRequest true "Feedback"
// @Success 204
// @Failure 400 {object} response.R
// @Failure 403 {object} response.R
// @Failure 500 {object} response.R
func postFeedback(c *gin.Context) {
	var dto dto.FeedbackRequest
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.RespondError(c, http.StatusBadRequest, "Invalid request")
		return
	}

	if !dbTypes.IsValidFeedbackRating(dto.Rating) {
		response.RespondError(c, http.StatusBadRequest, "Invalid rating")
		return
	}

	user := c.MustGet("x-user").(*dbTypes.User)
	feedback := &dbTypes.Feedback{
		SubmitterCID:   user.CID,
		SubmitterName:  user.FirstName + " " + user.LastName,
		SubmitterEmail: user.Email,
		ControllerID:   dto.Controller,
		FlightCallsign: dto.FlightCallsign,
		FlightDate:     dto.FlightDate,
		Comments:       dto.Comments,
		Status:         dbTypes.FeedbackStatus["pending"],
		Rating:         dto.Rating,
	}

	if err := database.DB.Create(feedback).Error; err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// @TODO Send to rabbitmq for bot messaging
	response.RespondBlank(c, http.StatusCreated)
}

// Patch Pilot Feedback
// @Summary Patch Pilot Feedback
// @Description Patch feedback for a pilot -- currently only the status field can be patched
// @Tags Feedback
// @Param id path int true "Feedback ID"
// @Param data body dto.FeedbackPatchRequest true "Feedback"
// @Success 204
// @Failure 400 {object} response.R
// @Failure 403 {object} response.R
// @Failure 500 {object} response.R
func patchFeedback(c *gin.Context) {
	var dtoFeedback dto.FeedbackPatchRequest
	if err := c.ShouldBindJSON(&dtoFeedback); err != nil {
		response.RespondError(c, http.StatusBadRequest, "Invalid request")
		return
	}

	if !dbTypes.IsValidFeedbackStatus(dtoFeedback.Status) {
		response.RespondError(c, http.StatusBadRequest, "Invalid status")
		return
	}

	feedback := &dbTypes.Feedback{}
	if err := database.DB.Where("id = ?", c.Param("id")).First(feedback).Error; err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if err := database.DB.Model(feedback).Update("status", dtoFeedback.Status).Error; err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	response.RespondBlank(c, http.StatusNoContent)
}
