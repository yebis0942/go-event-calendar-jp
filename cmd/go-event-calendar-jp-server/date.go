package main

import "time"

// GetMonthsRange returns an array of year-month strings in "YYYYMM" format for 3 months before and after the specified year and month
func GetMonthsRange(year int, month int) []string {
	// Use the 1st day of the specified year and month as the base date
	baseDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)

	// 3 months before + current month + 3 months after = 7 months
	result := make([]string, 7)

	// From 3 months before to 3 months after
	for i := -3; i <= 3; i++ {
		targetDate := baseDate.AddDate(0, i, 0)
		result[i+3] = targetDate.Format("200601")
	}

	return result
}
