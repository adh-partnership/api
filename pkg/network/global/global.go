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

package global

import (
	"github.com/adh-partnership/api/pkg/network/vatsim"
	"github.com/adh-partnership/api/pkg/network/vatusa"
)

type Location struct {
	Region      string
	Division    string
	Subdivision string
}

// GetLocation gets a controller's region, division, and subdivision
// from the VATSIM API, and where applicable, the VATUSA API.
func GetLocation(cid string) (Location, error) {
	reg, div, sub, err := vatsim.GetLocation(cid)
	if err != nil {
		return Location{}, err
	}

	loc := Location{
		Region:      reg,
		Division:    div,
		Subdivision: sub,
	}

	// Check if user is in VATUSA, and if so, query the VATUSA API
	// for the specific facility. VATUSA and VATCAN, are the only two
	// that I know of that do not report subdivisions to VATSIM.
	//
	// We cannot query VATCAN without being a VATCAN FIR, so we will
	// just report VATCAN as "AMAS, CAN, ''" for Canadians.
	if loc.Region == "AMAS" && loc.Division == "USA" {
		sub, err := vatusa.GetUserFacility(cid)
		if err != nil {
			return Location{}, err
		}
		loc.Subdivision = sub
	}

	return loc, nil
}
