// Package ip provides IP address related functions
package ip

import (
	"net"
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
