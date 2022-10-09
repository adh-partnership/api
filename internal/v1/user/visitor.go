package user

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/adh-partnership/api/pkg/config"
	"github.com/adh-partnership/api/pkg/database"
	"github.com/adh-partnership/api/pkg/database/models"
	"github.com/adh-partnership/api/pkg/database/models/constants"
	"github.com/adh-partnership/api/pkg/discord"
	"github.com/adh-partnership/api/pkg/email"
	"github.com/adh-partnership/api/pkg/gin/response"
	"github.com/adh-partnership/api/pkg/network/vatsim"
	"github.com/adh-partnership/api/pkg/network/vatusa"
)

// Submit a Visitor Application
// @Summary Submit a Visitor Application
// @Description Submit a Visitor Application
// @Tags User
// @Success 204
// @Failure 401 {object} response.R
// @Failure 406 {object} response.R "Not Acceptable - Generally means doesn't meet requirements"
// @Failure 409 {object} response.R "Conflict - Generally means already applied"
// @Failure 500 {object} response.R
// @Router /user/visitor [post]
func postVisitor(c *gin.Context) {
	user := c.MustGet("x-user").(*models.User)

	if user.Status != constants.ControllerStatusNone {
		response.RespondError(c, http.StatusConflict, "You are already a controller")
		return
	}

	if !isEligibleVisiting(user) {
		response.RespondError(c, http.StatusNotAcceptable, "You are not eligible to apply for visiting")
		return
	}

	app := &models.VisitorApplication{}
	if err := database.DB.Find(&models.VisitorApplication{UserID: user.CID}).First(&app).Error; err == nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
			return
		}
	}

	if app != nil {
		response.RespondError(c, http.StatusConflict, "Already applied")
		return
	}

	app = &models.VisitorApplication{
		User: user,
	}

	go func() {
		_ = discord.SendWebhookMessage(
			config.Cfg.Facility.Visiting.DiscordWebhookName,
			"Web API",
			fmt.Sprintf("New visiting application from %s %s (%d) [%s]", user.FirstName, user.LastName, user.CID, user.Rating.Short),
		)
	}()

	response.Respond(c, http.StatusNoContent, nil)
}

// Handle Visitor Application
// @Summary Handle Visitor Application
// @Description Handle Visitor Application
// @Tags User
// @Param id path int true "Visitor CID"
// @Param action body string true "Action to take (accept, deny)"
// @Param reason body string false "Reason for action for denials"
// @Success 204
// @Failure 401 {object} response.R
// @Failure 403 {object} response.R
// @Failure 404 {object} response.R
// @Failure 500 {object} response.R
// @Router /user/visitor/{id} [put]
func putVisitor(c *gin.Context) {
	var app models.VisitorApplication
	if err := database.DB.Find(&models.VisitorApplication{UserID: database.Atou(c.Param("id"))}).First(&app).Error; err != nil {
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
		response.RespondError(c, http.StatusNotAcceptable, "Reason required for denials")
		return
	}

	switch act.Action {
	case "accept":
		go func() {
			err := email.Send(
				app.User.Email,
				"",
				"Visitor Application Accepted",
				[]string{},
				"visitor_accepted",
				map[string]interface{}{
					"FirstName": app.User.FirstName,
					"LastName":  app.User.LastName,
					"Rating":    app.User.Rating.Short,
				},
			)
			if err != nil {
				log.Errorf("Error sending visitor accepted email to %s: %s", app.User.Email, err)
			}
			err = discord.SendWebhookMessage(
				config.Cfg.Facility.Visiting.DiscordWebhookName,
				"Web API",
				fmt.Sprintf("Visitor application accepted for %s %s (%d) [%s]", app.User.FirstName, app.User.LastName, app.User.CID, app.User.Rating.Short),
			)
			if err != nil {
				log.Errorf("Error sending visitor accepted Discord message: %s", err)
			}
		}()
		status, err := vatusa.AddVisitingController(fmt.Sprint(app.User.CID))
		if err != nil || status > 299 {
			log.Errorf("Error adding visiting controller to VATUSA for %d: %s", app.User.CID, err)
			err = discord.SendWebhookMessage(
				config.Cfg.Facility.Visiting.DiscordWebhookName,
				"Web API",
				fmt.Sprintf("Error adding visiting controller to VATUSA for cid: %d: %s -- Please add manually!", app.User.CID, err),
			)
			if err != nil {
				log.Errorf("Error sending visiting VATUSA error Discord message: %s", err)
			}
			response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		app.User.ControllerType = constants.ControllerTypeVisitor
		if err := database.DB.Save(&app.User).Error; err != nil {
			log.Errorf("Error updating user controller type to visitor for %d: %s", app.User.CID, err)
		}
	case "deny":
		go func() {
			err := email.Send(
				app.User.Email,
				"",
				"Visitor Application Denied",
				[]string{},
				"visitor_denied",
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
			err = discord.SendWebhookMessage(
				config.Cfg.Facility.Visiting.DiscordWebhookName,
				"Web API",
				fmt.Sprintf("Visitor application denied for %s %s (%d) [%s]", app.User.FirstName, app.User.LastName, app.User.CID, app.User.Rating.Short),
			)
			if err != nil {
				log.Errorf("Error sending visitor denied Discord message: %s", err)
			}
		}()
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

	if rating, _ := database.FindRatingByShort(config.Cfg.Facility.Visiting.MinRating); user.Rating.ID < rating.ID {
		return false
	}

	ratechange, err := vatsim.GetDateOfRatingChange(fmt.Sprint(user.CID))
	if err != nil {
		log.Errorf("Error getting date of rating change: %s", err)
		return false
	}

	// Check that ratechange is more than 90 days ago
	if ratechange.After(time.Now().AddDate(0, 0, -90)) {
		return false
	}

	return true
}
