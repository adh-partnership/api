package event

import (
	"github.com/gin-gonic/gin"

	"github.com/adh-partnership/api/pkg/gin/middleware/auth"
	"github.com/adh-partnership/api/pkg/logger"
)

var log = logger.Logger.WithField("component", "events")

func Routes(r *gin.RouterGroup) {
	r.GET("", getEvents)
	r.GET("/:id", getEvent)
	r.POST("", auth.NotGuest, auth.InGroup("events"), postEvent)
	r.PATCH(":id", auth.NotGuest, auth.InGroup("events"), patchEvent)
	r.DELETE(":id", auth.NotGuest, auth.InGroup("events"), deleteEvent)

	r.GET("/:id/positions", getEventPositions)
	r.POST("/:id/positions", auth.NotGuest, auth.InGroup("events"), addEventPosition)
	r.PUT("/:id/positions/:position", auth.NotGuest, auth.InGroup("events"), updateEventPosition)
	r.DELETE("/:id/positions/:position", auth.NotGuest, auth.InGroup("events"), deleteEventPosition)

	r.POST("/:id/signup", auth.NotGuest, postEventSignup)
	r.DELETE("/:id/signup", auth.NotGuest, deleteEventSignup)
}
