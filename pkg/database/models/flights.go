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

type Flights struct {
	ID          int       `json:"id" gorm:"primaryKey"`
	Callsign    string    `json:"callsign" gorm:"index;type:varchar(10)"`
	CID         int       `json:"cid" gorm:"index"`
	Facility    string    `json:"facility" gorm:"type:varchar(4)"`
	Latitude    float32   `json:"latitude" gorm:"type:float(10,8)"`
	Longitude   float32   `json:"longitude" gorm:"type:float(11,8)"`
	Groundspeed int       `json:"groundspeed"`
	Heading     int       `json:"heading"`
	Altitude    int       `json:"altitude"`
	Aircraft    string    `json:"aircraft" gorm:"type:varchar(10)"`
	Departure   string    `json:"departure" gorm:"type:varchar(4)"`
	Arrival     string    `json:"arrival" gorm:"type:varchar(4)"`
	Route       string    `json:"route" gorm:"type:text"`
	UpdateID    string    `json:"update_id" gorm:"type:varchar(36)"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
