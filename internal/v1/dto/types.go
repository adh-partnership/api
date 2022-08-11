package dto

import dbTypes "github.com/kzdv/api/pkg/database/types"

type StandardResponse struct {
	Message string      `json:"message" yaml:"message" xml:"message"`
	Data    interface{} `json:"data" yaml:"data" xml:"data"`
}

type SSOUserResponse struct {
	Message string        `json:"message" yaml:"message" xml:"message"`
	User    *dbTypes.User `json:"user" yaml:"user" xml:"user"`
}
