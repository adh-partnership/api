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

// Get Specific Training Request
// @Summary Get Specific Training Request
// @Tags training
// @Param id path string true "Training Session ID"
// @Success 200 {object} dto.TrainingSessionRequest
// @Failure 403 {object} response.R
// @Failure 404 {object} response.R
// @Router /v1/training/requests/{id} [GET]
func getTrainingRequest(c *gin.Context) {
	request, err := database.FindTrainingSessionRequestByID(c.Param("id"))
	if err != nil {
		response.RespondError(c, http.StatusNotFound, "Not Found")
		return
	}

	if !auth.InGroup(c.MustGet("x-user").(*models.User), "training") && request.StudentID != c.MustGet("x-user").(*models.User).CID {
		response.RespondError(c, http.StatusForbidden, "Forbidden")
		return
	}

	response.Respond(c, http.StatusOK, dto.ConvertTrainingRequestToDTO(request))
}

// Get Training Requests
// @Summary Get Training Sessions
// @Tags training
// @Param cid query string false "Filter by CID"
// @Param status query string false "Filter by Status"
// @Success 200 {array} []dto.TrainingSessionRequest
// @Failure 403 {object} response.R
// @Router /v1/training/requests [GET]
func getTrainingRequests(c *gin.Context) {
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
// @Param data body dto.TrainingRequestCreateRequest true "Training Session Request"
// @Success 201 {object} dto.TrainingSessionRequest
// @Failure 400 {object} response.R
// @Failure 401 {object} response.R
// @Failure 403 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/training/sessions [POST]
func postTrainingRequest(c *gin.Context) {
	var request dto.TrainingRequestCreateRequest
	if err := c.ShouldBind(&request); err != nil {
		response.RespondError(c, http.StatusBadRequest, "Bad Request")
		return
	}

	user := c.MustGet("x-user").(*models.User)

	req := &models.TrainingRequest{
		Student:   user,
		StudentID: user.CID,
		Position:  request.Position,
		Notes:     request.Notes,
		Status:    constants.TrainingSessionStatusOpen,
	}

	// Check if Position is valid
	if !models.IsValidPosition(req.Position) {
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

	if !areSlotsValid(request.Slots) {
		response.RespondError(c, http.StatusBadRequest, "Invalid Training Slots")
		return
	}

	for _, slot := range request.Slots {
		s := &models.TrainingRequestSlot{
			TrainingRequestID: req.ID,
			StartTime:         slot.StartTime,
			EndTime:           slot.EndTime,
		}

		if err := database.DB.Create(s).Error; err != nil {
			response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
			return
		}
	}

	go func(r *models.TrainingRequest) {
		u, _ := database.FindUserByCID(fmt.Sprint(r.Student.CID))
		_ = discord.NewMessage().
			SetContent("New Training Request has been submitted").
			AddEmbed(
				discord.NewEmbed().
					AddField(
						discord.NewField().SetName("Controller").SetValue(
							fmt.Sprintf("%s %s (%d/%s)", r.Student.FirstName, r.Student.LastName, r.Student.CID, u.Rating.Short),
						).SetInline(false),
					).
					AddField(
						discord.NewField().SetName("Position").SetValue(r.Position).SetInline(true),
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
// @Param data body dto.TrainingRequestEditRequest true "Training Session Request"
// @Success 200 {object} dto.TrainingSessionRequest
// @Failure 400 {object} response.R
// @Failure 401 {object} response.R
// @Failure 403 {object} response.R
// @Failure 404 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/training/requests/{id} [PATCH]
func patchTrainingRequest(c *gin.Context) {
	var request dto.TrainingRequestEditRequest
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

	if request.Position != "" && !models.IsValidPosition(request.Position) {
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
		if request.Start != "" {
			start, err := time.Parse(time.RFC3339, request.Start)
			if err != nil {
				response.RespondError(c, http.StatusBadRequest, "Invalid Start Time")
				return
			}
			req.Start = &start
		}

		if request.End != "" {
			end, err := time.Parse(time.RFC3339, request.End)
			if err != nil {
				response.RespondError(c, http.StatusBadRequest, "Invalid End Time")
				return
			}
			req.End = &end
		}

		if request.Status != "" && request.Status != req.Status {
			if request.Status == constants.TrainingSessionStatusAccepted && req.Status == constants.TrainingSessionStatusOpen {
				go func(r *models.TrainingRequest) {
					u, _ := database.FindUserByCID(fmt.Sprint(r.Student.CID))
					_ = discord.NewMessage().
						SetContent("Training Scheduled!").
						AddEmbed(
							discord.NewEmbed().
								SetColor(discord.GetColor("00", "00", "ff")).
								AddField(
									discord.NewField().SetName("Controller").SetValue(
										fmt.Sprintf("%s %s (%d/%s)", r.Student.FirstName, r.Student.LastName, r.Student.CID, u.Rating.Short),
									).SetInline(false),
								).
								AddField(
									discord.NewField().SetName("Position").SetValue(r.Position).SetInline(true),
								).
								AddField(
									discord.NewField().SetName("Scheduled At").SetValue(r.Start.Format(time.RFC1123)).SetInline(false),
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
			req.Start = nil
			req.End = nil
			go func(r *models.TrainingRequest) {
				u, _ := database.FindUserByCID(fmt.Sprint(r.Student.CID))
				_ = discord.NewMessage().
					SetContent("Scheduled Training Session set to Cancelled").
					AddEmbed(
						discord.NewEmbed().
							SetColor(discord.GetColor("ff", "00", "00")).
							AddField(
								discord.NewField().SetName("Controller").SetValue(
									fmt.Sprintf("%s %s (%d/%s)", r.Student.FirstName, r.Student.LastName, r.Student.CID, u.Rating.Short),
								).SetInline(false),
							).
							AddField(
								discord.NewField().SetName("Position").SetValue(r.Position).SetInline(true),
							).
							AddField(
								discord.NewField().SetName("Scheduled At").SetValue(r.Start.Format(time.RFC1123)).SetInline(false),
							).
							AddField(
								discord.NewField().SetName("Notes").SetValue(r.Notes).SetInline(false),
							).SetURL(fmt.Sprintf("%s/training/sessions/%s", config.Cfg.Facility.FrontendURL, r.ID)),
					).Send(config.Cfg.Facility.TrainingRequests.Discord.TrainingStaff)
			}(req)
		}

		req.Status = request.Status
	}

	if request.Position != "" && req.Position != request.Position {
		req.Position = request.Position
	}

	if request.Notes != "" && req.Notes != request.Notes {
		req.Notes = request.Notes
	}

	if req.InstructorID != &request.Instructor {
		ins, err := database.FindUserByCID(fmt.Sprint(request.Instructor))
		if err != nil {
			response.RespondError(c, http.StatusBadRequest, "Invalid Instructor ID")
			return
		}
		req.Instructor = ins
		req.InstructorID = &request.Instructor
	}

	if err := database.DB.Save(req).Error; err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	response.Respond(c, http.StatusOK, dto.ConvertTrainingRequestToDTO(req))
}

// Add a slot to a request
// @Summary Add a slot to a request
// @Tags training
// @Param id path string true "Training Session Request ID"
// @Param data body dto.TrainingRequestSlot true "Training Session Request Slot"
// @Success 200 {object} dto.TrainingSessionRequest
// @Failure 400 {object} response.R
// @Failure 401 {object} response.R
// @Failure 403 {object} response.R
// @Failure 404 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/training/requests/{id}/slots [POST]
func postTrainingRequestSlot(c *gin.Context) {
	var request dto.TrainingRequestSlot
	if err := c.ShouldBind(&request); err != nil {
		response.RespondError(c, http.StatusBadRequest, "Bad Request")
		return
	}

	// Check if start time or end time is nil
	if request.StartTime == nil || request.EndTime == nil {
		response.RespondError(c, http.StatusBadRequest, "Bad Request")
		return
	}

	trainingrequest, err := database.FindTrainingSessionRequestByID(c.Param("id"))
	if err != nil {
		response.RespondError(c, http.StatusNotFound, "Training Session Request not found")
		return
	}

	if trainingrequest.StudentID != c.MustGet("x-user").(*models.User).CID {
		response.RespondError(c, http.StatusForbidden, "Forbidden")
		return
	}

	trainingrequest.Slots = append(trainingrequest.Slots, &models.TrainingRequestSlot{
		StartTime: request.StartTime,
		EndTime:   request.EndTime,
	})

	trd := dto.ConvertTrainingRequestToDTO(trainingrequest)
	if !areSlotsValid(trd.Slots) {
		response.RespondError(c, http.StatusBadRequest, "Invalid Slots")
		return
	}

	if err := database.DB.Save(trainingrequest).Error; err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	response.Respond(c, http.StatusOK, trd)
}

// Delete a slot from a request
// @Summary Delete a slot from a request
// @Tags training
// @Param id path string true "Training Session Request ID"
// @Param slot_id path string true "Training Session Request Slot ID"
// @Success 200 {object} dto.TrainingSessionRequest
// @Failure 400 {object} response.R
// @Failure 401 {object} response.R
// @Failure 403 {object} response.R
// @Failure 404 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/training/requests/{id}/slots/{slot_id} [DELETE]
func deleteTrainingRequestSlot(c *gin.Context) {
	var request dto.TrainingRequestSlot
	if err := c.ShouldBind(&request); err != nil {
		response.RespondError(c, http.StatusBadRequest, "Bad Request")
		return
	}

	// Check if start time or end time is nil
	if request.StartTime == nil || request.EndTime == nil {
		response.RespondError(c, http.StatusBadRequest, "Bad Request")
		return
	}

	trainingrequest, err := database.FindTrainingSessionRequestByID(c.Param("id"))
	if err != nil {
		response.RespondError(c, http.StatusNotFound, "Training Session Request not found")
		return
	}

	if trainingrequest.StudentID != c.MustGet("x-user").(*models.User).CID {
		response.RespondError(c, http.StatusForbidden, "Forbidden")
		return
	}

	found := false
	for _, slot := range trainingrequest.Slots {
		if slot.ID.String() == c.Param("slot_id") {
			found = true
		}
	}

	if !found {
		response.RespondError(c, http.StatusNotFound, "Slot not found")
		return
	}

	if err := database.DB.Delete(request.ID).Error; err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	for _, slot := range trainingrequest.Slots {
		if slot.ID.String() == c.Param("slot_id") {
			slot = nil
		}
	}

	response.Respond(c, http.StatusOK, dto.ConvertTrainingRequestToDTO(trainingrequest))
}

func areSlotsValid(slots []*dto.TrainingRequestSlot) bool {
	validSlots := true
	firstDate := slots[0].StartTime
	endDate := slots[len(slots)-1].EndTime // Just filler data

	/* Rules for valid slots:
	 *
	 * - The earliest slot must be at least 24 hours from now
	 * - The latest slot cannot start more than 14 days from the earliest slot
	 * - No single slot can be longer than 24 hours
	 * - No single slot can be less than 1 hour long
	 */
	for _, slot := range slots {
		if slot.StartTime.Before(time.Now().Add(time.Hour * 24)) {
			validSlots = false
			break
		}
		if slot.EndTime.Sub(*slot.StartTime) > time.Hour*24 || slot.EndTime.Sub(*slot.StartTime) < time.Hour {
			validSlots = false
			break
		}
		// Check if end is before start
		if slot.EndTime.Before(*slot.StartTime) {
			validSlots = false
			break
		}
		// Update firstDate if this slot's start time is before it
		if slot.StartTime.Before(*firstDate) {
			firstDate = slot.StartTime
		}
		// Update endDate if this slot's end time is after it
		if slot.StartTime.After(*endDate) {
			endDate = slot.StartTime
		}
	}
	// Check if endDate is within 14 days of firstDate
	if endDate.Sub(*firstDate) > time.Hour*24*14 {
		validSlots = false
	}

	return validSlots
}
