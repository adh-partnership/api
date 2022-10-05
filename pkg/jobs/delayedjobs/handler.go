package delayedjobs

import (
	"time"

	"github.com/go-co-op/gocron"

	"github.com/adh-partnership/api/pkg/database"
	dbTypes "github.com/adh-partnership/api/pkg/database/types"
	"github.com/adh-partnership/api/pkg/logger"
	"github.com/adh-partnership/api/pkg/messaging"
)

var log = logger.Logger.WithField("component", "job/delayedjobs")

func ScheduleJobs(s *gocron.Scheduler) error {
	_, err := s.Every(1).Minute().SingletonMode().Do(HandleDelayedJobs)
	if err != nil {
		log.Errorf("Error scheduling HandleDelayedJobs: %v", err)
		return err
	}

	return nil
}

func HandleDelayedJobs() {
	log.Debug("Handling delayed jobs")
	var jobs []dbTypes.DelayedJob
	if err := database.DB.Where("not_before", "<=", time.Now()).Find(&jobs).Error; err != nil {
		log.Errorf("Error handling delayed jobs: %v", err)
		return
	}

	for _, job := range jobs {
		log.Debugf("Handling delayed job: %v", job)
		err := messaging.PublishMessage(job.Queue, job.Body)
		if err != nil {
			log.Errorf("Error handling delayed job: %v", err)
			return
		}

		if err := database.DB.Delete(&job).Error; err != nil {
			log.Errorf("Error deleting delayed job: %v", err)
			return
		}
	}
}
