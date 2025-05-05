package gojpcal

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestConnpassClient_FetchEvents(t *testing.T) {
	t.Parallel()

	event1 := Event{
		ID:        12345,
		Title:     "Go Conference 001",
		URL:       "https://connpass.com/event/12345/",
		StartTime: time.Date(2023, 10, 15, 10, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2023, 10, 15, 17, 0, 0, 0, time.UTC),
		Place:     "Tokyo Conference Center",
		Address:   "1-2-3 Shibuya, Tokyo",
	}
	event2 := Event{
		ID:        12346,
		Title:     "Go Meetup 001",
		URL:       "https://connpass.com/event/12346/",
		StartTime: time.Date(2023, 10, 20, 19, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2023, 10, 20, 21, 0, 0, 0, time.UTC),
		Place:     "Co-working Space",
		Address:   "4-5-6 Minato, Tokyo",
	}

	tests := map[string]struct {
		subdomains   []string
		yms          []string
		responseCode int
		responseBody ConnpassResponse
		wantErr      bool
		wantEvents   []Event
	}{
		"successful fetch": {
			subdomains:   []string{"golang", "go-tokyo"},
			yms:          []string{"202310"},
			responseCode: http.StatusOK,
			responseBody: ConnpassResponse{
				ResultsReturned:  2,
				ResultsStart:     1,
				ResultsAvailable: 2,
				Events:           []Event{event1, event2},
			},
			wantErr:    false,
			wantEvents: []Event{event1, event2},
		},
		"empty result": {
			subdomains:   []string{"nonexistent"},
			yms:          []string{"202310"},
			responseCode: http.StatusOK,
			responseBody: ConnpassResponse{
				ResultsReturned:  0,
				ResultsStart:     1,
				ResultsAvailable: 0,
				Events:           []Event{},
			},
			wantErr:    false,
			wantEvents: []Event{},
		},
		"api error": {
			subdomains:   []string{"golang"},
			yms:          []string{"202310"},
			responseCode: http.StatusInternalServerError,
			wantErr:      true,
			wantEvents:   nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			isRequestReceived := false

			// Create a test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				isRequestReceived = true

				// Check request method
				require.Equal(t, http.MethodGet, r.Method, "Request method doesn't match")

				// Check API key header
				apiKey := r.Header.Get("X-API-Key")
				require.Equal(t, "test-api-key", apiKey, "API key doesn't match")

				// Check query parameters
				q := r.URL.Query()
				require.Equal(t, strings.Join(tc.yms, ","), q.Get("ym"), "Year-month parameter doesn't match")

				// Check subdomain parameter
				require.Equal(t, strings.Join(tc.subdomains, ","), q.Get("subdomain"), "Subdomain parameter doesn't match")

				// Set response status code
				w.WriteHeader(tc.responseCode)

				// Set response body
				if tc.responseCode == http.StatusOK {
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(tc.responseBody)
				}
			}))
			defer server.Close()

			// Create client with test API key
			client := NewConnpassClient("test-api-key")

			// Override the API URL to use our test server
			client.SetConnpassAPIURL(server.URL + "/")

			// Call the method
			events, err := client.FetchEvents(context.Background(), tc.subdomains, tc.yms)

			require.True(t, isRequestReceived, "Request was not received by the test server")

			// Check the error
			if tc.wantErr {
				require.Error(t, err, "Expected an error but got nil")
				return
			}
			require.NoError(t, err, "Expected no error but got one")

			// Check the events
			require.Equal(t, tc.wantEvents, events, "Events don't match")
		})
	}
}
