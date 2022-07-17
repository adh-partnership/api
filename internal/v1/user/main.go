package user

import (
	"github.com/gin-gonic/gin"
	"github.com/kzdv/api/internal/v1/router"
)

func init() {
	router.AddGroup("/user", routes)
}

func routes(router *gin.RouterGroup) {
	router.GET("/login", GetLogin)
	router.GET("/login/callback", GetLoginCallback)
}
