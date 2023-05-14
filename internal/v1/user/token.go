package user

import (
	"net/http"

	"github.com/adh-partnership/api/pkg/oauth"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

// Login to account
// @Summary Login to account
// @Tags user, oauth
// @Param redirect path string false "Redirect URL"
// @Success 307
// @Failure 500 {object} response.R
// @Router /v1/user/login [GET]
func getToken(c *gin.Context) {
	state, _ := gonanoid.Generate("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", 64)
	session := sessions.Default(c)
	session.Set("state", state)
	session.Set("redirect", c.Query("redirect"))
	_ = session.Save()

	c.Redirect(http.StatusTemporaryRedirect, oauth.OAuthConfig.AuthCodeURL(state))
}
