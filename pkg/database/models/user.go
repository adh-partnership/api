/*
 * Copyright ADH Partnership
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package models

import (
	"time"
)

type User struct {
	CID               uint   `json:"cid" gorm:"primaryKey" example:"876594"`
	FirstName         string `json:"first_name" gorm:"type:varchar(128)" example:"Daniel"`
	LastName          string `json:"last_name" gorm:"type:varchar(128)" example:"Hawton"`
	Email             string `json:"email" gorm:"type:varchar(128)" example:"wm@denartcc.org"`
	OperatingInitials string `json:"operating_initials" gorm:"type:varchar(2)" example:"DH"`
	// Must be one of: none, active, inactive, loa
	ControllerType       string `json:"controllerType" gorm:"type:varchar(10)" example:"home"`
	ExemptedFromActivity bool   `json:"exemptedFromActivity" gorm:"default:false" example:"false"`
	RatingID             int    `json:"-"`
	Rating               Rating `json:"rating"`
	// Must be one of: none, active, inactive, loa
	Status    string  `json:"status" gorm:"type:varchar(10)" example:"active"`
	Roles     []*Role `json:"roles" gorm:"many2many:user_roles"`
	DiscordID string  `json:"discord_id" gorm:"type:varchar(128)" example:"123456789012345678"`
	Region    string  `json:"region" gorm:"type:varchar(10)" example:"AMAS"`
	Division  string  `json:"division" gorm:"type:varchar(10)" example:"USA"`
	// This may be blank
	Subdivision string `json:"subdivision" gorm:"type:varchar(10)" example:"ZDV"`
	// Internally used identifier during scheduled updates for removals
	UpdateID       string     `json:"updateid" gorm:"type:varchar(32)"`
	RosterJoinDate *time.Time `json:"roster_join_date" example:"2020-01-01T00:00:00Z"`
	CreatedAt      time.Time  `json:"created_at" example:"2020-01-01T00:00:00Z"`
	UpdatedAt      time.Time  `json:"updated_at" example:"2020-01-01T00:00:00Z"`
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
