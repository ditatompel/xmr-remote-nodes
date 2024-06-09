package monero

import (
	"testing"
)

func TestQueryProbers_toSQL(t *testing.T) {
	tests := []struct {
		name              string
		query             QueryProbers
		wantArgs          []interface{}
		wantWhere         string
		wantSortBy        string
		wantSortDirection string
	}{
		// TODO: Add test cases.
		{
			name: "Default query",
			query: QueryProbers{
				Search:        "",
				SortBy:        "last_submit_ts",
				SortDirection: "desc",
			},
			wantArgs:          []interface{}{},
			wantWhere:         "",
			wantSortBy:        "last_submit_ts",
			wantSortDirection: "DESC",
		},
		{
			name: "With name or api_key query",
			query: QueryProbers{
				Search:        "test",
				SortBy:        "last_submit_ts",
				SortDirection: "desc",
			},
			wantArgs:          []interface{}{"%test%", "%test%"},
			wantWhere:         "WHERE (name LIKE ? OR api_key LIKE ?)",
			wantSortBy:        "last_submit_ts",
			wantSortDirection: "DESC",
		},
		{
			name: "With sort direction",
			query: QueryProbers{
				Search:        "test",
				SortBy:        "last_submit_ts",
				SortDirection: "asc",
			},
			wantArgs:          []interface{}{"%test%", "%test%"},
			wantWhere:         "WHERE (name LIKE ? OR api_key LIKE ?)",
			wantSortBy:        "last_submit_ts",
			wantSortDirection: "ASC",
		},
		{
			name: "With sort by ID",
			query: QueryProbers{
				Search:        "test",
				SortBy:        "id",
				SortDirection: "asc",
			},
			wantArgs:          []interface{}{"%test%", "%test%"},
			wantWhere:         "WHERE (name LIKE ? OR api_key LIKE ?)",
			wantSortBy:        "id",
			wantSortDirection: "ASC",
		},
		{
			name: "With invalid sort by name and direction",
			query: QueryProbers{
				Search:        "test",
				SortBy:        "invalid",
				SortDirection: "invalid",
			},
			wantArgs:          []interface{}{"%test%", "%test%"},
			wantWhere:         "WHERE (name LIKE ? OR api_key LIKE ?)",
			wantSortBy:        "last_submit_ts",
			wantSortDirection: "DESC",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := QueryProbers{
				Search:        tt.query.Search,
				SortBy:        tt.query.SortBy,
				SortDirection: tt.query.SortDirection,
			}
			gotArgs, gotWhere, gotSortBy, gotSortDirection := q.toSQL()
			if !equalArgs(gotArgs, tt.wantArgs) {
				t.Errorf("QueryNodes.toSQL() gotArgs = %v, want %v", gotArgs, tt.wantArgs)
			}
			if gotWhere != tt.wantWhere {
				t.Errorf("QueryProbers.toSQL() gotWhere = %v, want %v", gotWhere, tt.wantWhere)
			}
			if gotSortBy != tt.wantSortBy {
				t.Errorf("QueryProbers.toSQL() gotSortBy = %v, want %v", gotSortBy, tt.wantSortBy)
			}
			if gotSortDirection != tt.wantSortDirection {
				t.Errorf("QueryProbers.toSQL() gotSortDirection = %v, want %v", gotSortDirection, tt.wantSortDirection)
			}
		})
	}
}

// TODO: Add database test table and then clean it up

func TestProberRepo_CheckApi(t *testing.T) {
	if !testMySQL {
		t.Skip("Skip integration test, not connected to database")
	}
	tests := []struct {
		name    string
		apiKey  string
		want    Prober
		wantErr bool
	}{
		{
			name:    "Empty key",
			apiKey:  "",
			want:    Prober{},
			wantErr: true,
		},
		{
			name:    "Invalid key",
			apiKey:  "invalid",
			want:    Prober{},
			wantErr: true,
		},
	}

	repo := NewProber()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := repo.CheckAPI(tt.apiKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProberRepo.CheckApi() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func BenchmarkProberRepo_CheckApi(b *testing.B) {
	if !testMySQL {
		b.Skip("Skip bench, not connected to database")
	}
	repo := NewProber()
	for i := 0; i < b.N; i++ {
		repo.CheckAPI("")
	}
}
