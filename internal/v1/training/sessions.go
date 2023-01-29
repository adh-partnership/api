package training

import (
	"fmt"
	"net/http"
	"time"

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

// Get Specific Training Session
// @Summary Get Specific Training Session
// @Tags training
// @Param id path string true "Training Session ID"
// @Success 200 {object} dto.TrainingSessionRequest
// @Failure 403 {object} response.R
// @Failure 404 {object} response.R
// @Router /v1/training/sessions/{id} [GET]
func getSession(c *gin.Context) {
	request, err := database.FindTrainingSessionRequestByID(c.Param("id"))
	if err != nil {
		response.RespondError(c, http.StatusNotFound, "Not Found")
		return
	}

	if !auth.InGroup(c.MustGet("x-user").(*models.User), "training") && request.User.CID != c.MustGet("x-user").(*models.User).CID {
		response.RespondError(c, http.StatusForbidden, "Forbidden")
		return
	}

	response.Respond(c, http.StatusOK, dto.ConvertTrainingRequestToDTO(request))
}

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

	// Check if TrainingFor is valid
	if !models.IsValidPosition(req.TrainingFor) {
		response.RespondError(c, http.StatusBadRequest, "Invalid Training Position")
		return
	}

	if !models.IsValidTrainingStatus(req.Status) {
		response.RespondError(c, http.StatusBadRequest, "Invalid Training Type")
		return
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
					).SetURL(fmt.Sprintf("%s/training/sessions/%s", config.Cfg.Facility.FrontendURL, r.ID)),
			).Send(config.Cfg.Facility.TrainingRequests.Discord.Scheduled)
	}(req)

	response.Respond(c, http.StatusCreated, dto.ConvertTrainingRequestToDTO(req))
}

// Edit Training Session Request
// @Summary Edit Training Session Request
// @Tags training
// @Param id path string true "Training Session Request ID"
// @Param data body dto.TrainingSessionRequestEditRequest true "Training Session Request"
// @Success 200 {object} dto.TrainingSessionRequest
// @Failure 400 {object} response.R
// @Failure 401 {object} response.R
// @Failure 403 {object} response.R
// @Failure 404 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/training/sessions/{id} [PATCH]
func patchSession(c *gin.Context) {
	var request dto.TrainingSessionRequestEditRequest
	if err := c.ShouldBind(&request); err != nil {
		response.RespondError(c, http.StatusBadRequest, "Bad Request")
		return
	}

	user := c.MustGet("x-user").(*models.User)
	id := c.Param("id")

	req, err := database.FindTrainingSessionRequestByID(id)
	if err != nil {
		response.RespondError(c, http.StatusNotFound, "Not Found")
		return
	}

	if request.TrainingFor != "" && !models.IsValidPosition(request.TrainingFor) {
		response.RespondError(c, http.StatusBadRequest, "Invalid Training Position")
		return
	}

	if request.Status != "" && !models.IsValidTrainingStatus(request.Status) {
		response.RespondError(c, http.StatusBadRequest, "Invalid Status")
		return
	}

	if auth.InGroup(user, "training") {
		if request.InstructorNotes != "" {
			req.InstructorNotes = request.InstructorNotes
		}
		if request.Scheduled != nil {
			req.ScheduledSession = request.Scheduled
		}

		if request.Status != "" && request.Status != req.Status {
			if request.Status == constants.TrainingSessionStatusAccepted && req.Status == constants.TrainingSessionStatusOpen {
				go func(r *models.TrainingSessionRequest) {
					u, _ := database.FindUserByCID(fmt.Sprint(r.User.CID))
					_ = discord.NewMessage().
						SetContent("Training Scheduled!").
						AddEmbed(
							discord.NewEmbed().
								SetColor(discord.GetColor("00", "00", "ff")).
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
									discord.NewField().SetName("Scheduled At").SetValue(r.ScheduledSession.Format(time.RFC1123)).SetInline(false),
								),
						).Send(config.Cfg.Facility.TrainingRequests.Discord.TrainingStaff)
				}(req)
			}
			req.Status = request.Status
		}
	}

	if request.Status != "" && req.Status != request.Status {
		if request.Status != constants.TrainingSessionStatusOpen && request.Status != constants.TrainingSessionStatusCancelled {
			response.RespondError(c, http.StatusBadRequest, "Invalid Status")
			return
		}
		if request.Status == constants.TrainingSessionStatusCancelled && req.Status == constants.TrainingSessionStatusAccepted {
			req.ScheduledSession = nil
			go func(r *models.TrainingSessionRequest) {
				u, _ := database.FindUserByCID(fmt.Sprint(r.User.CID))
				_ = discord.NewMessage().
					SetContent("Scheduled Training Session set to Cancelled").
					AddEmbed(
						discord.NewEmbed().
							SetColor(discord.GetColor("ff", "00", "00")).
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
								discord.NewField().SetName("Scheduled At").SetValue(r.ScheduledSession.Format(time.RFC1123)).SetInline(false),
							).
							AddField(
								discord.NewField().SetName("Notes").SetValue(r.Notes).SetInline(false),
							).SetURL(fmt.Sprintf("%s/training/sessions/%s", config.Cfg.Facility.FrontendURL, r.ID)),
					).Send(config.Cfg.Facility.TrainingRequests.Discord.TrainingStaff)
			}(req)
		}

		req.Status = request.Status
	}

	if request.TrainingType != "" && req.TrainingType != request.TrainingType {
		req.TrainingType = request.TrainingType
	}

	if request.TrainingFor != "" && req.TrainingFor != request.TrainingFor {
		req.TrainingFor = request.TrainingFor
	}

	if request.Notes != "" && req.Notes != request.Notes {
		req.Notes = request.Notes
	}

	if request.InstructorID > 0 && *req.InstructorID != request.InstructorID {
		req.InstructorID = &request.InstructorID
	}

	if err := database.DB.Save(req).Error; err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	response.Respond(c, http.StatusOK, dto.ConvertTrainingRequestToDTO(req))
}
