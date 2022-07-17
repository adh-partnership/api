package user

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kzdv/api/pkg/auth"
	"github.com/kzdv/api/pkg/database"
	"github.com/kzdv/api/pkg/database/dto"
	"github.com/kzdv/api/pkg/gin/response"
	"github.com/kzdv/api/pkg/vatusa"
	dbTypes "github.com/kzdv/types/database"
)

// Get User Information
// @Summary Get User Information
// @Tags user
// @Success 200 {object} dto.UserResponse
// @Failure 403 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/user [GET]
// @Router /v1/user/:cid [GET]
func getUser(c *gin.Context) {
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

	response.Respond(c, http.StatusOK, dto.ConvUserToUserResponse(user))
}

// Patch User Information
// @Summary Patch User Information
// @Tags user
// @Accept json
// @Produce json
// @Param user body dto.UserResponse true "User"
// @Param cid path string true "CID"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} response.R
// @Failure 403 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/user/{cid} [PATCH]
// @Router /v1/user [PATCH]
func patchUser(c *gin.Context) {
	user := c.MustGet("x-user").(*dbTypes.User)
	status := 200

	// If cid is set to 0 or not defined, user is patching their own information
	// which we allow for DiscordID
	cid := c.Param("cid")
	if cid == "0" || cid == "" {
		cid = fmt.Sprint(user.CID)
	}

	var req dto.UserResponse
	if err := c.ShouldBind(&req); err != nil {
		response.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	oldUser, err := database.FindUserByCID(cid)
	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Don't allow patching of some fields, so always force these to empty values
	req.CID = 0
	req.FirstName = ""
	req.LastName = ""

	if req.OperatingInitials != "" && !auth.InGroup(user, "admin") {
		response.RespondError(c, http.StatusForbidden, "Forbidden")
		return
	}

	if req.ControllerType != "" && !auth.InGroup(user, "admin") {
		response.RespondError(c, http.StatusForbidden, "Forbidden")
		return
	}

	if req.Certiciations.Delivery != "" && !auth.InGroup(user, "admin") && !auth.InGroup(user, "training") {
		response.RespondError(c, http.StatusForbidden, "Forbidden")
		return
	}

	if req.Certiciations.Ground != "" && !auth.InGroup(user, "admin") && !auth.InGroup(user, "training") {
		response.RespondError(c, http.StatusForbidden, "Forbidden")
		return
	}

	if req.Certiciations.Local != "" && !auth.InGroup(user, "admin") && !auth.InGroup(user, "training") {
		response.RespondError(c, http.StatusForbidden, "Forbidden")
		return
	}

	if req.Certiciations.Approach != "" && !auth.InGroup(user, "admin") && !auth.InGroup(user, "training") {
		response.RespondError(c, http.StatusForbidden, "Forbidden")
		return
	}

	if req.Certiciations.Enroute != "" && !auth.InGroup(user, "admin") && !auth.InGroup(user, "training") {
		response.RespondError(c, http.StatusForbidden, "Forbidden")
		return
	}

	if req.Rating != "" && !auth.InGroup(user, "admin") && !auth.HasRoleList(user, []string{"ta", "ins"}) {
		response.RespondError(c, http.StatusForbidden, "Forbidden")
		return
	}

	if req.ControllerType != oldUser.ControllerType {
		if (oldUser.ControllerType == dbTypes.ControllerTypeOptions["home"] ||
			oldUser.ControllerType == dbTypes.ControllerTypeOptions["visitor"]) &&
			req.RemovalReason == "" {

			response.RespondError(c, http.StatusBadRequest, "Removal reason required")
			return
		}

		var status int
		var err error

		if oldUser.ControllerType == dbTypes.ControllerTypeOptions["home"] {
			status, err = vatusa.RemoveController(c.Param("cid"), user.CID, req.RemovalReason)
		} else if oldUser.ControllerType == dbTypes.ControllerTypeOptions["visitor"] {
			status, err = vatusa.RemoveVisitingController(c.Param("cid"), user.CID, req.RemovalReason)
		} else if req.ControllerType == dbTypes.ControllerTypeOptions["visitor"] {
			status, err = vatusa.AddVisitingController(c.Param("cid"))
		} else {
			log.Errorf("Unknown controller type: %s", oldUser.ControllerType)
		}

		if err != nil {
			log.Errorf("Error setting controller type %s for %s: %s", req.ControllerType, c.Param("cid"), err)
			response.RespondError(c, http.StatusInternalServerError, "error changing controller type on vatusa")
			return
		}

		if status > 299 {
			log.Errorf("Got invalid status code from VATUSA changing controller type %s for %s: %d", req.ControllerType, c.Param("cid"), status)
			response.RespondError(c, http.StatusInternalServerError, "error changing controller type on vatusa")
			return
		}
	}

	if req.DiscordID != "" {
		// User can patch their own DiscordID
		if oldUser.CID != user.CID && !auth.InGroup(user, "admin") {
			response.RespondError(c, http.StatusForbidden, "Forbidden")
			return
		}
	}

	errors := dto.PatchUserFromUserResponse(oldUser, req)
	if len(errors) > 0 {
		response.RespondError(c, http.StatusBadRequest, strings.Join(errors, ", "))
		return
	}

	if err := database.DB.Save(oldUser).Error; err != nil {
		response.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Respond(c, status, dto.ConvUserToUserResponse(oldUser))
}
