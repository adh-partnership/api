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

type OnlineController struct {
	ID        int       `json:"id" example:"1"`
	UserID    uint      `json:"cid" example:"1"`
	User      *User     `json:"user"`
	Position  string    `json:"position" example:"ANC_TWR" gorm:"index"`
	Frequency string    `json:"frequency" example:"118.000"`
	LogonTime time.Time `json:"logon_time" example:"2020-01-01T00:00:00Z"`
	UpdateID  string    `json:"update_id" example:"1"`
	CreatedAt time.Time `json:"created_at" example:"2020-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2020-01-01T00:00:00Z"`
}
