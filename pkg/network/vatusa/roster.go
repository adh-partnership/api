package vatusa

import (
	"encoding/json"
	"fmt"

	"github.com/adh-partnership/api/pkg/config"
)

type VATUSAController struct {
	CID         int    `json:"cid"`
	FirstName   string `json:"fname"`
	LastName    string `json:"lname"`
	Email       string `json:"email"`
	Rating      int    `json:"rating"`
	Facility    string `json:"facility"`
	Membership  string `json:"membership"`
	RatingShort string `json:"rating_short"`
}

type VATUSAFacility struct {
	Info  *VATUSAFacilityInfo    `json:"info"`
	Roles []*VATUSAFacilityRole  `json:"roles"`
	Stats map[string]interface{} `json:"stats"`
}

type VATUSAFacilityInfo struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	URL    string `json:"url"`
	Region int    `json:"region"`
	ATM    int    `json:"atm"`
	DATM   int    `json:"datm"`
	TA     int    `json:"ta"`
	EC     int    `json:"ec"`
	FE     int    `json:"fe"`
	WM     int    `json:"wm"`
}

type VATUSAFacilityRole struct {
	ID       int    `json:"id"`
	CID      int    `json:"cid"`
	Facility string `json:"facility"`
	Role     string `json:"role"`
}

// RemoveController removes a home controller from the roster at VATUSA.
func RemoveController(cid string, by uint, reason string) (int, error) {
	status, _, err := handle("DELETE", "/facility/"+config.Cfg.VATUSA.Facility+"/roster/"+cid, map[string]string{
		"by":     fmt.Sprint(by),
		"reason": reason,
	})

	return status, err
}

// RemoveVisitingController removes a controller from the visiting roster at VATUSA.
func RemoveVisitingController(cid string, by uint, reason string) (int, error) {
	status, _, err := handle("DELETE", "/facility/"+config.Cfg.VATUSA.Facility+"/roster/visiting/"+cid, map[string]string{
		"by":     fmt.Sprint(by),
		"reason": reason,
	})

	return status, err
}

// AddVisitingController adds a visitor to the visiting roster at VATUSA.
func AddVisitingController(cid string) (int, error) {
	status, _, err := handle("POST", "/facility/"+config.Cfg.VATUSA.Facility+"/roster/manageVisitor/"+cid, nil)

	return status, err
}

func GetFacility(id string) (*VATUSAFacility, error) {
	status, content, err := handle("GET", "/facility/"+id, nil)
	if err != nil {
		return nil, err
	}

	if status > 299 {
		log.Warnf("Failed to get facility %s: %s", id, content)
		return nil, fmt.Errorf("invalid status code: %d", status)
	}

	var facility VATUSAFacility

	err = json.Unmarshal(content, &facility)
	if err != nil {
		return nil, err
	}

	return &facility, nil
}

// GetFacilityRoster will grab the facility roster from VATUSA. If membership is not specified, or it is not "home" or
// "visit", it will return both rosters.
func GetFacilityRoster(membership string) ([]VATUSAController, error) {
	if membership == "" || (membership != "home" && membership != "visit") {
		membership = "both"
	}
	status, content, err := handle("GET", "/facility/"+config.Cfg.VATUSA.Facility+"/roster/"+membership, nil)
	if err != nil {
		return nil, err
	}

	if status > 299 {
		log.Warnf("Failed to get facility roster: %s", content)
		return nil, fmt.Errorf("invalid status code: %d", status)
	}

	type response struct {
		Data []VATUSAController `json:"data"`
	}
	r := response{}

	err = json.Unmarshal(content, &r)
	if err != nil {
		return nil, err
	}

	return r.Data, nil
}

// GetFacility returns the Facility for a user in VATUSA.
//
// If a user is not in VATUSA's roster, it will return an error of "user not found". Other errors
// will be returned as-is.
//
// VATUSA-ism: a controller not assigned to a facility will be in "ZAE". These can be controllers that are
// in other divisions/regions, so you should also check the division/region from the VATSIM API pkg's GetLocation method.
func GetUserFacility(cid string) (string, error) {
	status, content, err := handle("GET", "/user/"+cid, nil)
	if err != nil {
		return "", err
	}

	if status == 404 {
		return "", fmt.Errorf("user not found")
	}

	if status > 299 {
		log.Warnf("Failed to get facility for %s: %s", cid, content)
		return "", fmt.Errorf("invalid status code: %d", status)
	}

	var facility struct {
		Facility string `json:"facility"`
	}

	err = json.Unmarshal(content, &facility)
	if err != nil {
		return "", err
	}

	return facility.Facility, nil
}
