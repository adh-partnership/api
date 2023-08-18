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

package vatsim

import "time"

type VATSIMData struct {
	Controllers []*VATSIMController `json:"controllers"`
	Flights     []*VATSIMFlight     `json:"pilots"`
}

type VATSIMController struct {
	CID       int        `json:"cid"`
	Callsign  string     `json:"callsign"`
	Name      string     `json:"name"`
	Frequency string     `json:"frequency"`
	Facility  int        `json:"facility"`
	Rating    int        `json:"rating"`
	LogonTime *time.Time `json:"logon_time"`
}

type VATSIMFlight struct {
	CID         int              `json:"cid"`
	Callsign    string           `json:"callsign"`
	Latitude    float64          `json:"latitude"`
	Longitude   float64          `json:"longitude"`
	Altitude    int              `json:"altitude"`
	Groundspeed int              `json:"groundspeed"`
	Heading     int              `json:"heading"`
	FlightPlan  VATSIMFlightPlan `json:"flight_plan"`
}

type VATSIMFlightPlan struct {
	Aircraft  string `json:"aircraft_faa"`
	Departure string `json:"departure"`
	Arrival   string `json:"arrival"`
	Route     string `json:"route"`
}
