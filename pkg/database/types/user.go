package database

import (
	"time"
)

type User struct {
	CID               uint   `json:"cid" gorm:"primaryKey" example:"876594"`
	FirstName         string `json:"firstname" gorm:"type:varchar(128)" example:"Daniel"`
	LastName          string `json:"lastname" gorm:"type:varchar(128)" example:"Hawton"`
	Email             string `json:"email" gorm:"type:varchar(128)" example:"wm@denartcc.org"`
	OperatingInitials string `json:"oi" gorm:"type:varchar(2)" example:"DH"`
	// Must be one of: none, active, inactive, loa
	ControllerType string `json:"controllerType" gorm:"type:varchar(10)" example:"home"`
	// Must be one of : none, training, solo, certified, major, cantrain
	DelCertification string `json:"delCertification" gorm:"type:varchar(15)" example:"certified"`
	// Must be one of : none, training, solo, certified, major, cantrain
	GndCertification string `json:"gndCertification" gorm:"type:varchar(15)" example:"certified"`
	// Must be one of : none, training, solo, certified, major, cantrain
	LclCertification string `json:"lclCertification" gorm:"type:varchar(15)" example:"certified"`
	// Must be one of : none, training, solo, certified, major, cantrain
	AppCertification string `json:"appCertification" gorm:"type:varchar(15)" example:"certified"`
	// Must be one of : none, training, solo, certified, major, cantrain
	CtrCertification string `json:"ctrCertification" gorm:"type:varchar(15)" example:"none"`
	RatingID         int    `json:"-"`
	Rating           Rating `json:"rating"`
	// Must be one of: none, active, inactive, loa
	Status    string  `json:"status" gorm:"type:varchar(10)" example:"active"`
	Roles     []*Role `json:"roles" gorm:"many2many:user_roles"`
	DiscordID string  `json:"discord_id" gorm:"type:varchar(128)" example:"123456789012345678"`
	Region    string  `json:"region" gorm:"type:varchar(10)" example:"AMAS"`
	Division  string  `json:"division" gorm:"type:varchar(10)" example:"USA"`
	// This may be blank
	Subdivision string `json:"subdivision" gorm:"type:varchar(10)" example:"ZDV"`
	// Internally used identifier during scheduled updates for removals
	UpdateID  string    `json:"updateid" gorm:"type:varchar(32)"`
	CreatedAt time.Time `json:"created_at" example:"2020-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2020-01-01T00:00:00Z"`
}

var CertificationOptions = map[string]string{
	"none":      "none",
	"training":  "training",
	"solo":      "solo",
	"certified": "certified",
	"major":     "major",
	"cantrain":  "cantrain",
}

var ControllerTypeOptions = map[string]string{
	"none":    "none",
	"visitor": "visitor",
	"home":    "home",
}

var ControllerStatusOptions = map[string]string{
	"none":     "none",
	"active":   "active",
	"inactive": "inactive",
	"loa":      "loa",
}
