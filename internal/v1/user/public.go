package user

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"

	"github.com/kzdv/api/pkg/database"
	"github.com/kzdv/api/pkg/database/dto"
	dbTypes "github.com/kzdv/api/pkg/database/types"
	"github.com/kzdv/api/pkg/gin/response"
	"github.com/kzdv/api/pkg/memcache"
)

// Get Full Roster
// @Summary Get Full Roster
// @Tags user
// @Success 200 {object} []dto.UserResponse
// @Failure 403 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/user/all [GET]
func getFullRoster(c *gin.Context) {
	users := []dbTypes.User{}
	ret := []*dto.UserResponse{}

	if err := database.DB.Preload(clause.Associations).Find(&users).Error; err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	for _, user := range users {
		ret = append(ret, dto.ConvUserToUserResponse(&user))
	}

	response.Respond(c, http.StatusOK, ret)
}

// Get Facility Roster
// @Summary Get Facility Roster
// @Tags user
// @Success 200 {object} []dto.UserResponse
// @Failure 403 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/user/roster [GET]
func getRoster(c *gin.Context) {
	users := []dbTypes.User{}
	ret := []*dto.UserResponse{}

	if err := database.DB.Preload(clause.Associations).Not(&dbTypes.User{Status: dbTypes.ControllerStatusOptions["none"]}).Find(&users).Error; err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	for _, user := range users {
		ret = append(ret, dto.ConvUserToUserResponse(&user))
	}

	response.Respond(c, http.StatusOK, ret)
}

// Get Facility Staff
// @Summary Get Facility Staff
// @Tags user
// @Success 200 {object} []dto.FacilityStaffResponse
// @Failure 403 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/user/staff [GET]
func getStaff(c *gin.Context) {
	if memcache.Cache.Get("staff") != nil {
		response.Respond(c, http.StatusOK, memcache.Cache.Get("staff"))
		return
	}

	staff, err := dto.GetStaffResponse()
	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// Cache for an hour, we do delete this cache if roles change.. so this should be okay to do
	// without worrying about stale lists
	memcache.Cache.Set("staff", staff, 1*time.Hour)

	response.Respond(c, http.StatusOK, staff)
}
