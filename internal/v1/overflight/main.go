package overflight

import (
	"github.com/gin-gonic/gin"

	"github.com/kzdv/api/pkg/logger"
)

var log = logger.Logger.WithField("component", "overflight")

func Routes(r *gin.RouterGroup) {
	r.GET("", getOverflights)
	r.GET("/:facility", getOverflights)
}
