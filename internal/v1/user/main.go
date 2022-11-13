package user

import (
	"github.com/gin-gonic/gin"

	"github.com/adh-partnership/api/pkg/gin/middleware/auth"
	"github.com/adh-partnership/api/pkg/logger"
)

var log = logger.Logger.WithField("component", "v1/user")

func Routes(r *gin.RouterGroup) {
	r.GET("/discord/link", getDiscordLink)
	r.GET("/discord/callback", auth.NotGuest, getDiscordCallback)

	r.GET("/login", getLogin)
	r.GET("/login/callback", getLoginCallback)
	r.GET("/logout", auth.NotGuest, getLogout)

	r.GET("/", auth.NotGuest, getUser)
	r.GET("/:cid", getUser)
	r.PATCH("/", auth.NotGuest, patchUser)
	r.PATCH("/:cid", auth.NotGuest, patchUser)

	r.POST("/visitor", auth.NotGuest, postVisitor)
	r.PUT("/visitor/:id", auth.NotGuest, auth.InGroup("admin"), putVisitor)

	r.GET("/all", getFullRoster)
	r.GET("/roster", getRoster)
	r.GET("/staff", getStaff)

	r.GET("/roles", auth.NotGuest, getUserRoles)
	r.GET("/:cid/roles", getUserRoles)
	r.PUT("/:cid/roles/:role", auth.NotGuest, putUserRoles)
	r.DELETE("/:cid/roles/:role", auth.NotGuest, deleteUserRoles)
}
