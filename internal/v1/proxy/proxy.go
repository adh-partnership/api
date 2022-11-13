package proxy

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/adh-partnership/api/pkg/gin/response"
)

// Proxy METAR Data
// @Summary Proxy METAR Data
// @Description Proxy METAR Data
// @Tags Proxy
// @Param icao path string true "ICAO, multiple ICAOs can be separated by a comma"
// @Success 200 {object} string
// @Failure 400 {object} response.R
// @Failure 403 {object} response.R
// @Failure 500 {object} response.R
// @Router /proxy/metar/{icao} [get]
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
