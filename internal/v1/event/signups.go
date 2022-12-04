package event

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/adh-partnership/api/pkg/database"
	"github.com/adh-partnership/api/pkg/database/dto"
	"github.com/adh-partnership/api/pkg/database/models"
	"github.com/adh-partnership/api/pkg/gin/response"
)

// Create/Edit User Signup for Event
// @Summary Create/Edit User Signup for Event
// @Description Create/Edit User Signup for Event. This will only work for the logged in user.
// @Tags Events
// @Param id path string true "Event ID"
// @Param signup body dto.EventSignupRequest true "Signup"
// @Success 200 {object} models.Event
// @Failure 400 {object} response.R
// @Failure 403 {object} response.R
// @Failure 404 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/events/{id}/signup [post]
func postEventSignup(c *gin.Context) {
	event, err := database.GetEvent(c.Param("id"))
	if err != nil {
		log.Errorf("Error getting event: %s", err)
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if event == nil {
		response.RespondError(c, http.StatusNotFound, "Not Found")
		return
	}

	data := &dto.EventSignupRequest{}
	if err := c.ShouldBindJSON(&data); err != nil {
		response.RespondError(c, http.StatusBadRequest, "Invalid request")
		return
	}

	user := c.MustGet("x-user").(*models.User)

	signup := &models.EventSignup{
		EventID: event.ID,
		User:    user,
		Choice1: data.Choice1,
		Choice2: data.Choice2,
		Choice3: data.Choice3,
		Notes:   data.Notes,
	}

	log.Debugf("Signups: %+v", event.Signups)

	found := false
	for i := range event.Signups {
		if event.Signups[i].User.CID == user.CID {
			event.Signups[i] = signup
			found = true
		}
	}

	log.Debugf("found: %t", found)

	if !found {
		event.Signups = append(event.Signups, signup)
	}

	if err := database.DB.Save(&event).Error; err != nil {
		log.Errorf("Error creating event signup: %s", err)
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	response.Respond(c, http.StatusOK, event)
}

// Delete User Signup
// @Summary Delete User Signup
// @Description Delete User Signup. This will only work for the logged in user.
// @Tags Events
// @Param id path string true "Event ID"
// @Success 204
// @Failure 403 {object} response.R
// @Failure 404 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/events/{id}/signup [delete]
func deleteEventSignup(c *gin.Context) {
	event, err := database.GetEvent(c.Param("id"))
	if err != nil {
		log.Errorf("Error getting event: %s", err)
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if event == nil {
		response.RespondError(c, http.StatusNotFound, "Not Found")
		return
	}

	user := c.MustGet("x-user").(*models.User)

	found := false
	for i := range event.Signups {
		if *event.Signups[i].UserID == user.CID {
			if err := database.DB.Delete(event.Signups[i]).Error; err != nil {
				log.Errorf("Error deleting event signup: %s", err)
				response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
				return
			}
			found = true
		}
	}

	if !found {
		response.RespondError(c, http.StatusNotFound, "Not Found")
		return
	}

	response.RespondBlank(c, http.StatusNoContent)
}
