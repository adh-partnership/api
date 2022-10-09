package activity

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron"

	"github.com/adh-partnership/api/pkg/config"
	"github.com/adh-partnership/api/pkg/database"
	"github.com/adh-partnership/api/pkg/database/models"
	"github.com/adh-partnership/api/pkg/database/models/constants"
	"github.com/adh-partnership/api/pkg/email"
	"github.com/adh-partnership/api/pkg/logger"
)

var log = logger.Logger.WithField("component", "job/activity")

func ScheduleJobs(s *gocron.Scheduler) error {
	t := fmt.Sprintf(
		"%d %d %d * *",
		config.Cfg.Facility.Activity.RunAtMinute,
		config.Cfg.Facility.Activity.RunAtHour,
		config.Cfg.Facility.Activity.RunOnDay,
	)
	_, err := s.Cron(t).SingletonMode().Do(handle)
	return err
}

func handle() {
	log.Info("Starting activity job")

	lastMonth := time.Now().AddDate(0, -1, 0)
	log.Debugf("Last Month=%+v", lastMonth)

	// Get all active controllers
	var controllers []*models.User
	if err := database.DB.Where(&models.User{Status: constants.ControllerStatusActive}).Find(&controllers).Error; err != nil {
		log.Errorf("Failed to get active controllers: %s", err)
		return
	}

OUTER:
	for _, controller := range controllers {
		var sum float32

		for i := 0; i < config.Cfg.Facility.Activity.Period; i++ {
			month := lastMonth.AddDate(0, -i, 0)
			cab, tracon, enroute, err := database.GetStatsForUserAndMonth(controller, int(lastMonth.Month()), lastMonth.Year())
			if err != nil {
				log.Errorf("Failed to get stats for %s: %s", controller.CID, err)
				continue OUTER
			}
			sum = cab + tracon + enroute
			log.Debugf("Controller %s - %s: %f", controller.CID, month.Format("2006-01"), sum)
			if sum >= float32(config.Cfg.Facility.Activity.MinHours) {
				log.Debugf("Controller %s has %f hours, enough for activity", controller.CID, sum)
				continue OUTER
			}
		}

		// They did not have enough hours for the past (config period) months
		log.Debugf("Controller %s has %f hours, not enough for activity", controller.CID, sum)
		controller.Status = constants.ControllerStatusInactive
		if err := database.DB.Save(controller).Error; err != nil {
			log.Errorf("Failed to save controller: %s", err)
			continue
		}

		log.Infof("Controller %s is now inactive", controller.CID)
		err := email.Send(
			controller.Email,
			"",
			"",
			[]string{},
			"activity_warning",
			map[string]interface{}{
				"FirstName": controller.FirstName,
				"LastName":  controller.LastName,
			},
		)
		if err != nil {
			log.Errorf("Failed to send email to %s: %s", controller.Email, err)
		}
	}

	log.Info("Finished activity job")
}
