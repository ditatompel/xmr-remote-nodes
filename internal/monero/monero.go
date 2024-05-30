package monero

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"net"
	"slices"
	"strings"
	"time"
	"xmr-remote-nodes/internal/database"
	"xmr-remote-nodes/internal/geo"

	"github.com/jmoiron/sqlx/types"
)

type MoneroRepository interface {
	Node(id int) (Node, error)
	Add(protocol string, host string, port uint) error
	Nodes(QueryNodes) (Nodes, error)
	GiveJob(acceptTor int) (Node, error)
	ProcessJob(report ProbeReport, proberId int64) error
	NetFee() []NetFee
	Countries() ([]Countries, error)
	Logs(QueryLogs) (FetchLogs, error)
}

type MoneroRepo struct {
	db *database.DB
}

func New() MoneroRepository {
	return &MoneroRepo{db: database.GetDB()}
}

// Node represents a single remote node
type Node struct {
	ID              uint           `json:"id,omitempty" db:"id"`
	Hostname        string         `json:"hostname" db:"hostname"`
	IP              string         `json:"ip" db:"ip_addr"`
	Port            uint           `json:"port" db:"port"`
	Protocol        string         `json:"protocol" db:"protocol"`
	IsTor           bool           `json:"is_tor" db:"is_tor"`
	IsAvailable     bool           `json:"is_available" db:"is_available"`
	Nettype         string         `json:"nettype" db:"nettype"`
	Height          uint           `json:"height" db:"height"`
	AdjustedTime    uint           `json:"adjusted_time" db:"adjusted_time"`
	DatabaseSize    uint           `json:"database_size" db:"database_size"`
	Difficulty      uint           `json:"difficulty" db:"difficulty"`
	Version         string         `json:"version" db:"version"`
	Status          string         `json:"status,omitempty"`
	Uptime          float64        `json:"uptime" db:"uptime"`
	EstimateFee     uint           `json:"estimate_fee" db:"estimate_fee"`
	ASN             uint           `json:"asn" db:"asn"`
	ASNName         string         `json:"asn_name" db:"asn_name"`
	CountryCode     string         `json:"cc" db:"country"`
	CountryName     string         `json:"country_name" db:"country_name"`
	City            string         `json:"city" db:"city"`
	Latitude        float64        `json:"latitude" db:"lat"`
	Longitude       float64        `json:"longitude" db:"lon"`
	DateEntered     uint           `json:"date_entered,omitempty" db:"date_entered"`
	LastChecked     uint           `json:"last_checked" db:"last_checked"`
	FailedCount     uint           `json:"failed_count,omitempty" db:"failed_count"`
	LastCheckStatus types.JSONText `json:"last_check_statuses" db:"last_check_status"`
	CORSCapable     bool           `json:"cors" db:"cors_capable"`
}

// Get node from database by id
func (repo *MoneroRepo) Node(id int) (Node, error) {
	var node Node
	err := repo.db.Get(&node, `SELECT * FROM tbl_node WHERE id = ?`, id)
	if err != nil && err != sql.ErrNoRows {
		fmt.Println("WARN:", err)
		return node, errors.New("Can't get node information")
	}
	if err == sql.ErrNoRows {
		return node, errors.New("Node not found")
	}
	return node, err
}

// Nodes represents a list of nodes
type Nodes struct {
	TotalRows   int     `json:"total_rows"`
	RowsPerPage int     `json:"rows_per_page"`
	Items       []*Node `json:"items"`
}

// QueryNodes represents database query parameters
type QueryNodes struct {
	Host     string
	Nettype  string
	Protocol string
	CC       string // 2 letter country code
	Status   int
	CORS     int

	// pagination
	RowsPerPage   int
	Page          int
	SortBy        string
	SortDirection string
}

