package feedback

import (
	"github.com/gin-gonic/gin"

	"github.com/adh-partnership/api/pkg/gin/middleware/auth"
)

func Routes(r *gin.RouterGroup) {
	r.GET("", getFeedback)
	r.POST("", auth.NotGuest, postFeedback)
	r.PATCH("", auth.NotGuest, auth.InGroup("admin"), patchFeedback)
}
