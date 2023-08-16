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
	"github.com/adh-partnership/api/pkg/database"
	"github.com/adh-partnership/api/pkg/database/models"
)

type UserResponse struct {
	CID                  uint                       `json:"cid" yaml:"cid" xml:"cid"`
	FirstName            string                     `json:"first_name" yaml:"first_name" xml:"first_name"`
	LastName             string                     `json:"last_name" yaml:"last_name" xml:"last_name"`
	OperatingInitials    string                     `json:"operating_initials" yaml:"operating_initials" xml:"operating_initials"`
	ControllerType       string                     `json:"controller_type" yaml:"controller_type" xml:"controller_type"`
	Certifications       UserResponseCertifications `json:"certifications" yaml:"certifications" xml:"certifications"`
	Rating               string                     `json:"rating" yaml:"rating" xml:"rating"`
	Status               string                     `json:"status" yaml:"status" xml:"status"`
	Roles                []string                   `json:"roles" yaml:"roles" xml:"roles"`
	Region               string                     `json:"region" yaml:"region" xml:"region"`
	Division             string                     `json:"division" yaml:"division" xml:"division"`
	Subdivision          string                     `json:"subdivision" yaml:"subdivision" xml:"subdivision"`
	DiscordID            string                     `json:"discord_id" yaml:"discord_id" xml:"discord_id"`
	RosterJoinDate       string                     `json:"roster_join_date" yaml:"roster_join_date" xml:"roster_join_date"`
	ExemptedFromActivity *bool                      `json:"exempted_from_activity" yaml:"exempted_from_activity" xml:"exempted_from_activity"`
	CreatedAt            string                     `json:"created_at" yaml:"created_at" xml:"created_at"`
	UpdatedAt            string                     `json:"updated_at" yaml:"updated_at" xml:"updated_at"`
}

type UserResponseAdmin struct {
	CID                  uint                       `json:"cid" yaml:"cid" xml:"cid"`
	FirstName            string                     `json:"first_name" yaml:"first_name" xml:"first_name"`
	LastName             string                     `json:"last_name" yaml:"last_name" xml:"last_name"`
	OperatingInitials    string                     `json:"operating_initials" yaml:"operating_initials" xml:"operating_initials"`
	ControllerType       string                     `json:"controller_type" yaml:"controller_type" xml:"controller_type"`
	Certifications       UserResponseCertifications `json:"certifications" yaml:"certifications" xml:"certifications"`
	RemovalReason        string                     `json:"removal_reason" yaml:"removal_reason" xml:"removal_reason"`
	Rating               string                     `json:"rating" yaml:"rating" xml:"rating"`
	Status               string                     `json:"status" yaml:"status" xml:"status"`
	Roles                []string                   `json:"roles" yaml:"roles" xml:"roles"`
	Region               string                     `json:"region" yaml:"region" xml:"region"`
	Division             string                     `json:"division" yaml:"division" xml:"division"`
	Subdivision          string                     `json:"subdivision" yaml:"subdivision" xml:"subdivision"`
	DiscordID            string                     `json:"discord_id" yaml:"discord_id" xml:"discord_id"`
	ExemptedFromActivity *bool                      `json:"exempted_from_activity" yaml:"exempted_from_activity" xml:"exempted_from_activity"`
	RosterJoinDate       string                     `json:"roster_join_date" yaml:"roster_join_date" xml:"roster_join_date"`
	CreatedAt            string                     `json:"created_at" yaml:"created_at" xml:"created_at"`
	UpdatedAt            string                     `json:"updated_at" yaml:"updated_at" xml:"updated_at"`
}

type VisitorResponse struct {
	ID        uint          `json:"id" yaml:"id" xml:"id"`
	User      *UserResponse `json:"user" yaml:"user" xml:"user"`
	CreatedAt string        `json:"created_at" yaml:"created_at" xml:"created_at"`
	UpdatedAt string        `json:"updated_at" yaml:"updated_at" xml:"updated_at"`
}

