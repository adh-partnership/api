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

package user

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/adh-partnership/api/pkg/config"
	"github.com/adh-partnership/api/pkg/database"
	"github.com/adh-partnership/api/pkg/database/dto"
	"github.com/adh-partnership/api/pkg/database/models"
	"github.com/adh-partnership/api/pkg/database/models/constants"
	"github.com/adh-partnership/api/pkg/discord"
	"github.com/adh-partnership/api/pkg/email"
	"github.com/adh-partnership/api/pkg/gin/response"
	"github.com/adh-partnership/api/pkg/network/global"
	"github.com/adh-partnership/api/pkg/network/vatsim"
	"github.com/adh-partnership/api/pkg/network/vatusa"
)

type VisitingEligibility struct {
	Eligible bool `json:"eligible"`
}

// Get visiting eligibility
// @Summary Get visiting eligibility
// @Description Get visiting eligibility
// @Tags user
// @Success 200 {object} VisitingEligibility
// @Failure 401 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/user/visitor/eligible [get]
func getVisitorEligibility(c *gin.Context) {
	user := c.MustGet("x-user").(*models.User)

	response.Respond(c, http.StatusOK, &VisitingEligibility{
		Eligible: user.ControllerType == constants.ControllerTypeNone && isEligibleVisiting(user),
	})
}

