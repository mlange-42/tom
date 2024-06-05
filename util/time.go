package util

import "time"

func WithLocation(t time.Time, loc *time.Location) time.Time {
	return time.Date(
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second(),
		0, loc,
	)
}
