package vmpusher

import (
	"bytes"
	"fmt"
	"net/http"
)

const (
	jsonStreamingContentType = "application/x-ndjson; charset=utf-8" // https://github.com/wardi/jsonlines/issues/9#issuecomment-143550948
)

func (c *Controller) push(payload string) (err error) {
	// Prepare request
	req, err := http.NewRequestWithContext(c.ctx, "POST", c.vmURL, bytes.NewBufferString(payload))
	if err != nil {
		return fmt.Errorf("failed to build the http request: %w", err)
	}
	req.Header.Set("Content-Type", jsonStreamingContentType)
	if c.vmUser != "" && c.vmPass != "" {
		req.SetBasicAuth(c.vmUser, c.vmPass)
	}
	// Send it
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute http request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected HTTP code: %s", resp.Status)
	}
	return
}
