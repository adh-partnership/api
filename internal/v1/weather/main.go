package weather

import (
	"github.com/gin-gonic/gin"

	"github.com/adh-partnership/api/pkg/gin/middleware/auth"
)

func Routes(r *gin.RouterGroup) {
	r.GET("/metar/:icao", getMetar)
	r.GET("/taf/:icao", getTaf)
	r.GET("/populate", auth.NotGuest, auth.InGroup("admin"), populate)
}
