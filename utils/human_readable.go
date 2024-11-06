package utils

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

// TimeSince converts an int64 timestamp to a relative time string
func TimeSince(timestamp int64) string {
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

// Convert the float to a string, trimming unnecessary zeros
func FormatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

// Formats bytes as a human-readable string with the specified number of decimal places.
func FormatBytes(bytes, decimals int) string {
	if bytes == 0 {
		return "0 Bytes"
	}

	const k float64 = 1024
	sizes := []string{"Bytes", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB", "ZiB", "YiB"}

	i := int(math.Floor(math.Log(float64(bytes)) / math.Log(k)))
	dm := decimals
	if dm < 0 {
		dm = 0
	}

	value := float64(bytes) / math.Pow(k, float64(i))
	return fmt.Sprintf("%.*f %s", dm, value, sizes[i])
}

// Formats a hash value (h) into human readable format.
//
// This function was adapted from jtgrassie/monero-pool project.
// Source: https://github.com/jtgrassie/monero-pool/blob/master/src/webui-embed.html
//
// Copyright (c) 2018, The Monero Project
func FormatHashes(h float64) string {
	switch {
	case h < 1e-12:
		return "0 H"
	case h < 1e-9:
		return fmt.Sprintf("%.0f pH", maxPrecision(h*1e12, 0))
	case h < 1e-6:
		return fmt.Sprintf("%.0f nH", maxPrecision(h*1e9, 0))
	case h < 1e-3:
		return fmt.Sprintf("%.0f μH", maxPrecision(h*1e6, 0))
	case h < 1:
		return fmt.Sprintf("%.0f mH", maxPrecision(h*1e3, 0))
	case h < 1e3:
		return fmt.Sprintf("%.0f H", h)
	case h < 1e6:
		return fmt.Sprintf("%.2f KH", maxPrecision(h*1e-3, 2))
	case h < 1e9:
		return fmt.Sprintf("%.2f MH", maxPrecision(h*1e-6, 2))
	default:
		return fmt.Sprintf("%.2f GH", maxPrecision(h*1e-9, 2))
	}
}

// Returns a number with a maximum precision.
//
// This function was adapted from jtgrassie/monero-pool project.
// Source: https://github.com/jtgrassie/monero-pool/blob/master/src/webui-embed.html
//
// Copyright (c) 2018, The Monero Project
func maxPrecision(n float64, p int) float64 {
	format := "%." + strconv.Itoa(p) + "f"
	result, _ := strconv.ParseFloat(fmt.Sprintf(format, n), 64)
	return result
}
