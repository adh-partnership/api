package dto

import dbTypes "github.com/kzdv/types/database"

type StandardResponse struct {
	Message string      `json:"message" yaml:"message" xml:"message"`
	Data    interface{} `json:"data" yaml:"data" xml:"data"`
}

type SSOUserResponse struct {
	Message string        `json:"message" yaml:"message" xml:"message"`
	User    *dbTypes.User `json:"user" yaml:"user" xml:"user"`
}
