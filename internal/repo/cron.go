package repo

import (
	"fmt"
	"math"
	"time"

	"github.com/ditatompel/xmr-nodes/internal/database"
)

type CronRepository interface {
	RunCronProcess()
	Crons() ([]CronTask, error)
}

type CronRepo struct {
	db *database.DB
}

type CronTask struct {
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
		fmt.Println("Running cron...")
		list, err := repo.queueList()
		if err != nil {
			fmt.Println("Error parsing to struct:", err)
			continue
		}
		for _, task := range list {
			startTime := time.Now()
			currentTs := startTime.Unix()
			delayedTask := currentTs - task.NextRun
			if task.CronState == 1 && delayedTask <= int64(rerunTimeout) {
				fmt.Println("SKIP STATE 1:", task.Slug)
				continue
			}

			repo.preRunTask(task.Id, currentTs)

			repo.execCron(task.Slug)

			runTime := math.Ceil(time.Since(startTime).Seconds()*1000) / 1000
			fmt.Println("Runtime:", runTime)
			nextRun := currentTs + int64(task.RunEvery)

			repo.postRunTask(task.Id, nextRun, runTime)
		}
		fmt.Println("Cron done!")
	}
}

func (repo *CronRepo) Crons() ([]CronTask, error) {
	tasks := []CronTask{}
	query := `SELECT * FROM tbl_cron`
	err := repo.db.Select(&tasks, query)
	return tasks, err
}

func (repo *CronRepo) queueList() ([]CronTask, error) {
	tasks := []CronTask{}
	query := `SELECT id, run_every, last_run, slug, next_run, cron_state FROM tbl_cron
    WHERE is_enabled = ? AND next_run <= ?`
	err := repo.db.Select(&tasks, query, 1, time.Now().Unix())

	return tasks, err
}

func (repo *CronRepo) preRunTask(id int, lastRunTs int64) {
	query := `UPDATE tbl_cron SET cron_state = ?, last_run = ? WHERE id = ?`
	row, err := repo.db.Query(query, 1, lastRunTs, id)
	if err != nil {
		fmt.Println("ERROR PRERUN:", err)
	}
	defer row.Close()
}

func (repo *CronRepo) postRunTask(id int, nextRun int64, runtime float64) {
	query := `UPDATE tbl_cron SET cron_state = ?, next_run = ?, run_time = ? WHERE id = ?`
	row, err := repo.db.Query(query, 0, nextRun, runtime, id)
	if err != nil {
		fmt.Println("ERROR PRERUN:", err)
	}
	defer row.Close()
}

func (repo *CronRepo) execCron(slug string) {
	switch slug {
	case "something":
		fmt.Println("Running task", slug)
		// do task
		break
	}
}
