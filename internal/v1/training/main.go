package training

import (
	"github.com/gin-gonic/gin"

	"github.com/adh-partnership/api/pkg/config"
	"github.com/adh-partnership/api/pkg/gin/middleware/auth"
	"github.com/adh-partnership/api/pkg/logger"
)

var log = logger.Logger.WithField("component", "training")

func Routes(r *gin.RouterGroup) {
	r.GET("/:cid", auth.NotGuest, getTraining)
	r.POST("/:cid", auth.NotGuest, auth.InGroup("training"), postTraining)
	r.PUT("/:cid/:id", auth.NotGuest, auth.InGroup("training"), putTraining)
	r.DELETE("/:cid/:id", auth.NotGuest, auth.InGroup("training"), deleteTraining)

	if config.Cfg.Facility.TrainingRequests.Enabled {
		r.GET("/sessions", auth.NotGuest, getSessions)
		r.POST("/sessions", auth.NotGuest, postSessions)
	}
}
