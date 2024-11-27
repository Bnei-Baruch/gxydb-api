package httputil

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// HTTPGetWithAuth sends a GET request to a specified URL with authorization header
func HTTPGetWithAuth(ctx context.Context, url string, auth string) (string, error) {
	if url == "" {
		return "", errors.New("URL cannot be empty")
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	if auth != "" {
		req.Header.Set("Authorization", auth)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer func() {
		if resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var body bytes.Buffer
		_, _ = io.Copy(&body, resp.Body)
		return "", fmt.Errorf("unexpected response status: %d, body: %s", resp.StatusCode, body.String())
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(body), nil
}
