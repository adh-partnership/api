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

package stats

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"

	"github.com/adh-partnership/api/pkg/database"
	"github.com/adh-partnership/api/pkg/database/dto"
	"github.com/adh-partnership/api/pkg/database/models"
	"github.com/adh-partnership/api/pkg/database/models/constants"
	"github.com/adh-partnership/api/pkg/gin/response"
)

// Get Online Controllers
// @Summary Get Online Controllers
// @Description Get Online Controllers
// @Tags Stats
// @Success 200 {object} []dto.OnlineController
// @Failure 400 {object} response.R
// @Failure 403 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/stats/online [get]
func getOnlineATC(c *gin.Context) {
	var controllers []models.OnlineController

	if err := database.DB.Preload(clause.Associations).Find(&controllers).Error; err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	response.Respond(c, http.StatusOK, dto.ConvertOnlineToDTOs(controllers))
}

// Get Historical Stats
// @Summary Get Historical Stats
// @Description Get Historical Stats
// @Tags Stats
// @Param year path int false "Year"
// @Param month path int false "Month"
// @Success 200 {object} []dto.ControllerStats
// @Failure 400 {object} response.R
// @Failure 403 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/stats/historical/{year}/{month} [get]
func getHistoricalStats(c *gin.Context) {
	var users []models.User
	var year, month int

	pathYear := c.Param("year")
	pathMonth := c.Param("month")
	if pathYear == "" {
		year = time.Now().Year()
		month = int(time.Now().Month())
	} else {
		year = database.Atoi(pathYear)
		month = database.Atoi(pathMonth)
	}

	var ret []dto.ControllerStats
	if err := database.DB.Preload(clause.Associations).Not(&models.User{ControllerType: constants.ControllerTypeNone}).Find(&users).Error; err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	for _, user := range users {
		stat, err := dto.GetDTOForUserAndMonth(&user, month, year)
		if err != nil {
			response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		ret = append(ret, *stat)
	}

	response.Respond(c, http.StatusOK, ret)
}

// Get Facility Report
// @Summary Get Facility Stats Report
// @Description Get Facility Stats Report
// @Tags Stats
// @Param prefix query string false "Prefix, ie ANC"
// @Param suffix query string false "Suffix, ie CTR"
// @Param from query string false "From, ie 2020-01-01"
// @Param to query string false "To, ie 2020-01-31"
// @Success 200 {object} []dto.FacilityReportDTO
// @Failure 400 {object} response.R
// @Failure 403 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/stats/reports/facility [get]
func getFacilityReport(c *gin.Context) {
	results := []*models.ControllerStat{}
	tx := database.DB
	if c.Query("prefix") != "" {
		tx = tx.Where("position LIKE ?", c.Query("prefix")+"_%")
	}
	if c.Query("suffix") != "" {
		tx = tx.Where("position LIKE ?", "%_"+c.Query("suffix"))
	}
	if c.Query("from") != "" {
		tx = tx.Where("logon_time >= ?", c.Query("from"))
	}
	if c.Query("to") != "" {
		tx = tx.Where("logon_time <= ?", c.Query("to"))
	}

	if err := tx.Find(&results).Error; err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	response.Respond(c, http.StatusOK, dto.ConvertControllerStatsToFacilityReport(results))
}
