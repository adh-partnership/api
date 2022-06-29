package auth

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/kzdv/api/pkg/database"
	"github.com/kzdv/api/pkg/gin/response"
)

func Auth(c *gin.Context) {
	session := sessions.Default(c)
	cid := session.Get("cid")
	if cid == nil {
		c.Set("x-guest", true)
		c.Next()
		return
	}

	user, err := database.FindUserByCID(cid.(string))
	if err != nil {
		c.Set("x-guest", false)
		c.Set("x-cid", cid)
		c.Set("x-user", user)
		c.Set("x-auth-type", "cookie")
		c.Next()
		return
	}

	// If we get here, they had a cookie with an invalid user
	// so delete it.
	session.Delete("cid")
	session.Save()
	c.Set("x-guest", true)
	c.Next()
}

func NotGuest(c *gin.Context) {
	if c.GetBool("x-guest") {
		response.RespondError(c, http.StatusForbidden, "Forbidden")
		c.Abort()
		return
	}
	c.Next()
}
