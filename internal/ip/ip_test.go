package ip

import (
	"net"
	"testing"
)

// Single test: go test ./internal/ip -bench TestIsIPv6Only -benchmem -run=^$ -v
func TestIsIPv6Only(t *testing.T) {
	tests := []struct {
		name string
		ips  []net.IP
		want bool
	}{
		{
			name: "IPv4",
			ips: []net.IP{
				net.ParseIP("1.1.1.1"),
			},
			want: false,
		},
		{
			name: "IPv6",
			ips: []net.IP{
				net.ParseIP("2606:4700::6810:85e5"),
			},
			want: true,
		},
		{
			name: "IPv6 and IPv4",
			ips: []net.IP{
				net.ParseIP("1.1.1.1"),
				net.ParseIP("2606:4700::6810:84e5"),
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsIPv6Only(tt.ips); got != tt.want {
				t.Errorf("IsIPv6Only() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Single test: go test ./internal/ip -bench TestSliceToString -benchmem -run=^$ -v
func TestSliceToString(t *testing.T) {
	tests := []struct {
		name string
		ips  []net.IP
		want string
	}{
		{
			name: "IPv4",
			ips: []net.IP{
				net.ParseIP("1.1.1.1"),
			},
			want: "1.1.1.1",
		},
		{
			name: "IPv6",
			ips: []net.IP{
				net.ParseIP("2606:4700::6810:85e5"),
			},
			want: "2606:4700::6810:85e5",
		},
		{
			name: "IPv6 and IPv4",
			ips: []net.IP{
				net.ParseIP("1.1.1.1"),
				net.ParseIP("2606:4700::6810:85e5"),
			},
			want: "1.1.1.1,2606:4700::6810:85e5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SliceToString(tt.ips); got != tt.want {
				t.Errorf("SliceToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
