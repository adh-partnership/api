package vatusa

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/kzdv/api/pkg/config"
)

const (
	baseUrl = "https://api.vatusa.net/v2"
)

func handle(method, endpoint string, formdata map[string]string) (int, []byte, error) {
	data := url.Values{}
	for k, v := range formdata {
		data.Set(k, v)
	}

	u, err := url.Parse(baseUrl + endpoint)
	if err != nil {
		return 0, nil, err
	}
	q := u.Query()
	q.Set("apikey", config.Cfg.VATUSA.APIKey)
	u.RawQuery = q.Encode()

	r, err := http.NewRequest(method, u.String(), strings.NewReader(data.Encode()))
	if err != nil {
		return 0, nil, err
	}
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("User-Agent", "kzdv-api")
	r.Header.Set("Accept", "application/json")

	defer r.Body.Close()
	contents, err := io.ReadAll(r.Body)
	if err != nil {
		return 0, nil, err
	}

	return r.Response.StatusCode, contents, nil
}
