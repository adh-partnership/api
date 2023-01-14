package dto

import (
	"reflect"
	"testing"
	"time"

	"github.com/adh-partnership/api/pkg/database/models"
	"github.com/adh-partnership/api/pkg/database/models/constants"
)

func TestConvUserToUserResponse(t *testing.T) {
	tim, _ := time.Parse("2006-01-02 15:04:05", "2020-01-01 00:00:00")
	user := &models.User{
		CID:                   123,
		FirstName:             "John",
		LastName:              "Doe",
		OperatingInitials:     "JD",
		ControllerType:        constants.ControllerTypeHome,
		GndCertification:      constants.CertificationCertified,
		MajorGndCertification: constants.CertificationCertified,
		LclCertification:      constants.CertificationSolo,
		MajorLclCertification: constants.CertificationSolo,
		AppCertification:      constants.CertificationCertified,
		MajorAppCertification: constants.CertificationCertified,
		CtrCertification:      constants.CertificationCanTrain,
		Rating:                models.Rating{Short: "C1", Long: "Controller"},
		Status:                constants.ControllerStatusActive,
		Roles: []*models.Role{
			{Name: "admin"},
			{Name: "user"},
		},
		DiscordID:      "123456789",
		RosterJoinDate: &tim,
		CreatedAt:      tim,
		UpdatedAt:      tim,
	}

	expectedResponse := &UserResponse{
		CID:               123,
		FirstName:         "John",
		LastName:          "Doe",
		OperatingInitials: "JD",
		ControllerType:    constants.ControllerTypeHome,
		Certifications: UserResponseCertifications{
			Ground:        constants.CertificationCertified,
			MajorGround:   constants.CertificationCertified,
			Local:         constants.CertificationSolo,
			MajorLocal:    constants.CertificationSolo,
			Approach:      constants.CertificationCertified,
			MajorApproach: constants.CertificationCertified,
			Enroute:       constants.CertificationCanTrain,
		},
		Rating:         "C1",
		Roles:          []string{"admin", "user"},
		Status:         constants.ControllerStatusActive,
		DiscordID:      "123456789",
		RosterJoinDate: "2020-01-01T00:00:00Z",
		CreatedAt:      "2020-01-01T00:00:00Z",
		UpdatedAt:      "2020-01-01T00:00:00Z",
	}

	userResponse := ConvUserToUserResponse(user)
	if !reflect.DeepEqual(userResponse, expectedResponse) {
		t.Errorf("ConvUserToUserResponse() = %+v\nwant %+v", userResponse, expectedResponse)
	}
}

func TestPatchUserFromUserResponse(t *testing.T) {
	tests := []struct {
		Name           string
		BaseUser       models.User
		Patch          UserResponseAdmin
		ExpectedUser   models.User
		ExpectedErrors []string
	}{
		{
			Name: "Patch OI",
			BaseUser: models.User{
				OperatingInitials: "JD",
			},
			Patch: UserResponseAdmin{
				OperatingInitials: "FB",
			},
			ExpectedUser: models.User{
				OperatingInitials: "FB",
			},
			ExpectedErrors: []string{},
		},
		{
			Name: "Invalid OI",
			BaseUser: models.User{
				OperatingInitials: "JD",
			},
			Patch: UserResponseAdmin{
				OperatingInitials: "ABC",
			},
			ExpectedUser: models.User{
				OperatingInitials: "JD",
			},
			ExpectedErrors: []string{ErrInvalidOperatingInitials},
		},
		{
			Name: "Patch ControllerType",
			BaseUser: models.User{
				ControllerType: constants.ControllerTypeHome,
			},
			Patch: UserResponseAdmin{
				ControllerType: constants.ControllerTypeNone,
			},
			ExpectedUser: models.User{
				ControllerType: constants.ControllerTypeNone,
			},
			ExpectedErrors: []string{},
		},
		{
			Name: "Invalid ControllerType",
			BaseUser: models.User{
				ControllerType: constants.ControllerTypeHome,
			},
			Patch: UserResponseAdmin{
				ControllerType: "invalid",
			},
			ExpectedUser: models.User{
				ControllerType: constants.ControllerTypeHome,
			},
			ExpectedErrors: []string{ErrInvalidControllerType},
		},
		{
			Name: "Patch Certification",
			BaseUser: models.User{
				GndCertification:      constants.CertificationCertified,
				MajorGndCertification: constants.CertificationCertified,
				LclCertification:      constants.CertificationSolo,
				MajorLclCertification: constants.CertificationSolo,
				AppCertification:      constants.CertificationCertified,
				MajorAppCertification: constants.CertificationCertified,
				CtrCertification:      constants.CertificationCanTrain,
			},
			Patch: UserResponseAdmin{
				Certifications: UserResponseCertifications{
					Ground:        constants.CertificationNone,
					MajorGround:   constants.CertificationNone,
					Local:         constants.CertificationNone,
					MajorLocal:    constants.CertificationNone,
					Approach:      constants.CertificationNone,
					MajorApproach: constants.CertificationNone,
					Enroute:       constants.CertificationNone,
				},
			},
			ExpectedUser: models.User{
				GndCertification:      constants.CertificationNone,
				MajorGndCertification: constants.CertificationNone,
				LclCertification:      constants.CertificationNone,
				MajorLclCertification: constants.CertificationNone,
				AppCertification:      constants.CertificationNone,
				MajorAppCertification: constants.CertificationNone,
				CtrCertification:      constants.CertificationNone,
			},
			ExpectedErrors: []string{},
		},
		{
			Name: "Invalid Certification",
			BaseUser: models.User{
				GndCertification:      constants.CertificationCertified,
				MajorGndCertification: constants.CertificationCertified,
				LclCertification:      constants.CertificationSolo,
				MajorLclCertification: constants.CertificationSolo,
				AppCertification:      constants.CertificationCertified,
				MajorAppCertification: constants.CertificationCertified,
				CtrCertification:      constants.CertificationCanTrain,
			},
			Patch: UserResponseAdmin{
				Certifications: UserResponseCertifications{
					Ground:   "invalid",
					Local:    "invalid",
					Approach: "invalid",
					Enroute:  "invalid",
				},
			},
			ExpectedUser: models.User{
				GndCertification:      constants.CertificationCertified,
				MajorGndCertification: constants.CertificationCertified,
				LclCertification:      constants.CertificationSolo,
				MajorLclCertification: constants.CertificationSolo,
				AppCertification:      constants.CertificationCertified,
				MajorAppCertification: constants.CertificationCertified,
				CtrCertification:      constants.CertificationCanTrain,
			},
			ExpectedErrors: []string{ErrInvalidCertification, ErrInvalidCertification, ErrInvalidCertification, ErrInvalidCertification},
		},
		{
			Name: "Patch DiscordID",
			BaseUser: models.User{
				DiscordID: "123456789",
			},
			Patch: UserResponseAdmin{
				DiscordID: "987654321",
			},
			ExpectedUser: models.User{
				DiscordID: "987654321",
			},
			ExpectedErrors: []string{},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			errors := PatchUserFromUserResponse(&test.BaseUser, test.Patch)
			if len(errors) != len(test.ExpectedErrors) {
				t.Errorf("Test: %s, Expected %d errors, got %d", test.Name, len(test.ExpectedErrors), len(errors))
			}
			if !reflect.DeepEqual(test.BaseUser, test.ExpectedUser) {
				t.Errorf("Test: %s, Expected %+v, got %+v", test.Name, test.ExpectedUser, test.BaseUser)
			}
		})
	}
}
