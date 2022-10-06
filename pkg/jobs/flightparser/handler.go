package flightparser

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"time"

	"github.com/go-co-op/gocron"
	"gorm.io/gorm"

	"github.com/adh-partnership/api/pkg/database"
	"github.com/adh-partnership/api/pkg/database/models"
	"github.com/adh-partnership/api/pkg/geo"
	"github.com/adh-partnership/api/pkg/logger"
	"github.com/adh-partnership/api/pkg/network/vatsim"
)

var fac []Facility

var log = logger.Logger.WithField("component", "job/flightparser")

func Initialize(cron *gocron.Scheduler) error {
	_, err := cron.Every(1).Minute().SingletonMode().Do(HandleParseFlights)
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

func HandleParseFlights() {
	log.Debug("Handling parse flights")

	vatsimData, err := vatsim.GetData()
	if err != nil {
		log.Errorf("Error getting vatsim data: %v", err)
		return
	}

	go database.DB.Where("updated_at < ?", time.Now().Add((time.Minute*5)*-1)).Delete(&models.Flights{})

	for i := 0; i < len(vatsimData.Flights); i++ {
		flight := vatsimData.Flights[i]
		f := &models.Flights{}
		if err := database.DB.Where("callsign = ?", flight.Callsign).First(f).Error; err != nil {
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
}
