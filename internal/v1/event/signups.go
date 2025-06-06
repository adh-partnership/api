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

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

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
	if err := c.ShouldBind(&data); err != nil {
		response.RespondError(c, http.StatusBadRequest, "Invalid request")
		return
	}

	user := c.MustGet("x-user").(*models.User)

	signup := &models.EventSignup{}
	if err := database.DB.Where("event_id = ? AND user_id = ?", event.ID, user.CID).First(signup).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			signup := &models.EventSignup{
				EventID: event.ID,
				User:    user,
				Choice1: data.Choice1,
				Choice2: data.Choice2,
				Choice3: data.Choice3,
				Notes:   data.Notes,
			}
			if err := database.DB.Create(signup).Error; err != nil {
				log.Errorf("Error creating event signup: %s", err)
				response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
				return
			}
		} else {
			log.Errorf("Error getting event signup: %s", err)
			response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
			return
		}
	} else {
		signup.Choice1 = data.Choice1
		signup.Choice2 = data.Choice2
		signup.Choice3 = data.Choice3
		signup.Notes = data.Notes

		if err := database.DB.Save(&signup).Error; err != nil {
			log.Errorf("Error updating event signup: %s", err)
			response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
			return
		}
	}

	event, _ = database.GetEvent(c.Param("id"))
	response.Respond(c, http.StatusOK, dto.ConvEventToEventsResponse(event))
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
