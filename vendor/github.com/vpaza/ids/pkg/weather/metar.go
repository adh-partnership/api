/*
 * Copyright Daniel Hawton
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

package weather

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
)

var ErrorTryAgain = fmt.Errorf("try again")

type addsResponse struct {
	METARs []METAR `xml:"data>METAR"`
}

// This is the AviationWeather ADDS data structure for the METAR block
type METAR struct {
	RawText            string         `xml:"raw_text"`
	StationID          string         `xml:"station_id"`
	ObservationTime    string         `xml:"observation_time"`
	Latitude           float64        `xml:"latitude"`
	Longitude          float64        `xml:"longitude"`
	Temperature        float64        `xml:"temp_c"`
	Dewpoint           float64        `xml:"dewpoint_c"`
	WindDirection      int            `xml:"wind_dir_degrees"`
	WindSpeed          int            `xml:"wind_speed_kt"`
	WindGust           int            `xml:"wind_gust_kt"`
	Visibility         float64        `xml:"visibility_statute_mi"`
	Altimeter          float64        `xml:"altim_in_hg"`
	SeaLevelPressure   float64        `xml:"sea_level_pressure_mb"`
	WxString           string         `xml:"wx_string"`
	SkyConditions      []SkyCondition `xml:"sky_condition"`
	FlightCategory     string         `xml:"flight_category"`
	VerticalVisibility int            `xml:"vert_vis_ft"`
	StationElevation   float64        `xml:"elevation_m"`
}

type SkyCondition struct {
	SkyCover  string `xml:"sky_cover,attr"`
	CloudBase int    `xml:"cloud_base_ft_agl,attr"`
}

func GetMetar(station string) (*METAR, error) {
	u := url.URL{
		Scheme: "https",
		Host:   "aviationweather.gov",
		Path:   "adds/dataserver_current/httpparam",
	}
	q := u.Query()
	q.Set("dataSource", "metars")
	q.Set("requestType", "retrieve")
	q.Set("format", "xml")
	q.Set("hoursBeforeNow", "1")
	q.Set("stationString", station)
	q.Set("mostRecent", "true")
	u.RawQuery = q.Encode()

	response, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode >= 500 {
		return nil, ErrorTryAgain
	} else if response.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP %d", response.StatusCode)
	}

	resp := addsResponse{}
	err = xml.NewDecoder(response.Body).Decode(&resp)
	if err != nil {
		return nil, err
	}

	if len(resp.METARs) == 0 {
		return nil, fmt.Errorf("no METARs returned")
	}

	return &resp.METARs[0], nil
}
