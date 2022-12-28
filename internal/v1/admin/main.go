package admin

import (
	"github.com/gin-gonic/gin"

	"github.com/adh-partnership/api/pkg/gin/middleware/auth"
	"github.com/adh-partnership/api/pkg/logger"
)

var log = logger.Logger.WithField("component", "admin")

func Routes(r *gin.RouterGroup) {
	r.GET("/logging", auth.NotGuest, auth.InGroup("admin"), getLogLevel)
}
