package monero

import (
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/ditatompel/xmr-remote-nodes/internal/config"
	"github.com/ditatompel/xmr-remote-nodes/internal/database"
	"github.com/ditatompel/xmr-remote-nodes/internal/paging"
)

var testMySQL = true

// TODO: Add database test table and then clean it up
func init() {
	// load test db config from OS environment variable
	//
	// Example:
	// TEST_DB_HOST=127.0.0.1 \
	// TEST_DB_PORT=3306 \
	// TEST_DB_USER=testuser \
	// TEST_DB_PASSWORD=testpass \
	// TEST_DB_NAME=testdb go test ./... -v
	//
	// To run benchmark only, add `-bench=. -run=^#` to the `go test` command
	config.DBCfg().Host = os.Getenv("TEST_DB_HOST")
	config.DBCfg().Port, _ = strconv.Atoi(os.Getenv("TEST_DB_PORT"))
	config.DBCfg().User = os.Getenv("TEST_DB_USER")
	config.DBCfg().Password = os.Getenv("TEST_DB_PASSWORD")
	config.DBCfg().Name = os.Getenv("TEST_DB_NAME")

	if err := database.ConnectDB(); err != nil {
		testMySQL = false
	}
}

// Single test:
// go test -race ./internal/monero -run=TestQueryNodes_toSQL -v
func TestQueryNodes_toSQL(t *testing.T) {
	tests := []struct {
		name        string
		query       QueryNodes
		wantArgs    []interface{}
		wantWhere   string
		wantSortBy  string
		wantSortDir string
	}{
		{
			name: "Default query",
			query: QueryNodes{
				Paging: paging.Paging{
					Limit:         10,
					Page:          1,
					SortBy:        "last_checked",
					SortDirection: "desc",
				},
				Host:     "",
				Nettype:  "any",
				Protocol: "any",
				CC:       "any",
				Status:   -1,
				CORS:     "",
			},
			wantArgs:    []interface{}{},
			wantWhere:   "",
			wantSortBy:  "last_checked",
			wantSortDir: "DESC",
		},
		{
			name: "With host query",
			query: QueryNodes{
				Paging: paging.Paging{
					Limit:         10,
					Page:          1,
					SortBy:        "last_checked",
					SortDirection: "desc",
				},
				Host:     "test",
				Nettype:  "any",
				Protocol: "any",
				CC:       "any",
				Status:   -1,
				CORS:     "",
			},
			wantArgs:    []interface{}{"%test%", "%test%"},
			wantWhere:   "WHERE (hostname LIKE ? OR ip_addr LIKE ?)",
			wantSortBy:  "last_checked",
			wantSortDir: "DESC",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotArgs, gotWhere := tt.query.toSQL()
			if !equalArgs(gotArgs, tt.wantArgs) {
				t.Errorf("QueryNodes.toSQL() gotArgs = %v, want %v", gotArgs, tt.wantArgs)
			}
			if gotWhere != tt.wantWhere {
				t.Errorf("QueryNodes.toSQL() gotWhere = %v, want %v", gotWhere, tt.wantWhere)
			}
			if tt.query.SortBy != tt.wantSortBy {
				t.Errorf("QueryNodes.toSQL() gotSortBy = %v, want %v", tt.query.SortBy, tt.wantSortBy)
			}
			if tt.query.SortDirection != tt.wantSortDir {
				t.Errorf("QueryNodes.toSQL() gotSortDir = %v, want %v", tt.query.SortDirection, tt.wantSortDir)
			}
		})
	}
}

// Single bench test:
// go test ./internal/monero -bench QueryNodes_toSQL -benchmem -run=^$ -v
func Benchmark_QueryNodes_toSQL(b *testing.B) {
	q := QueryNodes{
		Paging: paging.Paging{
			Limit:         10,
			Page:          1,
			SortBy:        "last_checked",
			SortDirection: "desc",
		},
		Host:     "test",
		Nettype:  "any",
		Protocol: "any",
		CC:       "any",
		Status:   -1,
		CORS:     "",
	}
	for i := 0; i < b.N; i++ {
		_, _ = q.toSQL()
	}
}

// Single test:
// go test -race ./internal/monero -run=TestValidTorHostname -v
func TestValidTorHostname(t *testing.T) {
	tests := []struct {
		name      string
		host      string
		wantValid bool
	}{
		{
			name:      "Empty host",
			host:      "",
			wantValid: false,
		},
		{
			name:      "Valid tor host",
			host:      "cakexmrl7bonq7ovjka5kuwuyd3f7qnkz6z6s6dmsy3uckwra7bvggyd.onion",
			wantValid: true,
		},
		{
			name:      "Valid tor host with subdomain",
			host:      "just-test.cakexmrl7bonq7ovjka5kuwuyd3f7qnkz6z6s6dmsy3uckwra7bvggyd.onion",
			wantValid: true,
		},
		{
			name:      "Invalid host",
			host:      "test.com",
			wantValid: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if r := validTorHostname(tt.host); r != tt.wantValid {
				t.Errorf("ValidTorHostname() error = %v, wantValid %v", r, tt.wantValid)
			}
		})
	}
}

// Single bench test:
// go test ./internal/monero -bench validTorHostname -benchmem -run=^$ -v
func Benchmark_validTorHostname(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = validTorHostname("cakexmrl7bonq7ovjka5kuwuyd3f7qnkz6z6s6dmsy3uckwra7bvggyd.onion")
	}
}

// Single test:
// go test -race ./internal/monero -run=TestValidI2PHostname -v
// TODO: Validate new format and allow naming service hostnames
func TestValidI2PHostname(t *testing.T) {
	tests := []struct {
		name      string
		host      string
		wantValid bool
	}{
		{
			name:      "Empty host",
			host:      "",
			wantValid: false,
		},
		{
			name:      "Valid b32 i2p host (old format)",
			host:      "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrst234567.b32.i2p",
			wantValid: true,
		},
		{
			name:      "Invalid b32 i2p host (old format)",
			host:      "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrst123456.b32.i2p",
			wantValid: false,
		},
		{
			name:      "clearnet domain",
			host:      "test.com",
			wantValid: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if r := validI2PHostname(tt.host); r != tt.wantValid {
				t.Errorf("ValidTorHostname() error = %v, wantValid %v", r, tt.wantValid)
			}
		})
	}
}

// Single bench test:
// go test ./internal/monero -bench validI2PHostname -benchmem -run=^$ -v
func Benchmark_validI2PHostname(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = validTorHostname("abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrst234567.b32.i2p")
	}
}

// equalArgs is helper function for testing.
//
// This returns true if two slices of interface{} are equal.
func equalArgs(a, b []interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if !reflect.DeepEqual(v, b[i]) {
			return false
		}
	}
	return true
}
