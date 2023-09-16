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
		r.GET("/zdv/requests/schedules", auth.NotGuest, getTrainingSchedules)
		r.POST("/zdv/requests/schedules", auth.NotGuest, auth.InGroup("training"), postTrainingSchedule)
		r.PUT("/zdv/requests/schedules/:id", auth.NotGuest, auth.InGroup("training"), putTrainingSchedule)
		r.GET("/zdv/requests/sessions", auth.NotGuest, getTrainingSessions)
		r.POST("/zdv/requests/sessions", auth.NotGuest, postTrainingSession)
		r.PUT("/zdv/requests/sessions/:id", auth.NotGuest, putTrainingSession)
		r.GET("/zdv/requests/ratings/:id", auth.NotGuest, auth.HasRole("ta"), getTeacherTrainingRating)
		r.PUT("/zdv/requests/ratings/:id", auth.NotGuest, auth.HasRole("ta"), putTeacherTrainingRating)
	}
}
