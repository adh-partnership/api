package models

import (
	"encoding/json"
	"time"
)

type OAuthClient struct {
	ID           uint   `json:"id" gorm:"primaryKey"`
	Name         string `json:"name" gorm:"primaryKey"`
	ClientID     string `json:"client_id" gorm:"type:varchar(128)"`
	ClientSecret string `json:"-" gorm:"type:varchar(255)"`   // Will only be presented on creation, should not be otherwise exposed and must be regenerated
	RedirectURIs string `json:"return_uris" gorm:"type:text"` // Stringified JSON array of redirect URIs, no wildcards permitted
	TTL          int    `json:"ttl" gorm:"type:int"`          // How long can a token live in seconds
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (c *OAuthClient) ValidURI(uri string) (bool, error) {
	uris := []string{}
	err := json.Unmarshal([]byte(c.RedirectURIs), &uris)
	if err != nil {
		return false, err
	}
	for _, v := range uris {
		if uri == v {
			return true, nil
		}
	}

	return false, nil
}
