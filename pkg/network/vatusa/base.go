/*
 * Copyright ADH Partnership
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package vatusa

import (
	"encoding/json"
	"net/url"

	"github.com/adh-partnership/api/pkg/config"
	"github.com/adh-partnership/api/pkg/logger"
	"github.com/adh-partnership/api/pkg/network"
)

const (
	baseURL = "https://api.vatusa.net"
)

var log = logger.Logger.WithField("component", "network/vatusa")

func handle(method, endpoint string, formdata map[string]string) (int, []byte, error) {
	u, err := url.Parse(baseURL + endpoint)
	if err != nil {
		return 0, nil, err
	}
	data := url.Values{}
	q := u.Query()
	q.Set("apikey", config.Cfg.VATUSA.APIKey)
	if config.Cfg.VATUSA.TestMode {
		q.Set("test", "true")
	}

	if method == "DELETE" {
		// VATUSA seems to have a problem with form data on delete requests... so add to query
		for k, v := range formdata {
			q.Set(k, v)
		}
	} else {
		for k, v := range formdata {
			data.Set(k, v)
		}
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
	j, _ := json.Marshal(formdata)

	return network.Handle(method, u.String(), "application/json", string(j))
}
