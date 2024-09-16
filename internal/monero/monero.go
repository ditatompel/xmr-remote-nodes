package monero

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"slices"
	"strings"
	"time"

	"github.com/ditatompel/xmr-remote-nodes/internal/database"
	"github.com/ditatompel/xmr-remote-nodes/internal/ip"
	"github.com/jmoiron/sqlx/types"
)

type moneroRepo struct {
	db *database.DB
}

func New() *moneroRepo {
	return &moneroRepo{db: database.GetDB()}
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
	IPv6Only        bool           `json:"ipv6_only" db:"ipv6_only"`
	IPAddresses     string         `json:"ip_addresses" db:"ip_addresses"`
}

// Get node from database by id
func (r *moneroRepo) Node(id int) (Node, error) {
	var node Node
	err := r.db.Get(&node, `SELECT * FROM tbl_node WHERE id = ?`, id)
	if err != nil && err != sql.ErrNoRows {
		slog.Error(err.Error())
		return node, errors.New("Can't get node information")
	}
	if err == sql.ErrNoRows {
		return node, errors.New("Node not found")
	}
	return node, err
}

// QueryNodes represents database query parameters
type QueryNodes struct {
	Host     string
	Nettype  string // Can be "any", mainnet, stagenet, testnet. Default: "any"
	Protocol string // Can be "any", tor, http, https. Default: "any"
	CC       string // 2 letter country code
	Status   int
	CORS     int

	// pagination
	RowsPerPage   int
	Page          int
	SortBy        string
	SortDirection string
}

// toSQL generates SQL query from query parameters
func (q *QueryNodes) toSQL() (args []interface{}, where string) {
	wq := []string{}

	if q.Host != "" {
		wq = append(wq, "(hostname LIKE ? OR ip_addr LIKE ?)")
		args = append(args, "%"+q.Host+"%", "%"+q.Host+"%")
	}
	if slices.Contains([]string{"mainnet", "stagenet", "testnet"}, q.Nettype) {
		wq = append(wq, "nettype = ?")
		args = append(args, q.Nettype)
	}
	if q.Protocol != "any" && slices.Contains([]string{"tor", "http", "https"}, q.Protocol) {
		if q.Protocol == "tor" {
			wq = append(wq, "is_tor = ?")
			args = append(args, 1)
		} else {
			wq = append(wq, "(protocol = ? AND is_tor = ?)")
			args = append(args, q.Protocol, 0)
		}
	}
	if q.CC != "any" {
		wq = append(wq, "country = ?")
		if q.CC == "UNKNOWN" {
			args = append(args, "")
		} else {
			args = append(args, q.CC)
		}
	}
	if q.Status != -1 {
		wq = append(wq, "is_available = ?")
		args = append(args, q.Status)
	}
	if q.CORS != -1 {
		wq = append(wq, "cors_capable = ?")
		args = append(args, q.CORS)
	}

	if len(wq) > 0 {
		where = "WHERE " + strings.Join(wq, " AND ")
	}

	if !slices.Contains([]string{"last_checked", "uptime"}, q.SortBy) {
		q.SortBy = "last_checked"
	}
	if q.SortDirection != "asc" {
		q.SortDirection = "DESC"
	}

	return args, where
}

// Nodes represents a list of nodes
type Nodes struct {
	TotalRows   int     `json:"total_rows"`
	RowsPerPage int     `json:"rows_per_page"`
	Items       []*Node `json:"items"`
}

// Get nodes from database
func (r *moneroRepo) Nodes(q QueryNodes) (Nodes, error) {
	args, where := q.toSQL()

	var nodes Nodes

	nodes.RowsPerPage = q.RowsPerPage

	qTotal := fmt.Sprintf(`
		SELECT
			COUNT(id) AS total_rows
		FROM
			tbl_node
		%s`, where)

	err := r.db.QueryRow(qTotal, args...).Scan(&nodes.TotalRows)
	if err != nil {
		return nodes, err
	}
	args = append(args, q.RowsPerPage, (q.Page-1)*q.RowsPerPage)

	query := fmt.Sprintf(`
		SELECT
			*
		FROM
			tbl_node
		%s
		ORDER BY
			%s
			%s
		LIMIT ?
		OFFSET ?`, where, q.SortBy, q.SortDirection)
	err = r.db.Select(&nodes.Items, query, args...)

	return nodes, err
}

func (r *moneroRepo) Add(protocol string, hostname string, port uint) error {
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

	ipAddr := ""
	ips := ""
	ipv6_only := false

	if !is_tor {
		hostIps, err := net.LookupIP(hostname)
		if err != nil {
			return err
		}

		ipv6_only = ip.IsIPv6Only(hostIps)

		hostIp := hostIps[0]
		if hostIp.IsPrivate() {
			return errors.New("IP address is private")
		}
		if hostIp.IsLoopback() {
			return errors.New("IP address is loopback address")
		}

		ipAddr = hostIp.String()
		ips = ip.SliceToString(hostIps)
	} else {
		if strings.HasPrefix(hostname, "http://") || strings.HasPrefix(hostname, "https://") {
			return errors.New("Don't start hostname with http:// or https://, just put your hostname")
		}
	}

	row, err := r.db.Query(`
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
	_, err = r.db.Exec(`
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
			last_check_status,
			ip_addresses,
			ipv6_only
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
			?,
			?,
			?
		)`,
		protocol,
		hostname,
		port,
		is_tor,
		"",
		ipAddr,
		0,
		0,
		time.Now().Unix(),
		0,
		string(statusDb),
		ips,
		ipv6_only)
	if err != nil {
		return err
	}

	return nil
}

func (r *moneroRepo) Delete(id uint) error {
	if _, err := r.db.Exec(`DELETE FROM tbl_node WHERE id = ?`, id); err != nil {
		return err
	}
	if _, err := r.db.Exec(`DELETE FROM tbl_probe_log WHERE node_id = ?`, id); err != nil {
		return err
	}

	return nil
}

type NetFee struct {
	Nettype     string `json:"nettype" db:"nettype"`
	EstimateFee uint   `json:"estimate_fee" db:"estimate_fee"`
	NodeCount   int    `json:"node_count" db:"node_count"`
}

// Get majority net fee from table tbl_fee
func (r *moneroRepo) NetFees() []NetFee {
	var netFees []NetFee
	err := r.db.Select(&netFees, `
		SELECT
			nettype,
			estimate_fee,
			node_count
		FROM
			tbl_fee
		`)
	if err != nil {
		slog.Error(fmt.Sprintf("[MONERO] Failed to get net fees: %s", err))
	}
	return netFees
}

// Countries represents list of countries
type Countries struct {
	TotalNodes int    `json:"total_nodes" db:"total_nodes"`
	CC         string `json:"cc" db:"country"` // country code
	Name       string `json:"name" db:"country_name"`
}

// Get list of countries (count by nodes)
func (r *moneroRepo) Countries() ([]Countries, error) {
	var c []Countries
	err := r.db.Select(&c, `
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
	return c, err
}
