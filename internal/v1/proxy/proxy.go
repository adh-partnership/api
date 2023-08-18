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

package proxy

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/adh-partnership/api/pkg/gin/response"
)

// Proxy METAR Data
// @Summary [Deprecated] Proxy METAR Data
// @Description Proxy METAR Data. Deprecated. VATSIM's Cloudflare front may throttle this, so use our job fetched data where possible (/v1/weather/metar/:icao)
// @Tags Proxy
// @Param icao path string true "ICAO, multiple ICAOs can be separated by a comma"
// @Success 200 {object} string
// @Failure 400 {object} response.R
// @Failure 403 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/proxy/metar/{icao} [get]
func getMetar(c *gin.Context) {
	icao := c.Param("icao")
	if icao == "" {
		response.RespondError(c, http.StatusBadRequest, "Missing ICAO")
		return
	}

	// Get METAR Data
	resp, err := http.Get("https://metar.vatsim.net/" + icao)
	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	defer resp.Body.Close()

	// Read Body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	response.Respond(c, http.StatusOK, string(body))
}

// Proxy TAF Data
// @Summary [Deprecated] Proxy TAF Data
// @Description Proxy TAF Data. Deprecated. VATSIM's Cloudflare front may throttle this, so use our job fetched data where possible (/v1/weather/taf/:icao)
// @Tags Proxy
// @Param icao path string true "ICAO"
// @Success 200 {object} string
// @Failure 400 {object} response.R
// @Failure 404 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/proxy/taf/{icao} [get]
func getTaf(c *gin.Context) {
	icao := c.Param("icao")
	if icao == "" {
		response.RespondError(c, http.StatusBadRequest, "Missing ICAO")
		return
	}

	// Get TAF Data
	resp, err := http.Get("https://tgftp.nws.noaa.gov/data/forecasts/taf/stations/" + icao + ".TXT")
	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		response.RespondError(c, http.StatusNotFound, "Not Found")
		return
	}

	// Read Body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	response.Respond(c, http.StatusOK, string(body))
}
