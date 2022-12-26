package vatusa

import (
	"encoding/json"
	"net/url"

	"github.com/adh-partnership/api/pkg/config"
	"github.com/adh-partnership/api/pkg/logger"
	"github.com/adh-partnership/api/pkg/network"
)

const (
	baseURL = "https://api.vatusa.net/v2"
)

var log = logger.Logger.WithField("component", "network/vatusa")

func handle(method, endpoint string, formdata map[string]string) (int, []byte, error) {
	data := url.Values{}
	for k, v := range formdata {
		data.Set(k, v)
	}

	u, err := url.Parse(baseURL + endpoint)
	if err != nil {
		return 0, nil, err
	}
	q := u.Query()
	q.Set("apikey", config.Cfg.VATUSA.APIKey)
	if config.Cfg.VATUSA.TestMode {
		q.Set("test", "true")
	}
	u.RawQuery = q.Encode()

	return network.Handle(method, u.String(), "application/x-www-form-urlencode", data.Encode())
}

func handleJSON(method, endpoint string, formdata map[string]string) (int, []byte, error) {
	u, err := url.Parse(baseURL + endpoint)
	if err != nil {
		return 0, nil, err
	}
	q := u.Query()
	q.Set("apikey", config.Cfg.VATUSA.APIKey)
	if config.Cfg.VATUSA.TestMode {
		q.Set("test", "true")
	}
	u.RawQuery = q.Encode()

	// Encode formdata to json
	data := url.Values{}
	for k, v := range formdata {
		data.Set(k, v)
	}
	j, _ := json.Marshal(data)

	return network.Handle(method, u.String(), "application/json", string(j))
}
