package monero

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"slices"
	"strings"
	"time"

	"github.com/ditatompel/xmr-remote-nodes/internal/ip/geo"
	"github.com/ditatompel/xmr-remote-nodes/internal/paging"
)

type QueryLogs struct {
	paging.Paging
	NodeID       int    `url:"node_id,omitempty"`       // 0 for all, >0 for specific node
	Status       int    `url:"status"`                  // -1 for all, 0 for failed, 1 for success
	FailedReason string `url:"failed_reason,omitempty"` // empty for all, non empty string will be used as search
}

func (q QueryLogs) toSQL() (args []interface{}, where, sortBy, sortDirection string) {
	wq := []string{}
	if q.NodeID != 0 {
		wq = append(wq, "node_id = ?")
		args = append(args, q.NodeID)
	}
	if q.Status != -1 {
		wq = append(wq, "is_available = ?")
		args = append(args, q.Status)
	}
	if q.FailedReason != "" {
		wq = append(wq, "failed_reason LIKE ?")
		args = append(args, "%"+q.FailedReason+"%")
	}

	if len(wq) > 0 {
		where = "WHERE " + strings.Join(wq, " AND ")
	}

	as := []string{"date_checked", "fetch_runtime"}
	sortBy = "id"
	if slices.Contains(as, q.SortBy) {
		sortBy = q.SortBy
	}
	sortDirection = "DESC"
	if q.SortDirection == "asc" {
		sortDirection = "ASC"
	}

	return args, where, sortBy, sortDirection
}

type FetchLog struct {
	ID           int     `db:"id" json:"id,omitempty"`
	NodeID       int     `db:"node_id" json:"node_id"`
	ProberID     int     `db:"prober_id" json:"prober_id"`
	Status       int     `db:"is_available" json:"status"`
	Height       int     `db:"height" json:"height"`
	AdjustedTime int64   `db:"adjusted_time" json:"adjusted_time"`
	DatabaseSize int     `db:"database_size" json:"database_size"`
	Difficulty   int     `db:"difficulty" json:"difficulty"`
	EstimateFee  int     `db:"estimate_fee" json:"estimate_fee"`
	DateChecked  int64   `db:"date_checked" json:"date_checked"`
	FailedReason string  `db:"failed_reason" json:"failed_reason"`
	FetchRuntime float64 `db:"fetch_runtime" json:"fetch_runtime"`
}

type FetchLogs struct {
	TotalRows   int         `json:"total_rows"`
	TotalPages  int         `json:"total_pages"` // total pages
	RowsPerPage int         `json:"rows_per_page"`
	Items       []*FetchLog `json:"items"`
}

// Logs returns list of fetched log result for given query
func (r *moneroRepo) Logs(q QueryLogs) (FetchLogs, error) {
	args, where, sortBy, sortDirection := q.toSQL()

	var fetchLogs FetchLogs
	fetchLogs.RowsPerPage = q.Limit

	qTotal := fmt.Sprintf(`SELECT COUNT(id) FROM tbl_probe_log %s`, where)
	err := r.db.QueryRow(qTotal, args...).Scan(&fetchLogs.TotalRows)
	if err != nil {
		return fetchLogs, err
	}

	fetchLogs.TotalPages = int(math.Ceil(float64(fetchLogs.TotalRows) / float64(q.Limit)))
	args = append(args, q.Limit, (q.Page-1)*q.Limit)

	fmt.Printf("%+v", fetchLogs)
	query := fmt.Sprintf(`
		SELECT
			*
		FROM
			tbl_probe_log
		%s -- where query
		ORDER BY
			%s
			%s
		LIMIT ?
		OFFSET ?`, where, sortBy, sortDirection)
	err = r.db.Select(&fetchLogs.Items, query, args...)

	return fetchLogs, err
}

// GiveJob returns node that should be probed for the next time
func (r *moneroRepo) GiveJob(acceptTor, acceptIPv6 int) (Node, error) {
	args := []interface{}{}
	wq := []string{}
	where := ""

	if acceptTor != 1 {
		wq = append(wq, "is_tor = ?")
		args = append(args, 0)
	}
	if acceptIPv6 != 1 {
		wq = append(wq, "ipv6_only = ?")
		args = append(args, 0)
	}

	if len(wq) > 0 {
		where = "WHERE " + strings.Join(wq, " AND ")
	}

	var node Node

	query := fmt.Sprintf(`
		SELECT
			id,
			hostname,
			port,
			protocol,
			is_tor,
			last_check_status
		FROM
			tbl_node
		%s -- where query if any
		ORDER BY
			last_checked ASC
		LIMIT 1`, where)
	err := r.db.QueryRow(query, args...).Scan(
		&node.ID,
		&node.Hostname,
		&node.Port,
		&node.Protocol,
		&node.IsTor,
		&node.LastCheckStatus)
	if err != nil {
		return node, err
	}

	_, err = r.db.Exec(`
		UPDATE tbl_node
		SET last_checked = ?
		WHERE id = ?`, time.Now().Unix(), node.ID)
	if err != nil {
		return node, err
	}

	return node, nil
}

type ProbeReport struct {
	TookTime float64 `json:"took_time"`
	Message  string  `json:"message"`
	Node     Node    `json:"node"`
}

