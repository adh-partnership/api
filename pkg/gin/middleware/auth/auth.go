package auth

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/adh-partnership/api/pkg/auth"
	"github.com/adh-partnership/api/pkg/database"
	"github.com/adh-partnership/api/pkg/database/models"
	"github.com/adh-partnership/api/pkg/gin/response"
	"github.com/adh-partnership/api/pkg/logger"
	"github.com/adh-partnership/api/pkg/utils"
)

var log = logger.Logger.WithField("component", "middleware/auth")

func Auth(c *gin.Context) {
	session := sessions.Default(c)
	cid := session.Get("cid")
	log.Debugf("Cookie cid: %v", utils.DumpToJSON(cid))
	if cid == nil {
		log.Debug("In Auth as Guest")
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
