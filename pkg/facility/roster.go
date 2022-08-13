package facility

import (
	"fmt"

	"gorm.io/gorm"

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
		if err != nil && err != gorm.ErrRecordNotFound {
			log.Errorf("Error finding user by CID %d: %s", controller.CID, err)
			continue
		}
		if err == gorm.ErrRecordNotFound {
			log.Infof("New user on roster: %d", controller.CID)
			create = true
			user = &dbTypes.User{
				CID:              uint(controller.CID),
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
			location, err := global.GetLocation(fmt.Sprint(controller.CID))
			if err != nil {
				log.Errorf("Error getting location for %d: %s", controller.CID, err.Error())
			} else {
				user.Region = location.Region
				user.Division = location.Division
				user.Subdivision = location.Subdivision
			}
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
