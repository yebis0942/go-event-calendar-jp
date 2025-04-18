package gojpcal

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestConnpassClient_FetchEvents(t *testing.T) {
	t.Parallel()

	event1 := Event{
		ID:          12345,
		Title:       "Go Conference 001",
		URL:         "https://connpass.com/event/12345/",
		Description: "A conference for Gophers",
		StartTime:   time.Date(2023, 10, 15, 10, 0, 0, 0, time.UTC),
		EndTime:     time.Date(2023, 10, 15, 17, 0, 0, 0, time.UTC),
		Place:       "Tokyo Conference Center",
		Address:     "1-2-3 Shibuya, Tokyo",
	}
	event2 := Event{
		ID:          12346,
		Title:       "Go Meetup 001",
		URL:         "https://connpass.com/event/12346/",
		Description: "A monthly Go meetup",
		StartTime:   time.Date(2023, 10, 20, 19, 0, 0, 0, time.UTC),
		EndTime:     time.Date(2023, 10, 20, 21, 0, 0, 0, time.UTC),
		Place:       "Co-working Space",
		Address:     "4-5-6 Minato, Tokyo",
	}

	tests := map[string]struct {
		subdomains   []string
		year         string
		month        string
		responseCode int
		responseBody ConnpassResponse
		wantErr      bool
		wantEvents   []Event
	}{
		"successful fetch": {
			subdomains:   []string{"golang", "go-tokyo"},
			year:         "2023",
			month:        "10",
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
			year:         "2023",
			month:        "10",
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
			year:         "2023",
			month:        "10",
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
				if r.Method != http.MethodGet {
					t.Errorf("Expected method %s, got %s", http.MethodGet, r.Method)
				}

				// Check API key header
				apiKey := r.Header.Get("X-API-Key")
				if apiKey != "test-api-key" {
					t.Errorf("Expected API key %s, got %s", "test-api-key", apiKey)
				}

				// Check query parameters
				q := r.URL.Query()
				if q.Get("ym") != tc.year+tc.month {
					t.Errorf("Expected ym %s, got %s", tc.year+tc.month, q.Get("ym"))
				}

				// Check subdomain parameter
				if strings.Join(tc.subdomains, ",") != q.Get("subdomain") {
					t.Errorf("Expected subdomain %s, got %s", strings.Join(tc.subdomains, ","), q.Get("subdomain"))
				}

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
			events, err := client.FetchEvents(tc.subdomains, tc.year, tc.month)

			if !isRequestReceived {
				t.Fatal("Request was not received by the test server")
			}

			// Check results
			if tc.wantErr {
				if err == nil {
					t.Error("Expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			if len(events) != len(tc.wantEvents) {
				t.Errorf("Expected %d events, got %d", len(tc.wantEvents), len(events))
			}

			for i, wantEvent := range tc.wantEvents {
				if wantEvent.ID != events[i].ID {
					t.Errorf("Event ID mismatch at index %d: expected %d, got %d", i, wantEvent.ID, events[i].ID)
				}
				if wantEvent.Title != events[i].Title {
					t.Errorf("Event Title mismatch at index %d: expected %s, got %s", i, wantEvent.Title, events[i].Title)
				}
				if wantEvent.URL != events[i].URL {
					t.Errorf("Event URL mismatch at index %d: expected %s, got %s", i, wantEvent.URL, events[i].URL)
				}
				if wantEvent.Description != events[i].Description {
					t.Errorf("Event Description mismatch at index %d", i)
				}
				if !wantEvent.StartTime.Equal(events[i].StartTime) {
					t.Errorf("Event StartTime mismatch at index %d: expected %s, got %s",
						i, wantEvent.StartTime.Format(time.RFC3339), events[i].StartTime.Format(time.RFC3339))
				}
				if !wantEvent.EndTime.Equal(events[i].EndTime) {
					t.Errorf("Event EndTime mismatch at index %d: expected %s, got %s",
						i, wantEvent.EndTime.Format(time.RFC3339), events[i].EndTime.Format(time.RFC3339))
				}
				if wantEvent.Place != events[i].Place {
					t.Errorf("Event Place mismatch at index %d: expected %s, got %s", i, wantEvent.Place, events[i].Place)
				}
				if wantEvent.Address != events[i].Address {
					t.Errorf("Event Address mismatch at index %d: expected %s, got %s", i, wantEvent.Address, events[i].Address)
				}
				if wantEvent.GroupID != events[i].GroupID {
					t.Errorf("Event GroupID mismatch at index %d: expected %s, got %s", i, wantEvent.GroupID, events[i].GroupID)
				}
				if wantEvent.GroupTitle != events[i].GroupTitle {
					t.Errorf("Event GroupTitle mismatch at index %d: expected %s, got %s", i, wantEvent.GroupTitle, events[i].GroupTitle)
				}
			}
		})
	}
}
