// Copyright (c) 2026 Karel Hanák
// SPDX-License-Identifier: MIT

package common

import (
	"fmt"
	"time"
)

// GetThisWeek returns the start and end time of the current week (Monday 00:00:00 to Sunday 23:59:59) in UTC.
func GetThisWeek() (start, end time.Time) {
	now := time.Now().UTC()
	daysToSubtract := int(now.Weekday()) - 1
	if daysToSubtract < 0 {
		daysToSubtract = 6
	}
	start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, -daysToSubtract)
	end = start.AddDate(0, 0, 7).Add(-time.Nanosecond)
	return start, end
}

// GetThisMonth returns the start and end time of the current month in UTC.
func GetThisMonth() (start, end time.Time) {
	now := time.Now().UTC()
	start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	end = start.AddDate(0, 1, 0).Add(-time.Nanosecond)
	return start, end
}

func FormatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}
