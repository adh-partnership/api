package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kzdv/api/pkg/logger"
)

var routeGroups map[string]func(*gin.RouterGroup)

var log = logger.Logger.WithField("component", "router/v1")

func init() {
	routeGroups = make(map[string]func(*gin.RouterGroup))
}

func SetupRoutes(router *gin.Engine) {
	v1 := router.Group("/v1")
	for prefix, f := range routeGroups {
		grp := v1.Group(prefix)
		f(grp)
	}

	// Setup redirect for old overflight endpoint
	router.GET("/live/{fac}", func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, "/v1/overflight/"+c.Param("fac"))
	})
}

func AddGroup(prefix string, f func(*gin.RouterGroup)) {
	if _, exists := routeGroups[prefix]; exists {
		log.Warnf("Route group prefix %s defined already but got request to add again... this will overwrite the previous definition", prefix)
	}
	routeGroups[prefix] = f
}
