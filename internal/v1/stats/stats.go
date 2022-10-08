package stats

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"

	"github.com/adh-partnership/api/pkg/database"
	"github.com/adh-partnership/api/pkg/database/dto"
	"github.com/adh-partnership/api/pkg/database/models"
	"github.com/adh-partnership/api/pkg/gin/response"
)

// Get Online Controllers
// @Summary Get Online Controllers
// @Description Get Online Controllers
// @Tags Stats
// @Success 200 {object} models.OnlineController[]
// @Failure 400 {object} response.R
// @Failure 403 {object} response.R
// @Failure 500 {object} response.R
// @Router /stats/online [get]
func getOnlineATC(c *gin.Context) {
	var controllers []models.OnlineController

	if err := database.DB.Find(&controllers).Error; err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	response.Respond(c, http.StatusOK, controllers)
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
// @Router /stats/historical/{year}/{month} [get]
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
	if err := database.DB.Preload(clause.Associations).Not(&models.User{Status: models.ControllerStatusOptions["none"]}).Find(&users).Error; err != nil {
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
