package repo

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net"
	"slices"
	"strings"
	"time"

	"github.com/ditatompel/xmr-nodes/internal/database"

	"github.com/jmoiron/sqlx/types"
)

type MoneroRepository interface {
	Add(protocol string, host string, port uint) error
	Nodes(q MoneroQueryParams) (MoneroNodes, error)
	GiveJob(acceptTor int) (MoneroNode, error)
	ProcessJob(report ProbeReport, proberId int64) error
}

type MoneroRepo struct {
	db *database.DB
}

func NewMoneroRepo(db *database.DB) MoneroRepository {
	return &MoneroRepo{db}
}

type MoneroNode struct {
	Id              uint           `json:"id,omitempty" db:"id"`
	Hostname        string         `json:"hostname" db:"hostname"`
	Ip              string         `json:"ip" db:"ip_addr"`
	Port            uint           `json:"port" db:"port"`
	Protocol        string         `json:"protocol" db:"protocol"`
	IsTor           bool           `json:"is_tor" db:"is_tor"`
	IsAvailable     bool           `json:"is_available" db:"is_available"`
	NetType         string         `json:"nettype" db:"nettype"`
	Height          uint           `json:"height" db:"height"`
	AdjustedTime    uint           `json:"adjusted_time" db:"adjusted_time"`
	DatabaseSize    uint           `json:"database_size" db:"database_size"`
	Difficulty      uint           `json:"difficulty" db:"difficulty"`
	Version         string         `json:"version" db:"version"`
	Status          string         `json:"status,omitempty"`
	Uptime          float64        `json:"uptime" db:"uptime"`
	EstimateFee     uint           `json:"estimate_fee" db:"estimate_fee"`
	Asn             uint           `json:"asn" db:"asn"`
	AsnName         string         `json:"asn_name" db:"asn_name"`
	CountryCode     string         `json:"cc" db:"country"`
	CountryName     string         `json:"country_name" db:"country_name"`
	City            string         `json:"city" db:"city"`
	Lat             float64        `json:"latitude" db:"lat"`
	Lon             float64        `json:"longitude" db:"lon"`
	DateEntered     uint           `json:"date_entered,omitempty" db:"date_entered"`
	LastChecked     uint           `json:"last_checked" db:"last_checked"`
	FailedCount     uint           `json:"failed_count,omitempty" db:"failed_count"`
	LastCheckStatus types.JSONText `json:"last_check_statuses" db:"last_check_status"`
	CorsCapable     bool           `json:"cors" db:"cors_capable"`
}

type MoneroNodes struct {
	TotalRows   int           `json:"total_rows"`
	RowsPerPage int           `json:"rows_per_page"`
	CurrentPage int           `json:"current_page"`
	NextPage    int           `json:"next_page"`
	Items       []*MoneroNode `json:"items"`
}

type MoneroQueryParams struct {
	Host          string
	RowsPerPage   int
	Page          int
	SortBy        string
	SortDirection string
}

