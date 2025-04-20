package gojpcal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const connpassAPIURL = "https://connpass.com/api/v2/"

type ConnpassClient struct {
	apiKey         string
	httpClient     *http.Client
	connpassAPIURL string
}

func NewConnpassClient(apiKey string) *ConnpassClient {
	return &ConnpassClient{
		apiKey:         apiKey,
		httpClient:     http.DefaultClient,
		connpassAPIURL: connpassAPIURL,
	}
}

func (c *ConnpassClient) SetHTTPClient(httpClient *http.Client) {
	c.httpClient = httpClient
}

func (c *ConnpassClient) SetConnpassAPIURL(connpassAPIURL string) {
	c.connpassAPIURL = connpassAPIURL
}

// FetchEvents fetches events from Connpass API
func (c *ConnpassClient) FetchEvents(groupSubdomains, yearMonths []string) ([]Event, error) {
	// Build query parameters
	params := url.Values{}
	params.Add("subdomain", strings.Join(groupSubdomains, ","))
	params.Add("ym", strings.Join(yearMonths, ","))

	// Build request URL
	reqURL := c.connpassAPIURL + "events/" + "?" + params.Encode()

	// Create a new request
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add API key header
	req.Header.Add("X-API-Key", c.apiKey)

	// Make HTTP request
	resp, err := c.httpClient.Do(req)
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
