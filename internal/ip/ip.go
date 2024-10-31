// Package ip provides IP address related functions
package ip

import (
	"net"
	"strings"
)

// IsIPv6Only returns true if all given IPs are IPv6
func IsIPv6Only(ips []net.IP) bool {
	for _, ip := range ips {
		if ip.To4() != nil {
			return false
		}
	}
	return true
}

// SliceToString converts []net.IP to a string separated by comma.
// If the separator is empty, it defaults to ",".
func SliceToString(ips []net.IP) string {
	r := make([]string, len(ips))
	for i, j := range ips {
		r[i] = j.String()
	}

	return strings.Join(r, ",")
}

// Add brackets based on whether the given string is IPv6 or not.
// If the input is an IPv6 address, wraps it in square brackets `[ ]`.
// Otherwise, it returns the input string as-is (for domain names or IPv4
// addresses).
func FormatHostname(hostname string) string {
	ip := net.ParseIP(hostname)
	if ip != nil && ip.To4() == nil {
		return "[" + hostname + "]"
	}

	return hostname
}