func (repo *MoneroRepo) Nodes(q MoneroQueryParams) (MoneroNodes, error) {
	queryParams := []interface{}{}
	whereQueries := []string{}
	where := ""

	if q.Host != "" {
		whereQueries = append(whereQueries, "(hostname LIKE ? OR ip_addr LIKE ?)")
		queryParams = append(queryParams, "%"+q.Host+"%")
		queryParams = append(queryParams, "%"+q.Host+"%")
	}

	if len(whereQueries) > 0 {
		where = "WHERE " + strings.Join(whereQueries, " AND ")
	}

	nodes := MoneroNodes{}

	queryTotalRows := fmt.Sprintf("SELECT COUNT(id) AS total_rows FROM tbl_node %s", where)

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

	query := fmt.Sprintf("SELECT id, protocol, hostname, port, is_tor, is_available, nettype, height, adjusted_time, database_size, difficulty, version, uptime, estimate_fee, ip_addr, asn, asn_name, country, country_name, city, lat, lon, date_entered, last_checked, last_check_status, cors_capable FROM tbl_node %s ORDER BY %s %s LIMIT ? OFFSET ?", where, sortBy, sortDirection)

	row, err := repo.db.Query(query, queryParams...)
	if err != nil {
		return nodes, err
	}
	defer row.Close()

	nodes.RowsPerPage = q.RowsPerPage
	nodes.CurrentPage = q.Page
	nodes.NextPage = q.Page + 1

	for row.Next() {
		node := MoneroNode{}
		err = row.Scan(&node.Id, &node.Protocol, &node.Hostname, &node.Port, &node.IsTor, &node.IsAvailable, &node.NetType, &node.Height, &node.AdjustedTime, &node.DatabaseSize, &node.Difficulty, &node.Version, &node.Uptime, &node.EstimateFee, &node.Ip, &node.Asn, &node.AsnName, &node.CountryName, &node.CountryCode, &node.City, &node.Lat, &node.Lon, &node.DateEntered, &node.LastChecked, &node.LastCheckStatus, &node.CorsCapable)
		if err != nil {
			return nodes, err
		}
		nodes.Items = append(nodes.Items, &node)
	}

	return nodes, nil
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

	check := `SELECT id FROM tbl_node WHERE protocol = ? AND hostname = ? AND port = ? LIMIT 1`
	row, err := repo.db.Query(check, protocol, hostname, port)
	if err != nil {
		return err
	}
	defer row.Close()

	if row.Next() {
		return errors.New("Node already monitored")
	}
	statusDb, _ := json.Marshal([5]int{2, 2, 2, 2, 2})

	query := `INSERT INTO tbl_node (protocol, hostname, port, is_tor, nettype, ip_addr, lat, lon, date_entered, last_checked, last_check_status) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = repo.db.Exec(query, protocol, hostname, port, is_tor, "", ip, 0, 0, time.Now().Unix(), 0, string(statusDb))
	if err != nil {
		return err
	}

	return nil
}

func (repo *MoneroRepo) GiveJob(acceptTor int) (MoneroNode, error) {
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

	node := MoneroNode{}

	query := fmt.Sprintf(`SELECT id, hostname, port, protocol, is_tor, last_check_status FROM tbl_node %s ORDER BY last_checked ASC LIMIT 1`, where)
	err := repo.db.QueryRow(query, queryParams...).Scan(&node.Id, &node.Hostname, &node.Port, &node.Protocol, &node.IsTor, &node.LastCheckStatus)
	if err != nil {
		return node, err
	}

	update := `UPDATE tbl_node SET last_checked = ? WHERE id = ?`
	_, err = repo.db.Exec(update, time.Now().Unix(), node.Id)
	if err != nil {
		return node, err
	}

	return node, nil
}

type ProbeReport struct {
	TookTime float64    `json:"took_time"`
	Message  string     `json:"message"`
	NodeInfo MoneroNode `json:"node_info"`
}

func (repo *MoneroRepo) ProcessJob(report ProbeReport, proberId int64) error {
	if report.NodeInfo.Id == 0 {
		return errors.New("Invalid node")
	}

	qInsertLog := `INSERT INTO tbl_probe_log (node_id, prober_id, is_available, height, adjusted_time, database_size, difficulty, estimate_fee, date_checked, failed_reason, fetch_runtime) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := repo.db.Exec(qInsertLog, report.NodeInfo.Id, proberId, report.NodeInfo.IsAvailable, report.NodeInfo.Height, report.NodeInfo.AdjustedTime, report.NodeInfo.DatabaseSize, report.NodeInfo.Difficulty, report.NodeInfo.EstimateFee, time.Now().Unix(), report.Message, report.TookTime)
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

	qstats := `SELECT
	SUM(if(is_available='1',1,0)) AS online,
	SUM(if(is_available='0',1,0)) AS offline,
	SUM(if(id='0',0,1)) AS total_fetched FROM
	tbl_probe_log WHERE node_id = ? AND date_checked > ?`
	repo.db.Get(&nodeStats, qstats, report.NodeInfo.Id, limitTs)

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
	// TODO: Fill the data using GeoIP

	// if report.NodeInfo.Ip != "" {
	// 	ipInfo, errGeoIp := GetGeoIpInfo(report.NodeInfo.Ip)
	// 	if errGeoIp == nil {
	// 		report.NodeInfo.Asn = ipInfo.Asn
	// 		report.NodeInfo.AsnName = ipInfo.AsnOrg
	// 		report.NodeInfo.CountryCode = ipInfo.CountryCode
	// 		report.NodeInfo.CountryName = ipInfo.CountryName
	// 		report.NodeInfo.City = ipInfo.City
	// 	}
	// }

	update := `UPDATE tbl_node SET
        is_available = ?, nettype = ?, height = ?, adjusted_time = ?,
        database_size = ?, difficulty = ?, version = ?, uptime = ?,
        estimate_fee = ?, ip_addr = ?, asn = ?, asn_name = ?, country = ?,
      country_name = ?, city = ?, last_checked = ?, last_check_status = ?,
      cors_capable = ?
    WHERE id = ?`

	_, err = repo.db.Exec(update,
		nodeAvailable, report.NodeInfo.NetType, report.NodeInfo.Height, report.NodeInfo.AdjustedTime, report.NodeInfo.DatabaseSize, report.NodeInfo.Difficulty, report.NodeInfo.Version, report.NodeInfo.Uptime, report.NodeInfo.EstimateFee, report.NodeInfo.Ip, report.NodeInfo.Asn, report.NodeInfo.AsnName, report.NodeInfo.CountryCode, report.NodeInfo.CountryName, report.NodeInfo.City, now.Unix(), string(statuesValueToDb), report.NodeInfo.CorsCapable, report.NodeInfo.Id)

	return err
}
