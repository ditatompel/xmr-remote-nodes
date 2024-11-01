package views

import (
	"fmt"
	"strconv"
	"time"
)

// Convert the float to a string, trimming unnecessary zeros
func formatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

// TimeSince converts an int64 timestamp to a relative time string
func timeSince(timestamp int64) string {
	var duration time.Duration
	var suffix string

	t := time.Unix(timestamp, 0)

	if t.After(time.Now()) {
		duration = time.Until(t)
		suffix = "from now"
	} else {
		duration = time.Since(t)
		suffix = "ago"
	}

	switch {
	case duration < time.Minute:
		return fmt.Sprintf("%ds %s", int(duration.Seconds()), suffix)
	case duration < time.Hour:
		return fmt.Sprintf("%dm %s", int(duration.Minutes()), suffix)
	case duration < time.Hour*24:
		return fmt.Sprintf("%dh %s", int(duration.Hours()), suffix)
	case duration < time.Hour*24*7:
		return fmt.Sprintf("%dd %s", int(duration.Hours()/24), suffix)
	case duration < time.Hour*24*30:
		return fmt.Sprintf("%dw %s", int(duration.Hours()/(24*7)), suffix)
	default:
		months := int(duration.Hours() / (24 * 30))
		if months == 1 {
			return fmt.Sprintf("1 month %s", suffix)
		}
		return fmt.Sprintf("%d months %s", months, suffix)
	}
}
