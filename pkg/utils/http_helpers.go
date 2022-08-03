package utils

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Get is a helper function to make a GET request to the specified URL and return the response body as a string.
func Get(ctx context.Context, url string, headers map[string]string, timeout time.Duration) (string, error) {
	client := http.Client{
		Transport: http.DefaultTransport,
		Timeout:   timeout,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	for header, value := range headers {
		req.Header.Add(header, value)
	}

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if res.StatusCode != 200 {
		return "", fmt.Errorf("status code %d trying to GET %s", res.StatusCode, url)
	}

	defer res.Body.Close()
	all, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("error while reading response from oraclecloud metadata endpoint: %s", err)
	}

	return string(all), nil
}
