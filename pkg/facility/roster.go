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

package facility

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/adh-partnership/api/pkg/database"
	"github.com/adh-partnership/api/pkg/database/models"
	"github.com/adh-partnership/api/pkg/database/models/constants"
	"github.com/adh-partnership/api/pkg/discord"
	"github.com/adh-partnership/api/pkg/logger"
	"github.com/adh-partnership/api/pkg/network/global"
	"github.com/adh-partnership/api/pkg/network/vatusa"
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
			user = &models.User{
				CID:            uint(controller.CID),
				FirstName:      controller.FirstName,
				LastName:       controller.LastName,
				ControllerType: constants.ControllerTypeNone,
			}
			oi, err := database.FindOI(user)
			if err != nil {
				log.Infof("Error generating new OI: %s", err.Error())
				oi = ""
			}
			if oi == "" {
				go func() {
					_ = discord.NewMessage().
						SetContent(fmt.Sprintf("New user on roster, %s %s (%d), needs to be assigned an OI", user.FirstName, user.LastName, user.CID)).
						Send("seniorstaff")
				}()
			}
			user.OperatingInitials = oi
		}

		// Check if they are new on roster, previously in the table, and don't have an OI set
		if user.OperatingInitials == "" &&
			user.ControllerType == constants.ControllerTypeNone &&
			(controller.Membership == "home" || controller.Membership == "visit") {
			oi, err := database.FindOI(user)
			if err != nil {
				log.Infof("Error generating new OI: %s", err.Error())
				oi = ""
			}
			if oi == "" {
				go func() {
					_ = discord.NewMessage().
						SetContent(
							fmt.Sprintf(
								"User %s %s (%d) is back on the roster, but auto-generated OI failed as their first initial + last initial was already in use. Please assign one.",
								user.FirstName,
								user.LastName,
								user.CID,
							),
						).
						Send("seniorstaff")
				}()
			}
			user.OperatingInitials = oi
		}

		user.FirstName = controller.FirstName
		user.LastName = controller.LastName
		user.Email = controller.Email
		rating, _ := database.FindRatingByShort(controller.RatingShort)
		user.Rating = *rating
		user.RatingID = rating.ID
		user.UpdateID = updateid
		user.RosterJoinDate = &(controller.FacilityJoin)

		// If their status is none or empty, set it to active
		if user.Status == constants.ControllerStatusNone || user.Status == "" {
			_ = discord.NewMessage().SetContent(
				fmt.Sprintf("User %s %s (%d) is on our roster with no status set. Assuming active.",
					user.FirstName,
					user.LastName,
					user.CID)).
				Send("seniorstaff")
			user.Status = constants.ControllerStatusActive
		}

		switch controller.Membership {
		case "visit":
			user.ControllerType = constants.ControllerTypeVisitor
			if controller.Facility != "ZZN" {
				user.Region = "AMAS"
				user.Division = "USA"
				user.Subdivision = controller.Facility

				if controller.Facility == "ZAE" && isInDailyCheck() {
					_ = discord.NewMessage().SetContent(
						fmt.Sprintf("User %s %s (%d) is a visitor, but is in ZAE. Verify eligibility as this should not happen.",
							user.FirstName,
							user.LastName,
							user.CID)).
						Send("seniorstaff")
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

						_ = discord.NewMessage().SetContent(
							fmt.Sprintf("%s %s (%d) (%s) is a visitor, VATSIM API indicates they are in %s, %s, %s "+
								"but VATUSA has them in a non-member facility (ZZN) -- verify eligibility and raise to VATUSA's Tech Manager as this should not happen (unless they "+
								"JUST transferred into VATUSA and the div sync job hasn't run yet)",
								user.FirstName, user.LastName, user.CID, controller.RatingShort, user.Region, user.Division, user.Subdivision)).
							Send("seniorstaff")
					}
				}
			}
		case "home":
			user.Region = "AMAS"
			user.Division = "USA"
			user.Subdivision = controller.Facility
			user.ControllerType = constants.ControllerTypeHome
		default:
			// This shouldn't happen... but...
			user.ControllerType = constants.ControllerTypeNone
		}

		if create {
			if err := database.DB.Create(user).Error; err != nil {
				log.Errorf("Error creating user: %s", err.Error())
				continue
			}
		} else {
			if err := database.DB.Session(&gorm.Session{FullSaveAssociations: true}).Save(user).Error; err != nil {
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
