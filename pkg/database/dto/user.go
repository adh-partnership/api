package dto

import (
	"github.com/kzdv/api/pkg/database"
	dbTypes "github.com/kzdv/api/pkg/database/types"
)

type UserResponse struct {
	CID               uint                       `json:"cid" yaml:"cid" xml:"cid"`
	FirstName         string                     `json:"first_name" yaml:"first_name" xml:"first_name"`
	LastName          string                     `json:"last_name" yaml:"last_name" xml:"last_name"`
	OperatingInitials string                     `json:"operating_initials" yaml:"operating_initials" xml:"operating_initials"`
	ControllerType    string                     `json:"controller_type" yaml:"controller_type" xml:"controller_type"`
	RemovalReason     string                     `json:"removal_reason" yaml:"removal_reason" xml:"removal_reason"`
	Certiciations     UserResponseCertifications `json:"certifications" yaml:"certifications" xml:"certifications"`
	Rating            string                     `json:"rating" yaml:"rating" xml:"rating"`
	Status            string                     `json:"status" yaml:"status" xml:"status"`
	Roles             []string                   `json:"roles" yaml:"roles" xml:"roles"`
	DiscordID         string                     `json:"discord_id" yaml:"discord_id" xml:"discord_id"`
	CreatedAt         string                     `json:"created_at" yaml:"created_at" xml:"created_at"`
	UpdatedAt         string                     `json:"updated_at" yaml:"updated_at" xml:"updated_at"`
}

type UserResponseCertifications struct {
	Delivery string `json:"delivery" yaml:"delivery" xml:"delivery"`
	Ground   string `json:"ground" yaml:"ground" xml:"ground"`
	Local    string `json:"local" yaml:"local" xml:"local"`
	Approach string `json:"approach" yaml:"approach" xml:"approach"`
	Enroute  string `json:"enroute" yaml:"enroute" xml:"enroute"`
}

type FacilityStaffResponse struct {
	ATM  []*UserResponse `json:"atm" yaml:"atm" xml:"atm"`
	DATM []*UserResponse `json:"datm" yaml:"datm" xml:"datm"`
	TA   []*UserResponse `json:"ta" yaml:"ta" xml:"ta"`
	EC   []*UserResponse `json:"ec" yaml:"ec" xml:"ec"`
	FE   []*UserResponse `json:"fe" yaml:"fe" xml:"fe"`
	WM   []*UserResponse `json:"wm" yaml:"wm" xml:"wm"`
}

func ConvUserToUserResponse(user *dbTypes.User) *UserResponse {
	roles := []string{}
	if user.Roles != nil {
		for _, role := range user.Roles {
			roles = append(roles, role.Name)
		}
	}

	return &UserResponse{
		CID:               user.CID,
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		OperatingInitials: user.OperatingInitials,
		ControllerType:    user.ControllerType,
		Certiciations: UserResponseCertifications{
			Delivery: user.DelCertification,
			Ground:   user.GndCertification,
			Local:    user.LclCertification,
			Approach: user.AppCertification,
			Enroute:  user.CtrCertification,
		},
		Roles:     roles,
		Rating:    user.Rating.Short,
		Status:    user.Status,
		DiscordID: user.DiscordID,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

const (
	ErrInvalidOperatingInitials = "invalid operating initials"
	ErrInvalidControllerType    = "invalid controller type"
	ErrInvalidCertification     = "invalid certification"
	ErrInvalidRating            = "invalid rating"
	ErrInvalidStatus            = "invalid status"
)

func PatchUserFromUserResponse(user *dbTypes.User, userResponse UserResponse) []string {
	var errs []string

	if len(userResponse.OperatingInitials) != 2 && userResponse.OperatingInitials != "" {
		errs = append(errs, ErrInvalidOperatingInitials)
	} else {
		user.OperatingInitials = userResponse.OperatingInitials
	}

	if userResponse.ControllerType != "" {
		if _, ok := dbTypes.ControllerTypeOptions[userResponse.ControllerType]; !ok {
			errs = append(errs, ErrInvalidControllerType)
		} else {
			user.ControllerType = userResponse.ControllerType
		}
	}

	if userResponse.Certiciations.Delivery != "" {
		if _, ok := dbTypes.CertificationOptions[userResponse.Certiciations.Delivery]; !ok {
			errs = append(errs, ErrInvalidCertification)
		} else {
			user.DelCertification = userResponse.Certiciations.Delivery
		}
	}

	if userResponse.Certiciations.Ground != "" {
		if _, ok := dbTypes.CertificationOptions[userResponse.Certiciations.Ground]; !ok {
			errs = append(errs, ErrInvalidCertification)
		} else {
			user.GndCertification = userResponse.Certiciations.Ground
		}
	}

	if userResponse.Certiciations.Local != "" {
		if _, ok := dbTypes.CertificationOptions[userResponse.Certiciations.Local]; !ok {
			errs = append(errs, ErrInvalidCertification)
		} else {
			user.LclCertification = userResponse.Certiciations.Local
		}
	}

	if userResponse.Certiciations.Approach != "" {
		if _, ok := dbTypes.CertificationOptions[userResponse.Certiciations.Approach]; !ok {
			errs = append(errs, ErrInvalidCertification)
		} else {
			user.AppCertification = userResponse.Certiciations.Approach
		}
	}

	if userResponse.Certiciations.Enroute != "" {
		if _, ok := dbTypes.CertificationOptions[userResponse.Certiciations.Enroute]; !ok {
			errs = append(errs, ErrInvalidCertification)
		} else {
			user.CtrCertification = userResponse.Certiciations.Enroute
		}
	}

	if userResponse.DiscordID != "" {
		user.DiscordID = userResponse.DiscordID
	}

	if userResponse.Status != "" {
		if _, ok := dbTypes.ControllerStatusOptions[userResponse.Status]; !ok {
			errs = append(errs, ErrInvalidStatus)
		} else {
			user.Status = userResponse.Status
		}
	}

	return errs
}

func GetUsersByRole(role string) ([]*UserResponse, error) {
	var users []*UserResponse

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
	roles := []string{"atm", "datm", "ta", "ec", "fe", "wm"}
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
		}
	}

	return staff, nil
}
