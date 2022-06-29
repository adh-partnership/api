package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kzdv/api/internal/v1/user"
)

func userGroup(router *gin.RouterGroup) {
	userGroup := router.Group("/user")
	{
		userGroup.GET("/login", user.GetLogin)
		userGroup.GET("/login/callback", user.GetLoginCallback)
	}
}