type UserResponseCertifications struct {
	Ground        string `json:"ground" yaml:"ground" xml:"ground"`
	MajorGround   string `json:"major_ground" yaml:"major_ground" xml:"major_ground"`
	Local         string `json:"local" yaml:"local" xml:"local"`
	MajorLocal    string `json:"major_local" yaml:"major_local" xml:"major_local"`
	Approach      string `json:"approach" yaml:"approach" xml:"approach"`
	MajorApproach string `json:"major_approach" yaml:"major_approach" xml:"major_approach"`
	Enroute       string `json:"enroute" yaml:"enroute" xml:"enroute"`
	Oceanic       string `json:"oceanic" yaml:"oceanic" xml:"oceanic"`
}

type FacilityStaffResponse struct {
	ATM        []*UserResponse `json:"atm" yaml:"atm" xml:"atm"`
	DATM       []*UserResponse `json:"datm" yaml:"datm" xml:"datm"`
	TA         []*UserResponse `json:"ta" yaml:"ta" xml:"ta"`
	EC         []*UserResponse `json:"ec" yaml:"ec" xml:"ec"`
	FE         []*UserResponse `json:"fe" yaml:"fe" xml:"fe"`
	WM         []*UserResponse `json:"wm" yaml:"wm" xml:"wm"`
	Events     []*UserResponse `json:"events" yaml:"events" xml:"events"`
	Facilities []*UserResponse `json:"facilities" yaml:"facilities" xml:"facilities"`
	Web        []*UserResponse `json:"web" yaml:"web" xml:"web"`
	Instructor []*UserResponse `json:"instructor" yaml:"instructor" xml:"instructor"`
	Mentor     []*UserResponse `json:"mentor" yaml:"mentor" xml:"mentor"`
}

