package vatsim

import (
	"encoding/json"
	"fmt"
)

var Ratings = [13]string{
	"",
	"OBS",
	"S1",
	"S2",
	"S3",
	"C1",
	"C2",
	"C3",
	"I1",
	"I2",
	"I3",
	"SUP",
	"ADM",
}

// GetRating returns the rating of a VATSIM CID from the VATSIM API.
// The rating is a integar that represents the rating "id".
// Typical: OBS=1, S1=2, S2=3, S3=4, C1=5, C2=6, C3=7, I1=8, I2=9, I3=10, SUP=11, ADM=12
func GetRating(cid string) (int, error) {
	status, contents, err := handle("GET", "/ratings/"+cid+"/", nil)
	if err != nil {
		return 0, err
	}

	if status > 299 {
		log.Warnf("Failed to get rating for %s: %s", cid, contents)
		return 0, fmt.Errorf("invalid status code: %d", status)
	}

	var rating struct {
		Rating int `json:"rating"`
	}
	err = json.Unmarshal(contents, &rating)
	if err != nil {
		return 0, err
	}

	return rating.Rating, nil
}

// GetLocation returns the Region, Division and Subdivision of a VATSIM CID
// from the VATSIM API.
// Returns: Region, Division, Subdivision, error
func GetLocation(cid string) (string, string, string, error) {
	status, contents, err := handle("GET", "/ratings/"+cid+"/", nil)
	if err != nil {
		return "", "", "", err
	}

	if status > 299 {
		log.Warnf("Failed to get division for %s: %s", cid, contents)
		return "", "", "", fmt.Errorf("invalid status code: %d", status)
	}

	var division struct {
		Region      string `json:"region"`
		Division    string `json:"division"`
		Subdivision string `json:"subdivision"`
	}
	err = json.Unmarshal(contents, &division)
	if err != nil {
		return "", "", "", err
	}

	return division.Region, division.Division, division.Subdivision, nil
}

func GetData() (*VATSIMData, error) {
	status, contents, err := handleData()
	if err != nil {
		return nil, err
	}

	if status > 299 {
		log.Warnf("Failed to get data: %s", contents)
		return nil, fmt.Errorf("invalid status code: %d", status)
	}

	data := &VATSIMData{}
	err = json.Unmarshal(contents, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
