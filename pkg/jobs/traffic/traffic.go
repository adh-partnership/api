package traffic

import (
	"fmt"
	"math"
	"regexp"

	"golang.org/x/exp/slices"

	"github.com/adh-partnership/api/pkg/config"
	"github.com/adh-partnership/api/pkg/database"
	"github.com/adh-partnership/api/pkg/database/models"
	"github.com/adh-partnership/api/pkg/logger"
	"github.com/adh-partnership/api/pkg/network/vatsim"
	"github.com/go-co-op/gocron"
)

var log = logger.Logger.WithField("component", "job/traffic")

const MAX_DISTANCE = 5.0
const ALERT_THRESHOLD = 10

func ScheduleJobs(s *gocron.Scheduler) error {
	// TODO pick better time
	_, err := s.Every(60).Seconds().SingletonMode().Do(handleTraffic)
	if err != nil {
		return fmt.Errorf("failed to schedule traffic job: %s", err)
	}
	return nil
}

func handleTraffic() {
	cfg_entries := config.Cfg.Facility.TrafficAlerts.Entries
	if len(cfg_entries) == 0 {
		return
	}
	log.Info("Checking traffic")
	airports := []models.Airport{}
	if err := database.DB.Where("twr_type_code != ?", "NON-ATCT").Find(&airports).Error; err != nil {
		log.Errorf("Failed to get airports: %s", err)
		return
	}
	network_data, err := vatsim.GetData()
	if err != nil {
		log.Errorf("Failed to get network data: %s", err)
		return
	}
	log.Infof("Looking at %d airports", len(airports))
	alerts := []string{}
	for _, airport := range airports {
		entry_idx := slices.IndexFunc(cfg_entries, func(e config.ConfigFacilityAlertsEntry) bool { return e.Airport == airport.ICAO })
		if entry_idx == -1 {
			continue
		}
		log.Infof("Airport found in DB and config: %s", airport.ICAO)
		var count uint = 0
		for _, flight := range network_data.Flights {
			distance := haversineDistance(flight.Latitude, flight.Longitude, float64(airport.Latitude), float64(airport.Longitude))
			if distance <= MAX_DISTANCE {
				count += 1
			}
		}
		log.Infof("Pilots near %s: %d", airport.ICAO, count)
		if count < cfg_entries[entry_idx].Threshold {
			continue
		}
		isCovered := false
		for _, could_cover := range cfg_entries[entry_idx].CoveredBy {
			for _, position := range network_data.Controllers {
				re := regexp.MustCompile(could_cover)
				if re.MatchString(position.Callsign) {
					isCovered = true
					break
				}
			}
			if isCovered {
				break
			}
		}
		if isCovered {
			log.Infof("%s airport is being controlled", airport.ICAO)
		} else {
			alerts = append(alerts, airport.ICAO)
		}
	}
	log.Infof("Traffic alerts: %v", alerts)
}

// <https://www.movable-type.co.uk/scripts/latlong.html>
func haversineDistance(lat1, lon1, lat2, lon2 float64) int {
	R := 6371e3
	φ1 := (lat1 * math.Pi) / 180
	φ2 := (lat2 * math.Pi) / 180
	Δφ := ((lat2 - lat1) * math.Pi) / 180
	Δλ := ((lon2 - lon1) * math.Pi) / 180
	a := math.Sin(Δφ/2)*math.Sin(Δφ/2) +
		math.Cos(φ1)*math.Cos(φ2)*math.Sin(Δλ/2)*math.Sin(Δλ/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	d := R * c
	return int(math.Round(d * 0.00054))
}
