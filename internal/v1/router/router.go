package router

import "github.com/gin-gonic/gin"

func SetupRoutes(router *gin.Engine) {
	v1 := router.Group("/v1")
	{
		userGroup(v1)
	}
}
