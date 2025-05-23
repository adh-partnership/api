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

type Weather struct {
	METAR string `json:"metar"`
	TAF   string `json:"taf"`
}

type response struct {
	METARs []METAR `xml:"data>METAR"`
	TAFs   []TAF   `xml:"data>TAF"`
}

type METAR struct {
	StationID string `xml:"station_id"`
	RawText   string `xml:"raw_text"`
}

type TAF struct {
	StationID string `xml:"station_id"`
	RawText   string `xml:"raw_text"`
}
