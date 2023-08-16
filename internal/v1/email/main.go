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

package email

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/adh-partnership/api/pkg/database"
	"github.com/adh-partnership/api/pkg/email"
	"github.com/adh-partnership/api/pkg/gin/middleware/auth"
	"github.com/adh-partnership/api/pkg/gin/response"
	"github.com/adh-partnership/api/pkg/logger"
)

var log = logger.Logger.WithField("component", "email")

func Routes(r *gin.RouterGroup) {
	r.GET("/test", auth.NotGuest, auth.InGroup("admin"), getTest)
}

func getTest(c *gin.Context) {
	// Get name of template
	name := c.Query("template")
	// Get target user
	cid := c.Query("cid")

	// Lookup user
	user, err := database.FindUserByCID(cid)
	if err != nil {
		response.RespondError(c, http.StatusNotFound, "User not found")
		return
	}

	// All of the templates currently only use the first and last name
	err = email.Send(
		user.Email,
		"",
		"",
		name,
		map[string]interface{}{
			"FirstName": user.FirstName,
			"LastName":  user.LastName,
		},
	)
	if err != nil {
		log.Errorf("Failed to send email: %v", err)
		response.RespondError(c, http.StatusInternalServerError, "Failed to send email")
		return
	}

	response.RespondMessage(c, http.StatusOK, "Email sent")
}
