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
