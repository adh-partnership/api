package overflight

import (
	"github.com/gin-gonic/gin"

	"github.com/kzdv/api/internal/v1/router"
	"github.com/kzdv/api/pkg/logger"
)

var log = logger.Logger.WithField("component", "overflight")

func init() {
	router.AddGroup("/overflight", routes)
}

func routes(router *gin.RouterGroup) {
	router.GET("", getOverflights)
	router.GET("/:facility", getOverflights)
}
