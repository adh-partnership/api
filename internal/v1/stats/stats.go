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

// Get Historical Stats
// @Summary Get Historical Stats
// @Description Get Historical Stats
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
	if c.Query("prefix") != "" {
		database.DB.Where("position LIKE ?", c.Query("prefix")+"_%")
	}
	if c.Query("suffix") != "" {
		database.DB.Where("position LIKE ?", "%_"+c.Query("suffix"))
	}
	if c.Query("from") != "" {
		database.DB.Where("logon_time >= ?", c.Query("from"))
	}
	if c.Query("to") != "" {
		database.DB.Where("logon_time <= ?", c.Query("to"))
	}

	if err := database.DB.Find(&results).Error; err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	response.Respond(c, http.StatusOK, dto.ConvertControllerStatsToFacilityReport(results))
}
