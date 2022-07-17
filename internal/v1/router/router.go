package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kzdv/api/pkg/logger"
)

var routeGroups map[string](func(*gin.RouterGroup))

var log = logger.Logger.WithField("component", "router/v1")

func init() {
	routeGroups = make(map[string]func(*gin.RouterGroup))
}

func SetupRoutes(router *gin.Engine) {
	router.Group("/v1")
	{
		for prefix, f := range routeGroups {
			grp := router.Group(prefix)
			{
				f(grp)
			}
		}
	}
}

func AddGroup(prefix string, f func(*gin.RouterGroup)) {
	if _, exists := routeGroups[prefix]; exists {
		log.Warnf("Route group prefix %s defined already but got request to add again... this will overwrite the previous definition", prefix)
	}
	routeGroups[prefix] = f
}
