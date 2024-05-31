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
	"xmr-remote-nodes/internal/geo"
)

type QueryLogs struct {
	NodeID       int    // 0 fpr all, >0 for specific node
	Status       int    // -1 for all, 0 for failed, 1 for success
	FailedReason string // empty for all, if not empty, will be used as search from failed_reaso

	RowsPerPage   int
	Page          int
	SortBy        string
	SortDirection string
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
	AdjustedTime int     `db:"adjusted_time" json:"adjusted_time"`
	DatabaseSize int     `db:"database_size" json:"database_size"`
	Difficulty   int     `db:"difficulty" json:"difficulty"`
	EstimateFee  int     `db:"estimate_fee" json:"estimate_fee"`
	DateChecked  int     `db:"date_checked" json:"date_checked"`
	FailedReason string  `db:"failed_reason" json:"failed_reason"`
	FetchRuntime float64 `db:"fetch_runtime" json:"fetch_runtime"`
}

type FetchLogs struct {
	TotalRows   int         `json:"total_rows"`
	RowsPerPage int         `json:"rows_per_page"`
	Items       []*FetchLog `json:"items"`
}

// Logs returns list of fetched log result for given query
func (r *MoneroRepo) Logs(q QueryLogs) (FetchLogs, error) {
	args, where, sortBy, sortDirection := q.toSQL()

	var fetchLogs FetchLogs
	fetchLogs.RowsPerPage = q.RowsPerPage

	qTotal := fmt.Sprintf(`SELECT COUNT(id) FROM tbl_probe_log %s`, where)
	err := r.db.QueryRow(qTotal, args...).Scan(&fetchLogs.TotalRows)
	if err != nil {
		return fetchLogs, err
	}
	args = append(args, q.RowsPerPage, (q.Page-1)*q.RowsPerPage)

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
func (r *MoneroRepo) GiveJob(acceptTor int) (Node, error) {
	args := []interface{}{}
	wq := []string{}
	where := ""

	if acceptTor != 1 {
		wq = append(wq, "is_tor = ?")
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
	NodeInfo Node    `json:"node_info"`
}

func (r *MoneroRepo) ProcessJob(report ProbeReport, proberId int64) error {
	if report.NodeInfo.ID == 0 {
		return errors.New("Invalid node")
	}

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
		report.NodeInfo.ID,
		proberId,
		report.NodeInfo.IsAvailable,
		report.NodeInfo.Height,
		report.NodeInfo.AdjustedTime,
		report.NodeInfo.DatabaseSize,
		report.NodeInfo.Difficulty,
		report.NodeInfo.EstimateFee,
		time.Now().Unix(),
		report.Message,
		report.TookTime)
	if err != nil {
		return err
	}

	now := time.Now()
	limitTs := now.AddDate(0, -1, 0).Unix()

	nodeStats := struct {
		OnlineCount  uint `db:"online"`
		OfflineCount uint `db:"offline"`
		TotalFetched uint `db:"total_fetched"`
	}{}

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
	if err := r.db.Get(&nodeStats, qstats, report.NodeInfo.ID, limitTs); err != nil {
		slog.Warn(err.Error())
	}

	avgUptime := (float64(nodeStats.OnlineCount) / float64(nodeStats.TotalFetched)) * 100
	report.NodeInfo.Uptime = math.Ceil(avgUptime*100) / 100

	var statuses [5]int
	errUnmarshal := report.NodeInfo.LastCheckStatus.Unmarshal(&statuses)
	if errUnmarshal != nil {
		fmt.Println("Warning", errUnmarshal.Error())
		statuses = [5]int{2, 2, 2, 2, 2}
	}

	nodeAvailable := 0

	if report.NodeInfo.IsAvailable {
		nodeAvailable = 1
	}
	newStatuses := statuses[1:]
	newStatuses = append(newStatuses, nodeAvailable)
	statuesValueToDb, errMarshalStatus := json.Marshal(newStatuses)
	if errMarshalStatus != nil {
		fmt.Println("WARN", errMarshalStatus.Error())
	}

	// recheck IP
	if report.NodeInfo.IP != "" {
		if ipInfo, errGeoIp := geo.Info(report.NodeInfo.IP); errGeoIp != nil {
			fmt.Println("WARN:", errGeoIp.Error())
		} else {
			report.NodeInfo.ASN = ipInfo.ASN
			report.NodeInfo.ASNName = ipInfo.ASNOrg
			report.NodeInfo.CountryCode = ipInfo.CountryCode
			report.NodeInfo.CountryName = ipInfo.CountryName
			report.NodeInfo.City = ipInfo.City
			report.NodeInfo.Longitude = ipInfo.Longitude
			report.NodeInfo.Latitude = ipInfo.Latitude
		}
	}

	if report.NodeInfo.IsAvailable {
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
			cors_capable = ?
		WHERE
			id = ?`
		_, err := r.db.Exec(update,
			nodeAvailable,
			report.NodeInfo.Nettype,
			report.NodeInfo.Height,
			report.NodeInfo.AdjustedTime,
			report.NodeInfo.DatabaseSize,
			report.NodeInfo.Difficulty,
			report.NodeInfo.Version,
			report.NodeInfo.Uptime,
			report.NodeInfo.EstimateFee,
			report.NodeInfo.IP,
			report.NodeInfo.ASN,
			report.NodeInfo.ASNName,
			report.NodeInfo.CountryCode,
			report.NodeInfo.CountryName,
			report.NodeInfo.City,
			now.Unix(),
			string(statuesValueToDb),
			report.NodeInfo.CORSCapable,
			report.NodeInfo.ID)
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
			last_check_status = ?
		WHERE
			id = ?`
		if _, err := r.db.Exec(u, nodeAvailable, report.NodeInfo.Uptime, now.Unix(), string(statuesValueToDb), report.NodeInfo.ID); err != nil {
			slog.Warn(err.Error())
		}
	}

	if avgUptime <= 0 && nodeStats.TotalFetched > 300 {
		fmt.Println("Deleting Monero node (0% uptime from > 300 records)")
		if err := r.Delete(report.NodeInfo.ID); err != nil {
			slog.Warn(err.Error())
		}
	}

	_, err = r.db.Exec(`
		UPDATE tbl_prober
		SET last_submit_ts = ?
		WHERE id = ?`, now.Unix(), proberId)

	return err
}
