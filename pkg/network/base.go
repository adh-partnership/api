package network

import (
	"io"
	"net/http"
	"strings"
)

// UserAgent is the user agent to pass in request headers
var UserAgent = "kzdv-network"

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
	if err != nil {
		return 0, nil, err
	}
	r.Header.Set("Content-Type", contenttype)
	r.Header.Set("User-Agent", UserAgent)
	r.Header.Set("Accept", "application/json")
	for k, v := range headers {
		r.Header.Set(k, v)
	}

	defer func() {
		_ = r.Body.Close()
	}()

	contents, err := io.ReadAll(r.Body)
	if err != nil {
		return 0, nil, err
	}

	return r.Response.StatusCode, contents, nil
}
