package proxy

import "github.com/gin-gonic/gin"

func Routes(r *gin.RouterGroup) {
	r.GET("/metar/:icao", getMetar)
}
