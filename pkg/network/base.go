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

package network

import (
	"io"
	"net/http"
	"strings"

	"moul.io/http2curl"

	"github.com/adh-partnership/api/pkg/logger"
)

// UserAgent is the user agent to pass in request headers
var UserAgent = "adh-partnership-network"

var log = logger.Logger.WithField("component", "network")

// Handle will make the request as presented in a structured and expected way. It adds appropriate headers, including
// a user agent so our requests can be known to come from us.
//
// Handle is an alias of HandleWithHeaders that passes no extra headers
func Handle(method, url, contenttype, body string) (int, []byte, error) {
	return HandleWithHeaders(method, url, contenttype, body, map[string]string{})
}

// HandleWithHeaders will make the request as presented in a structured and expected way. It adds appropriate headers, including
// a user agent so our requests can be known to come from us.
func HandleWithHeaders(method, url, contenttype, body string, headers map[string]string) (int, []byte, error) {
	r, err := http.NewRequest(method, url, strings.NewReader(body))
	log.Debugf("Making request: %s %s", method, url)
	log.Debugf("with Body: %s", body)
	if err != nil {
		return 0, nil, err
	}
	r.Header.Set("Content-Type", contenttype)
	r.Header.Set("User-Agent", UserAgent)
	r.Header.Set("Accept", "application/json")
	for k, v := range headers {
		r.Header.Set(k, v)
	}

	curl, _ := http2curl.GetCurlCommand(r)
	log.Debugf("Request curl equivalent: %s", curl)

	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		return 0, nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	contents, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}

	log.Debugf("Response from %s %s: %d", method, url, resp.StatusCode)
	log.Debugf("Response body: %s", string(contents))

	return resp.StatusCode, contents, nil
}
