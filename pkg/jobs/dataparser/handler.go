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

package dataparser

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/go-co-op/gocron"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"

	"github.com/adh-partnership/api/pkg/config"
	"github.com/adh-partnership/api/pkg/database"
	"github.com/adh-partnership/api/pkg/database/models"
	"github.com/adh-partnership/api/pkg/database/models/constants"
	"github.com/adh-partnership/api/pkg/discord"
	"github.com/adh-partnership/api/pkg/geo"
	"github.com/adh-partnership/api/pkg/logger"
	"github.com/adh-partnership/api/pkg/network/vatsim"
	"github.com/adh-partnership/api/pkg/server"
)

var fac []Facility

var log = logger.Logger.WithField("component", "job/flightparser")

func Initialize(cron *gocron.Scheduler) error {
	_, err := cron.Every(1).Minute().SingletonMode().Do(handle)
	if err != nil {
		log.Errorf("Error scheduling ParseFlights: %v", err)
		return err
	}

	jsonfile, err := os.Open("boundaries.json")
	if err != nil {
		log.Errorf("Error opening boundaries.json: %v", err)
		return err
	}

	data, err := io.ReadAll(jsonfile)
	if err != nil {
		log.Errorf("Error reading boundaries.json: %v", err)
		return err
	}

	defer func() {
		_ = jsonfile.Close()
	}()

	if err := json.Unmarshal(data, &fac); err != nil {
		log.Errorf("Error unmarshalling boundaries.json: %v", err)
		return err
	}

	for i := 0; i < len(fac); i++ {
		var points []geo.Point
		for j := 0; j < len(fac[i].Boundary); j++ {
			points = append(points, geo.Point{X: fac[i].Boundary[j][0], Y: fac[i].Boundary[j][1]})
		}
		fac[i].Polygon = geo.Polygon{Points: points}
	}

	return nil
}

func handle() {
	log.Debug("Handling parse flights")

	vatsimData, err := vatsim.GetData()
	if err != nil {
		log.Errorf("Error getting vatsim data: %v", err)
		return
	}

	flightDone := make(chan bool)
	atcDone := make(chan bool)

	go parseFlights(flightDone, vatsimData.Flights)
	go parseATC(atcDone, vatsimData.Controllers)

	<-flightDone
	<-atcDone
}

