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

package dto

import (
	"time"

	"github.com/adh-partnership/api/pkg/database"
	"github.com/adh-partnership/api/pkg/database/models"
)

type FacilityReportDTO struct {
	Position  string
	Duration  int
	LogonTime *time.Time `json:"logon_time"`
}

type ControllerStats struct {
	CID               uint    `json:"cid" example:"1"`
	FirstName         string  `json:"first_name" example:"Daniel"`
	LastName          string  `json:"last_name" example:"Hawton"`
	OperatingInitials string  `json:"operating_initials" example:"DAN"`
	ControllerType    string  `json:"controllerType" example:"home"`
	Rating            string  `json:"rating" example:"S1"`
	Cab               float32 `json:"cab" example:"0.5"`
	Terminal          float32 `json:"terminal" example:"0.5"`
	Enroute           float32 `json:"enroute" example:"0.5"`
}

type OnlineController struct {
	CID         uint          `json:"cid" example:"1"`
	Controller  *UserResponse `json:"controller"`
	Position    string        `json:"position" example:"ANC_00_CTR"`
	Frequency   string        `json:"frequency" example:"118.000"`
	OnlineSince string        `json:"online_since" example:"2020-01-01T00:00:00Z"`
}

func ConvertControllerStatsToFacilityReport(stats []*models.ControllerStat) []*FacilityReportDTO {
	var facilityReport []*FacilityReportDTO
	for _, stat := range stats {
		facilityReport = append(facilityReport, &FacilityReportDTO{
			Position:  stat.Position,
			Duration:  stat.Duration,
			LogonTime: &stat.LogonTime,
		})
	}
	return facilityReport
}

func GetDTOForUserAndMonth(user *models.User, month int, year int) (*ControllerStats, error) {
	cab, terminal, enroute, err := database.GetStatsForUserAndMonth(user, month, year)
	if err != nil {
		return nil, err
	}
	return &ControllerStats{
		CID:               user.CID,
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		OperatingInitials: user.OperatingInitials,
		ControllerType:    user.ControllerType,
		Rating:            user.Rating.Short,
		Cab:               cab,
		Terminal:          terminal,
		Enroute:           enroute,
	}, nil
}

func ConvertOnlineToDTO(online *models.OnlineController) *OnlineController {
	return &OnlineController{
		CID:         online.User.CID,
		Controller:  ConvUserToUserResponse(online.User),
		Position:    online.Position,
		Frequency:   online.Frequency,
		OnlineSince: online.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func ConvertOnlineToDTOs(online []models.OnlineController) []OnlineController {
	ret := []OnlineController{}

	for _, o := range online {
		ret = append(ret, *ConvertOnlineToDTO(&o))
	}
	return ret
}
