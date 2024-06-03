package monero

import (
	"testing"

	"github.com/jmoiron/sqlx/types"
)

func TestNode_parseStatuses(t *testing.T) {
	tests := []struct {
		name   string
		report Node
		want   string
	}{
		{
			name: "Invalid initial LastCheckStatus type",
			report: Node{
				IsAvailable:     true,
				LastCheckStatus: types.JSONText("-invalid source-"),
			},
			want: "[2,2,2,2,1]",
		},
		{
			name: "Node goes online",
			report: Node{
				IsAvailable:     true,
				LastCheckStatus: types.JSONText("[0,1,0,0,0]"),
			},
			want: "[1,0,0,0,1]",
		},
		{
			name: "Node goes offline",
			report: Node{
				IsAvailable:     false,
				LastCheckStatus: types.JSONText("[0,1,0,1,1]"),
			},
			want: "[1,0,1,1,0]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &ProbeReport{
				Node: Node{
					IsAvailable:     tt.report.IsAvailable,
					LastCheckStatus: tt.report.LastCheckStatus,
				},
			}
			if got := n.parseStatuses(); got != tt.want {
				t.Errorf("Node.parseStatuses() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueryLogs_toSQL(t *testing.T) {
	tests := []struct {
		name              string
		fields            QueryLogs
		wantArgs          []interface{}
		wantWhere         string
		wantSortBy        string
		wantSortDirection string
	}{
		{
			name: "Default query",
			fields: QueryLogs{
				NodeID:        0,
				Status:        -1,
				FailedReason:  "",
				RowsPerPage:   10,
				Page:          1,
				SortBy:        "date_checked",
				SortDirection: "desc",
			},
			wantArgs:          []interface{}{},
			wantWhere:         "",
			wantSortBy:        "date_checked",
			wantSortDirection: "DESC",
		},
		{
			name: "With node_id query",
			fields: QueryLogs{
				NodeID:        1,
				Status:        -1,
				FailedReason:  "",
				RowsPerPage:   10,
				Page:          1,
				SortBy:        "date_checked",
				SortDirection: "desc",
			},
			wantArgs:          []interface{}{1},
			wantWhere:         "WHERE node_id = ?",
			wantSortBy:        "date_checked",
			wantSortDirection: "DESC",
		},
		{
			name: "All possible query",
			fields: QueryLogs{
				NodeID:        1,
				Status:        0,
				FailedReason:  "test",
				RowsPerPage:   10,
				Page:          1,
				SortBy:        "date_checked",
				SortDirection: "asc",
			},
			wantArgs:          []interface{}{1, 0, "%test%"},
			wantWhere:         "WHERE node_id = ? AND is_available = ? AND failed_reason LIKE ?",
			wantSortBy:        "date_checked",
			wantSortDirection: "ASC",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := QueryLogs{
				NodeID:        tt.fields.NodeID,
				Status:        tt.fields.Status,
				FailedReason:  tt.fields.FailedReason,
				RowsPerPage:   tt.fields.RowsPerPage,
				Page:          tt.fields.Page,
				SortBy:        tt.fields.SortBy,
				SortDirection: tt.fields.SortDirection,
			}
			gotArgs, gotWhere, gotSortBy, gotSortDirection := q.toSQL()
			if !equalArgs(gotArgs, tt.wantArgs) {
				t.Errorf("QueryNodes.toSQL() gotArgs = %v, want %v", gotArgs, tt.wantArgs)
			}
			if gotWhere != tt.wantWhere {
				t.Errorf("QueryLogs.toSQL() gotWhere = %v, want %v", gotWhere, tt.wantWhere)
			}
			if gotSortBy != tt.wantSortBy {
				t.Errorf("QueryLogs.toSQL() gotSortBy = %v, want %v", gotSortBy, tt.wantSortBy)
			}
			if gotSortDirection != tt.wantSortDirection {
				t.Errorf("QueryLogs.toSQL() gotSortDirection = %v, want %v", gotSortDirection, tt.wantSortDirection)
			}
		})
	}
}
