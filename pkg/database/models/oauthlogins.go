package models

import "time"

type OAuthLogin struct {
	ID                  uint        `json:"id" gorm:"primaryKey"`
	Token               string      `json:"token" gorm:"type:varchar(128)"`
	Code                string      `json:"code" gorm:"varchar(48)"`
	UserAgent           string      `json:"ua" gorm:"type:varchar(255)"`
	IP                  string      `json:"ip" gorm:"type:varchar(128)"`
	RedirectURI         string      `json:"url" gorm:"type:varchar(255)"`
	ClientID            uint        `json:"-"`
	Client              OAuthClient `json:"-"`
	State               string      `json:"state"`
	CodeChallenge       string      `json:"-"`
	CodeChallengeMethod string      `json:"-"`
	Scope               string      `json:"-"`
	Nonce               string      `json:"type:varchar(255)"`
	CID                 uint        `json:"cid"`
	ExpiresAt           time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
