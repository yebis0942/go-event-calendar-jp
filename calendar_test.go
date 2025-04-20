package gojpcal_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	gojpcal "github.com/yebis0942/golang-jp-event-calendar"
)

func TestGenerateICalendar(t *testing.T) {
	events := []gojpcal.Event{
		{
			ID:          1,
			Title:       "Sample Event",
			Description: "This is a sample event.",
			StartTime:   time.Date(2023, 10, 1, 10, 0, 0, 0, time.UTC),
			EndTime:     time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC),
			URL:         "http://example.com",
			Place:       "Sample Place",
			Address:     "123 Sample Street",
		},
	}

	now := time.Date(2023, 10, 1, 9, 0, 0, 0, time.UTC)
	got, err := gojpcal.GenerateICalendar(events, now)
	require.NoError(t, err)

	want := "BEGIN:VCALENDAR\r\n" +
		"VERSION:2.0\r\n" +
		"PRODID:github.com/yebis0942/golang-jp-event-calendar\r\n" +
		"CALSCALE:GREGORIAN\r\n" +
		"METHOD:PUBLISH\r\n" +
		"X-WR-CALNAME;VALUE=TEXT:Goコミュニティのイベント\r\n" +
		"BEGIN:VEVENT\r\n" +
		"UID:connpass-1@golang-jp-event-calendar.yebis0942.workers.dev\r\n" +
		"DTSTART:20231001T100000Z\r\n" +
		"DTEND:20231001T120000Z\r\n" +
		"DTSTAMP:20231001T090000Z\r\n" +
		"SUMMARY:Sample Event\r\n" +
		"DESCRIPTION:This is a sample event.\\n\\nURL: http://example.com\r\n" +
		"LOCATION:Sample Place\\, 123 Sample Street\r\n" +
		"URL:http://example.com\r\n" +
		"END:VEVENT\r\n" +
		"END:VCALENDAR\r\n"
	require.Equal(t, want, got)
}
