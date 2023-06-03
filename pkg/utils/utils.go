package utils

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"
	"strings"
)

func ArrayContains(array []string, item string) bool {
	for _, a := range array {
		if a == item {
			return true
		}
	}
	return false
}

func StringToSlug(s string) string {
	if len(s) > 100 {
		s = s[:99]
	}
	s = strings.ToLower(s)
	re := regexp.MustCompile(`[^a-z0-9 -]+`)
	s = re.ReplaceAllString(s, "")
	s = strings.Replace(s, " ", "-", -1)
	s = strings.Replace(s, "--", "-", -1)
	s = strings.TrimRight(s, "-")
	s = strings.TrimSpace(s)
	return s
}

func DumpToJSON(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}

func GetAirportTAF(icao string) ([]byte, error) {
	// Get TAF Data
	resp, err := http.Get("https://tgftp.nws.noaa.gov/data/forecasts/taf/stations/" + icao + ".TXT")
	if err != nil {
		return nil, errors.New("internal server error")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New("not found")
	}

	// Read Body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("internal server error")
	}

	return body, nil
}
