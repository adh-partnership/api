package airport

import (
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/adh-partnership/api/pkg/database"
	"github.com/adh-partnership/api/pkg/database/dto"
	_ "github.com/adh-partnership/api/pkg/database/models" // for docs
	"github.com/adh-partnership/api/pkg/gin/response"
)

// Get airports in ARTCC
// @Summary Get airports in ARTCC
// @Description Get airports in ARTCC
// @Tags Airports
// @Param center path string true "ARTCC identifier, ie ZAN"
// @Success 200 {object} []models.Airport
// @Router /v1/airports/:center [get]
func getCenter(c *gin.Context) {
	airports, err := database.FindAirportsByARTCC(c.Param("center"))
	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if len(airports) == 0 {
		response.RespondError(c, http.StatusNotFound, "Not Found")
		return
	}

	response.Respond(c, http.StatusOK, airports)
}

// Get Airport
// @Summary Get Airport
// @Description Get Airport
// @Tags Airports
// @Param center path string true "ARTCC identifier, ie ZAN"
// @Param id path string true "Airport identifier, ie KATL [FAA Identifier or ICAO]"
// @Success 200 {object} models.Airport
// @Router /v1/airports/:center/:id [get]
func getAirport(c *gin.Context) {
	airport, err := database.FindAirportByID(c.Param("id"))
	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if airport == nil {
		response.RespondError(c, http.StatusNotFound, "Not Found")
		return
	}

	response.Respond(c, http.StatusOK, airport)
}

// Get Airport ATC
// @Summary Get Airport ATC
// @Description Get Airport ATC
// @Tags Airports
// @Param center path string true "ARTCC identifier, ie ZAN"
// @Param id path string true "Airport identifier, ie KATL [FAA Identifier or ICAO]"
// @Success 200 {object} []models.AirportATC
// @Router /v1/airports/:center/:id/atc [get]
func getAirportATC(c *gin.Context) {
	airport, err := database.FindAirportByID(c.Param("id"))
	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	atc, err := database.FindAirportATCByID(airport.ID)
	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if atc == nil {
		response.RespondError(c, http.StatusNotFound, "Not Found")
		return
	}

	response.Respond(c, http.StatusOK, atc)
}

// Get Airport Charts
// @Summary Get Airport Charts
// @Description Get Airport Charts
// @Tags Airports
// @Param center path string true "ARTCC identifier, ie ZAN"
// @Param id path string true "Airport identifier, ie KATL [FAA Identifier or ICAO]"
// @Success 200 {object} []models.AirportChart
// @Router /v1/airports/:center/:id/charts [get]
func getAirportCharts(c *gin.Context) {
	airport, err := database.FindAirportByID(c.Param("id"))
	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	charts, err := database.FindAirportChartsByID(airport.ID)
	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if charts == nil {
		response.RespondError(c, http.StatusNotFound, "Not Found")
		return
	}

	response.Respond(c, http.StatusOK, charts)
}

// Get Airport Weather
// @Summary Get Airport Weather
// @Description Get Airport Weather
// @Tags Airports
// @Param center path string true "ARTCC identifier, ie ZAN"
// @Param id path string true "Airport identifier, ie KATL [FAA Identifier or ICAO]"
// @Success 200 {object} dto.AirportWeatherDTO
// @Router /v1/airports/:center/:id/weather [get]
func getAirportWeather(c *gin.Context) {
	airport, err := database.FindAirportByID(c.Param("id"))
	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	d := &dto.AirportWeatherDTO{
		ID: airport.ICAO,
	}

	type Result struct {
		data string
		err  error
	}

	metar := make(chan Result)
	taf := make(chan Result)

	go func(icao string) {
		// Get METAR Data
		resp, err := http.Get("https://metar.vatsim.net/" + icao)
		if err != nil {
			metar <- Result{err: err}
			return
		}
		defer resp.Body.Close()

		// Read Body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			metar <- Result{err: err}
			return
		}

		metar <- Result{data: string(body)}
	}(d.ID)

	go func(icao string) {
		// Get TAF Data
		resp, err := http.Get("https://tgftp.nws.noaa.gov/data/forecasts/taf/stations/" + icao + ".TXT")
		if err != nil {
			taf <- Result{err: err}
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNotFound {
			taf <- Result{err: errors.New("not found")}
			return
		}

		// Read Body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			taf <- Result{err: err}
			return
		}
		taf <- Result{data: string(body)}
	}(d.ID)

	// Wait for both to finish
	metarResult := <-metar
	tafResult := <-taf

	if metarResult.err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if tafResult.err != nil && tafResult.err.Error() != "not found" {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	response.Respond(c, http.StatusOK, dto.AirportWeatherDTO{
		ID:    d.ID,
		METAR: metarResult.data,
		TAF:   tafResult.data,
	})
}