// Get visiting applications
// @Summary Get visiting applications
// @Description Get visiting applications
// @Tags user
// @Success 200 {object} []models.VisitorApplication
// @Failure 401 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/user/visitor [get]
func getVisitor(c *gin.Context) {
	apps := []models.VisitorApplication{}
	if err := database.DB.Preload("User.Rating").Preload(clause.Associations).Find(&apps).Error; err != nil {
		log.Errorf("Error getting visitor applications: %s", err)
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	response.Respond(c, http.StatusOK, dto.ConvVisitorApplicationsToResponse(apps))
}

// Submit a Visitor Application
// @Summary Submit a Visitor Application
// @Description Submit a Visitor Application
// @Tags user
// @Success 204
// @Failure 401 {object} response.R
// @Failure 406 {object} response.R "Not Acceptable - Generally means doesn't meet requirements"
// @Failure 409 {object} response.R "Conflict - Generally means already applied"
// @Failure 500 {object} response.R
// @Router /v1/user/visitor [post]
func postVisitor(c *gin.Context) {
	user := c.MustGet("x-user").(*models.User)

	location, err := global.GetLocation(fmt.Sprint(user.CID))
	changed := false
	if err != nil {
		log.Errorf("Error getting location: %s", err)
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if location.Subdivision != user.Subdivision || location.Division != user.Division || location.Region != user.Region {
		user.Region = location.Region
		user.Division = location.Division
		user.Subdivision = location.Subdivision
		changed = true
	}

	rating, err := vatsim.GetRating(fmt.Sprint(user.CID))
	if err != nil {
		log.Errorf("Error getting rating: %s", err)
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if user.RatingID != rating {
		user.RatingID = rating
		changed = true
	}

	if changed {
		if err := database.DB.Save(&user).Error; err != nil {
			log.Errorf("Error updating user location and rating: %s", err)
			response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
			return
		}
	}

	if user.ControllerType != constants.ControllerTypeNone {
		response.RespondError(c, http.StatusConflict, "You are already a controller")
		return
	}

	if !isEligibleVisiting(user) {
		response.RespondError(c, http.StatusNotAcceptable, "You are not eligible to apply for visiting")
		return
	}

	app, err := database.FindVisitorApplicationByCID(fmt.Sprint(user.CID))
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Errorf("Error getting visitor application: %s", err)
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if app != nil {
		log.Infof("Visitor application already exists for %d", user.CID)
		log.Infof("Application=%+v", app)
		response.RespondError(c, http.StatusConflict, "Already applied")
		return
	}

	app = &models.VisitorApplication{
		UserID: user.CID,
		User:   user,
	}

	if err := database.DB.Create(&app).Error; err != nil {
		log.Errorf("Error creating visitor application: %s", err)
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	err = discord.NewMessage().
		SetContent("New Visiting Application").
		AddEmbed(
			discord.NewEmbed().
				SetTitle("New Visiting Application").
				SetColor(
					discord.GetColor("00", "00", "ff"),
				).
				AddField(
					discord.NewField().SetName("Name").SetValue(fmt.Sprintf("%s %s", user.FirstName, user.LastName)).SetInline(true),
				).
				AddField(
					discord.NewField().SetName("CID").SetValue(fmt.Sprintf("%d", user.CID)).SetInline(true),
				).
				AddField(
					discord.NewField().SetName("Rating").SetValue(user.Rating.Short).SetInline(true),
				).
				AddField(
					discord.NewField().SetName("Visiting From").SetValue(fmt.Sprintf("%s/%s/%s", user.Region, user.Division, user.Subdivision)).SetInline(true),
				),
		).Send("visiting_application")
	if err != nil {
		log.Errorf("Error sending discord message: %s", err)
	}

	response.Respond(c, http.StatusNoContent, nil)
}

// Handle Visitor Application
// @Summary Handle Visitor Application
// @Description Handle Visitor Application
// @Tags user
// @Param id path int true "Visitor CID"
// @Param action body string true "Action to take (accept, deny)"
// @Param reason body string false "Reason for action for denials"
// @Success 204
// @Failure 401 {object} response.R
// @Failure 403 {object} response.R
// @Failure 404 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/user/visitor/{id} [put]
func putVisitor(c *gin.Context) {
	var app models.VisitorApplication
	if err := database.DB.Preload(clause.Associations).Find(&app, database.Atou(c.Param("id"))).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.RespondError(c, http.StatusNotFound, "Not Found")
			return
		}

		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	type action struct {
		Action string `json:"action"`
		Reason string `json:"reason"`
	}
	act := &action{}
	if err := c.ShouldBind(&act); err != nil {
		response.RespondError(c, http.StatusBadRequest, "Bad Request")
		return
	}

	if act.Action == "deny" && act.Reason == "" {
		act.Reason = "No reason provided"
	}

	switch act.Action {
	case "accept":
		if config.Cfg.Facility.Visiting.SendWelcome {
			go func() {
				err := email.Send(
					app.User.Email,
					"",
					"",
					email.Templates["visiting_added"],
					map[string]interface{}{
						"FirstName": app.User.FirstName,
						"LastName":  app.User.LastName,
						"Rating":    app.User.Rating.Short,
					},
				)
				if err != nil {
					log.Errorf("Error sending visitor accepted email to %s: %s", app.User.Email, err)
				}
			}()
		}
		status, err := vatusa.AddVisitingController(fmt.Sprint(app.User.CID))
		if err != nil || status > 299 {
			log.Errorf("Error adding visiting controller to VATUSA for %d: %s", app.User.CID, err)
			err = discord.NewMessage().SetContent(
				fmt.Sprintf("Error adding visiting controller %s %s (%d) to VATUSA roster", app.User.FirstName, app.User.LastName, app.User.CID),
			).Send("visiting_application")
			if err != nil {
				log.Errorf("Error sending discord message: %s", err)
			}
			response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		app.User.ControllerType = constants.ControllerTypeVisitor
		app.User.Status = constants.ControllerStatusActive
		t := time.Now()
		app.User.RosterJoinDate = &t
		if err := database.DB.Save(&app.User).Error; err != nil {
			log.Errorf("Error updating user controller type to visitor for %d: %s", app.User.CID, err)
		}
	case "deny":
		if config.Cfg.Facility.Visiting.SendRejected {
			go func() {
				err := email.Send(
					app.User.Email,
					"",
					"",
					email.Templates["visiting_rejected"],
					map[string]interface{}{
						"FirstName": app.User.FirstName,
						"LastName":  app.User.LastName,
						"Rating":    app.User.Rating.Short,
						"Reason":    act.Reason,
					},
				)
				if err != nil {
					log.Errorf("Error sending visitor rejected email to %s: %s", app.User.Email, err)
				}
			}()
		}
	default:
		response.RespondError(c, http.StatusBadRequest, "Bad Request")
		return
	}

	if err := database.DB.Delete(&app).Error; err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	response.Respond(c, http.StatusNoContent, nil)
}

// Check if user is eligible
func isEligibleVisiting(user *models.User) bool {
	if user.ControllerType != constants.ControllerTypeNone {
		return false
	}

	// This is now hardcoded as it is defined by the division's T
	if rating, _ := database.FindRatingByShort("S3"); user.Rating.ID < rating.ID {
		return false
	}

	// Is user apart of VATUSA and in the ZAE subdivision? If so, they are not eligible
	if user.Region == "AMAS" && user.Division == "USA" && user.Subdivision == "ZAE" {
		return false
	}

	// Check VATUSA eligibility
	eligible, _, err := vatusa.IsVisitorEligible(fmt.Sprint(user.CID))
	if err != nil {
		log.Errorf("Error checking VATUSA eligibility: %s", err)
		return false
	}

	return eligible
}
