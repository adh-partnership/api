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

type Airport struct {
	ID          string      `json:"arpt_id" example:"FAI" gorm:"primary_key"`
	ICAO        string      `json:"icao_id" example:"PAFA" gorm:"index"`
	State       string      `json:"state_code" example:"AK"`
	City        string      `json:"city" example:"Fairbanks"`
	Name        string      `json:"arpt_name" example:"FAIRBANKS INTL AIRPORT"`
	ARTCC       string      `json:"resp_artcc_id" example:"ZNY" gorm:"index"`
	Status      string      `json:"arpt_status" example:"O"`
	TwrTypeCode string      `json:"twr_type_code" example:"T"`
	Elevation   float32     `json:"elevation" example:"13.0"`
	Latitude    float32     `json:"latitude" example:"40.639801"`
	Longitude   float32     `json:"longitude" example:"-73.778900"`
	ATC         *AirportATC `json:"atc,omitempty" gorm:"foreignKey:ID;references:ID"`
	HasMETAR    bool        `json:"-"`
	HasTAF      bool        `json:"-"`
	METAR       string      `json:"metar,omitempty"`
	TAF         string      `json:"taf,omitempty"`
	CreatedAt   time.Time   `json:"created_at" example:"2021-09-01T00:00:00Z"`
	UpdatedAt   time.Time   `json:"updated_at" example:"2021-09-01T00:00:00Z"`
}

type AirportATC struct {
	ID                     string    `json:"arpt_id" example:"JFK" gorm:"index"`
	FacilityType           string    `json:"facility_type" example:"ATCT-A/C"`
	TwrOperatorCode        string    `json:"twr_operator_code" example:"A"`
	TwrCall                string    `json:"twr_call" example:"Fairbanks"`
	TwrHrs                 string    `json:"twr_hrs" example:"1200-0100 local time"`
	PrimaryApchRadioCall   string    `json:"primary_apch_radio_call" example:"Fairbanks"`
	ApchPProvider          string    `json:"apch_p_provider" example:"FAI"`
	ApchPProvTypeCD        string    `json:"apch_p_prov_type_cd" example:"A"`
	SecondaryApchRadioCall string    `json:"secondary_apch_radio_call" example:"Anchorage"`
	ApchSProvider          string    `json:"apch_s_provider" example:"ANC"`
	ApchSProvTypeCD        string    `json:"apch_s_prov_type_cd" example:"A"`
	PrimaryDepRadioCall    string    `json:"primary_dep_radio_call" example:"Fairbanks"`
	DepPProvider           string    `json:"dep_p_provider" example:"FAI"`
	DepPProvTypeCD         string    `json:"dep_p_prov_type_cd" example:"A"`
	SecondaryDepRadioCall  string    `json:"secondary_dep_radio_call" example:"Anchorage"`
	DepSProvider           string    `json:"dep_s_provider" example:"ANC"`
	DepSProvTypeCD         string    `json:"dep_s_prov_type_cd" example:"A"`
	CreatedAt              time.Time `json:"created_at" example:"2021-09-01T00:00:00Z"`
	UpdatedAt              time.Time `json:"updated_at" example:"2021-09-01T00:00:00Z"`
}

type AirportChart struct {
	ID        string    `json:"-" gorm:"primary_key"`
	AirportID string    `json:"arpt_id" example:"FAI" gorm:"index"`
	Cycle     int       `json:"cycle" example:"2213" gorm:"index"`
	FromDate  time.Time `json:"from_date" example:"2021-09-01T00:00:00Z"`
	ToDate    time.Time `json:"to_date" example:"2021-09-30T00:00:00Z"`
	ChartCode string    `json:"chart_code" example:"DP"`
	ChartName string    `json:"chart_name" example:"RDFLG FOUR (RNAV)"`
	ChartURL  string    `json:"chart_url" example:"https://aeronav.faa.gov/d-tpp/2212/01234RDFLG.PDF"`
	CreatedAt time.Time `json:"created_at" example:"2021-09-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2021-09-01T00:00:00Z"`
}
