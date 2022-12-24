package roster

import (
	"fmt"

	"github.com/go-co-op/gocron"
	gonanoid "github.com/matoous/go-nanoid/v2"

	"github.com/adh-partnership/api/pkg/database"
	"github.com/adh-partnership/api/pkg/database/models"
	"github.com/adh-partnership/api/pkg/database/models/constants"
	"github.com/adh-partnership/api/pkg/discord"
	"github.com/adh-partnership/api/pkg/facility"
	"github.com/adh-partnership/api/pkg/logger"
	"github.com/adh-partnership/api/pkg/network/global"
	"github.com/adh-partnership/api/pkg/network/vatusa"
)

var log = logger.Logger.WithField("component", "job/roster")

func ScheduleJobs(s *gocron.Scheduler) error {
	_, err := s.Cron("1,11,21,31,41,51 * * * *").SingletonMode().Do(UpdateRoster)
	if err != nil {
		log.Errorf("Error scheduling UpdateRoster: %s", err)
		return err
	}

	_, err = s.Every(1).Day().At("00:00").Do(UpdateForeignRoster)
	if err != nil {
		log.Errorf("Error scheduling UpdateForeignRoster: %s", err)
		return err
	}

	_, err = s.Every(1).Day().At("08:00").Do(NagJob)
	if err != nil {
		log.Errorf("Error scheduling NagJob: %s", err)
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

	// Cleanup operating initials from controllers that are gone
	if err := database.DB.Model(&models.User{}).Where(models.User{ControllerType: constants.ControllerTypeNone}).
		Updates(map[string]interface{}{
			"operating_initials": "",
		}).Error; err != nil {
		return err
	}

	// Users not part of the VATUSA roster will be removed from our roster
	if err := database.DB.Model(&models.User{}).Not(models.User{UpdateID: updateid}).
		Or("update_id IS NULL").
		Updates(models.User{
			ControllerType: constants.ControllerTypeNone,
			UpdateID:       updateid,
		}).Error; err != nil {
		return err
	}

	return nil
}

func UpdateForeignRoster() {
	// Update foreign visitors
	var users []models.User
	if err := database.DB.
		Where(models.User{ControllerType: models.ControllerTypeOptions["visit"]}).
		Not(models.User{Region: "AMAS", Division: "USA"}).Find(&users).Error; err != nil {
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

func NagJob() {
	// Get users that are home/visitor and have no operating initials
	var users []models.User
	if err := database.DB.
		Where("controller_type = ? OR controller_type = ?", constants.ControllerTypeHome, constants.ControllerTypeVisitor).
		Where("operating_initials = ?", "").
		Find(&users).Error; err != nil {
		log.Errorf("Error getting users to nag: %s", err)
		return
	}

	// Send nag message to senior staff discord webhook
	for _, user := range users {
		_ = discord.NewMessage().SetContent(
			fmt.Sprintf("User %s %s (%d) has no operating initials", user.FirstName, user.LastName, user.CID),
		).Send("seniorstaff")
	}
}
