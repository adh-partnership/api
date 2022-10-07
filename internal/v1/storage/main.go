package storage

import (
	"github.com/gin-gonic/gin"

	"github.com/adh-partnership/api/pkg/gin/middleware/auth"
	"github.com/adh-partnership/api/pkg/logger"
)

var log = logger.Logger.WithField("component", "storage")

func Routes(r *gin.RouterGroup) {
	r.GET("/*category", getStorage)
	r.POST("", auth.NotGuest, auth.InGroup("files"), postStorage)
	r.PUT("/:id", auth.NotGuest, auth.InGroup("files"), putStorage)
	r.DELETE("/:id", auth.NotGuest, auth.InGroup("files"), deleteStorage)

	r.PUT("/:id/file", auth.NotGuest, auth.InGroup("files"), putStorageFile)
}

func SetBase(b string) {
	base = b
}
