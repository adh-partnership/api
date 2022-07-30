package vatusa

import (
	"net/url"

	"github.com/kzdv/api/pkg/config"
	"github.com/kzdv/api/pkg/logger"
	"github.com/kzdv/api/pkg/network"
)

const (
	baseURL = "https://api.vatusa.net/v2"
)

var (
	FacilityID = "ZDV"
	log        = logger.Logger.WithField("component", "network/vatusa")
)

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
	u.RawQuery = q.Encode()

	return network.Handle(method, u.String(), "application/x-www-form-urlencode", data.Encode())
}
