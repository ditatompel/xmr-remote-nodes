package monero

import (
	"net"
	"net/netip"
	"testing"
)

// Single test:
// go test -race ./internal/monero -run=TestParseIPCIDRToPrefix -v
func TestParseIPCIDRToPrefix(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantOK   bool
		wantPref netip.Prefix
	}{
		{
			name:     "Plain IPv4",
			input:    "192.168.1.1",
			wantOK:   true,
			wantPref: netip.MustParsePrefix("192.168.1.1/32"),
		},
		{
			name:     "IPv4 CIDR",
			input:    "192.168.1.0/24",
			wantOK:   true,
			wantPref: netip.MustParsePrefix("192.168.1.0/24"),
		},
		{
			name:     "Plain IPv6",
			input:    "2001:db8::1",
			wantOK:   true,
			wantPref: netip.MustParsePrefix("2001:db8::1/128"),
		},
		{
			name:     "IPv6 CIDR",
			input:    "2001:db8::/32",
			wantOK:   true,
			wantPref: netip.MustParsePrefix("2001:db8::/32"),
		},
		{
			name:   "Invalid IP",
			input:  "999.1.1.1",
			wantOK: false,
		},
		{
			name:   "Invalid CIDR",
			input:  "192.168.1.1/33",
			wantOK: false,
		},
		{
			name:   "Commented string",
			input:  "# Monero Research Lab spy node ban list",
			wantOK: false,
		},
		{
			name:   "Empty string",
			input:  "",
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := parseIPCIDRToPrefix(tt.input)

			if ok != tt.wantOK {
				t.Fatalf("expected ok=%v, got %v", tt.wantOK, ok)
			}

			if !ok {
				return
			}

			if got != tt.wantPref {
				t.Errorf("expected prefix %v, got %v", tt.wantPref, got)
			}
		})
	}
}

// Single test:
// go test -race ./internal/monero -run=TestIsBannedIP -v
func TestIsBannedIP(t *testing.T) {
	tests := []struct {
		name     string
		banList  []string
		inputIPs []net.IP
		want     bool
	}{
		{
			name:     "Commented string",
			banList:  []string{"# Monero Research Lab spy node ban list"},
			inputIPs: []net.IP{net.ParseIP("192.168.1.123")},
			want:     false,
		},
		{
			name:     "Empty ban list",
			banList:  []string{},
			inputIPs: []net.IP{net.ParseIP("192.168.1.123")},
			want:     false,
		},
		{
			name:     "Exact IP match",
			banList:  []string{"192.168.1.123"},
			inputIPs: []net.IP{net.ParseIP("192.168.1.123")},
			want:     true,
		},
		{
			name:     "IP in CIDR",
			banList:  []string{"10.0.0.0/8"},
			inputIPs: []net.IP{net.ParseIP("10.1.2.3")},
			want:     true,
		},
		{
			name:     "No match",
			banList:  []string{"192.168.1.0/24", "10.0.0.0/8"},
			inputIPs: []net.IP{net.ParseIP("192.168.2.1")},
			want:     false,
		},
		{
			name:    "Multiple IPs, one match",
			banList: []string{"10.0.0.0/8", "172.16.0.0/12"},
			inputIPs: []net.IP{
				net.ParseIP("192.168.1.1"),
				net.ParseIP("8.8.8.8"),
				net.ParseIP("172.16.5.10"),
			},
			want: true,
		},
		{
			name:     "IPv6 match",
			banList:  []string{"2001:db8::/32"},
			inputIPs: []net.IP{net.ParseIP("2001:db8::1")},
			want:     true,
		},
		{
			name:     "IPv6 no match",
			banList:  []string{"2001:db8::/32", "10.0.0.0/8", "172.16.0.0/12", "8.8.8.8"},
			inputIPs: []net.IP{net.ParseIP("2001:dead::1")},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isBannedIP(tt.banList, tt.inputIPs)
			if got != tt.want {
				t.Errorf("isIPBanned() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Single bench test:
// go test ./internal/monero -bench IsBannedIP -benchmem -run=^$ -v
func Benchmark_IsBannedIP(b *testing.B) {
	banList := []string{
		"192.168.0.0/16", "10.0.0.0/8", "172.16.0.0/12",
	}

	inputIPs := []net.IP{
		net.ParseIP("192.168.1.1"),
		net.ParseIP("10.0.0.1"),
		net.ParseIP("172.16.99.99"),
		net.ParseIP("8.8.8.8"),
	}

	for i := 0; i < b.N; i++ {
		_ = isBannedIP(banList, inputIPs)
	}
}
