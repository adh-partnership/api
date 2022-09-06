package feedback

import (
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.RouterGroup) {
	r.GET("", getFeedback)
	r.POST("", postFeedback)
	r.PATCH("", patchFeedback)
}
