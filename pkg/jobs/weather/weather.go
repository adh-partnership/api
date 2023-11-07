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

	"github.com/adh-partnership/api/pkg/database"
	"github.com/adh-partnership/api/pkg/database/models"
	"github.com/adh-partnership/api/pkg/logger"
	"github.com/adh-partnership/api/pkg/weather"
)

var log = logger.Logger.WithField("component", "job/weather")

func ScheduleJobs(s *gocron.Scheduler) error {
	_, err := s.Every(2).Minutes().SingletonMode().Do(handleWeather)
	if err != nil {
		return fmt.Errorf("failed to schedule weather job: %s", err)
	}

	return nil
}

func handleWeather() {
	err := weather.UpdateWeatherCache()
	if err != nil {
		log.Errorf("Failed to update weather cache: %s", err)
		return
	}

	airports := []models.Airport{}
	if err := database.DB.Find(&airports).Error; err != nil {
		log.Errorf("Failed to get airports: %s", err)
		return
	}

	for _, airport := range airports {
		if airport.HasMETAR || airport.HasTAF {
			wx, err := weather.GetWeather(airport.ICAO)
			if err != nil {
				log.Errorf("Failed to get weather for %s: %s", airport.ICAO, err)
				continue
			}
			if airport.HasMETAR {
				airport.METAR = wx.METAR
			}
			if airport.HasTAF {
				airport.TAF = wx.TAF
			}
			if err := database.DB.Save(&airport).Error; err != nil {
				log.Errorf("Failed to save airport %s: %s", airport.ICAO, err)
				continue
			}
		}
	}
}
