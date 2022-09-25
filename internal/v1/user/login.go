package user

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	gonanoid "github.com/matoous/go-nanoid/v2"

	"github.com/kzdv/api/internal/v1/dto"
	"github.com/kzdv/api/pkg/config"
	"github.com/kzdv/api/pkg/gin/response"
	"github.com/kzdv/api/pkg/oauth"
	"github.com/kzdv/api/pkg/utils"
)

// Login to account
// @Summary Login to account
// @Tags user, oauth
// @Param redirect path string false "Redirect URL"
// @Success 307
// @Failure 500 {object} response.R
// @Router /v1/user/login [GET]
func getLogin(c *gin.Context) {
	state, _ := gonanoid.Generate("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", 64)
	session := sessions.Default(c)
	session.Set("state", state)
	session.Set("redirect", c.Query("redirect"))
	_ = session.Save()

	c.Redirect(http.StatusTemporaryRedirect, oauth.OAuthConfig.AuthCodeURL(state))
}

// Login callback
// @Summary Login callback
// @Tags user, oauth
// @Success 307
// @Success 200 {object} response.R
// @Failure 400 {object} response.R
// @Failure 403 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/user/login/callback [GET]
func getLoginCallback(c *gin.Context) {
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
	token, err := oauth.OAuthConfig.Exchange(c.Request.Context(), c.Query("code"))
	if err != nil {
		log.Warnf("Error exchanging code for token: %s", err.Error())
		response.RespondError(c, http.StatusBadRequest, "Bad Request")
		return
	}
	res, err := http.NewRequest("GET", fmt.Sprintf("%s%s", config.Cfg.OAuth.BaseURL, config.Cfg.OAuth.Endpoints.UserInfo), nil)
	res.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	res.Header.Add("Accept", "application/json")
	res.Header.Add("User-Agent", "kzdv-api")
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

	user := &dto.SSOUserResponse{}
	if err := json.Unmarshal(contents, &user); err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	session.Delete("state")
	session.Set("cid", fmt.Sprint(user.User.CID))
	log.Tracef("User %s logged in", utils.DumpToJSON(user.User.CID))
	_ = session.Save()

	redirect := session.Get("redirect").(string)
	if redirect != "" {
		c.Redirect(http.StatusTemporaryRedirect, redirect)
		return
	}
	response.RespondMessage(c, http.StatusOK, "Logged In")
}
