package training

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/adh-partnership/api/pkg/auth"
	"github.com/adh-partnership/api/pkg/config"
	"github.com/adh-partnership/api/pkg/database"
	"github.com/adh-partnership/api/pkg/database/dto"
	"github.com/adh-partnership/api/pkg/database/models"
	"github.com/adh-partnership/api/pkg/database/models/constants"
	"github.com/adh-partnership/api/pkg/discord"
	"github.com/adh-partnership/api/pkg/gin/response"
)

// Get Training Sessions
// @Summary Get Training Sessions
// @Tags training
// @Param cid query string false "Filter by CID"
// @Param status query string false "Filter by Status"
// @Success 200 {array} []dto.TrainingSessionRequest
// @Failure 403 {object} response.R
// @Router /v1/training/sessions [GET]
func getSessions(c *gin.Context) {
	user := c.MustGet("x-user").(*models.User)
	filter := &database.TrainingSessionRequestFilter{}

	// Only training staff should be able to see all training requests
	// So lock the filter to their CID unless in the training group,
	// where we can then respect the value of the cid query param
	if auth.InGroup(user, "training") {
		filter.CID = fmt.Sprint(user.CID)
	} else {
		filter.CID = c.Query("cid")
	}
	filter.Status = c.Query("status")

	requests, err := database.FindTrainingSessionRequestWithFilter(filter)
	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	response.Respond(c, http.StatusOK, dto.ConvertTrainingRequestsToDTO(requests))
}

// Create new training session request for user
// @Summary Create new training session request for user
// @Tags training
// @Param data body dto.TrainingSessionRequestCreateRequest true "Training Session Request"
// @Success 201 {object} dto.TrainingSessionRequest
// @Failure 400 {object} response.R
// @Failure 401 {object} response.R
// @Failure 403 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/training/sessions [POST]
func postSessions(c *gin.Context) {
	var request dto.TrainingSessionRequestCreateRequest
	if err := c.ShouldBind(&request); err != nil {
		response.RespondError(c, http.StatusBadRequest, "Bad Request")
		return
	}

	user := c.MustGet("x-user").(*models.User)

	req := &models.TrainingSessionRequest{
		User:         user,
		TrainingType: request.TrainingType,
		TrainingFor:  request.TrainingFor,
		Notes:        request.Notes,
		Status:       constants.TrainingSessionStatusOpen,
	}

	if err := database.DB.Create(req).Error; err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	go func(r *models.TrainingSessionRequest) {
		u, _ := database.FindUserByCID(fmt.Sprint(r.User.CID))
		_ = discord.NewMessage().
			SetContent("New Training Request has been submitted").
			AddEmbed(
				discord.NewEmbed().
					AddField(
						discord.NewField().SetName("Controller").SetValue(
							fmt.Sprintf("%s %s (%d/%s)", r.User.FirstName, r.User.LastName, r.User.CID, u.Rating.Short),
						).SetInline(false),
					).
					AddField(
						discord.NewField().SetName("Type").SetValue(r.TrainingType).SetInline(true),
					).
					AddField(
						discord.NewField().SetName("Position").SetValue(r.TrainingFor).SetInline(true),
					).
					AddField(
						discord.NewField().SetName("Notes").SetValue(r.Notes).SetInline(false),
					),
			).Send(config.Cfg.Facility.TrainingRequests.Discord.Scheduled)
	}(req)

	response.Respond(c, http.StatusCreated, dto.ConvertTrainingRequestToDTO(req))
}