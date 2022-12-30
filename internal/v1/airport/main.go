package airport

import (
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.RouterGroup) {
	r.GET("/:center", getCenter)
	r.GET("/:center/:id", getAirport)
	r.GET("/:center/:id/atc", getAirportATC)
	r.GET("/:center/:id/charts", getAirportCharts)
	r.GET("/:center/:id/weather", getAirportWeather)
}
