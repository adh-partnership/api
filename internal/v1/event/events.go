package event

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/adh-partnership/api/pkg/database"
	"github.com/adh-partnership/api/pkg/database/dto"
	"github.com/adh-partnership/api/pkg/database/models"
	"github.com/adh-partnership/api/pkg/gin/response"
)

// Get Events
// @Summary Get Events
// @Description Get Upcoming Events
// @Tags Events
// @Param limit query number false "Limit to X events, default 5 (max 10)"
// @Success 200 {object} []dto.EventsResponse
// @Failure 500 {object} response.R
// @Router /v1/events [get]
func getEvents(c *gin.Context) {
	var limit uint
	if c.Query("limit") == "" {
		limit = 5
	} else {
		limit = database.Atou(c.Query("limit"))
		if limit > 10 {
			limit = 10
		}
	}

	events, err := database.GetEvents(int(limit))
	if err != nil {
		log.Errorf("Error getting events: %s", err)
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	response.Respond(c, http.StatusOK, dto.ConvEventsToEventsResponse(events))
}

// Get Event
// @Summary Get Event
// @Description Get an event
// @Tags Events
// @Param id path string true "Event ID"
// @Success 200 {object} dto.EventsResponse
// @Failure 404 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/events/{id} [get]
func getEvent(c *gin.Context) {
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

	response.Respond(c, http.StatusOK, dto.ConvEventToEventsResponse(event))
}

// Create Event
// @Summary Create Event
// @Description Create an event
// @Tags Events
// @Param data body dto.EventRequest true "Event Data"
// @Success 201
// @Failure 400 {object} response.R
// @Failure 403 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/events [post]
func postEvent(c *gin.Context) {
	var dto dto.EventRequest
	if err := c.ShouldBindJSON(&dto); err != nil {
		log.Debugf("Error binding dto: %s", err)
		response.RespondError(c, http.StatusBadRequest, "Invalid request")
		return
	}

	event := models.Event{
		Title:       dto.Title,
		Description: dto.Description,
		Banner:      dto.Banner,
		StartDate:   *dto.StartDate,
		EndDate:     *dto.EndDate,
	}

	if err := database.DB.Create(&event).Error; err != nil {
		log.Errorf("Error creating event: %s", err)
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	response.Respond(c, http.StatusCreated, event)
}

// Patch Event
// @Summary Patch Event
// @Description Patch an event
// @Tags Events
// @Param id path string true "Event ID"
// @Param data body dto.EventRequest true "Event Data"
// @Success 200 {object} models.Event
// @Failure 400 {object} response.R
// @Failure 403 {object} response.R
// @Failure 404 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/events/{id} [patch]
func patchEvent(c *gin.Context) {
	var data dto.EventRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		log.Debugf("Error binding dto: %s", err)
		response.RespondError(c, http.StatusBadRequest, "Invalid request")
		return
	}

	event, err := database.GetEvent(c.Param("id"))
	if err != nil {
		log.Errorf("Error getting event: %s", err)
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if event == nil {
		response.RespondError(c, http.StatusNotFound, "Event not found")
		return
	}

	patchedEvent := dto.PatchEventRequest(event, data)

	if err := database.DB.Save(patchedEvent).Error; err != nil {
		log.Errorf("Error updating event: %s", err)
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	response.Respond(c, http.StatusOK, patchedEvent)
}

// Delete Event
// @Summary Delete Event
// @Description Delete an event
// @Tags Events
// @Param id path string true "Event ID"
// @Success 204
// @Failure 403 {object} response.R
// @Failure 404 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/events/{id} [delete]
func deleteEvent(c *gin.Context) {
	event, err := database.GetEvent(c.Param("id"))
	if err != nil {
		log.Errorf("Error getting event: %s", err)
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if event == nil {
		response.RespondError(c, http.StatusNotFound, "Event not found")
		return
	}

	if err := database.DB.Delete(event.Positions).Error; err != nil {
		log.Errorf("Error deleting positions: %s", err)
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if err := database.DB.Delete(event.Signups).Error; err != nil {
		log.Errorf("Error deleting signups: %s", err)
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if err := database.DB.Delete(event).Error; err != nil {
		log.Errorf("Error deleting event: %s", err)
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	response.Respond(c, http.StatusNoContent, nil)
}
