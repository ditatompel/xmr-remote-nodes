package repo

import (
	"fmt"
	"slices"
	"strings"

	"github.com/ditatompel/xmr-nodes/internal/database"

	"github.com/google/uuid"
)

type ProberRepository interface {
	AddProber(name string) error
	Probers(q ProbersQueryParams) (Probers, error)
}

type ProberRepo struct {
	db *database.DB
}

type Prober struct {
	Id           int64     `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	ApiKey       uuid.UUID `json:"api_key" db:"api_key"`
	LastSubmitTs int64     `json:"last_submit_ts" db:"last_submit_ts"`
}

type ProbersQueryParams struct {
	Name   string
	ApiKey string

	RowsPerPage   int
	Page          int
	SortBy        string
	SortDirection string
}

type Probers struct {
	TotalRows   int       `json:"total_rows"`
	RowsPerPage int       `json:"rows_per_page"`
	CurrentPage int       `json:"current_page"`
	NextPage    int       `json:"next_page"`
	Items       []*Prober `json:"items"`
}

func NewProberRepo(db *database.DB) ProberRepository {
	return &ProberRepo{db}
}

func (repo *ProberRepo) AddProber(name string) error {
	query := `INSERT INTO tbl_prober (name, api_key, last_submit_ts) VALUES (?, ?, ?)`
	_, err := repo.db.Exec(query, name, uuid.New(), 0)
	if err != nil {
		return err
	}
	return nil
}

func (repo *ProberRepo) Probers(q ProbersQueryParams) (Probers, error) {
	queryParams := []interface{}{}
	whereQueries := []string{}
	where := ""

	if q.Name != "" {
		whereQueries = append(whereQueries, "name LIKE ?")
		queryParams = append(queryParams, "%"+q.Name+"%")
	}
	if q.ApiKey != "" {
		whereQueries = append(whereQueries, "api_key LIKE ?")
		queryParams = append(queryParams, "%"+q.ApiKey+"%")
	}

	if len(whereQueries) > 0 {
		where = "WHERE " + strings.Join(whereQueries, " AND ")
	}

	probers := Probers{}
	queryTotalRows := fmt.Sprintf("SELECT COUNT(id) AS total_rows FROM tbl_prober %s", where)

	err := repo.db.QueryRow(queryTotalRows, queryParams...).Scan(&probers.TotalRows)
	if err != nil {
		return probers, err
	}
	queryParams = append(queryParams, q.RowsPerPage, (q.Page-1)*q.RowsPerPage)

	allowedSort := []string{"id", "last_submit_ts"}
	sortBy := "id"
	if slices.Contains(allowedSort, q.SortBy) {
		sortBy = q.SortBy
	}
	sortDirection := "DESC"
	if q.SortDirection == "asc" {
		sortDirection = "ASC"
	}

	query := fmt.Sprintf("SELECT id, name, api_key, last_submit_ts FROM tbl_prober %s ORDER BY %s %s LIMIT ? OFFSET ?", where, sortBy, sortDirection)

	row, err := repo.db.Query(query, queryParams...)
	if err != nil {
		return probers, err
	}
	defer row.Close()

	probers.RowsPerPage = q.RowsPerPage
	probers.CurrentPage = q.Page
	probers.NextPage = q.Page + 1

	for row.Next() {
		prober := Prober{}
		err = row.Scan(&prober.Id, &prober.Name, &prober.ApiKey, &prober.LastSubmitTs)
		if err != nil {
			return probers, err
		}
		probers.Items = append(probers.Items, &prober)
	}
	return probers, nil
}
