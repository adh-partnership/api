package staffing

import (
	"github.com/gin-gonic/gin"

	"github.com/adh-partnership/api/pkg/gin/middleware/auth"
	"github.com/adh-partnership/api/pkg/logger"
)

var log = logger.Logger.WithField("component", "staffing")

func Routes(r *gin.RouterGroup) {
	r.POST("", auth.NotGuest, requestStaffing)
}
