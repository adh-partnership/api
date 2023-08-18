/*
 * Copyright ADH Partnership
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package event

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/adh-partnership/api/pkg/database"
	"github.com/adh-partnership/api/pkg/database/dto"
	"github.com/adh-partnership/api/pkg/database/models"
	"github.com/adh-partnership/api/pkg/gin/response"
)

// Get Event Positions
// @Summary Get Event Positions
// @Description Get Positions for an event
// @Tags Events
// @Param id path string true "Event ID"
// @Success 200 {object} models.EventPosition[]
// @Failure 404 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/events/{id}/positions [get]
func getEventPositions(c *gin.Context) {
	event, err := database.GetEvent(c.Param("id"))
	if err != nil {
		log.Errorf("Error getting event positions: %s", err)
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if event == nil {
		response.RespondError(c, http.StatusNotFound, "Not Found")
		return
	}

	response.Respond(c, http.StatusOK, event.Positions)
}

// Add Event Position
// @Summary Add Event Position
// @Description Add a position to an event
// @Tags Events
// @Param id path string true "Event ID"
// @Param position body dto.EventPositionRequest true "Position. CID 0 means unassigned."
// @Success 200 {object} models.Event
// @Failure 400 {object} response.R
// @Failure 404 {object} response.R
// @Failure 409 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/events/{id}/positions [post]
func addEventPosition(c *gin.Context) {
	event, err := database.GetEvent(c.Param("id"))
	if err != nil {
		log.Errorf("Error getting event: %s", err)
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	data := &dto.EventPositionRequest{}
	if err := c.ShouldBind(&data); err != nil || data.Position == "" {
		response.RespondError(c, http.StatusBadRequest, "Bad Request")
		return
	}

	if event == nil {
		response.RespondError(c, http.StatusNotFound, "Not Found")
		return
	}

	// Check if position already exists
	for _, position := range event.Positions {
		if position.Position == data.Position {
			response.RespondError(c, http.StatusConflict, "Position already exists")
			return
		}
	}

	var user *models.User
	if data.UserID != 0 {
		user, err = database.FindUserByCID(fmt.Sprint(data.UserID))
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Errorf("Error getting user: %s", err)
			response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
			return
		}
	}

	position := &models.EventPosition{
		Position: data.Position,
		User:     user,
	}

	event.Positions = append(event.Positions, position)
	if err := database.DB.Save(&event).Error; err != nil {
		log.Errorf("Error adding event position: %s", err)
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	response.Respond(c, http.StatusOK, event)
}

// Update Event Position
// @Summary Update Event Position
// @Description Update a position for an event.
// @Tags Events
// @Param id path string true "Event ID"
// @Param position path string true "Position Name"
// @Param position body dto.EventPositionRequest true "Position. CID 0 means unassigned."
// @Success 200 {object} models.Event
// @Failure 400 {object} response.R
// @Failure 404 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/events/{id}/positions/{position} [put]
func updateEventPosition(c *gin.Context) {
	event, err := database.GetEvent(c.Param("id"))
	if err != nil {
		log.Errorf("Error getting event: %s", err)
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	data := &dto.EventPositionRequest{}
	if err := c.ShouldBind(&data); err != nil {
		response.RespondError(c, http.StatusBadRequest, "Bad Request")
		return
	}

	if event == nil {
		response.RespondError(c, http.StatusNotFound, "Not Found")
		return
	}

	var cid *uint
	var user *models.User
	if data.UserID != 0 {
		user, err = database.FindUserByCID(fmt.Sprint(data.UserID))
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Errorf("Error getting user: %s", err)
			response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		cid = &user.CID
	}

	for _, position := range event.Positions {
		if position.Position == c.Param("position") {
			// Never let position be "", assume they meant to keep the same position name
			if data.Position == "" {
				data.Position = position.Position
			}
			position.Position = data.Position
			position.User = user
			position.UserID = cid
			if err := database.DB.Save(&position).Error; err != nil {
				log.Errorf("Error updating event position: %s", err)
				response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
				return
			}
			response.Respond(c, http.StatusOK, event)
			return
		}
	}

	response.RespondError(c, http.StatusNotFound, "Not Found")
}

// Delete Event Position
// @Summary Delete Event Position
// @Description Delete a position from an event
// @Tags Events
// @Param id path string true "Event ID"
// @Param position path string true "Position Name"
// @Success 204
// @Failure 404 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/events/{id}/positions/{position} [delete]
func deleteEventPosition(c *gin.Context) {
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

	for _, position := range event.Positions {
		if position.Position == c.Param("position") {
			if err := database.DB.Delete(position).Error; err != nil {
				log.Errorf("Error deleting event position: %s", err)
				response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
				return
			}
			response.RespondBlank(c, http.StatusNoContent)
			return
		}
	}

	response.RespondError(c, http.StatusNotFound, "Not Found")
}
