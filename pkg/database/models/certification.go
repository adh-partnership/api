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

import "time"

type Certification struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Order       uint   `json:"order"`
	Hidden      bool   `json:"hidden"`
}

type UserCertification struct {
	ID        uint      `json:"id"`
	CID       uint      `json:"cid" gorm:"index:idx_cid,index:idx_cid_certification"`
	Name      string    `json:"name" gorm:"index:idx_cid_certification"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var CertificationOptions = map[string]string{
	"none":      "none",
	"training":  "training",
	"solo":      "solo",
	"certified": "certified",
	"cantrain":  "cantrain",
}

func IValidUserCertificationValue(value string) bool {
	_, ok := CertificationOptions[value]
	return ok
}
