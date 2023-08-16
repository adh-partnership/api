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

package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"

	"github.com/adh-partnership/api/pkg/database"
	"github.com/adh-partnership/api/pkg/database/dto"
	"github.com/adh-partnership/api/pkg/database/models"
	"github.com/adh-partnership/api/pkg/database/models/constants"
	"github.com/adh-partnership/api/pkg/gin/response"
)

// Get Full Roster
// @Summary Get Full Roster
// @Tags user
// @Success 200 {object} []dto.UserResponse
// @Failure 403 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/user/all [GET]
func getFullRoster(c *gin.Context) {
	users := []models.User{}
	ret := []*dto.UserResponse{}

	if err := database.DB.Preload(clause.Associations).Find(&users).Error; err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	for _, user := range users {
		ret = append(ret, dto.ConvUserToUserResponse(&user))
	}

	response.Respond(c, http.StatusOK, ret)
}

// Get Facility Roster
// @Summary Get Facility Roster
// @Tags user
// @Success 200 {object} []dto.UserResponse
// @Failure 403 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/user/roster [GET]
func getRoster(c *gin.Context) {
	users := []models.User{}
	ret := []*dto.UserResponse{}

	if err := database.DB.Preload(clause.Associations).Not(&models.User{ControllerType: constants.ControllerTypeNone}).Find(&users).Error; err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	for _, user := range users {
		ret = append(ret, dto.ConvUserToUserResponse(&user))
	}

	response.Respond(c, http.StatusOK, ret)
}

// Get Facility Staff
// @Summary Get Facility Staff
// @Tags user
// @Success 200 {object} []dto.FacilityStaffResponse
// @Failure 403 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/user/staff [GET]
func getStaff(c *gin.Context) {
	staff, err := dto.GetStaffResponse()
	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	response.Respond(c, http.StatusOK, staff)
}
