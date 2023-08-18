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
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	gonanoid "github.com/matoous/go-nanoid/v2"

	"github.com/adh-partnership/api/pkg/database"
	"github.com/adh-partnership/api/pkg/database/models"
	"github.com/adh-partnership/api/pkg/gin/response"
	"github.com/adh-partnership/api/pkg/oauth"
)

// Link account and Discord account
// @Summary Link account and Discord account
// @Tags user, oauth
// @Param redirect path string false "Redirect URL"
// @Success 307
// @Failure 500 {object} response.R
// @Router /v1/user/discord/link [GET]
func getDiscordLink(c *gin.Context) {
	state, _ := gonanoid.Generate("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", 64)
	session := sessions.Default(c)
	session.Set("state", state)
	session.Set("discord-redirect", c.Query("redirect"))
	_ = session.Save()

	u := c.Request.URL

	// They are a guest, so we'll redirect to Login with a redirect back to here
	if c.GetBool("x-guest") {
		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/v1/user/login?redirect=%s/v1/user/discord/link", u.Scheme+"://"+u.Host, u.Scheme+"://"+u.Host))
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, oauth.DiscordOAuthConfig.AuthCodeURL(state))
}

// This is not the entirety of the user structure, but compromises the only field(s) we need
type DiscordUser struct {
	ID string `json:"id"`
}

// Discord Callback
// @Summary Discord Callback
// @Tags user, oauth
// @Success 307
// @Success 200 {object} response.R
// @Failure 400 {object} response.R
// @Failure 403 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/user/discord/callback [GET]
func getDiscordCallback(c *gin.Context) {
	session := sessions.Default(c)
	state := session.Get("state")
	if state == nil {
		log.Warn("State is nil")
		response.RespondError(c, http.StatusForbidden, "Forbidden")
		return
	}
	if state != c.Query("state") {
		log.Warnf("State is not equal: %s != %s", state, c.Query("state"))
		response.RespondError(c, http.StatusForbidden, "Forbidden")
		return
	}
	token, err := oauth.DiscordOAuthConfig.Exchange(c.Request.Context(), c.Query("code"))
	if err != nil {
		log.Warnf("Error exchanging code for token: %s", err.Error())
		response.RespondError(c, http.StatusBadRequest, "Bad Request")
		return
	}
	res, err := http.NewRequest("GET", "https://discord.com/api/v10/users/@me", nil)
	res.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	res.Header.Add("Accept", "application/json")
	res.Header.Add("User-Agent", "adh-partnership-api")
	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	client := &http.Client{}
	resp, err := client.Do(res)
	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	contents, err := io.ReadAll(resp.Body)
	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if resp.StatusCode >= 299 {
		log.Warnf("Error getting user info: %s, %s", resp.Status, string(contents))
		response.RespondError(c, http.StatusForbidden, "Forbidden")
		return
	}

	discorduser := &DiscordUser{}
	if err := json.Unmarshal(contents, &discorduser); err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	user := c.MustGet("x-user").(*models.User)
	user.DiscordID = discorduser.ID
	if err := database.DB.Save(user).Error; err != nil {
		log.Errorf("Error saving user %d discord id: %+v", user.CID, err)
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	redirect := session.Get("discord-redirect").(string)
	if redirect != "" {
		c.Redirect(http.StatusTemporaryRedirect, redirect)
		return
	}
	response.RespondMessage(c, http.StatusOK, "Account Linked")
}
