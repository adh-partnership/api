package stats

import "github.com/gin-gonic/gin"

func Routes(r *gin.RouterGroup) {
	r.GET("/online", getOnlineATC)
	r.GET("/historical", getHistoricalStats)
	r.GET("/historical/:year/:month", getHistoricalStats)
}
