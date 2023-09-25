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
	"net/http"

	"github.com/adh-partnership/api/pkg/database"
	"github.com/adh-partnership/api/pkg/database/dto"
	"github.com/adh-partnership/api/pkg/database/models"
	"github.com/adh-partnership/api/pkg/gin/response"
	"github.com/gin-gonic/gin"
)

// Get event tracking stats
// @Summary Get event tracking stats
// @Description Get event tracking stats for a user
// @Tags Events
// @param id path string true "CID"
// @Success 200 {object} dto.EventsResponse
// @Failure 404 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/events/user/{id}/stats [get]
func getEventTracking(c *gin.Context) {
	cid := uint(database.Atoi(c.Param("id")))
	data, err := database.GetEventTracking(cid)
	if err != nil {
		log.Errorf("Error getting event tracking for %d: %s", cid, err)
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if data == nil {
		blank := &models.EventTracking{}
		response.Respond(c, http.StatusOK, blank)
		return
	}

	response.Respond(c, http.StatusOK, data)
}

// Update event tracking stats
// @Summary Update event tracking stats
// @Description Update event tracking stats for a user
// @Tags Events
// @param id path string true "CID"
// @Param data body dto.EventStatsRequest true "Stats data"
// @Success 201
// @Failure 400 {object} response.R
// @Failure 403 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/events/user/{id}/stats [put]
func updateEventTracking(c *gin.Context) {
	cid := uint(database.Atoi(c.Param("id")))
	record, err := database.GetEventTracking(cid)
	if err != nil {
		log.Errorf("Error getting event tracking for %s: %s", c.Param("id"), err)
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if record == nil {
		record = &models.EventTracking{CID: cid}
	}

	data := &dto.EventStatsRequest{}
	if err := c.ShouldBind(&data); err != nil {
		response.RespondError(c, http.StatusBadRequest, "Invalid request")
		return
	}

	record.Tardies = data.Tardies
	record.NoShows = data.NoShows

	if err := database.DB.Save(&record).Error; err != nil {
		log.Errorf("Error updating events stats for %d: %s", cid, err)
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	response.RespondBlank(c, http.StatusOK)
}
