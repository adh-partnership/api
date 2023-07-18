package staffing

import (
	"fmt"
	"net/http"

	"github.com/adh-partnership/api/pkg/database/dto"
	"github.com/adh-partnership/api/pkg/database/models"
	"github.com/adh-partnership/api/pkg/discord"
	"github.com/adh-partnership/api/pkg/gin/middleware/auth"
	"github.com/adh-partnership/api/pkg/gin/response"
	"github.com/gin-gonic/gin"
	"istio.io/pkg/log"
)

func Routes(r *gin.RouterGroup) {
	r.POST("", auth.NotGuest, requestStaffing)
}

// Submit a staffing request.
// @Summary Submit a staffing request
// @Description Submit a staffing request
// @Tags Staffing
// @Param data body dto.StaffingRequest true "Request Data"
// @Success 202
// @Failure 400 {object} response.R "Invalid form submission"
// @Failure 401 {object} response.R "Not logged in"
// @Failure 500 {object} response.R
// @Router /v1/staffing/ [post]
func requestStaffing(c *gin.Context) {
	user := c.MustGet("x-user").(*models.User)

	var dto dto.StaffingRequest
	if err := c.ShouldBindJSON(&dto); err != nil {
		log.Debugf("Error binding dto: %s", err)
		response.RespondError(c, http.StatusBadRequest, "Invalid request")
		return
	}

	_ = discord.NewMessage().
		SetContent("New staffing request").
		AddEmbed(discord.NewEmbed().
			AddField(discord.NewField().SetName("Requester").SetValue(fmt.Sprintf("%s %s", user.FirstName, user.LastName))).
			AddField(discord.NewField().SetName("CID").SetValue(fmt.Sprintf("%d", user.CID))).
			AddField(discord.NewField().SetName("Date").SetValue(dto.Date)).
			AddField(discord.NewField().SetName("Start").SetValue(dto.Start)).
			AddField(discord.NewField().SetName("End").SetValue(dto.End)).
			AddField(discord.NewField().SetName("DepartureAirport").SetValue(dto.DepartureAirport)).
			AddField(discord.NewField().SetName("ArrivalAirport").SetValue(dto.ArrivalAirport)).
			AddField(discord.NewField().SetName("Pilots").SetValue(dto.Pilots)).
			AddField(discord.NewField().SetName("Comments").SetValue(dto.Comments)),
		).Send("staffing_request")

	response.RespondBlank(c, http.StatusAccepted)
}