type nodeStats struct {
	Online       uint `db:"online"`  // total count online
	Offline      uint `db:"offline"` // total count offline
	TotalFetched uint `db:"total_fetched"`
}

// Create new status indicator based from LastCheckStatus and recent IsAvailable status
func (p *ProbeReport) parseStatuses() string {
	var s [5]int
	if err := p.Node.LastCheckStatus.Unmarshal(&s); err != nil {
		slog.Warn(err.Error())
		s = [5]int{2, 2, 2, 2, 2}
	}

	si := 0 // set default "status indicator" to offline
	if p.Node.IsAvailable {
		si = 1
	}

	ns := s[1:]
	ns = append(ns, si)
	j, err := json.Marshal(ns)
	if err != nil {
		slog.Warn(err.Error())
	}

	return string(j)
}

// Process report data from probers
func (r *moneroRepo) ProcessJob(report ProbeReport, proberId int64) error {
	if report.Node.ID == 0 {
		return errors.New("Invalid node")
	}

	now := time.Now()

	qInsertLog := `
		INSERT INTO tbl_probe_log (
			node_id,
			prober_id,
			is_available,
			height,
			adjusted_time,
			database_size,
			difficulty,
			estimate_fee,
			date_checked,
			failed_reason,
			fetch_runtime
		) VALUES (
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?
		)`
	_, err := r.db.Exec(qInsertLog,
		report.Node.ID,
		proberId,
		report.Node.IsAvailable,
		report.Node.Height,
		report.Node.AdjustedTime,
		report.Node.DatabaseSize,
		report.Node.Difficulty,
		report.Node.EstimateFee,
		now.Unix(),
		report.Message,
		report.TookTime)
	if err != nil {
		return err
	}

	limitTs := now.AddDate(0, -1, 0).Unix()

	var stats nodeStats

	qstats := `
		SELECT
			SUM(if(is_available='1',1,0)) AS online,
			SUM(if(is_available='0',1,0)) AS offline,
			SUM(if(id='0',0,1)) AS total_fetched
		FROM
			tbl_probe_log
		WHERE
			node_id = ?
			AND date_checked > ?`
	if err := r.db.Get(&stats, qstats, report.Node.ID, limitTs); err != nil {
		slog.Warn(err.Error())
	}

	avgUptime := (float64(stats.Online) / float64(stats.TotalFetched)) * 100
	report.Node.Uptime = math.Ceil(avgUptime*100) / 100

	statuses := report.parseStatuses()

	// recheck IP
	if report.Node.IP != "" {
		if ipInfo, errGeoIp := geo.Info(report.Node.IP); errGeoIp != nil {
			fmt.Println("WARN:", errGeoIp.Error())
		} else {
			report.Node.ASN = ipInfo.ASN
			report.Node.ASNName = ipInfo.ASNOrg
			report.Node.CountryCode = ipInfo.CountryCode
			report.Node.CountryName = ipInfo.CountryName
			report.Node.City = ipInfo.City
			report.Node.Longitude = ipInfo.Longitude
			report.Node.Latitude = ipInfo.Latitude
		}
	}

	if report.Node.IsAvailable {
		update := `
		UPDATE tbl_node
		SET
			is_available = ?,
			nettype = ?,
			height = ?,
			adjusted_time = ?,
			database_size = ?,
			difficulty = ?,
			version = ?,
			uptime = ?,
			estimate_fee = ?,
			ip_addr = ?,
			asn = ?,
			asn_name = ?,
			country = ?,
			country_name = ?,
			city = ?,
			last_checked = ?,
			last_check_status = ?,
			cors_capable = ?,
			ip_addresses = ?,
			ipv6_only = ?
		WHERE
			id = ?`
		_, err := r.db.Exec(update,
			1,
			report.Node.Nettype,
			report.Node.Height,
			report.Node.AdjustedTime,
			report.Node.DatabaseSize,
			report.Node.Difficulty,
			report.Node.Version,
			report.Node.Uptime,
			report.Node.EstimateFee,
			report.Node.IP,
			report.Node.ASN,
			report.Node.ASNName,
			report.Node.CountryCode,
			report.Node.CountryName,
			report.Node.City,
			now.Unix(),
			statuses,
			report.Node.CORSCapable,
			report.Node.IPAddresses,
			report.Node.IPv6Only,
			report.Node.ID)
		if err != nil {
			slog.Warn(err.Error())
		}
	} else {
		u := `
		UPDATE tbl_node
		SET
			is_available = ?,
			uptime = ?,
			last_checked = ?,
			last_check_status = ?,
			ip_addresses = ?,
			ipv6_only = ?
		WHERE
			id = ?`
		if _, err := r.db.Exec(u, 0, report.Node.Uptime, now.Unix(), statuses, report.Node.IPAddresses, report.Node.IPv6Only, report.Node.ID); err != nil {
			slog.Warn(err.Error())
		}
	}

	if avgUptime <= 0 && stats.TotalFetched > 300 {
		fmt.Println("Deleting Monero node (0% uptime from > 300 records)")
		if err := r.Delete(report.Node.ID); err != nil {
			slog.Warn(err.Error())
		}
	}

	_, err = r.db.Exec(`
		UPDATE tbl_prober
		SET last_submit_ts = ?
		WHERE id = ?`, now.Unix(), proberId)

	return err
}
