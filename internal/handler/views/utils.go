package views

import "strconv"

// Convert the float to a string, trimming unnecessary zeros
func formatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}
