package monero

import (
	"fmt"
	"slices"
	"strings"
	"xmr-remote-nodes/internal/database"

	"github.com/google/uuid"
)

type ProberRepository interface {
	Add(name string) (Prober, error)
	Edit(id int, name string) error
	Probers(q ProbersQueryParams) ([]Prober, error)
	CheckApi(key string) (Prober, error)
	Delete(id int) error
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

func NewProberRepo(db *database.DB) ProberRepository {
	return &ProberRepo{db}
}

func (repo *ProberRepo) Add(name string) (Prober, error) {
	apiKey := uuid.New()
	query := `
		INSERT INTO tbl_prober (
			name,
			api_key,
			last_submit_ts
		) VALUES (
			?,
			?,
			?
		)`
	_, err := repo.db.Exec(query, name, apiKey, 0)
	if err != nil {
		return Prober{}, err
	}
	return Prober{Name: name, ApiKey: apiKey}, nil
}

func (repo *ProberRepo) Edit(id int, name string) error {
	query := `UPDATE tbl_prober SET name = ? WHERE id = ?`
	res, err := repo.db.Exec(query, name, id)
	if err != nil {
		return err
	}
	row, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if row == 0 {
		return fmt.Errorf("no rows affected")
	}
	return err
}

func (repo *ProberRepo) Delete(id int) error {
	res, err := repo.db.Exec(`DELETE FROM tbl_prober WHERE id = ?`, id)
	if err != nil {
		return err
	}
	row, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if row == 0 {
		return fmt.Errorf("no rows affected")
	}
	return err
}

type ProbersQueryParams struct {
	Search        string
	SortBy        string
	SortDirection string
}

func (repo *ProberRepo) Probers(q ProbersQueryParams) ([]Prober, error) {
	queryParams := []interface{}{}
	whereQueries := []string{}
	where := ""

	if q.Search != "" {
		whereQueries = append(whereQueries, "(name LIKE ? OR api_key LIKE ?)")
		queryParams = append(queryParams, "%"+q.Search+"%")
		queryParams = append(queryParams, "%"+q.Search+"%")
	}

	if len(whereQueries) > 0 {
		where = "WHERE " + strings.Join(whereQueries, " AND ")
	}

	probers := []Prober{}

	allowedSort := []string{"id", "last_submit_ts"}
	sortBy := "last_submit_ts"
	if slices.Contains(allowedSort, q.SortBy) {
		sortBy = q.SortBy
	}
	sortDirection := "DESC"
	if q.SortDirection == "asc" {
		sortDirection = "ASC"
	}

	query := fmt.Sprintf(`
		SELECT
			id,
			name,
			api_key,
			last_submit_ts
		FROM
			tbl_prober
		%s -- where clause if any
		ORDER BY %s %s`, where, sortBy, sortDirection)

	row, err := repo.db.Query(query, queryParams...)
	if err != nil {
		return probers, err
	}
	defer row.Close()

	for row.Next() {
		prober := Prober{}
		err = row.Scan(&prober.Id, &prober.Name, &prober.ApiKey, &prober.LastSubmitTs)
		if err != nil {
			return probers, err
		}
		probers = append(probers, prober)
	}
	return probers, nil
}

func (repo *ProberRepo) CheckApi(key string) (Prober, error) {
	prober := Prober{}
	query := `
		SELECT
			id,
			name,
			api_key,
			last_submit_ts
		FROM
			tbl_prober
		WHERE
			api_key = ?
		LIMIT 1`
	err := repo.db.QueryRow(query, key).Scan(&prober.Id, &prober.Name, &prober.ApiKey, &prober.LastSubmitTs)
	return prober, err
}
