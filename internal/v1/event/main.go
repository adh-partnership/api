/*
 * Copyright ADH Partnership
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

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

	r.GET("/user/:id/stats", auth.NotGuest, getEventTracking)
	r.PUT("/user/:id/stats", auth.NotGuest, auth.HasRole("ec"), updateEventTracking)
}