// Get nodes from database
func (repo *MoneroRepo) Nodes(q QueryNodes) (Nodes, error) {
	queryParams := []interface{}{}
	whereQueries := []string{}
	where := ""

	if q.Host != "" {
		whereQueries = append(whereQueries, "(hostname LIKE ? OR ip_addr LIKE ?)")
		queryParams = append(queryParams, "%"+q.Host+"%")
		queryParams = append(queryParams, "%"+q.Host+"%")
	}
	if q.Nettype != "any" {
		if q.Nettype != "mainnet" && q.Nettype != "stagenet" && q.Nettype != "testnet" {
			return Nodes{}, errors.New("Invalid nettype, must be one of 'mainnet', 'stagenet', 'testnet' or 'any'")
		}
		whereQueries = append(whereQueries, "nettype = ?")
		queryParams = append(queryParams, q.Nettype)
	}
	if q.Protocol != "any" {
		allowedProtocols := []string{"tor", "http", "https"}
		if !slices.Contains(allowedProtocols, q.Protocol) {
			return Nodes{}, errors.New("Invalid protocol, must be one of '" + strings.Join(allowedProtocols, "', '") + "' or 'any'")
		}
		if q.Protocol == "tor" {
			whereQueries = append(whereQueries, "is_tor = ?")
			queryParams = append(queryParams, 1)
		} else {
			whereQueries = append(whereQueries, "(protocol = ? AND is_tor = ?)")
			queryParams = append(queryParams, q.Protocol)
			queryParams = append(queryParams, 0)
		}
	}
	if q.CC != "any" {
		whereQueries = append(whereQueries, "country = ?")
		if q.CC == "UNKNOWN" {
			queryParams = append(queryParams, "")
		} else {
			queryParams = append(queryParams, q.CC)
		}
	}
	if q.Status != -1 {
		whereQueries = append(whereQueries, "is_available = ?")
		queryParams = append(queryParams, q.Status)
	}
	if q.CORS != -1 {
		whereQueries = append(whereQueries, "cors_capable = ?")
		queryParams = append(queryParams, 1)
	}

	if len(whereQueries) > 0 {
		where = "WHERE " + strings.Join(whereQueries, " AND ")
	}

	nodes := Nodes{}

	queryTotalRows := fmt.Sprintf(`
		SELECT
			COUNT(id) AS total_rows
		FROM
			tbl_node
		%s`, where)

	err := repo.db.QueryRow(queryTotalRows, queryParams...).Scan(&nodes.TotalRows)
	if err != nil {
		return nodes, err
	}
	queryParams = append(queryParams, q.RowsPerPage, (q.Page-1)*q.RowsPerPage)

	allowedSort := []string{"last_checked", "uptime"}
	sortBy := "last_checked"
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
			protocol,
			hostname,
			port,
			is_tor,
			is_available,
			nettype,
			height,
			adjusted_time,
			database_size,
			difficulty,
			version,
			uptime,
			estimate_fee,
			ip_addr,
			asn,
			asn_name,
			country,
			country_name,
			city,
			lat,
			lon,
			date_entered,
			last_checked,
			last_check_status,
			cors_capable
		FROM
			tbl_node
		%s -- where query if any
		ORDER BY
			%s
			%s
		LIMIT ?
		OFFSET ?`, where, sortBy, sortDirection)

	row, err := repo.db.Query(query, queryParams...)
	if err != nil {
		return nodes, err
	}
	defer row.Close()

	nodes.RowsPerPage = q.RowsPerPage

	for row.Next() {
		var node Node
		err = row.Scan(
			&node.ID,
			&node.Protocol,
			&node.Hostname,
			&node.Port,
			&node.IsTor,
			&node.IsAvailable,
			&node.Nettype,
			&node.Height,
			&node.AdjustedTime,
			&node.DatabaseSize,
			&node.Difficulty,
			&node.Version,
			&node.Uptime,
			&node.EstimateFee,
			&node.IP,
			&node.ASN,
			&node.ASNName,
			&node.CountryCode,
			&node.CountryName,
			&node.City,
			&node.Latitude,
			&node.Longitude,
			&node.DateEntered,
			&node.LastChecked,
			&node.LastCheckStatus,
			&node.CORSCapable)
		if err != nil {
			return nodes, err
		}
		nodes.Items = append(nodes.Items, &node)
	}

	return nodes, nil
}

type QueryLogs struct {
	NodeID       int    // 0 fpr all, >0 for specific node
	WorkerID     int    // 0 for all, >0 for specific worker
	Status       int    // -1 for all, 0 for failed, 1 for success
	FailedReason string // empty for all, if not empty, will be used as search from failed_reaso

	RowsPerPage   int
	Page          int
	SortBy        string
	SortDirection string
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
func (repo *MoneroRepo) Logs(q QueryLogs) (FetchLogs, error) {
	queryParams := []interface{}{}
	whereQueries := []string{}
	where := ""

	if q.NodeID != 0 {
		whereQueries = append(whereQueries, "node_id = ?")
		queryParams = append(queryParams, q.NodeID)
	}
	if q.Status != -1 {
		whereQueries = append(whereQueries, "is_available = ?")
		queryParams = append(queryParams, q.Status)
	}
	if q.FailedReason != "" {
		whereQueries = append(whereQueries, "failed_reason LIKE ?")
		queryParams = append(queryParams, "%"+q.FailedReason+"%")
	}

	if len(whereQueries) > 0 {
		where = "WHERE " + strings.Join(whereQueries, " AND ")
	}

	var fetchLogs FetchLogs

	queryTotalRows := fmt.Sprintf("SELECT COUNT(id) FROM tbl_probe_log %s", where)
	err := repo.db.QueryRow(queryTotalRows, queryParams...).Scan(&fetchLogs.TotalRows)
	if err != nil {
		return fetchLogs, err
	}
	queryParams = append(queryParams, q.RowsPerPage, (q.Page-1)*q.RowsPerPage)

	allowedSort := []string{"date_checked", "fetch_runtime"}
	sortBy := "id"
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
		FROM
			tbl_probe_log
		%s -- where query
		ORDER BY
			%s
			%s
		LIMIT ?
		OFFSET ?`, where, sortBy, sortDirection)

	row, err := repo.db.Query(query, queryParams...)
	if err != nil {
		return fetchLogs, err
	}
	defer row.Close()

	fetchLogs.RowsPerPage = q.RowsPerPage

	for row.Next() {
		var fl FetchLog
		err = row.Scan(
			&fl.ID,
			&fl.NodeID,
			&fl.ProberID,
			&fl.Status,
			&fl.Height,
			&fl.AdjustedTime,
			&fl.DatabaseSize,
			&fl.Difficulty,
			&fl.EstimateFee,
			&fl.DateChecked,
			&fl.FailedReason,
			&fl.FetchRuntime)
		if err != nil {
			return fetchLogs, err
		}

		fetchLogs.Items = append(fetchLogs.Items, &fl)
	}

	return fetchLogs, nil
}

