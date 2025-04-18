package gojpcal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const connpassAPIURL = "https://connpass.com/api/v2/events/"

// FetchEvents fetches events from Connpass API for the specified year and month
func FetchEvents(apiKey string, groupSubdomains []string, year, month string) ([]Event, error) {
	// Build query parameters
	params := url.Values{}
	params.Add("subdomain", strings.Join(groupSubdomains, ","))
	params.Add("ym", year+month)

	// Build request URL
	reqURL := connpassAPIURL + "?" + params.Encode()

	// Create a new request
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add API key header
	req.Header.Add("X-API-Key", apiKey)

	// Make HTTP request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-OK status: %s", resp.Status)
	}

	// Parse response
	var connpassResp ConnpassResponse
	err = json.NewDecoder(resp.Body).Decode(&connpassResp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse API response: %w", err)
	}

	return connpassResp.Events, nil
}
