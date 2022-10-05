package vatsim

import (
	"encoding/json"
	"net/url"

	"github.com/adh-partnership/api/pkg/logger"
	"github.com/adh-partnership/api/pkg/network"
)

const (
	baseURL = "https://api.vatsim.net/api"
	dataURL = "https://data.vatsim.net/v3/vatsim-data.json"
)

var log = logger.Logger.WithField("component", "network/vatsim")

func handle(method, endpoint string, formdata map[string]string) (int, []byte, error) {
	return handleFull(method, baseURL+endpoint, formdata)
}

func handleData() (int, []byte, error) {
	return handleFull("GET", dataURL, nil)
}

func handleFull(method, fullURL string, formdata map[string]string) (int, []byte, error) {
	data, err := json.Marshal(formdata)
	if err != nil {
		return 0, nil, err
	}

	u, err := url.Parse(fullURL)
	if err != nil {
		return 0, nil, err
	}
	q := u.Query()
	u.RawQuery = q.Encode()

	return network.Handle(method, u.String(), "application/json", string(data))
}