func (repo *MoneroRepo) Add(protocol string, hostname string, port uint) error {
	if protocol != "http" && protocol != "https" {
		return errors.New("Invalid protocol, must one of or HTTP/HTTPS")
	}

	if port > 65535 || port < 1 {
		return errors.New("Invalid port number")
	}

	is_tor := false
	if strings.HasSuffix(hostname, ".onion") {
		is_tor = true
	}
	ip := ""

	if !is_tor {
		hostIps, err := net.LookupIP(hostname)
		if err != nil {
			return err
		}

		hostIp := hostIps[0].To4()
		if hostIp == nil {
			return errors.New("Host IP is not IPv4")
		}
		if hostIp.IsPrivate() {
			return errors.New("IP address is private")
		}
		if hostIp.IsLoopback() {
			return errors.New("IP address is loopback address")
		}

		ip = hostIp.String()
	}

	row, err := repo.db.Query(`
		SELECT
			id
		FROM
			tbl_node
		WHERE
			protocol = ?
			AND hostname = ?
			AND port = ?
		LIMIT 1`, protocol, hostname, port)
	if err != nil {
		return err
	}
	defer row.Close()

	if row.Next() {
		return errors.New("Node already monitored")
	}
	statusDb, _ := json.Marshal([5]int{2, 2, 2, 2, 2})
	_, err = repo.db.Exec(`
		INSERT INTO tbl_node (
			protocol,
			hostname,
			port,
			is_tor,
			nettype,
			ip_addr,
			lat,
			lon,
			date_entered,
			last_checked,
			last_check_status
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
		)`,
		protocol,
		hostname,
		port,
		is_tor,
		"",
		ip,
		0,
		0,
		time.Now().Unix(),
		0,
		string(statusDb))
	if err != nil {
		return err
	}

	return nil
}

