package user

import (
	"github.com/gin-gonic/gin"

	"github.com/kzdv/api/pkg/gin/middleware/auth"
	"github.com/kzdv/api/pkg/logger"
)

var log = logger.Logger.WithField("component", "user")

func Routes(r *gin.RouterGroup) {
	r.GET("/login", getLogin)
	r.GET("/login/callback", getLoginCallback)

	r.GET("/", auth.NotGuest, getUser)
	r.GET("/:cid", getUser)
	r.PATCH("/", auth.NotGuest, patchUser)
	r.PATCH("/:cid", auth.NotGuest, patchUser)

	r.GET("/roles", auth.NotGuest, getUserRoles)
	r.GET("/:cid/roles", getUserRoles)
	r.PUT("/:cid/roles/:role", auth.NotGuest, putUserRoles)
	r.DELETE("/:cid/roles/:role", auth.NotGuest, deleteUserRoles)
}
