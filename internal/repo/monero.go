package repo

import (
	"encoding/json"
	"errors"
	"fmt"
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
	LastHeight      uint           `json:"last_height" db:"last_height"`
	AdjustedTime    uint           `json:"adjusted_time" db:"adjusted_time"`
	DatabaseSize    uint           `json:"database_size" db:"database_size"`
	Difficulty      uint           `json:"difficulty" db:"difficulty"`
	NodeVersion     string         `json:"node_version" db:"node_version"`
	Uptime          float32        `json:"uptime" db:"uptime"`
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

	query := fmt.Sprintf("SELECT id, protocol, hostname, port, is_tor, is_available, nettype, last_height, adjusted_time, database_size, difficulty, node_version, uptime, estimate_fee, ip_addr, asn, asn_name, country, country_name, city, lat, lon, date_entered, last_checked, last_check_status, cors_capable FROM tbl_node %s ORDER BY %s %s LIMIT ? OFFSET ?", where, sortBy, sortDirection)

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
		err = row.Scan(&node.Id, &node.Protocol, &node.Hostname, &node.Port, &node.IsTor, &node.IsAvailable, &node.NetType, &node.LastHeight, &node.AdjustedTime, &node.DatabaseSize, &node.Difficulty, &node.NodeVersion, &node.Uptime, &node.EstimateFee, &node.Ip, &node.Asn, &node.AsnName, &node.CountryName, &node.CountryCode, &node.City, &node.Lat, &node.Lon, &node.DateEntered, &node.LastChecked, &node.LastCheckStatus, &node.CorsCapable)
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
