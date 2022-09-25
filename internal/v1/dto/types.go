package dto

import (
	"github.com/kzdv/api/pkg/database/dto"
	dbTypes "github.com/kzdv/api/pkg/database/types"
)

type StandardResponse struct {
	Message string      `json:"message" yaml:"message" xml:"message"`
	Data    interface{} `json:"data" yaml:"data" xml:"data"`
}

type SSOUserResponse struct {
	Message string        `json:"message" yaml:"message" xml:"message"`
	User    *dbTypes.User `json:"user" yaml:"user" xml:"user"`
}

type FacilityStaffResponse struct {
	ATM  []*dto.UserResponse `json:"atm" yaml:"atm" xml:"atm"`
	DATM []*dto.UserResponse `json:"datm" yaml:"datm" xml:"datm"`
	TA   []*dto.UserResponse `json:"ta" yaml:"ta" xml:"ta"`
	EC   []*dto.UserResponse `json:"ec" yaml:"ec" xml:"ec"`
	FE   []*dto.UserResponse `json:"fe" yaml:"fe" xml:"fe"`
	WM   []*dto.UserResponse `json:"wm" yaml:"wm" xml:"wm"`
}