func (repo *MoneroRepo) Delete(id uint) error {
	if _, err := repo.db.Exec(`DELETE FROM tbl_node WHERE id = ?`, id); err != nil {
		return err
	}
	if _, err := repo.db.Exec(`DELETE FROM tbl_probe_log WHERE node_id = ?`, id); err != nil {
		return err
	}

	return nil
}

func (repo *MoneroRepo) GiveJob(acceptTor int) (Node, error) {
	queryParams := []interface{}{}
	whereQueries := []string{}
	where := ""

	if acceptTor != 1 {
		whereQueries = append(whereQueries, "is_tor = ?")
		queryParams = append(queryParams, 0)
	}

	if len(whereQueries) > 0 {
		where = "WHERE " + strings.Join(whereQueries, " AND ")
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
	err := repo.db.QueryRow(query, queryParams...).Scan(
		&node.ID,
		&node.Hostname,
		&node.Port,
		&node.Protocol,
		&node.IsTor,
		&node.LastCheckStatus)
	if err != nil {
		return node, err
	}

	_, err = repo.db.Exec(`
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

func (repo *MoneroRepo) ProcessJob(report ProbeReport, proberId int64) error {
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
	_, err := repo.db.Exec(qInsertLog,
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
	if err := repo.db.Get(&nodeStats, qstats, report.NodeInfo.ID, limitTs); err != nil {
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
		_, err := repo.db.Exec(update,
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
		if _, err := repo.db.Exec(u, nodeAvailable, report.NodeInfo.Uptime, now.Unix(), string(statuesValueToDb), report.NodeInfo.ID); err != nil {
			slog.Warn(err.Error())
		}
	}

	if avgUptime <= 0 && nodeStats.TotalFetched > 300 {
		fmt.Println("Deleting Monero node (0% uptime from > 300 records)")
		if err := repo.Delete(report.NodeInfo.ID); err != nil {
			slog.Warn(err.Error())
		}
	}

	_, err = repo.db.Exec(`
		UPDATE tbl_prober
		SET last_submit_ts = ?
		WHERE id = ?`, now.Unix(), proberId)

	return err
}

type NetFee struct {
	Nettype     string `json:"nettype" db:"nettype"`
	EstimateFee uint   `json:"estimate_fee" db:"estimate_fee"`
	NodeCount   int    `json:"node_count" db:"node_count"`
}

func (repo *MoneroRepo) NetFee() []NetFee {
	netTypes := [3]string{"mainnet", "stagenet", "testnet"}
	netFees := []NetFee{}

	for _, net := range netTypes {
		fees := NetFee{}
		err := repo.db.Get(&fees, `
			SELECT
				COUNT(id) AS node_count,
				nettype,
				estimate_fee
			FROM
				tbl_node
			WHERE
				nettype = ?
			GROUP BY
				estimate_fee
			ORDER BY
				node_count DESC
			LIMIT 1`, net)
		if err != nil {
			fmt.Println("WARN:", err.Error())
			continue
		}
		netFees = append(netFees, fees)
	}

	return netFees
}

type Countries struct {
	TotalNodes int    `json:"total_nodes" db:"total_nodes"`
	CC         string `json:"cc" db:"country"` // country code
	Name       string `json:"name" db:"country_name"`
}

func (repo *MoneroRepo) Countries() ([]Countries, error) {
	countries := []Countries{}
	err := repo.db.Select(&countries, `
		SELECT
			COUNT(id) AS total_nodes,
			country,
			country_name
		FROM
			tbl_node
		GROUP BY
			country
		ORDER BY
			country ASC`)
	return countries, err
}
