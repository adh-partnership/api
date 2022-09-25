package facility

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	internalDTO "github.com/kzdv/api/internal/v1/dto"
	"github.com/kzdv/api/pkg/database"
	"github.com/kzdv/api/pkg/database/dto"
	dbTypes "github.com/kzdv/api/pkg/database/types"
	"github.com/kzdv/api/pkg/gin/response"
	"github.com/kzdv/api/pkg/memcache"
)

// Get Facility Roster
// @Summary Get Facility Roster
// @Tags user
// @Success 200 {object} []dto.UserResponse
// @Failure 403 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/facility/roster [GET]
func getRoster(c *gin.Context) {
	users := []dbTypes.User{}
	ret := []*dto.UserResponse{}

	if err := database.DB.Not(&dbTypes.User{Status: dbTypes.ControllerStatusOptions["none"]}).Find(&users).Error; err != nil {
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
// @Success 200 {object} []internalDTO.FacilityStaffResponse
// @Failure 403 {object} response.R
// @Failure 500 {object} response.R
// @Router /v1/facility/staff [GET]
func getStaff(c *gin.Context) {
	if memcache.Cache.Get("staff") != nil {
		response.Respond(c, http.StatusOK, memcache.Cache.Get("staff"))
		return
	}

	atm, err := database.FindUsersWithRole("atm")
	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	datm, err := database.FindUsersWithRole("datm")
	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	ta, err := database.FindUsersWithRole("ta")
	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	ec, err := database.FindUsersWithRole("ec")
	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	fe, err := database.FindUsersWithRole("fe")
	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	wm, err := database.FindUsersWithRole("wm")
	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	staff := internalDTO.FacilityStaffResponse{
		ATM:  []*dto.UserResponse{},
		DATM: []*dto.UserResponse{},
		TA:   []*dto.UserResponse{},
		EC:   []*dto.UserResponse{},
		FE:   []*dto.UserResponse{},
		WM:   []*dto.UserResponse{},
	}

	for _, user := range atm {
		staff.ATM = append(staff.ATM, dto.ConvUserToUserResponse(&user))
	}

	for _, user := range datm {
		staff.DATM = append(staff.DATM, dto.ConvUserToUserResponse(&user))
	}

	for _, user := range ta {
		staff.TA = append(staff.TA, dto.ConvUserToUserResponse(&user))
	}

	for _, user := range ec {
		staff.EC = append(staff.EC, dto.ConvUserToUserResponse(&user))
	}

	for _, user := range fe {
		staff.FE = append(staff.FE, dto.ConvUserToUserResponse(&user))
	}

	for _, user := range wm {
		staff.WM = append(staff.WM, dto.ConvUserToUserResponse(&user))
	}

	// Cache for an hour, we do delete this cache if roles change.. so this should be okay to do
	// without worrying about stale lists
	memcache.Cache.Set("staff", staff, 1*time.Hour)

	response.Respond(c, http.StatusOK, staff)
}
