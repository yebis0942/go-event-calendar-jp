package main

import (
	"flag"
	"fmt"
	"os"

	gojpcal "github.com/yebis0942/golang-jp-event-calendar"
)

func main() {
	yyyymm := flag.String("yyyymm", "", "Year and month in format YYYYMM")
	flag.Parse()

	if *yyyymm == "" {
		fmt.Fprintf(os.Stderr, "Error: -yyyymm flag is required\n")
		flag.Usage()
		os.Exit(1)
	}

	if len(*yyyymm) != 6 {
		fmt.Fprintf(os.Stderr, "Error: -yyyymm must be in format YYYYMM\n")
		os.Exit(1)
	}

	apiKey := os.Getenv("CONNPASS_API_KEY")
	if apiKey == "" {
		fmt.Fprintf(os.Stderr, "Error: CONNPASS_API_KEY environment variable is required\n")
		os.Exit(1)
	}

	// Parse year and month
	year := (*yyyymm)[:4]
	month := (*yyyymm)[4:]

	// Load config
	config := gojpcal.LoadConfig()

	// Get events
	events, err := gojpcal.FetchEvents(apiKey, config.ConnpassGroups, year, month)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching events: %v\n", err)
		os.Exit(1)
	}

	// Generate iCalendar
	ical, err := gojpcal.GenerateICalendar(events)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating iCalendar: %v\n", err)
		os.Exit(1)
	}

	// Write to file
	outputFile := fmt.Sprintf("%s.ics", *yyyymm)
	err = os.WriteFile(outputFile, []byte(ical), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Calendar successfully written to %s\n", outputFile)
}
