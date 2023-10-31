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

package activity

import (
	"github.com/go-co-op/gocron"

	"github.com/adh-partnership/api/pkg/logger"
)

var log = logger.Logger.WithField("component", "job/activity")

// @TODO
func ScheduleJobs(s *gocron.Scheduler) error {
	/*
		_, err := s.Cron("0 1 * * *").SingletonMode().Do(handleWarning)
		if err != nil {
			return fmt.Errorf("failed to schedule inactivity job: %s", err)
		}

		months := "*"
		if len(config.Cfg.Facility.Activity.Inactive.Months) != 0 {
			var m []string
			for _, month := range config.Cfg.Facility.Activity.Inactive.Months {
				if month < 0 || month > 12 {
					log.Warnf("Invalid month %d, skipping", month)
					continue
				}
				m = append(m, fmt.Sprint(month))
			}
			months = strings.Join(m, ",")
		}

		_, err = s.Cron(fmt.Sprintf("0 2 1 %s *", months)).SingletonMode().Do(handleInactive)
		if err != nil {
			return fmt.Errorf("failed to schedule inactive job: %s", err)
		}
	*/
	log.Warnf("Due to changes in GCAP, Activity Jobs are not yet supported and will be reintroduced later down the line.")

	return nil
}

/*
func handleInactive() {
	if !config.Cfg.Facility.Activity.Inactive.Enabled {
		return
	}

	log.Info("Starting inactivity job")
	beginning := time.Now()
	if beginning.Day() != 1 {
		beginning = BeginningOfNextMonth(beginning)
	}

	lastMonth := beginning.AddDate(0, -1, 0)
	log.Debugf("Last Month=%+v", lastMonth)

	// Get all active controllers that have joined prior to config.Cfg.Facility.Activity.Period month(s) ago
	var controllers []*models.User
	if err := database.DB.Preload(clause.Associations).
		Where(&models.User{Status: constants.ControllerStatusActive}).
		Where("roster_join_date <= date_sub(date_sub(now(), interval ? month), interval 1 day)", config.Cfg.Facility.Activity.Inactive.Period).
		Find(&controllers).Error; err != nil {
		log.Errorf("Failed to get active controllers: %s", err)
		return
	}

	for _, controller := range controllers {
		if controller.ExemptedFromActivity {
			continue
		}

		if !IsInactive(controller, &lastMonth) {
			continue
		}

		log.Infof("Controller %d is inactive", controller.CID)
		controller.Status = constants.ControllerStatusInactive
		if err := database.DB.Save(controller).Error; err != nil {
			log.Errorf("Failed to save controller: %s", err)
			continue
		}

		go func(controller *models.User) {
			_ = discord.NewMessage().
				SetContent(
					fmt.Sprintf("Controller marked inactive: %s %s - %s (%d/%s)",
						controller.FirstName,
						controller.LastName,
						controller.OperatingInitials,
						controller.CID,
						controller.Rating.Short,
					)).
				Send("seniorstaff")
		}(controller)
	}
}

func handleWarning() {
	// Check if we need to run
	if !config.Cfg.Facility.Activity.Warning.Enabled {
		return
	}

	daysBefore := config.Cfg.Facility.Activity.Warning.DaysBefore
	// Is today at the daysBefore mark?
	if time.Now().Day() != DaysBefore(BeginningOfNextMonth(time.Now()), daysBefore).Day() {
		return
	}

	log.Info("Starting inactivity warning job")

	beginning := time.Now()
	if beginning.Day() != 1 {
		beginning = BeginningOfNextMonth(beginning)
	}

	lastMonth := beginning.AddDate(0, -1, 0)
	log.Debugf("Last Month=%+v", lastMonth)

	// Get all active controllers that have joined prior to config.Cfg.Facility.Activity.Period month(s) ago
	var controllers []*models.User
	if err := database.DB.Preload(clause.Associations).
		Where(&models.User{Status: constants.ControllerStatusActive}).
		Where("roster_join_date <= date_sub(date_sub(now(), interval ? month), interval 1 day)", config.Cfg.Facility.Activity.Inactive.Period).
		Find(&controllers).Error; err != nil {
		log.Errorf("Failed to get active controllers: %s", err)
		return
	}

	for _, controller := range controllers {
		if controller.ExemptedFromActivity {
			continue
		}

		if !IsInactive(controller, &lastMonth) {
			continue
		}

		log.Infof("Sending inactivity warning to controller %d", controller.CID)
		go func(controller *models.User) {
			err := email.Send(
				controller.Email,
				"",
				"",
				email.Templates["inactive_warning"],
				map[string]interface{}{
					"FirstName": controller.FirstName,
					"LastName":  controller.LastName,
				},
			)
			if err != nil {
				log.Errorf("Failed to send email to %s: %s", controller.Email, err)
			}

			_ = discord.NewMessage().
				SetContent(
					fmt.Sprintf("Sent inactivity warning to controlller %s %s - %s (%d/%s)",
						controller.FirstName,
						controller.LastName,
						controller.OperatingInitials,
						controller.CID,
						controller.Rating.Short,
					)).
				Send("seniorstaff")
		}(controller)
	}

	log.Info("Finished activity job")
}

// Checks
func IsInactive(controller *models.User, since *time.Time) bool {
	var sum float32

	if since == nil {
		since = func(t time.Time) *time.Time { return &t }(time.Now())
	}

	// If we are on the first, we need to check the previous month
	// otherwise include this month
	var lastMonth time.Time
	if since.Day() != 1 {
		lastMonth = since.AddDate(0, -1, 0)
	} else {
		lastMonth = *since
	}

	for i := 0; i < config.Cfg.Facility.Activity.Inactive.Period; i++ {
		month := lastMonth.AddDate(0, -i, 0)
		cab, tracon, enroute, err := database.GetStatsForUserAndMonth(controller, int(lastMonth.Month()), lastMonth.Year())
		if err != nil {
			log.Errorf("Failed to get stats for %d: %s", controller.CID, err)
			return false
		}
		sum = cab + tracon + enroute
		log.Debugf("Controller %d - %s: %f", controller.CID, month.Format("2006-01"), sum)
		if sum >= float32(config.Cfg.Facility.Activity.Inactive.MinHours) {
			log.Debugf("Controller %d has %f hours, enough for activity", controller.CID, sum)
			return false
		}
	}
	// They did not have enough hours for the past (config period) months
	log.Debugf("Controller %d has %f hours, not enough for activity", controller.CID, sum)
	return true
}

func BeginningOfNextMonth(t time.Time) time.Time {
	return t.AddDate(0, 1, -t.Day()+1)
}

func BeginningOfMonth(t time.Time) time.Time {
	return t.AddDate(0, 0, -t.Day()+1)
}

func EndOfMonth(t time.Time) time.Time {
	return t.AddDate(0, 1, -t.Day())
}

func DaysBefore(t time.Time, days int) time.Time {
	return t.AddDate(0, 0, -days)
}
*/
