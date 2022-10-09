package email

import (
	"github.com/gin-gonic/gin"

	"github.com/adh-partnership/api/pkg/gin/middleware/auth"
	"github.com/adh-partnership/api/pkg/logger"
)

var log = logger.Logger.WithField("component", "email")

func Routes(r *gin.RouterGroup) {
	r.GET("/templates", auth.NotGuest, getTemplate)
	r.POST("/templates/:name", auth.NotGuest, postTemplate)
}