func parseATC(atcDone chan bool, controllers []*vatsim.VATSIMController) {
	updateid, _ := gonanoid.New(24)

	for _, controller := range controllers {
		prefix := strings.Split(controller.Callsign, "_")[0]

		// Ignore observers
		if controller.Facility == 0 {
			continue
		}

		allowedSuffixes := []string{"_RMP", "_DEL", "_GND", "_TWR", "_APP", "_DEP", "_CTR", "_RDO", "_FSS", "_TMU", "_FMP"}
		isAllowed := false
		for _, suffix := range allowedSuffixes {
			if strings.HasSuffix(controller.Callsign, suffix) {
				isAllowed = true
				break
			}
		}
		if !isAllowed {
			log.Tracef("Skipping %s, disallowed suffix", controller.Callsign)
			continue
		}

		// Check if we are tracking this prefix, if so, add it to the database
		if _, ok := server.Server.TrackedPrefixes[prefix]; !ok {
			continue
		}

		c := &models.OnlineController{}
		if err := database.DB.Where("position = ?", controller.Callsign).First(&c).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				log.Errorf("Error looking up controller (%s): %v", controller.Callsign, err)
				continue
			}
		}

		if c.Position == "" {
			// Safe to assume this is a new controller
			go func(controller *vatsim.VATSIMController) {
				user, err := database.FindUserByCID(fmt.Sprint(controller.CID))
				if err != nil {
					log.Errorf("Error finding user with CID %d: %v", controller.CID, err)
					return
				}
				if user == nil || user.ControllerType == constants.ControllerStatusNone {
					_ = discord.NewMessage().AddEmbed(
						discord.NewEmbed().SetTitle("Not active controller is on position").SetColor(
							discord.GetColor("ff", "00", "00"),
						).AddField(
							discord.NewField().SetName("CID").SetValue(fmt.Sprint(controller.CID)).SetInline(true),
						).AddField(
							discord.NewField().SetName("Name").SetValue(controller.Name).SetInline(true),
						).AddField(
							discord.NewField().SetName("Position").SetValue(controller.Callsign).SetInline(true),
						),
					).Send("seniorstaff")

					return
				}

				if config.Cfg.Features.ControllerOnline {
					_ = discord.NewMessage().
						AddEmbed(
							discord.NewEmbed().SetTitle(fmt.Sprintf("%s is now online!", controller.Callsign)).SetColor(
								discord.GetColor("00", "00", "ff"),
							).
								SetDescription(fmt.Sprintf(
									"%s %s (%s) is now online as %s",
									user.FirstName,
									user.LastName,
									user.OperatingInitials,
									controller.Callsign,
								)),
						).Send("online")
				}
			}(controller)
		}

		c.UserID = uint(controller.CID)
		c.Position = controller.Callsign
		c.Frequency = controller.Frequency
		c.LogonTime = *controller.LogonTime
		c.UpdateID = updateid

		if err := database.DB.Save(&c).Error; err != nil {
			log.Error("Error saving controller information for " + c.Position + " to database: " + err.Error())
		}
	}

	var oldControllers []models.OnlineController
	if err := database.DB.Where("update_id != ?", updateid).Find(&oldControllers).Error; err != nil {
		log.Errorf("Error looking up old controllers: %v", err)
	}

	for _, controller := range oldControllers {
		stat := &models.ControllerStat{
			UserID:    controller.UserID,
			Position:  controller.Position,
			LogonTime: controller.LogonTime,
			Duration:  int(controller.UpdatedAt.Sub(controller.LogonTime) / time.Second),
		}
		if err := database.DB.Create(&stat).Error; err != nil {
			log.Errorf("Error creating controller stat: %v", err)
		}
		if err := database.DB.Delete(&controller).Error; err != nil {
			log.Errorf("Error deleting old controller entry: %v", err)
		}
	}

	atcDone <- true
}

func parseFlights(flightDone chan bool, flights []*vatsim.VATSIMFlight) {
	updateid, _ := gonanoid.New(24)

	for _, flight := range flights {
		f := &models.Flights{}
		if err := database.DB.Where("callsign = ?", flight.Callsign).First(&f).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				log.Errorf("Error looking up flight (%s): %v", flight.Callsign, err)
				continue
			}
		}

		f.Aircraft = flight.FlightPlan.Aircraft
		f.CID = flight.CID
		f.Callsign = flight.Callsign
		f.Latitude = float32(flight.Latitude)
		f.Longitude = float32(flight.Longitude)
		f.Altitude = flight.Altitude
		f.Heading = flight.Heading
		f.Groundspeed = flight.Groundspeed
		f.Departure = flight.FlightPlan.Departure
		f.Arrival = flight.FlightPlan.Arrival
		f.Route = flight.FlightPlan.Route
		f.Facility = ""
		f.UpdateID = updateid

		if f.Latitude < 75.0 && f.Latitude > 21.0 && f.Longitude < -50.0 && f.Longitude > -179.0 {
			for j := 0; j < len(fac); j++ {
				facID := fac[j].ID
				poly := fac[j].Polygon
				p := geo.Point{X: float64(f.Longitude), Y: float64(f.Latitude)}
				if geo.PointInPolygon(p, poly) {
					f.Facility = facID
				}
			}
		}

		if err := database.DB.Save(&f).Error; err != nil {
			log.Error("Error saving flight information for " + f.Callsign + " to database: " + err.Error())
		}
	}

	if err := database.DB.Where("update_id != ?", updateid).Delete(&models.Flights{}).Error; err != nil {
		log.Errorf("Error deleting old flights: %v", err)
	}

	flightDone <- true
}
