package monero

import (
	"fmt"
	"slices"
	"strings"

	"github.com/ditatompel/xmr-remote-nodes/internal/database"

	"github.com/google/uuid"
)

const ProberAPIKey = "X-Prober-Api-Key" // HTTP header key

type proberRepo struct {
	db *database.DB
}

type Prober struct {
	ID           int64     `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	APIKey       uuid.UUID `json:"api_key" db:"api_key"`
	LastSubmitTS int64     `json:"last_submit_ts" db:"last_submit_ts"`
}

// Initializes a new ProberRepository
//
// NOTE: This "prober" is different with "probe" which is used to fetch a new job
func NewProber() *proberRepo {
	return &proberRepo{db: database.GetDB()}
}

// Add a new prober machine
func (r *proberRepo) Add(name string) (Prober, error) {
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
	_, err := r.db.Exec(query, name, apiKey, 0)
	if err != nil {
		return Prober{}, err
	}
	return Prober{Name: name, APIKey: apiKey}, nil
}

// Edit an existing prober
func (r *proberRepo) Edit(id int, name string) error {
	query := `UPDATE tbl_prober SET name = ? WHERE id = ?`
	res, err := r.db.Exec(query, name, id)
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

// Delete an existing prober
func (r *proberRepo) Delete(id int) error {
	res, err := r.db.Exec(`DELETE FROM tbl_prober WHERE id = ?`, id)
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

type QueryProbers struct {
	Search        string
	SortBy        string
	SortDirection string
}

func (q QueryProbers) toSQL() (args []interface{}, where, sortBy, sortDirection string) {
	wq := []string{}
	if q.Search != "" {
		wq = append(wq, "(name LIKE ? OR api_key LIKE ?)")
		args = append(args, "%"+q.Search+"%", "%"+q.Search+"%")
	}
	if len(wq) > 0 {
		where = "WHERE " + strings.Join(wq, " AND ")
	}

	as := []string{"id", "last_submit_ts"}
	sortBy = "last_submit_ts"
	if slices.Contains(as, q.SortBy) {
		sortBy = q.SortBy
	}
	sortDirection = "DESC"
	if q.SortDirection == "asc" {
		sortDirection = "ASC"
	}

	return args, where, sortBy, sortDirection
}

func (r *proberRepo) Probers(q QueryProbers) ([]Prober, error) {
	args, where, sortBy, sortDirection := q.toSQL()

	var probers []Prober

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

	row, err := r.db.Query(query, args...)
	if err != nil {
		return probers, err
	}
	defer row.Close()

	for row.Next() {
		var p Prober
		err = row.Scan(&p.ID, &p.Name, &p.APIKey, &p.LastSubmitTS)
		if err != nil {
			return probers, err
		}
		probers = append(probers, p)
	}
	return probers, nil
}

func (r *proberRepo) CheckAPI(key string) (Prober, error) {
	var p Prober
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
	err := r.db.QueryRow(query, key).Scan(&p.ID, &p.Name, &p.APIKey, &p.LastSubmitTS)
	return p, err
}
