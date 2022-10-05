package roster

import (
	"fmt"

	"github.com/go-co-op/gocron"
	gonanoid "github.com/matoous/go-nanoid/v2"

	"github.com/adh-partnership/api/pkg/database"
	dbTypes "github.com/adh-partnership/api/pkg/database/types"
	"github.com/adh-partnership/api/pkg/facility"
	"github.com/adh-partnership/api/pkg/logger"
	"github.com/adh-partnership/api/pkg/network/global"
	"github.com/adh-partnership/api/pkg/network/vatusa"
)

var log = logger.Logger.WithField("component", "job/roster")

func ScheduleJobs(s *gocron.Scheduler) error {
	_, err := s.Cron("1,11,21,31,41,51 * * * *").Do(UpdateRoster)
	if err != nil {
		log.Errorf("Error scheduling UpdateRoster: %s", err)
		return err
	}

	_, err = s.Every(1).Day().At("00:00").Do(UpdateForeignRoster)
	if err != nil {
		log.Errorf("Error scheduling UpdateForeignRoster: %s", err)
		return err
	}

	return nil
}

func UpdateRoster() error {
	controllers, err := vatusa.GetFacilityRoster("both")
	if err != nil {
		return err
	}

	updateid, _ := gonanoid.New(24)
	err = facility.UpdateControllerRoster(controllers, updateid)
	if err != nil {
		return err
	}

	// Users not part of the VATUSA roster will be removed from our roster
	if err := database.DB.Model(&dbTypes.User{}).Not(dbTypes.User{UpdateID: updateid}).
		Updates(dbTypes.User{
			ControllerType: dbTypes.ControllerTypeOptions["none"],
			UpdateID:       updateid,
		}).Error; err != nil {
		return err
	}

	return nil
}

func UpdateForeignRoster() {
	// Update foreign visitors
	var users []dbTypes.User
	if err := database.DB.
		Where(dbTypes.User{ControllerType: dbTypes.ControllerTypeOptions["visit"]}).
		Not(dbTypes.User{Region: "AMAS", Division: "USA"}).Find(&users).Error; err != nil {
		log.Errorf("Error getting foreign visitors: %s", err)
	}
	for _, user := range users {
		location, err := global.GetLocation(fmt.Sprint(user.CID))
		if err != nil {
			log.Errorf("Error getting location for user %d: %s", user.CID, err)
			continue
		}
		user.Region = location.Region
		user.Division = location.Division
		user.Subdivision = location.Subdivision
		if err := database.DB.Save(&user).Error; err != nil {
			log.Errorf("Error saving user %d: %s", user.CID, err)
		}
	}
}
