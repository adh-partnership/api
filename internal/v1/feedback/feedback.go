package feedback

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"

	"github.com/adh-partnership/api/pkg/auth"
	"github.com/adh-partnership/api/pkg/database"
	"github.com/adh-partnership/api/pkg/database/dto"
	"github.com/adh-partnership/api/pkg/database/models"
	"github.com/adh-partnership/api/pkg/database/models/constants"
	"github.com/adh-partnership/api/pkg/discord"
	"github.com/adh-partnership/api/pkg/gin/response"
)

// Get Pilot Feedback
// @Summary Get Pilot Feedback
// @Description Get pilot feedback
// @Tags Feedback
// @Param cid query string false "Controller CID filter"
// @Param status query string false "Status filter, valid values: pending, approved, rejected. Default: pending"
// @Success 200 {object} []dto.FeedbackResponse
// @Failure 400 {object} response.R
// @Failure 403 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/feedback [get]
func getFeedback(c *gin.Context) {
	var feedback []*models.Feedback

	res := database.DB
	if c.Query("cid") != "" {
		res = res.Where("controller_id = ?", c.Query("cid"))
	}
	if c.Query("status") != "" {
		res = res.Where("status = ?", c.Query("status"))
	} else {
		res = res.Where("status = 'pending'")
	}
	if err := res.Preload(clause.Associations).Find(&feedback).Error; err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	includeEmail := false
	if auth.InGroup(c.MustGet("x-user").(*models.User), "admin") {
		includeEmail = true
	}

	response.Respond(c, http.StatusOK, dto.ConvertFeedbacktoResponse(feedback, includeEmail))
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
// @Router /v1/feedback [post]
func postFeedback(c *gin.Context) {
	var dto dto.FeedbackRequest
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.RespondError(c, http.StatusBadRequest, "Invalid request")
		return
	}

	if !models.IsValidFeedbackRating(dto.Rating) {
		response.RespondError(
			c,
			http.StatusBadRequest,
			fmt.Sprintf("Invalid rating. Valid values: %s, %s, %s, %s",
				constants.FeedbackRatingExcellent,
				constants.FeedbackRatingGood,
				constants.FeedbackRatingFair,
				constants.FeedbackRatingPoor,
			),
		)
		return
	}

	user := c.MustGet("x-user").(*models.User)
	var controller *models.User
	var err error
	if dto.Controller != "" {
		controller, err = database.FindUserByCID(dto.Controller)
		if err != nil {
			response.RespondError(c, http.StatusBadRequest, "Invalid controller")
			return
		}
	}
	feedback := &models.Feedback{
		Submitter:    user,
		Controller:   controller,
		Rating:       dto.Rating,
		Comments:     dto.Comments,
		Position:     dto.Position,
		Callsign:     dto.Callsign,
		Status:       constants.FeedbackStatusPending,
		ContactEmail: user.Email,
	}

	if err := database.DB.Create(feedback).Error; err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	_ = discord.SendWebhookMessage(
		"pending_feedback",
		"Web API",
		fmt.Sprintf(
			"New feedback submitted for %s on %s by %s (%d) - Rating: %s Comments: %s",
			feedback.Controller.FirstName+" "+feedback.Controller.LastName,
			feedback.Position,
			feedback.Submitter.FirstName+" "+feedback.Submitter.LastName,
			feedback.Submitter.CID,
			feedback.Rating,
			feedback.Comments,
		),
	)

	response.RespondBlank(c, http.StatusNoContent)
}

// Patch Pilot Feedback
// @Summary Patch Pilot Feedback
// @Description Patch feedback for a pilot -- currently only the status field can be patched. Accepted values: approved, rejected
// @Tags Feedback
// @Param id path int true "Feedback ID"
// @Param data body dto.FeedbackPatchRequest true "Feedback"
// @Success 204
// @Failure 400 {object} response.R
// @Failure 403 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/feedback/{id} [patch]
func patchFeedback(c *gin.Context) {
	var dtoFeedback dto.FeedbackPatchRequest
	if err := c.ShouldBindJSON(&dtoFeedback); err != nil {
		response.RespondError(c, http.StatusBadRequest, "Invalid request")
		return
	}

	if !models.IsValidFeedbackStatus(dtoFeedback.Status) {
		response.RespondError(c, http.StatusBadRequest, "Invalid status")
		return
	}

	feedback := &models.Feedback{}
	if err := database.DB.Preload(clause.Associations).Where("id = ?", c.Param("id")).First(feedback).Error; err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if err := database.DB.Model(feedback).Update("status", dtoFeedback.Status).Error; err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if shouldBroadcastFeedback(feedback) {
		_ = discord.SendWebhookMessage(
			"broadcast_feedback",
			"Web API",
			fmt.Sprintf(
				"New feedback received! Controller %s on %s received a %s rating. Comments: %s",
				feedback.Controller.FirstName+" "+feedback.Controller.LastName,
				feedback.Position,
				feedback.Rating,
				feedback.Comments,
			),
		)
	}

	response.RespondBlank(c, http.StatusNoContent)
}

// shouldBroadcastFeedback returns true if the feedback should be broadcast to the public on approval
func shouldBroadcastFeedback(feedback *models.Feedback) bool {
	if feedback.Status != constants.FeedbackStatusApproved {
		return false
	}

	return feedback.Rating == constants.FeedbackRatingExcellent || feedback.Rating == constants.FeedbackRatingGood
}
