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

package user

import (
	"github.com/gin-gonic/gin"

	"github.com/adh-partnership/api/pkg/gin/middleware/auth"
	"github.com/adh-partnership/api/pkg/logger"
)

var log = logger.Logger.WithField("component", "v1/user")

func Routes(r *gin.RouterGroup) {
	r.GET("/discord/link", getDiscordLink)
	r.GET("/discord/callback", auth.NotGuest, getDiscordCallback)

	r.GET("/login", getLogin)
	r.GET("/login/callback", getLoginCallback)
	r.GET("/logout", auth.NotGuest, getLogout)

	r.GET("/", auth.NotGuest, getUser)
	r.GET("/:cid", getUser)
	r.PATCH("/", auth.NotGuest, patchUser)
	r.PATCH("/:cid", auth.NotGuest, patchUser)

	r.GET("/visitor", auth.NotGuest, getVisitor)
	r.POST("/visitor", auth.NotGuest, postVisitor)
	r.PUT("/visitor/:id", auth.NotGuest, auth.InGroup("admin"), putVisitor)
	r.GET("/visitor/eligible", auth.NotGuest, getVisitorEligibility)

	r.GET("/all", getFullRoster)
	r.GET("/roster", getRoster)
	r.GET("/staff", getStaff)

	r.GET("/roles", auth.NotGuest, getUserRoles)
	r.GET("/:cid/roles", getUserRoles)
	r.PUT("/:cid/roles/:role", auth.NotGuest, putUserRoles)
	r.DELETE("/:cid/roles/:role", auth.NotGuest, deleteUserRoles)
}
