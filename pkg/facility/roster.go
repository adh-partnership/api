package facility

import (
	"fmt"
	"time"

	"github.com/kzdv/api/pkg/database"
	dbTypes "github.com/kzdv/api/pkg/database/types"
	"github.com/kzdv/api/pkg/discord"
	"github.com/kzdv/api/pkg/logger"
	"github.com/kzdv/api/pkg/network/global"
	"github.com/kzdv/api/pkg/network/vatusa"
)

var log = logger.Logger.WithField("component", "facility")

func UpdateControllerRoster(controllers []vatusa.VATUSAController, updateid string) error {
	for _, controller := range controllers {
		create := false
		user, err := database.FindUserByCID(fmt.Sprint(controller.CID))
		if err != nil {
			log.Errorf("Error finding user by CID %d: %s", controller.CID, err)
			continue
		}
		if user == nil {
			log.Infof("New user on roster: %d", controller.CID)
			create = true
			user = &dbTypes.User{
				CID:              uint(controller.CID),
				FirstName:        controller.FirstName,
				LastName:         controller.LastName,
				ControllerType:   dbTypes.ControllerTypeOptions["none"],
				DelCertification: dbTypes.CertificationOptions["none"],
				GndCertification: dbTypes.CertificationOptions["none"],
				LclCertification: dbTypes.CertificationOptions["none"],
				AppCertification: dbTypes.CertificationOptions["none"],
				CtrCertification: dbTypes.CertificationOptions["none"],
			}
			oi, err := database.FindOI(user)
			if err != nil {
				log.Infof("Error generating new OI: %s", err.Error())
				oi = ""
			}
			if oi == "" {
				go func() {
					msg := fmt.Sprintf("New user on roster, %s %s (%d), needs to be assigned an OI", user.FirstName, user.LastName, user.CID)
					err := discord.SendWebhookMessage(
						"seniorstaff",
						"KZDV Web API",
						msg,
					)
					if err != nil {
						log.Errorf("Error sending discord message (%s): %s", msg, err.Error())
						return
					}
				}()
			}
			user.OperatingInitials = oi
		}
		user.FirstName = controller.FirstName
		user.LastName = controller.LastName
		user.Email = controller.Email
		user.RatingID = controller.Rating
		user.UpdateID = updateid

		if controller.Membership == "visit" {
			if controller.Facility != "ZZN" {
				user.Region = "AMAS"
				user.Division = "USA"
				user.Subdivision = controller.Facility

				if controller.Facility == "ZAE" && isInDailyCheck() {
					err := discord.SendWebhookMessage("seniorstaff", "KZDV Web API", fmt.Sprintf("%s %s (%d) (%s) is a visitor, but is in %s, %s, %s -- verify eligibility",
						user.FirstName, user.LastName, user.CID, controller.RatingShort, user.Region, user.Division, user.Subdivision))
					if err != nil {
						log.Errorf("Error sending discord message: %s", err.Error())
						return err
					}
				}
			} else {
				location, err := global.GetLocation(fmt.Sprint(controller.CID))
				if err != nil {
					log.Errorf("Error getting location for %d: %s", controller.CID, err.Error())
				} else {
					user.Region = location.Region
					user.Division = location.Division
					user.Subdivision = location.Subdivision

					// This in theory shouldn't happen, but flag if it does so we can raise to division
					if user.Region == "AMAS" && user.Division == "USA" && controller.Facility == "ZZN" && isInDailyCheck() {
						log.Infof("%s %s (%d) (%s) is a visitor, VATSIM API indicates they are in %s, %s, %s "+
							"but VATUSA has them in a non-member facility (ZZN) -- verify eligibility and raise to VATUSA's Tech Manager as this should not happen (unless they "+
							"JUST transferred into VATUSA and the div sync job hasn't run yet)",
							user.FirstName, user.LastName, user.CID, controller.RatingShort, user.Region, user.Division, user.Subdivision)

						err := discord.SendWebhookMessage("seniorstaff", "KZDV Web API", fmt.Sprintf("%s %s (%d) (%s) is a visitor, VATSIM API indicates they are in %s, %s, %s "+
							"but VATUSA has them in a non-member facility (ZZN) -- verify eligibility and raise to VATUSA's Tech Manager as this should not happen (unless they "+
							"JUST transferred into VATUSA and the div sync job hasn't run yet)",
							user.FirstName, user.LastName, user.CID, controller.RatingShort, user.Region, user.Division, user.Subdivision))
						if err != nil && err != discord.ErrWebhookNotConfigured && err != discord.ErrUsedDefaultWebhook {
							log.Errorf("Error sending discord message: %s", err.Error())
							return err
						} else if err != nil {
							log.Warnf("Error sending discord message: %s", err.Error())
						}
					}
				}
			}
			user.Status = dbTypes.ControllerTypeOptions["visitor"]
		} else if controller.Membership == "home" {
			user.Status = dbTypes.ControllerTypeOptions["home"]
		} else {
			// This shouldn't happen... but...
			user.Status = dbTypes.ControllerTypeOptions["none"]
		}

		if create {
			if err := database.DB.Create(user).Error; err != nil {
				log.Errorf("Error creating user: %s", err.Error())
				continue
			}
		} else {
			if err := database.DB.Save(user).Error; err != nil {
				log.Errorf("Error saving user: %s", err.Error())
				continue
			}
		}
	}

	return nil
}

// This will be used to check whether or not to send flags, we don't want to harass every 10 minutes so we'll only send once per day
func isInDailyCheck() bool {
	return time.Now().Hour() == 12 && time.Now().Minute() < 10
}
