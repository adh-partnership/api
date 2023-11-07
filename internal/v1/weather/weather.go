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

package weather

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	idsWeather "github.com/vpaza/ids/pkg/weather"

	"github.com/adh-partnership/api/pkg/database"
	"github.com/adh-partnership/api/pkg/database/models"
	"github.com/adh-partnership/api/pkg/gin/response"
	"github.com/adh-partnership/api/pkg/logger"
	"github.com/adh-partnership/api/pkg/utils"
	"github.com/adh-partnership/api/pkg/weather"
)

var log = logger.Logger.WithField("component", "weather")

func populate(c *gin.Context) {
	go func() {
		// We use charts here because charts is already filtered to the areas we may want, at least for AK and HI... ZDV might be another issue
		// since charts are generally filtered by state and ZDV has some partial state coverage... but if we cache weather for those areas, is
		// it really a problem?
		charts := []models.AirportChart{}
		if err := database.DB.Distinct("airport_id").Find(&charts).Error; err != nil {
			log.Errorf("Failed to get charts: %s", err.Error())
			return
		}

		for _, chart := range charts {
			airport, err := database.FindAirportByID(chart.AirportID)
			if err != nil {
				continue
			}

			metar, err := idsWeather.GetMetar(airport.ICAO)
			if err != nil {
				continue
			}

			airport.HasMETAR = true
			airport.METAR = metar.RawText
			body, err := utils.GetAirportTAF(airport.ICAO)
			if err != nil {
				continue
			}

			airport.HasTAF = true
			airport.TAF = string(body)

			if err := database.DB.Save(&airport).Error; err != nil {
				log.Errorf("Failed to save airport %s: %s", airport.ICAO, err.Error())
				continue
			}
			log.Infof("Populated airport %s with %v %v", airport.ICAO, airport.HasMETAR, airport.HasTAF)
			time.Sleep(500 * time.Millisecond) // Sleep for half a second
		}
	}()

	response.Respond(c, http.StatusOK, "Populating")
}

// Get METAR Data
// @Summary Get METAR Data
// @Description Get METAR Data
// @Tags Weather
// @Param icao path string true "ICAO"
// @Success 200 {object} string
// @Failure 400 {object} response.R
// @Failure 403 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/weather/metar/{icao} [get]
func getMetar(c *gin.Context) {
	icao := c.Param("icao")
	if icao == "" {
		response.RespondError(c, http.StatusBadRequest, "Missing ICAO")
		return
	}

	wx, err := weather.GetWeather(icao)
	if err != nil {
		response.RespondError(c, http.StatusNotFound, "Not Found")
		return
	}

	response.Respond(c, http.StatusOK, wx.METAR)
}

// Get TAF Data
// @Summary Get TAF Data
// @Description Get TAF Data
// @Tags Weather
// @Param icao path string true "ICAO"
// @Success 200 {object} string
// @Failure 400 {object} response.R
// @Failure 404 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/weather/taf/{icao} [get]
func getTaf(c *gin.Context) {
	icao := c.Param("icao")
	if icao == "" {
		response.RespondError(c, http.StatusBadRequest, "Missing ICAO")
		return
	}

	wx, err := weather.GetWeather(icao)
	if err != nil {
		response.RespondError(c, http.StatusNotFound, "Not Found")
		return
	}

	response.Respond(c, http.StatusOK, wx.TAF)
}
