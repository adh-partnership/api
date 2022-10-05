package dto

import (
	"reflect"
	"testing"
	"time"

	dbTypes "github.com/adh-partnership/api/pkg/database/types"
)

func TestConvUserToUserResponse(t *testing.T) {
	tim, _ := time.Parse("2006-01-02 15:04:05", "2020-01-01 00:00:00")
	user := &dbTypes.User{
		CID:               123,
		FirstName:         "John",
		LastName:          "Doe",
		OperatingInitials: "JD",
		ControllerType:    dbTypes.ControllerTypeOptions["home"],
		DelCertification:  dbTypes.CertificationOptions["training"],
		GndCertification:  dbTypes.CertificationOptions["major"],
		LclCertification:  dbTypes.CertificationOptions["solo"],
		AppCertification:  dbTypes.CertificationOptions["certified"],
		CtrCertification:  dbTypes.CertificationOptions["cantrain"],
		Rating:            dbTypes.Rating{Short: "C1", Long: "Controller"},
		Status:            dbTypes.ControllerStatusOptions["active"],
		Roles: []*dbTypes.Role{
			{Name: "admin"},
			{Name: "user"},
		},
		DiscordID: "123456789",
		CreatedAt: tim,
		UpdatedAt: tim,
	}

	expectedResponse := &UserResponse{
		CID:               123,
		FirstName:         "John",
		LastName:          "Doe",
		OperatingInitials: "JD",
		ControllerType:    dbTypes.ControllerTypeOptions["home"],
		Certiciations: UserResponseCertifications{
			Delivery: dbTypes.CertificationOptions["training"],
			Ground:   dbTypes.CertificationOptions["major"],
			Local:    dbTypes.CertificationOptions["solo"],
			Approach: dbTypes.CertificationOptions["certified"],
			Enroute:  dbTypes.CertificationOptions["cantrain"],
		},
		Rating:    "C1",
		Roles:     []string{"admin", "user"},
		Status:    dbTypes.ControllerStatusOptions["active"],
		DiscordID: "123456789",
		CreatedAt: "2020-01-01T00:00:00Z",
		UpdatedAt: "2020-01-01T00:00:00Z",
	}

	userResponse := ConvUserToUserResponse(user)
	if !reflect.DeepEqual(userResponse, expectedResponse) {
		t.Errorf("ConvUserToUserResponse() = %+v\nwant %+v", userResponse, expectedResponse)
	}
}

func TestPatchUserFromUserResponse(t *testing.T) {
	tests := []struct {
		Name           string
		BaseUser       dbTypes.User
		Patch          UserResponseAdmin
		ExpectedUser   dbTypes.User
		ExpectedErrors []string
	}{
		{
			Name: "Patch OI",
			BaseUser: dbTypes.User{
				OperatingInitials: "JD",
			},
			Patch: UserResponseAdmin{
				OperatingInitials: "FB",
			},
			ExpectedUser: dbTypes.User{
				OperatingInitials: "FB",
			},
			ExpectedErrors: []string{},
		},
		{
			Name: "Invalid OI",
			BaseUser: dbTypes.User{
				OperatingInitials: "JD",
			},
			Patch: UserResponseAdmin{
				OperatingInitials: "ABC",
			},
			ExpectedUser: dbTypes.User{
				OperatingInitials: "JD",
			},
			ExpectedErrors: []string{ErrInvalidOperatingInitials},
		},
		{
			Name: "Patch ControllerType",
			BaseUser: dbTypes.User{
				ControllerType: dbTypes.ControllerTypeOptions["home"],
			},
			Patch: UserResponseAdmin{
				ControllerType: dbTypes.ControllerTypeOptions["none"],
			},
			ExpectedUser: dbTypes.User{
				ControllerType: dbTypes.ControllerTypeOptions["none"],
			},
			ExpectedErrors: []string{},
		},
		{
			Name: "Invalid ControllerType",
			BaseUser: dbTypes.User{
				ControllerType: dbTypes.ControllerTypeOptions["home"],
			},
			Patch: UserResponseAdmin{
				ControllerType: "invalid",
			},
			ExpectedUser: dbTypes.User{
				ControllerType: dbTypes.ControllerTypeOptions["home"],
			},
			ExpectedErrors: []string{ErrInvalidControllerType},
		},
		{
			Name: "Patch Certification",
			BaseUser: dbTypes.User{
				DelCertification: dbTypes.CertificationOptions["training"],
				GndCertification: dbTypes.CertificationOptions["major"],
				LclCertification: dbTypes.CertificationOptions["solo"],
				AppCertification: dbTypes.CertificationOptions["certified"],
				CtrCertification: dbTypes.CertificationOptions["cantrain"],
			},
			Patch: UserResponseAdmin{
				Certiciations: UserResponseCertifications{
					Delivery: dbTypes.CertificationOptions["none"],
					Ground:   dbTypes.CertificationOptions["none"],
					Local:    dbTypes.CertificationOptions["none"],
					Approach: dbTypes.CertificationOptions["none"],
					Enroute:  dbTypes.CertificationOptions["none"],
				},
			},
			ExpectedUser: dbTypes.User{
				DelCertification: dbTypes.CertificationOptions["none"],
				GndCertification: dbTypes.CertificationOptions["none"],
				LclCertification: dbTypes.CertificationOptions["none"],
				AppCertification: dbTypes.CertificationOptions["none"],
				CtrCertification: dbTypes.CertificationOptions["none"],
			},
			ExpectedErrors: []string{},
		},
		{
			Name: "Invalid Certification",
			BaseUser: dbTypes.User{
				DelCertification: dbTypes.CertificationOptions["training"],
				GndCertification: dbTypes.CertificationOptions["major"],
				LclCertification: dbTypes.CertificationOptions["solo"],
				AppCertification: dbTypes.CertificationOptions["certified"],
				CtrCertification: dbTypes.CertificationOptions["cantrain"],
			},
			Patch: UserResponseAdmin{
				Certiciations: UserResponseCertifications{
					Delivery: "invalid",
					Ground:   "invalid",
					Local:    "invalid",
					Approach: "invalid",
					Enroute:  "invalid",
				},
			},
			ExpectedUser: dbTypes.User{
				DelCertification: dbTypes.CertificationOptions["training"],
				GndCertification: dbTypes.CertificationOptions["major"],
				LclCertification: dbTypes.CertificationOptions["solo"],
				AppCertification: dbTypes.CertificationOptions["certified"],
				CtrCertification: dbTypes.CertificationOptions["cantrain"],
			},
			ExpectedErrors: []string{ErrInvalidCertification, ErrInvalidCertification, ErrInvalidCertification, ErrInvalidCertification, ErrInvalidCertification},
		},
		{
			Name: "Patch DiscordID",
			BaseUser: dbTypes.User{
				DiscordID: "123456789",
			},
			Patch: UserResponseAdmin{
				DiscordID: "987654321",
			},
			ExpectedUser: dbTypes.User{
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
