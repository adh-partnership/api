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
