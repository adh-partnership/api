package vatsim

import "time"

type VATSIMData struct {
	Controllers []*VATSIMController `json:"controllers"`
	Flights     []*VATSIMFlight     `json:"pilots"`
}

type VATSIMController struct {
	CID       int        `json:"cid"`
	Callsign  string     `json:"callsign"`
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
