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

package auth

import (
	"net/http"
	"regexp"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/adh-partnership/api/pkg/auth"
	"github.com/adh-partnership/api/pkg/database"
	"github.com/adh-partnership/api/pkg/database/models"
	"github.com/adh-partnership/api/pkg/gin/response"
	"github.com/adh-partnership/api/pkg/logger"
	"github.com/adh-partnership/api/pkg/utils"
)

var log = logger.Logger.WithField("component", "middleware/auth")
var tokenHeader = regexp.MustCompile(`^Token\s+(.+)$`)

func Auth(c *gin.Context) {
	// Check for an Authorization header
	authHeader := c.GetHeader("Authorization")
	xApiToken := c.GetHeader("X-Api-Token")

	if authHeader != "" || xApiToken != "" {
		var apikey *models.APIKeys
		var err error
		// We have an API Key
		if xApiToken != "" {
			log.Debugf("X-Api-Token: %s", xApiToken)
			apikey, err = database.FindAPIKey(xApiToken)
		} else {
			if tokenHeader.MatchString(authHeader) {
				token := tokenHeader.FindStringSubmatch(authHeader)[1]
				log.Debugf("Authorization Token: %s", token)
				apikey, err = database.FindAPIKey(token)
			}
		}

		if err == gorm.ErrRecordNotFound {
			log.Debugf("API Key not found: %s", err)
			response.RespondError(c, http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		if err != nil {
			log.Errorf("Error finding API Key: %s", err)
			response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
			c.Abort()
			return
		}

		log.Debugf("API Key: %+v", apikey)
		rating, err := database.FindRating(1)
		if err != nil {
			rating = &models.Rating{
				ID:    1,
				Short: "OBS",
				Long:  "Observer",
			}
		}

		user := &models.User{
			CID:               1000,
			FirstName:         "Automation",
			LastName:          "User",
			Email:             "",
			OperatingInitials: "",
			Status:            "active",
			RatingID:          rating.ID,
			Rating:            *rating,
			DiscordID:         "",
			RosterJoinDate:    &apikey.CreatedAt,
			CreatedAt:         apikey.CreatedAt,
			UpdatedAt:         apikey.CreatedAt,
		}
		c.Set("x-guest", false)
		c.Set("x-user", user)
		c.Set("x-auth-type", "apikey")
		c.Set("x-cid", user.CID)
		c.Next()
		return
	}

	// @TODO Add JWT support here

	session := sessions.Default(c)
	cid := session.Get("cid")
	log.Debugf("Cookie cid: %v", utils.DumpToJSON(cid))
	if cid == nil {
		c.Set("x-guest", true)
		c.Next()
		return
	}

	user, err := database.FindUserByCID(cid.(string))
	if err == nil {
		log.Debugf("User: %+v", user)
		c.Set("x-guest", false)
		c.Set("x-cid", cid.(string))
		c.Set("x-user", user)
		c.Set("x-auth-type", "cookie")
		c.Next()
		return
	}

	log.Debugf("User does not exist??? %+v", err)
	// If we get here, they had a cookie with an invalid user
	// so delete it.
	session.Delete("cid")
	c.Set("x-guest", true)
	c.Next()
}

func NotGuest(c *gin.Context) {
	if c.GetBool("x-guest") {
		log.Debug("In NotGuest as Guest")
		response.RespondError(c, http.StatusUnauthorized, "Unauthorized")
		c.Abort()
		return
	}
	c.Next()
}

func IsAPIKey(c *gin.Context) bool {
	return c.GetString("x-auth-type") == "apikey"
}

func HasRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("x-user").(*models.User)
		if auth.HasRoleList(user, roles) || auth.InGroup(user, "admin") {
			c.Next()
			return
		}
		response.RespondError(c, http.StatusForbidden, "Forbidden")
		c.Abort()
	}
}

func InGroup(group ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("x-user").(*models.User)
		for _, g := range group {
			if auth.InGroup(user, g) || auth.InGroup(user, "admin") {
				c.Next()
				return
			}
		}

		response.RespondError(c, http.StatusForbidden, "Forbidden")
		c.Abort()
	}
}
