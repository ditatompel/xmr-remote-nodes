package repo

import (
	"fmt"
	"log/slog"
	"math"
	"slices"
	"strings"
	"time"
	"xmr-remote-nodes/internal/database"
)

type CronRepository interface {
	RunCronProcess()
	Crons(q CronQueryParams) (CronTasks, error)
}

type CronRepo struct {
	db *database.DB
}

type Cron struct {
	Id          int     `json:"id" db:"id"`
	Title       string  `json:"title" db:"title"`
	Slug        string  `json:"slug" db:"slug"`
	Description string  `json:"description" db:"description"`
	RunEvery    int     `json:"run_every" db:"run_every"`
	LastRun     int64   `json:"last_run" db:"last_run"`
	NextRun     int64   `json:"next_run" db:"next_run"`
	RunTime     float64 `json:"run_time" db:"run_time"`
	CronState   int     `json:"cron_state" db:"cron_state"`
	IsEnabled   int     `json:"is_enabled" db:"is_enabled"`
}

var rerunTimeout = 300

func NewCron(db *database.DB) CronRepository {
	return &CronRepo{db}
}

func (repo *CronRepo) RunCronProcess() {
	for {
		time.Sleep(60 * time.Second)
		slog.Info("[CRON] Running cron cycle...")
		list, err := repo.queueList()
		if err != nil {
			slog.Warn(fmt.Sprintf("[CRON] Error parsing queue list to struct: %s", err))
			continue
		}
		for _, task := range list {
			startTime := time.Now()
			currentTs := startTime.Unix()
			delayedTask := currentTs - task.NextRun
			if task.CronState == 1 && delayedTask <= int64(rerunTimeout) {
				slog.Debug(fmt.Sprintf("[CRON] Skipping task %s because it is already running", task.Slug))
				continue
			}

			repo.preRunTask(task.Id, currentTs)

			repo.execCron(task.Slug)

			runTime := math.Ceil(time.Since(startTime).Seconds()*1000) / 1000
			slog.Info(fmt.Sprintf("[CRON] Task %s done in %f seconds", task.Slug, runTime))
			nextRun := currentTs + int64(task.RunEvery)

			repo.postRunTask(task.Id, nextRun, runTime)
		}
		slog.Info("[CRON] Cron cycle done!")
	}
}

type CronQueryParams struct {
	Title         string
	Description   string
	IsEnabled     int
	CronState     int
	RowsPerPage   int
	Page          int
	SortBy        string
	SortDirection string
}

type CronTasks struct {
	TotalRows   int     `json:"total_rows"`
	RowsPerPage int     `json:"rows_per_page"`
	Items       []*Cron `json:"items"`
}

func (repo *CronRepo) Crons(q CronQueryParams) (CronTasks, error) {
	queryParams := []interface{}{}
	whereQueries := []string{}
	where := ""

	if q.Title != "" {
		whereQueries = append(whereQueries, "title LIKE ?")
		queryParams = append(queryParams, "%"+q.Title+"%")
	}
	if q.Description != "" {
		whereQueries = append(whereQueries, "description LIKE ?")
		queryParams = append(queryParams, "%"+q.Description+"%")
	}
	if q.IsEnabled != -1 {
		whereQueries = append(whereQueries, "is_enabled = ?")
		queryParams = append(queryParams, q.IsEnabled)
	}
	if q.CronState != -1 {
		whereQueries = append(whereQueries, "cron_state = ?")
		queryParams = append(queryParams, q.CronState)
	}
	if len(whereQueries) > 0 {
		where = "WHERE " + strings.Join(whereQueries, " AND ")
	}
	tasks := CronTasks{}

	queryTotalRows := fmt.Sprintf("SELECT COUNT(id) FROM tbl_cron %s", where)
	err := repo.db.QueryRow(queryTotalRows, queryParams...).Scan(&tasks.TotalRows)
	if err != nil {
		return tasks, err
	}
	queryParams = append(queryParams, q.RowsPerPage, (q.Page-1)*q.RowsPerPage)
	allowedSort := []string{"id", "run_every", "last_run", "next_run", "run_time"}
	sortBy := "id"
	if slices.Contains(allowedSort, q.SortBy) {
		sortBy = q.SortBy
	}
	sortDirection := "DESC"
	if q.SortDirection == "asc" {
		sortDirection = "ASC"
	}

	query := fmt.Sprintf("SELECT id, title, slug, description, run_every, last_run, next_run, run_time, cron_state, is_enabled FROM tbl_cron %s ORDER BY %s %s LIMIT ? OFFSET ?", where, sortBy, sortDirection)
	err = repo.db.Select(&tasks.Items, query, queryParams...)
	if err != nil {
		return tasks, err
	}
	tasks.RowsPerPage = q.RowsPerPage

	return tasks, nil
}

func (repo *CronRepo) queueList() ([]Cron, error) {
	tasks := []Cron{}
	query := `SELECT id, run_every, last_run, slug, next_run, cron_state FROM tbl_cron
    WHERE is_enabled = ? AND next_run <= ?`
	err := repo.db.Select(&tasks, query, 1, time.Now().Unix())

	return tasks, err
}

func (repo *CronRepo) preRunTask(id int, lastRunTs int64) {
	query := `UPDATE tbl_cron SET cron_state = ?, last_run = ? WHERE id = ?`
	row, err := repo.db.Query(query, 1, lastRunTs, id)
	if err != nil {
		slog.Error(fmt.Sprintf("[CRON] Failed to update pre cron state: %s", err))
	}
	defer row.Close()
}

func (repo *CronRepo) postRunTask(id int, nextRun int64, runtime float64) {
	query := `UPDATE tbl_cron SET cron_state = ?, next_run = ?, run_time = ? WHERE id = ?`
	row, err := repo.db.Query(query, 0, nextRun, runtime, id)
	if err != nil {
		slog.Error(fmt.Sprintf("[CRON] Failed to update post cron state: %s", err))
	}
	defer row.Close()
}

func (repo *CronRepo) execCron(slug string) {
	switch slug {
	case "delete_old_probe_logs":
		slog.Info(fmt.Sprintf("[CRON] Start running task: %s", slug))
		repo.deleteOldProbeLogs()
		break
	}
}

func (repo *CronRepo) deleteOldProbeLogs() {
	// for now, we only delete stats older than 1 month +2 days
	startTs := time.Now().AddDate(0, -1, -2).Unix()
	query := `DELETE FROM tbl_probe_log WHERE date_checked < ?`
	_, err := repo.db.Exec(query, startTs)
	if err != nil {
		slog.Error(fmt.Sprintf("[CRON] Failed to delete old probe logs: %s", err))
	}
}
