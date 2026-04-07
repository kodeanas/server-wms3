package utils

import "time"

// ParseDate parses input date value into time.Time with layout.
func ParseDate(value, layout string) (time.Time, error) {
	return time.Parse(layout, value)
}

// FormatDate formats time.Time with layout.
func FormatDate(t time.Time, layout string) string {
	return t.Format(layout)
}

// Today returns start of current local day.
func Today() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}

// NowUnixNano returns current unix nano timestamp
func NowUnixNano() int64 {
	return time.Now().UnixNano()
}

// StartOfDay returns start-of-day for specified time.
func StartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// EndOfDay returns end-of-day for specified time.
func EndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), t.Location())
}

// AddDays adds days to t.
func AddDays(t time.Time, days int) time.Time {
	return t.AddDate(0, 0, days)
}

// DurationDays returns number of full days between a and b.
func DurationDays(a, b time.Time) int {
	d := b.Sub(a)
	return int(d.Hours() / 24)
}
