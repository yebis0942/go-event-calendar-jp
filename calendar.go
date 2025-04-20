package gojpcal

import (
	"fmt"
	"strings"
	"time"
)

const (
	prodID      = "github.com/yebis0942/golang-jp-event-calendar"
	eventDomain = "golang-jp-event-calendar.yebis0942.workers.dev"
)

// GenerateICalendar generates an iCalendar (RFC 5545) format string from events
func GenerateICalendar(events []Event, now time.Time) (string, error) {
	var sb strings.Builder

	// iCalendar header
	sb.WriteString("BEGIN:VCALENDAR\r\n")
	sb.WriteString("VERSION:2.0\r\n")
	sb.WriteString(fmt.Sprintf("PRODID:%s\r\n", prodID))
	sb.WriteString("CALSCALE:GREGORIAN\r\n")
	sb.WriteString("METHOD:PUBLISH\r\n")
	sb.WriteString("X-WR-CALNAME;VALUE=TEXT:Goコミュニティのイベント\r\n")

	// Add each event
	for _, event := range events {
		sb.WriteString("BEGIN:VEVENT\r\n")

		// Generate a UID for the event
		uid := fmt.Sprintf("connpass-%d@%s", event.ID, eventDomain)
		sb.WriteString(fmt.Sprintf("UID:%s\r\n", uid))

		// Format dates according to iCalendar spec
		startStr := formatDateTimeUTC(event.StartTime)
		endStr := formatDateTimeUTC(event.EndTime)

		sb.WriteString(fmt.Sprintf("DTSTART:%s\r\n", startStr))
		sb.WriteString(fmt.Sprintf("DTEND:%s\r\n", endStr))

		// Add creation timestamp
		nowUTC := now.UTC()
		sb.WriteString(fmt.Sprintf("DTSTAMP:%s\r\n", formatDateTimeUTC(nowUTC)))

		// Event details
		sb.WriteString(fmt.Sprintf("SUMMARY:%s\r\n", escape(event.Title)))
		sb.WriteString(fmt.Sprintf("DESCRIPTION:%s\\n\\nURL: %s\r\n", escape(event.Description), escape(event.URL)))

		if event.Place != "" {
			location := event.Place
			if event.Address != "" {
				location += ", " + event.Address
			}
			sb.WriteString(fmt.Sprintf("LOCATION:%s\r\n", escape(location)))
		}

		sb.WriteString(fmt.Sprintf("URL:%s\r\n", escape(event.URL)))

		sb.WriteString("END:VEVENT\r\n")
	}

	sb.WriteString("END:VCALENDAR\r\n")
	return sb.String(), nil
}

// formatDateTimeUTC formats a time.Time as iCalendar format in UTC
func formatDateTimeUTC(t time.Time) string {
	return t.UTC().Format("20060102T150405Z")
}

// escape handles special characters in iCalendar text fields
func escape(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, ";", "\\;")
	s = strings.ReplaceAll(s, ",", "\\,")
	s = strings.ReplaceAll(s, "\n", "\\n")
	return s
}
