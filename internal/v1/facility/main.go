package facility

import (
	"github.com/gin-gonic/gin"

	"github.com/kzdv/api/pkg/logger"
)

var log = logger.Logger.WithField("component", "v1/facility")

func Routes(r *gin.RouterGroup) {
	r.GET("/roster", getRoster)
	r.GET("/staff", getStaff)
}
