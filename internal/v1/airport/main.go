package airport

import (
	"github.com/gin-gonic/gin"

	"github.com/adh-partnership/api/pkg/logger"
)

var log = logger.Logger.WithField("component", "admin")

func Routes(r *gin.RouterGroup) {
	r.GET("/:center", getCenter)
	r.GET("/:center/:id", getAirport)
	r.GET("/:center/:id/atc", getAirportATC)
	r.GET("/:center/:id/charts", getAirportCharts)
	r.GET("/:center/:id/weather", getAirportWeather)
}
