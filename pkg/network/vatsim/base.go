package vatsim

import (
	"encoding/json"
	"net/url"

	"github.com/kzdv/api/pkg/logger"
	"github.com/kzdv/api/pkg/network"
)

const (
	baseURL = "https://api.vatsim.net/api"
)

var log = logger.Logger.WithField("component", "network/vatsim")

func handle(method, endpoint string, formdata map[string]string) (int, []byte, error) {
	data, err := json.Marshal(formdata)
	if err != nil {
		return 0, nil, err
	}

	u, err := url.Parse(baseURL + endpoint)
	if err != nil {
		return 0, nil, err
	}
	q := u.Query()
	u.RawQuery = q.Encode()

	return network.Handle(method, u.String(), "application/json", string(data))
}
