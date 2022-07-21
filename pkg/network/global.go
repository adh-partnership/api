package network

import (
	"github.com/kzdv/api/pkg/network/vatsim"
	"github.com/kzdv/api/pkg/network/vatusa"
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
		sub, err := vatusa.GetFacility(cid)
		if err != nil {
			return Location{}, err
		}
		loc.Subdivision = sub
	}

	return loc, nil
}
