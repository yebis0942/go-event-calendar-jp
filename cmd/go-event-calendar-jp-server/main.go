package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/syumai/workers"
	"github.com/syumai/workers/cloudflare"
	gojpcal "github.com/yebis0942/go-event-calendar-jp"
)

func main() {
	// Get API key from environment
	apiKey := cloudflare.Getenv("CONNPASS_API_KEY")
	if apiKey == "" {
		log.Fatal("CONNPASS_API_KEY environment variable is required")
	}

	// Initialize cache
	cache, err := NewCache()
	if err != nil {
		log.Fatalf("Failed to initialize cache: %v", err)
	}

	// Load configuration
	config := gojpcal.LoadConfig()

	// Set up HTTP handler
	http.HandleFunc("/calendar.ics", func(w http.ResponseWriter, r *http.Request) {
		// Get current year and month
		now := time.Now()

		// Try to get cached calendar
		cachedCalendar, ok, err := cache.Lookup()
		if err != nil {
			http.Error(w, "Failed to get calendar from cache", http.StatusInternalServerError)
			return
		}
		if ok {
			// Cache hit
			w.Header().Set("Content-Type", "text/calendar; charset=utf-8")
			fmt.Fprint(w, cachedCalendar)
			return
		}

		// Cache miss, generate new calendar
		// Initialize the client
		client := gojpcal.NewConnpassClient(apiKey)
		client.SetHTTPClient(&http.Client{
			Transport: gojpcal.NewWorkerTransport(),
		})

		// Fetch events
		yearMonths := GetMonthsRange(now.Year(), int(now.Month()))
		events, err := client.FetchEvents(context.Background(), config.ConnpassGroups, yearMonths)
		if err != nil {
			log.Printf("Error fetching events: %v", err)
			http.Error(w, "Failed to fetch events", http.StatusInternalServerError)
			return
		}

		// Generate iCalendar
		calendar, err := gojpcal.GenerateICalendar(events, now)
		if err != nil {
			log.Printf("Error generating calendar: %v", err)
			http.Error(w, "Failed to generate calendar", http.StatusInternalServerError)
			return
		}

		// Cache the generated calendar
		err = cache.Put(calendar)
		if err != nil {
			log.Printf("Warning: Failed to cache calendar: %v", err)
		}

		// Set content type and headers
		w.Header().Set("Content-Type", "text/calendar; charset=utf-8")

		// Write the calendar data
		fmt.Fprint(w, calendar)
	})

	workers.Serve(nil)
}
