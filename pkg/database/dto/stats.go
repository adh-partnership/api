package dto

import (
	"github.com/adh-partnership/api/pkg/database"
	"github.com/adh-partnership/api/pkg/database/models"
)

type ControllerStats struct {
	CID            uint    `json:"cid" example:"1"`
	FirstName      string  `json:"firstname" example:"Daniel"`
	LastName       string  `json:"lastname" example:"Hawton"`
	ControllerType string  `json:"controllerType" example:"home"`
	Rating         string  `json:"rating" example:"S1"`
	Cab            float32 `json:"cab" example:"0.5"`
	Terminal       float32 `json:"terminal" example:"0.5"`
	Enroute        float32 `json:"enroute" example:"0.5"`
}

func GetDTOForUserAndMonth(user *models.User, month int, year int) (*ControllerStats, error) {
	cab, terminal, enroute, err := database.GetStatsForUserAndMonth(user, month, year)
	if err != nil {
		return nil, err
	}
	return &ControllerStats{
		CID:            user.CID,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		ControllerType: user.ControllerType,
		Rating:         user.Rating.Short,
		Cab:            cab,
		Terminal:       terminal,
		Enroute:        enroute,
	}, nil
}
