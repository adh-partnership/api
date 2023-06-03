package weather

import (
	"fmt"

	"github.com/go-co-op/gocron"
	idsWeather "github.com/vpaza/ids/pkg/weather"

	"github.com/adh-partnership/api/pkg/database"
	"github.com/adh-partnership/api/pkg/database/models"
	"github.com/adh-partnership/api/pkg/logger"
	"github.com/adh-partnership/api/pkg/utils"
)

var log = logger.Logger.WithField("component", "job/weather")

func ScheduleJobs(s *gocron.Scheduler) error {
	_, err := s.Every(2).Minutes().SingletonMode().Do(handleMETAR)
	if err != nil {
		return fmt.Errorf("failed to schedule metar job: %s", err)
	}

	_, err = s.Every(15).Minutes().SingletonMode().Do(handleTAF)
	if err != nil {
		return fmt.Errorf("failed to schedule taf job: %s", err)
	}

	return nil
}

func handleMETAR() {
	airports := []models.Airport{}
	if err := database.DB.Where(&models.Airport{HasMETAR: true}).Find(&airports).Error; err != nil {
		log.Errorf("Failed to get airports: %s", err)
		return
	}

	for _, airport := range airports {
		metar, err := idsWeather.GetMetar(airport.ICAO)
		if err != nil {
			log.Errorf("Failed to get METAR for %s: %s", airport.ICAO, err)
			continue
		}

		if airport.METAR == metar.RawText {
			continue
		}

		if airport.METAR != metar.RawText {
			airport.METAR = metar.RawText
			if err := database.DB.Save(&airport).Error; err != nil {
				log.Errorf("Failed to save airport %s: %s", airport.ICAO, err)
				continue
			}
		}
	}
}

func handleTAF() {
	airports := []models.Airport{}
	if err := database.DB.Where(&models.Airport{HasTAF: true}).Find(&airports).Error; err != nil {
		log.Errorf("Failed to get airports: %s", err)
		return
	}

	for _, airport := range airports {
		body, err := utils.GetAirportTAF(airport.ICAO)
		if err != nil {
			log.Errorf("Failed to get TAF for %s: %s", airport.ICAO, err)
			continue
		}

		if airport.TAF == string(body) {
			continue
		}

		if airport.TAF != string(body) {
			airport.TAF = string(body)
			if err := database.DB.Save(&airport).Error; err != nil {
				log.Errorf("Failed to save airport %s: %s", airport.ICAO, err)
				continue
			}
		}
	}
}
