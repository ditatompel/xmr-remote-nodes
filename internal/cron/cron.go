package cron

import (
	"fmt"
	"log/slog"
	"math"
	"time"
	"xmr-remote-nodes/internal/database"
)

type CronRepository interface {
	RunCronProcess(chan struct{})
	Crons() ([]Cron, error)
}

type CronRepo struct {
	db *database.DB
}

type Cron struct {
	ID          int     `json:"id" db:"id"`
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

func New() CronRepository {
	return &CronRepo{db: database.GetDB()}
}

func (r *CronRepo) RunCronProcess(c chan struct{}) {
	for {
		select {
		case <-time.After(60 * time.Second):
			slog.Info("[CRON] Running cron cycle...")
			list, err := r.queueList()
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

				r.preRunTask(task.ID, currentTs)
				r.execCron(task.Slug)

				runTime := math.Ceil(time.Since(startTime).Seconds()*1000) / 1000
				slog.Info(fmt.Sprintf("[CRON] Task %s done in %f seconds", task.Slug, runTime))
				nextRun := currentTs + int64(task.RunEvery)

				r.postRunTask(task.ID, nextRun, runTime)
			}
			slog.Info("[CRON] Cron cycle done!")
		case <-c:
			slog.Info("[CRON] Shutting down cron...")
			return
		}
	}
}

func (r *CronRepo) Crons() ([]Cron, error) {
	var tasks []Cron
	err := r.db.Select(&tasks, `
		SELECT
			id,
			title,
			slug,
			description,
			run_every,
			last_run,
			next_run,
			run_time,
			cron_state,
			is_enabled
		FROM
			tbl_cron`)
	return tasks, err
}

func (r *CronRepo) queueList() ([]Cron, error) {
	tasks := []Cron{}
	query := `
		SELECT
			id,
			run_every,
			last_run,
			slug,
			next_run,
			cron_state
		FROM
			tbl_cron
		WHERE
			is_enabled = ?
			AND next_run <= ?`
	err := r.db.Select(&tasks, query, 1, time.Now().Unix())

	return tasks, err
}

func (r *CronRepo) preRunTask(id int, lastRunTs int64) {
	query := `
		UPDATE tbl_cron
		SET
			cron_state = ?,
			last_run = ?
		WHERE
			id = ?`
	row, err := r.db.Query(query, 1, lastRunTs, id)
	if err != nil {
		slog.Error(fmt.Sprintf("[CRON] Failed to update pre cron state: %s", err))
	}
	defer row.Close()
}

func (r *CronRepo) postRunTask(id int, nextRun int64, runtime float64) {
	query := `
		UPDATE tbl_cron
		SET
			cron_state = ?,
			next_run = ?,
			run_time = ?
		WHERE
			id = ?`
	row, err := r.db.Query(query, 0, nextRun, runtime, id)
	if err != nil {
		slog.Error(fmt.Sprintf("[CRON] Failed to update post cron state: %s", err))
	}
	defer row.Close()
}

func (r *CronRepo) execCron(slug string) {
	switch slug {
	case "delete_old_probe_logs":
		slog.Info(fmt.Sprintf("[CRON] Start running task: %s", slug))
		r.deleteOldProbeLogs()
	}
}

func (r *CronRepo) deleteOldProbeLogs() {
	// for now, we only delete stats older than 1 month +2 days
	startTs := time.Now().AddDate(0, -1, -2).Unix()
	query := `DELETE FROM tbl_probe_log WHERE date_checked < ?`
	_, err := r.db.Exec(query, startTs)
	if err != nil {
		slog.Error(fmt.Sprintf("[CRON] Failed to delete old probe logs: %s", err))
	}
}
