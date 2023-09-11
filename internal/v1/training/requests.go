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
	"net/http"

	"github.com/adh-partnership/api/pkg/database"
	"github.com/adh-partnership/api/pkg/database/dto"
	"github.com/adh-partnership/api/pkg/gin/response"
	"github.com/gin-gonic/gin"
)

func getTrainingSchedules(c *gin.Context) {
	schedules, err := database.FindTrainingSchedules()
	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	response.Respond(c, http.StatusOK, dto.ConvertTrainingSchedulesToDTO(schedules))
}

func postTrainingSchedule(c *gin.Context) {
	//
}

func putTrainingSchedule(c *gin.Context) {
	//
}

func getTrainingSessions(c *gin.Context) {
	//
}

func postTrainingSession(c *gin.Context) {
	//
}

func putTrainingSession(c *gin.Context) {
	//
}

func getTeacherTrainingRating(c *gin.Context) {
	//
}

func putTeacherTrainingRating(c *gin.Context) {
	//
}
