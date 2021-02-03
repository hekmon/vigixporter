package hubeau

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	// RequestMaxSize is the maximum number of series that is queriable in one request
	RequestMaxSize = 20000
	baseURLStr     = "https://hubeau.eaufrance.fr/api/v1/"
	maxURL         = 2083
	devMode        = false
)

var (
	baseURL *url.URL
)

func init() {
	var err error
	baseURL, err = url.Parse(baseURLStr)
	if err != nil {
		panic(err)
	}
}

func (c *Controller) request(ctx context.Context, method, path string, queryParams url.Values, output interface{}) (err error) {
	// Build url
	URL := *baseURL
	URL.Path += path
	URL.RawQuery = queryParams.Encode()
	// Build query
	URLstr := URL.String()
	if len(URLstr) > maxURL {
		return fmt.Errorf("crafted URL is longer than %d which is the limit hubeau can process: %s", maxURL, URLstr)
	}
	req, err := http.NewRequestWithContext(ctx, method, URL.String(), nil)
	if err != nil {
		return fmt.Errorf("building HTTP query failed: %w", err)
	}
	// Execute query
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute http query: %w", err)
	}
	defer resp.Body.Close()
	// Handle answer
	switch resp.StatusCode {
	case http.StatusOK, http.StatusPartialContent:
		// noop (continue below)
	default:
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("unexpected HTTP code: %s. Failed to read body: %w", resp.Status, err)
		}
		return fmt.Errorf("unexpected HTTP code: %s. Body:\n%s", resp.Status, string(bodyBytes))
	}
	// Unmarshall
	jsonBodyDecoder := json.NewDecoder(resp.Body)
	if devMode {
		jsonBodyDecoder.DisallowUnknownFields()
	}
	if err = jsonBodyDecoder.Decode(output); err != nil {
		err = fmt.Errorf("answer payload json unmarshalling failed: %w", err)
	}
	return
}
