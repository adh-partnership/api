package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/adh-partnership/api/pkg/auth"
	"github.com/adh-partnership/api/pkg/database"
	dbTypes "github.com/adh-partnership/api/pkg/database/types"
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
	user := c.MustGet("x-user").(*dbTypes.User)

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
	reqUser := c.MustGet("x-user").(*dbTypes.User)

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
	reqUser := c.MustGet("x-user").(*dbTypes.User)

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

	// If roles change, invalidate staff cache
	memcache.Cache.Delete("staff")

	response.RespondBlank(c, http.StatusNoContent)
}
