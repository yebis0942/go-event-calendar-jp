package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	gojpcal "github.com/yebis0942/golang-jp-event-calendar"
)

func main() {
	// Get API key from environment
	apiKey := os.Getenv("CONNPASS_API_KEY")
	if apiKey == "" {
		log.Fatal("CONNPASS_API_KEY environment variable is required")
	}

	// Initialize the client
	client := gojpcal.NewConnpassClient(apiKey)

	// Load configuration
	config := gojpcal.LoadConfig()

	// Set up HTTP handler
	http.HandleFunc("/calendar.ics", func(w http.ResponseWriter, r *http.Request) {
		// Get current year and month
		now := time.Now()

		// Generate a range of months from 3 months before to 3 months after
		var yearMonths []string
		for i := -3; i <= 3; i++ {
			targetDate := now.AddDate(0, i, 0)
			yearMonth := fmt.Sprintf("%d%02d", targetDate.Year(), int(targetDate.Month()))
			yearMonths = append(yearMonths, yearMonth)
		}

		// Fetch events
		events, err := client.FetchEvents(config.ConnpassGroups, yearMonths)
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

		// Set content type and headers
		w.Header().Set("Content-Type", "text/calendar; charset=utf-8")

		// Write the calendar data
		fmt.Fprint(w, calendar)
	})

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
