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

package staffing

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/adh-partnership/api/pkg/database/dto"
	"github.com/adh-partnership/api/pkg/database/models"
	"github.com/adh-partnership/api/pkg/discord"
	"github.com/adh-partnership/api/pkg/gin/response"
)

// Submit a staffing request. [Feature Gated]
// @Summary Submit a staffing request [Feature Gated]
// @Description Submit a staffing request [Feature Gated]
// @Tags Staffing
// @Param data body dto.StaffingRequest true "Request Data"
// @Success 204
// @Failure 400 {object} response.R "Invalid form submission"
// @Failure 401 {object} response.R "Not logged in"
// @Failure 404 {object} response.R "Not Found -- feature disabled"
// @Failure 500 {object} response.R
// @Router /v1/staffing/ [post]
func requestStaffing(c *gin.Context) {
	user := c.MustGet("x-user").(*models.User)

	var dto dto.StaffingRequest
	if err := c.ShouldBind(&dto); err != nil {
		log.Debugf("Error binding dto: %s", err)
		response.RespondError(c, http.StatusBadRequest, "Invalid request")
		return
	}

	err := discord.NewMessage().
		SetContent("New staffing request").
		AddEmbed(discord.NewEmbed().
			AddField(discord.NewField().SetName("Requester").SetValue(fmt.Sprintf("%s %s (%d)", user.FirstName, user.LastName, user.CID))).
			AddField(discord.NewField().SetName("Start Date").SetValue(dto.StartDate)).
			AddField(discord.NewField().SetName("End Date").SetValue(dto.EndDate)).
			AddField(discord.NewField().SetName("DepartureAirport").SetValue(dto.DepartureAirport)).
			AddField(discord.NewField().SetName("ArrivalAirport").SetValue(dto.ArrivalAirport)).
			AddField(discord.NewField().SetName("Pilots").SetValue(strconv.Itoa(dto.Pilots))).
			AddField(discord.NewField().SetName("Contact").SetValue(dto.ContactInfo)).
			AddField(discord.NewField().SetName("Organization").SetValue(dto.Organization)).
			AddField(discord.NewField().SetName("Banner").SetValue(dto.BannerUrl)).
			AddField(discord.NewField().SetName("Comments").SetValue(dto.Comments)),
		).Send("staffing_request")
	if err != nil {
		log.Errorf("Error sending staffing request message to Discord: %s", err.Error())
	}

	response.RespondBlank(c, http.StatusNoContent)
}
