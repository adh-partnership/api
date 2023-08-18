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

package overflight

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/adh-partnership/api/pkg/database"
	"github.com/adh-partnership/api/pkg/database/models"
	"github.com/adh-partnership/api/pkg/gin/response"
)

type Flightsv1 struct {
	Callsign    string    `json:"callsign" example:"N462AW"`
	CID         int       `json:"cid" example:"876594"`
	Facility    string    `json:"facility" example:"ZDV"`
	Latitude    float32   `json:"lat" example:"-33.867"`
	Longitude   float32   `json:"lon" example:"151.206"`
	Groundspeed int       `json:"spd" example:"150"`
	Heading     int       `json:"hdg" example:"180"`
	Altitude    int       `json:"alt" example:"10500"`
	Aircraft    string    `json:"type" example:"C208"`
	Departure   string    `json:"dep" example:"KLMO"`
	Arrival     string    `json:"arr" example:"KLMO"`
	Route       string    `json:"route" example:"DCT"`
	UpdatedAt   time.Time `json:"lastSeen" example:"2020-01-01T00:00:00Z"`
}

// Get Overflights for Facility
// @Summary Get Overflights for Facility
// @Tags overflight
// @Param fac path string false "Facility, defaults to ZDV if no facility id provided"
// @Success 200 {object} []Flightsv1
// @Failure 500 {object} response.R
// @Router /v1/overflight [GET]
// @Router /v1/overflight/:facility [GET]
func getOverflights(c *gin.Context) {
	var flights []models.Flights

	facility := c.Param("fac")
	if facility == "" {
		facility = "ZDV"
	}

	if err := database.DB.Where(models.Flights{Facility: facility}).Find(&flights).Error; err != nil {
		log.Errorf("Error getting flights for facility %s: %s", facility, err)
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// This could be a DTO in the future
	var flightsv1 []Flightsv1
	for _, flight := range flights {
		flightsv1 = append(flightsv1, Flightsv1{
			Callsign:    flight.Callsign,
			CID:         flight.CID,
			Facility:    flight.Facility,
			Latitude:    flight.Latitude,
			Longitude:   flight.Longitude,
			Groundspeed: flight.Groundspeed,
			Heading:     flight.Heading,
			Altitude:    flight.Altitude,
			Aircraft:    flight.Aircraft,
			Departure:   flight.Departure,
			Arrival:     flight.Arrival,
			Route:       flight.Route,
			UpdatedAt:   flight.UpdatedAt,
		})
	}

	response.Respond(c, http.StatusOK, flightsv1)
}

// Get Overflights for Facility [Legacy/Deprecated]
// @Summary Get Overflights for Facility [Legacy/Deprecated]
// @Tags overflight
// @Param fac path string false "Facility, defaults to ZDV if no facility id provided"
// @Success 200 {object} []Flightsv1
// @Failure 500 {object} response.R
// @Router /live/:facility [GET]
func GetOverflightsLegacy(c *gin.Context) {
	var flights []models.Flights

	facility := c.Param("fac")
	if facility == "" {
		facility = "ZDV"
	}

	if err := database.DB.Where(models.Flights{Facility: facility}).Find(&flights).Error; err != nil {
		log.Errorf("Error getting flights for facility %s: %s", facility, err)
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// This could be a DTO in the future
	var flightsv1 []Flightsv1
	for _, flight := range flights {
		flightsv1 = append(flightsv1, Flightsv1{
			Callsign:    flight.Callsign,
			CID:         flight.CID,
			Facility:    flight.Facility,
			Latitude:    flight.Latitude,
			Longitude:   flight.Longitude,
			Groundspeed: flight.Groundspeed,
			Heading:     flight.Heading,
			Altitude:    flight.Altitude,
			Aircraft:    flight.Aircraft,
			Departure:   flight.Departure,
			Arrival:     flight.Arrival,
			Route:       flight.Route,
			UpdatedAt:   flight.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, flightsv1)
}
