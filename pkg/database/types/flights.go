package database

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
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
