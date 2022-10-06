package dto

import "github.com/adh-partnership/api/pkg/database/models"

type StandardResponse struct {
	Message string      `json:"message" yaml:"message" xml:"message"`
	Data    interface{} `json:"data" yaml:"data" xml:"data"`
}

type SSOUserResponse struct {
	Message string       `json:"message" yaml:"message" xml:"message"`
	User    *models.User `json:"user" yaml:"user" xml:"user"`
}
