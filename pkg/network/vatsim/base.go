package vatsim

import (
	"encoding/json"
	"github.com/kzdv/api/pkg/config"
	"github.com/kzdv/api/pkg/logger"
	"github.com/kzdv/api/pkg/network"
	"net/url"
)

const (
	baseUrl = "https://api.vatsim.net/api"
)

var log = logger.Logger.WithField("component", "network/vatsim")

func handle(method, endpoint string, formdata map[string]string) (int, []byte, error) {
	data, err := json.Marshal(formdata)
	if err != nil {
		return 0, nil, err
	}

	u, err := url.Parse(baseUrl + endpoint)
	if err != nil {
		return 0, nil, err
	}
	q := u.Query()
	q.Set("apikey", config.Cfg.VATUSA.APIKey)
	u.RawQuery = q.Encode()

	return network.Handle(method, u.String(), "application/json", string(data))
}
