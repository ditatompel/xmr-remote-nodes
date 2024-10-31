package paging

import (
	"reflect"

	"github.com/google/go-querystring/query"
)

type Paging struct {
	Limit   int    `url:"limit,omitempty"` // rows per page
	Page    int    `url:"page"`
	SortBy  string `url:"sort_by,omitempty"`
	SortDir string `url:"sort_dir,omitempty"`

	SortDirection string `url:"sort_direction,omitempty"` // DEPRECATED: use SortDir

	// Refresh interval
	Refresh int `url:"refresh,omitempty"`
}

// a-h templ helpers
func EncodedQuery(q interface{}, exclude interface{}) string {
	arr := reflect.ValueOf(exclude)
	v, _ := query.Values(q)

	for i := 0; i < arr.Len(); i++ {
		v.Del(arr.Index(i).String())
	}

	return v.Encode()
}

type Pagination struct {
	CurrentPage int
	TotalPages  int
	Pages       []int
}

func NewPagination(currentPage, totalPages int) Pagination {
	var pages []int
	const maxButtons = 5

	if totalPages <= maxButtons {
		for i := 1; i <= totalPages; i++ {
			pages = append(pages, i)
		}
	} else {
		start := max(1, currentPage-2)
		end := min(totalPages, currentPage+2)

		if currentPage <= 3 {
			end = maxButtons
		} else if currentPage > totalPages-3 {
			start = totalPages - (maxButtons - 1)
		}

		for i := start; i <= end; i++ {
			pages = append(pages, i)
		}
		if start > 1 {
			pages = append([]int{1, -1}, pages...) // -1 indicates ellipsis
		}
		if end < totalPages {
			pages = append(pages, -1, totalPages) // -1 indicates ellipsis
		}
	}

	return Pagination{
		CurrentPage: currentPage,
		TotalPages:  totalPages,
		Pages:       pages,
	}
}
