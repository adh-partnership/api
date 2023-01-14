package feedback

import (
	"github.com/gin-gonic/gin"

	"github.com/adh-partnership/api/pkg/gin/middleware/auth"
)

func Routes(r *gin.RouterGroup) {
	r.GET("", getFeedback)
	r.GET("/:id", getSingleFeedback)
	r.POST("", auth.NotGuest, postFeedback)
	r.PATCH("/:id", auth.NotGuest, auth.InGroup("admin"), patchFeedback)
}
