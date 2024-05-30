package monero

import (
	"os"
	"reflect"
	"strconv"
	"testing"
	"xmr-remote-nodes/internal/config"
	"xmr-remote-nodes/internal/database"
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

func TestQueryNodes_toSQL(t *testing.T) {
	tests := []struct {
		name              string
		query             QueryNodes
		wantArgs          []interface{}
		wantWhere         string
		wantSortBy        string
		wantSortDirection string
	}{
		{
			name: "Default query",
			query: QueryNodes{
				Host:          "",
				Nettype:       "any",
				Protocol:      "any",
				CC:            "any",
				Status:        -1,
				CORS:          -1,
				RowsPerPage:   10,
				Page:          1,
				SortBy:        "last_checked",
				SortDirection: "desc",
			},
			wantArgs:          []interface{}{},
			wantWhere:         "",
			wantSortBy:        "last_checked",
			wantSortDirection: "DESC",
		},
		{
			name: "With host query",
			query: QueryNodes{
				Host:          "test",
				Nettype:       "any",
				Protocol:      "any",
				CC:            "any",
				Status:        -1,
				CORS:          -1,
				RowsPerPage:   10,
				Page:          1,
				SortBy:        "last_checked",
				SortDirection: "desc",
			},
			wantArgs:          []interface{}{"%test%", "%test%"},
			wantWhere:         "WHERE (hostname LIKE ? OR ip_addr LIKE ?)",
			wantSortBy:        "last_checked",
			wantSortDirection: "DESC",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotArgs, gotWhere, gotSortBy, gotSortDirection := tt.query.toSQL()

			if !equalArgs(gotArgs, tt.wantArgs) {
				t.Errorf("QueryNodes.toSQL() gotArgs = %v, want %v", gotArgs, tt.wantArgs)
			}
			if gotWhere != tt.wantWhere {
				t.Errorf("QueryNodes.toSQL() gotWhere = %v, want %v", gotWhere, tt.wantWhere)
			}
			if gotSortBy != tt.wantSortBy {
				t.Errorf("QueryNodes.toSQL() gotSortBy = %v, want %v", gotSortBy, tt.wantSortBy)
			}
			if gotSortDirection != tt.wantSortDirection {
				t.Errorf("QueryNodes.toSQL() gotSortDirection = %v, want %v", gotSortDirection, tt.wantSortDirection)
			}
		})
	}
}

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
