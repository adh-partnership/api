package traffic

import (
	"fmt"
	"math"
	"regexp"
	"strings"

	"github.com/go-co-op/gocron"
	"golang.org/x/exp/slices"

	"github.com/adh-partnership/api/pkg/config"
	"github.com/adh-partnership/api/pkg/database"
	"github.com/adh-partnership/api/pkg/database/models"
	"github.com/adh-partnership/api/pkg/discord"
	"github.com/adh-partnership/api/pkg/logger"
	"github.com/adh-partnership/api/pkg/network/vatsim"
)

var log = logger.Logger.WithField("component", "job/traffic")

// Max distance for a pilot to be "counted" at an airport.
const maxDistance = 5.0

// Register the repeating job.
func ScheduleJobs(s *gocron.Scheduler) error {
	_, err := s.Every(15).Minutes().SingletonMode().Do(handleTraffic)
	if err != nil {
		return fmt.Errorf("failed to schedule traffic job: %s", err)
	}
	return nil
}

// Checks for unstaffed airports that have pilot traffic.
func handleTraffic() {
	cfgEntries := config.Cfg.Facility.TrafficAlerts.Entries
	if len(cfgEntries) == 0 {
		return
	}
	log.Info("Checking traffic for unstaffed airports")
	airports := []models.Airport{}
	if err := database.DB.Where("twr_type_code != ?", "NON-ATCT").Find(&airports).Error; err != nil {
		log.Errorf("Failed to get airports: %s", err)
		return
	}
	networkData, err := vatsim.GetData()
	if err != nil {
		log.Errorf("Failed to get network data: %s", err)
		return
	}

	needCoverage := make(map[string]uint)
	for _, airport := range airports {
		icao, count := checkAirport(cfgEntries, &airport, networkData)
		if count > 0 {
			needCoverage[icao] = count
		}
	}

	err = sendNotification(needCoverage)
	if err != nil {
		log.Errorf("Error sending uncovered traffic message to Discord: %s", err.Error())
	}
}

// Checks a single airport for pilot traffic and no top-down coverage, as configured.
func checkAirport(
	cfgEntries []config.ConfigFacilityAlertsEntry,
	airport *models.Airport,
	networkData *vatsim.VATSIMData,
) (string, uint) {
	entryIdx := slices.IndexFunc(cfgEntries, func(e config.ConfigFacilityAlertsEntry) bool {
		return e.Airport == airport.ICAO
	})
	if entryIdx == -1 {
		return "", 0
	}
	var count uint
	for _, flight := range networkData.Flights {
		distance := haversineDistance(
			flight.Latitude,
			flight.Longitude,
			float64(airport.Latitude),
			float64(airport.Longitude),
		)
		if distance <= maxDistance {
			count++
		}
	}
	if count < cfgEntries[entryIdx].Threshold {
		return "", 0
	}
	isCovered := false
	for _, couldCover := range cfgEntries[entryIdx].CoveredBy {
		for _, position := range networkData.Controllers {
			re := regexp.MustCompile(couldCover)
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
		return "", 0
	}
	return airport.ICAO, count
}

// Takes the uncovered data, constructs a string, and
// sends the string to the configured Discord webhook.
func sendNotification(needCoverage map[string]uint) error {
	builder := strings.Builder{}
	builder.WriteString("**Unstaffed traffic alerts**\n\n")
	for airport, count := range needCoverage {
		builder.WriteString(fmt.Sprintf("- %s: %d pilots\n", airport, count))
	}
	return discord.NewMessage().SetContent(builder.String()).Send("staffing_request")
}

// From: https://www.movable-type.co.uk/scripts/latlong.html
func haversineDistance(lat1, lon1, lat2, lon2 float64) int {
	R := 6371e3
	φ1 := (lat1 * math.Pi) / 180
	φ2 := (lat2 * math.Pi) / 180
	Δφ := ((lat2 - lat1) * math.Pi) / 180
	Δλ := ((lon2 - lon1) * math.Pi) / 180
	a := math.Sin(Δφ/2)*math.Sin(Δφ/2) + math.Cos(φ1)*math.Cos(φ2)*math.Sin(Δλ/2)*math.Sin(Δλ/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	d := R * c
	return int(math.Round(d * 0.00054))
}