func ConvUserToUserResponse(user *models.User) *UserResponse {
	if user == nil {
		return nil
	}

	roles := []string{}
	if user.Roles != nil {
		for _, role := range user.Roles {
			roles = append(roles, role.Name)
		}
	}

	u := &UserResponse{
		CID:               user.CID,
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		OperatingInitials: user.OperatingInitials,
		ControllerType:    user.ControllerType,
		Certifications: UserResponseCertifications{
			Ground:        user.GndCertification,
			MajorGround:   user.MajorGndCertification,
			Local:         user.LclCertification,
			MajorLocal:    user.MajorLclCertification,
			Approach:      user.AppCertification,
			MajorApproach: user.MajorAppCertification,
			Enroute:       user.CtrCertification,
			Oceanic:       user.OceanicCertification,
		},
		Roles:                roles,
		Rating:               user.Rating.Short,
		Status:               user.Status,
		DiscordID:            user.DiscordID,
		Region:               user.Region,
		Division:             user.Division,
		Subdivision:          user.Subdivision,
		ExemptedFromActivity: &user.ExemptedFromActivity,
		CreatedAt:            user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:            user.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	if user.RosterJoinDate != nil {
		u.RosterJoinDate = user.RosterJoinDate.Format("2006-01-02T15:04:05Z")
	}

	return u
}

func ConvVisitorApplicationsToResponse(applications []models.VisitorApplication) []*VisitorResponse {
	response := []*VisitorResponse{}
	for _, application := range applications {
		response = append(response, &VisitorResponse{
			ID:        application.ID,
			User:      ConvUserToUserResponse(application.User),
			CreatedAt: application.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt: application.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	return response
}

const (
	ErrInvalidOperatingInitials = "invalid operating initials"
	ErrInvalidControllerType    = "invalid controller type"
	ErrInvalidCertification     = "invalid certification"
	ErrInvalidRating            = "invalid rating"
	ErrInvalidStatus            = "invalid status"
)

func PatchUserFromUserResponse(user *models.User, userResponse UserResponseAdmin) []string {
	var errs []string

	if len(userResponse.OperatingInitials) != 2 && userResponse.OperatingInitials != "" {
		errs = append(errs, ErrInvalidOperatingInitials)
	} else if userResponse.OperatingInitials != "" {
		user.OperatingInitials = userResponse.OperatingInitials
	}

	if userResponse.ControllerType != "" {
		if _, ok := models.ControllerTypeOptions[userResponse.ControllerType]; !ok {
			errs = append(errs, ErrInvalidControllerType)
		} else {
			user.ControllerType = userResponse.ControllerType
		}
	}

	if userResponse.Certifications.Ground != "" {
		if _, ok := models.CertificationOptions[userResponse.Certifications.Ground]; !ok {
			errs = append(errs, ErrInvalidCertification)
		} else {
			user.GndCertification = userResponse.Certifications.Ground
		}
	}

	if userResponse.Certifications.MajorGround != "" {
		if _, ok := models.CertificationOptions[userResponse.Certifications.MajorGround]; !ok {
			errs = append(errs, ErrInvalidCertification)
		} else {
			user.MajorGndCertification = userResponse.Certifications.MajorGround
		}
	}

	if userResponse.Certifications.Local != "" {
		if _, ok := models.CertificationOptions[userResponse.Certifications.Local]; !ok {
			errs = append(errs, ErrInvalidCertification)
		} else {
			user.LclCertification = userResponse.Certifications.Local
		}
	}

	if userResponse.Certifications.MajorLocal != "" {
		if _, ok := models.CertificationOptions[userResponse.Certifications.MajorLocal]; !ok {
			errs = append(errs, ErrInvalidCertification)
		} else {
			user.MajorLclCertification = userResponse.Certifications.MajorLocal
		}
	}

	if userResponse.Certifications.Approach != "" {
		if _, ok := models.CertificationOptions[userResponse.Certifications.Approach]; !ok {
			errs = append(errs, ErrInvalidCertification)
		} else {
			user.AppCertification = userResponse.Certifications.Approach
		}
	}

	if userResponse.Certifications.MajorApproach != "" {
		if _, ok := models.CertificationOptions[userResponse.Certifications.MajorApproach]; !ok {
			errs = append(errs, ErrInvalidCertification)
		} else {
			user.MajorAppCertification = userResponse.Certifications.MajorApproach
		}
	}

	if userResponse.Certifications.Enroute != "" {
		if _, ok := models.CertificationOptions[userResponse.Certifications.Enroute]; !ok {
			errs = append(errs, ErrInvalidCertification)
		} else {
			user.CtrCertification = userResponse.Certifications.Enroute
		}
	}

	if userResponse.Certifications.Oceanic != "" {
		if _, ok := models.CertificationOptions[userResponse.Certifications.Oceanic]; !ok {
			errs = append(errs, ErrInvalidCertification)
		} else {
			user.OceanicCertification = userResponse.Certifications.Oceanic
		}
	}

	if userResponse.DiscordID != "" {
		// @TODO: This is a dirty hack.. revisit this at another point.
		if userResponse.DiscordID == "-1" {
			user.DiscordID = ""
		} else {
			user.DiscordID = userResponse.DiscordID
		}
	}

	if userResponse.ExemptedFromActivity != nil {
		user.ExemptedFromActivity = *userResponse.ExemptedFromActivity
	}

	if userResponse.Status != "" {
		if _, ok := models.ControllerStatusOptions[userResponse.Status]; !ok {
			errs = append(errs, ErrInvalidStatus)
		} else {
			user.Status = userResponse.Status
		}
	}

	return errs
}

func GetUsersByRole(role string) ([]*UserResponse, error) {
	users := []*UserResponse{}

	u, err := database.FindUsersWithRole(role)
	if err != nil {
		return nil, err
	}

	for _, user := range u {
		users = append(users, ConvUserToUserResponse(&user))
	}

	return users, nil
}

func GetStaffResponse() (*FacilityStaffResponse, error) {
	roles := []string{"atm", "datm", "ta", "ec", "fe", "wm", "events", "facilities", "web", "ins", "mtr"}
	staff := &FacilityStaffResponse{}

	for _, role := range roles {
		u, err := GetUsersByRole(role)
		if err != nil {
			return nil, err
		}

		switch role {
		case "atm":
			staff.ATM = u
		case "datm":
			staff.DATM = u
		case "ta":
			staff.TA = u
		case "ec":
			staff.EC = u
		case "fe":
			staff.FE = u
		case "wm":
			staff.WM = u
		case "events":
			staff.Events = u
		case "facilities":
			staff.Facilities = u
		case "web":
			staff.Web = u
		case "ins":
			staff.Instructor = u
		case "mtr":
			staff.Mentor = u
		}
	}

	return staff, nil
}
