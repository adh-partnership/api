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
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/adh-partnership/api/pkg/auth"
	"github.com/adh-partnership/api/pkg/database"
	models "github.com/adh-partnership/api/pkg/database/models"
	"github.com/adh-partnership/api/pkg/discord"
	"github.com/adh-partnership/api/pkg/gin/response"
	"github.com/adh-partnership/api/pkg/memcache"
)

// Get User Information
// @Summary Get User Information
// @Tags user
// @Param cid path string false "CID, if not provided, defaults to logged in user"
// @Success 200 {object} []string
// @Failure 403 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/user/:cid/roles [GET]
func getUserRoles(c *gin.Context) {
	var err error
	user := c.MustGet("x-user").(*models.User)

	if c.Param("cid") != "" {
		user, err = database.FindUserByCID(c.Param("cid"))
		if err != nil {
			response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		if user == nil {
			response.RespondError(c, http.StatusNotFound, "User not found")
			return
		}
	}

	roles := []string{}
	for _, role := range user.Roles {
		roles = append(roles, role.Name)
	}

	response.Respond(c, http.StatusOK, roles)
}

// Add User Role
// @Summary Add User Role
// @Tags user
// @Param role path string true "Role"
// @Param cid path string true "CID"
// @Success 204
// @Failure 309 {object} response.R
// @Failure 403 {object} response.R
// @Failure 404 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/user/:cid/roles/:role [PUT]
func putUserRoles(c *gin.Context) {
	reqUser := c.MustGet("x-user").(*models.User)

	user, err := database.FindUserByCID(c.Param("cid"))
	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if user == nil {
		response.RespondError(c, http.StatusNotFound, "User not found")
		return
	}

	role := c.Param("role")

	if _, ok := auth.Roles[role]; !ok {
		response.RespondError(c, http.StatusNotFound, "Role not found")
		return
	}

	if auth.HasRole(user, role) {
		response.RespondError(c, http.StatusConflict, "Conflict")
		return
	}

	if !auth.CanUserModifyRole(reqUser, role) {
		response.RespondError(c, http.StatusForbidden, "Forbidden")
		return
	}

	dbrole, err := database.FindRole(role)
	if err != nil {
		log.Errorf("Error finding role %s: %s", role, err)
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if err = database.AddRoleToUser(user, dbrole); err != nil {
		log.Errorf("Error adding role %s to %d: %s", role, user.CID, err)
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	_ = discord.NewMessage().
		SetContent(
			fmt.Sprintf("%s %s has added role %s to %s %s (%d)",
				reqUser.FirstName,
				reqUser.LastName,
				role,
				user.FirstName,
				user.LastName,
				user.CID)).
		Send("role")

	// If roles change, invalidate staff cache
	memcache.Cache.Delete("staff")

	response.RespondBlank(c, http.StatusNoContent)
}

// Remove User Role
// @Summary Remove User Role
// @Tags user
// @Param role path string true "Role"
// @Param cid path string true "CID"
// @Success 204
// @Failure 403 {object} response.R
// @Failure 404 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/user/:cid/roles/:role [DELETE]
func deleteUserRoles(c *gin.Context) {
	reqUser := c.MustGet("x-user").(*models.User)

	user, err := database.FindUserByCID(c.Param("cid"))
	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if user == nil {
		response.RespondError(c, http.StatusNotFound, "User not found")
		return
	}

	role := c.Param("role")

	if _, ok := auth.Roles[role]; !ok {
		response.RespondError(c, http.StatusNotFound, "Role not found")
		return
	}

	if !auth.HasRole(user, role) {
		response.RespondError(c, http.StatusNotFound, "Role not found")
		return
	}

	if !auth.CanUserModifyRole(reqUser, role) {
		response.RespondError(c, http.StatusForbidden, "Forbidden")
		return
	}

	dbrole, err := database.FindRole(role)
	if err != nil {
		log.Errorf("Error finding role %s: %s", role, err)
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if err = database.RemoveRoleFromUser(user, dbrole); err != nil {
		log.Errorf("Error removing role %s from %d: %s", role, user.CID, err)
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	_ = discord.NewMessage().
		SetContent(
			fmt.Sprintf("%s %s has removed role %s from %s %s (%d)",
				reqUser.FirstName,
				reqUser.LastName,
				role,
				user.FirstName,
				user.LastName,
				user.CID)).
		Send("role")

	// If roles change, invalidate staff cache
	memcache.Cache.Delete("staff")

	response.RespondBlank(c, http.StatusNoContent)
}
