package user

import (
	"github.com/gin-gonic/gin"
	"github.com/kzdv/api/internal/v1/router"
	"github.com/kzdv/api/pkg/gin/middleware/auth"
	"github.com/kzdv/api/pkg/logger"
)

var log = logger.Logger.WithField("component", "user")

func init() {
	router.AddGroup("/user", routes)
}

func routes(router *gin.RouterGroup) {
	router.GET("/login", getLogin)
	router.GET("/login/callback", getLoginCallback)

	router.GET("/", auth.NotGuest, getUser)
	router.GET("/:cid", getUser)
	router.PATCH("/", auth.NotGuest, patchUser)
	router.PATCH("/:cid", auth.NotGuest, patchUser)

	router.GET("/roles", auth.NotGuest, getUserRoles)
	router.GET("/:cid/roles", getUserRoles)
	router.PUT("/:cid/roles/:role", auth.NotGuest, putUserRoles)
	router.DELETE("/:cid/roles/:role", auth.NotGuest, deleteUserRoles)
}
