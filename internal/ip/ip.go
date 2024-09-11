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
