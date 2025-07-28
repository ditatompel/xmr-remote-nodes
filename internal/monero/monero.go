package monero

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"net"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/ditatompel/xmr-remote-nodes/internal/database"
	"github.com/ditatompel/xmr-remote-nodes/internal/ip"
	"github.com/ditatompel/xmr-remote-nodes/internal/paging"
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
	IsI2P           bool           `json:"is_i2p" db:"is_i2p"`
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
	DateEntered     int64          `json:"date_entered,omitempty" db:"date_entered"`
	SubmitterIPHash string         `json:"submitter_iphash,omitempty" db:"submitter_iphash"`
	LastChecked     int64          `json:"last_checked" db:"last_checked"`
	FailedCount     uint           `json:"failed_count,omitempty" db:"failed_count"`
	LastCheckStatus types.JSONText `json:"last_check_statuses" db:"last_check_status"`
	CORSCapable     bool           `json:"cors" db:"cors_capable"`
	IPv6Only        bool           `json:"ipv6_only" db:"ipv6_only"`
	IPAddresses     string         `json:"ip_addresses" db:"ip_addresses"`
	IsArchived      int            `json:"is_archived" db:"is_archived"`
	// Rucknium's node data
	IsSpyNode         int `json:"is_spy_node" db:"is_spy_node"`                   // 0 = no, 1 = yes, 2 = not applied
	MRLBanListEnabled int `json:"mrl_ban_list_enabled" db:"mrl_ban_list_enabled"` // 0 = no, 1 = yes, 2 = not applied
	DNSBanListEnabled int `json:"dns_ban_list_enabled" db:"dns_ban_list_enabled"` // 0 = no, 1 = yes, 2 = not applied
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
	paging.Paging
	Host       string `url:"host,omitempty"`
	Nettype    string `url:"nettype,omitempty"`  // Can be empty string, "any", mainnet, stagenet, testnet.
	Protocol   string `url:"protocol,omitempty"` // Can be "any", tor, http, https. Default: "any"
	CC         string `url:"cc,omitempty"`       // 2 letter country code
	Status     int    `url:"status"`
	CORS       string `url:"cors,omitempty"`
	IsArchived int    `url:"archived,omitempty"`
	IsSpyNode  int    `url:"spynode"`
	MRLBan     string `url:"mrlban,omitempty"`
	DNSBan     string `url:"dnsban,omitempty"`
}

// toSQL generates SQL query from query parameters
func (q *QueryNodes) toSQL() (args []interface{}, where string) {
	wq := []string{}

	if q.Host != "" {
		wq = append(wq, "(hostname LIKE ? OR ip_addresses LIKE ?)")
		args = append(args, "%"+q.Host+"%", "%"+q.Host+"%")
	}
	if slices.Contains([]string{"mainnet", "stagenet", "testnet"}, q.Nettype) {
		wq = append(wq, "nettype = ?")
		args = append(args, q.Nettype)
	}
	if q.Protocol != "any" && slices.Contains([]string{"tor", "i2p", "http", "https"}, q.Protocol) {
		switch q.Protocol {
		case "i2p":
			wq = append(wq, "is_i2p = ?")
			args = append(args, 1)
		case "tor":
			wq = append(wq, "is_tor = ?")
			args = append(args, 1)
		default:
			wq = append(wq, "(protocol = ? AND is_tor = ? AND is_i2p = ?)")
			args = append(args, q.Protocol, 0, 0)
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
	if q.CORS == "on" || q.CORS == "1" { // DEPRECATED: CORS = int is deprecated, use CORS = on" instead
		wq = append(wq, "cors_capable = ?")
		args = append(args, 1)
	}
	if q.IsArchived != -1 {
		wq = append(wq, "is_archived = ?")
		args = append(args, q.IsArchived)
	}
	if q.IsSpyNode != -1 {
		wq = append(wq, "is_spy_node = ?")
		args = append(args, q.IsSpyNode)
	}
	if q.MRLBan == "on" {
		wq = append(wq, "mrl_ban_list_enabled = ?")
		args = append(args, 1)
	}
	if q.DNSBan == "on" {
		wq = append(wq, "dns_ban_list_enabled = ?")
		args = append(args, 1)
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
	TotalPages  int     `json:"total_pages"` // total pages
	RowsPerPage int     `json:"rows_per_page"`
	Items       []*Node `json:"items"`
}

// Get nodes from database
func (r *moneroRepo) Nodes(q QueryNodes) (Nodes, error) {
	args, where := q.toSQL()

	var nodes Nodes

	nodes.RowsPerPage = q.Limit

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
	nodes.TotalPages = int(math.Ceil(float64(nodes.TotalRows) / float64(q.Limit)))
	args = append(args, q.Limit, (q.Page-1)*q.Limit)

	query := fmt.Sprintf(`
		SELECT
			id,
			hostname,
			ip_addr,
			port,
			protocol,
			is_tor,
			is_i2p,
			is_available,
			nettype,
			height,
			adjusted_time,
			database_size,
			difficulty,
			version,
			uptime,
			estimate_fee,
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
			cors_capable,
			ipv6_only,
			ip_addresses,
			is_archived,
			is_spy_node,
			mrl_ban_list_enabled,
			dns_ban_list_enabled
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

func (r *moneroRepo) Add(submitterIP, salt, protocol, hostname string, port uint) error {
	if protocol != "http" && protocol != "https" {
		return errors.New("Invalid protocol, must one of or HTTP/HTTPS")
	}

	if port > 65535 || port < 1 {
		return errors.New("Invalid port number")
	}

	is_tor := false
	if strings.HasSuffix(hostname, ".onion") {
		if !validTorHostname(hostname) {
			return errors.New("Invalid TOR v3 .onion hostname")
		}
		is_tor = true
	}

	is_i2p := false
	if strings.HasSuffix(hostname, ".i2p") {
		if !validI2PHostname(hostname) {
			return errors.New("Invalid I2P hostname")
		}
		is_i2p = true
	}

	ipAddr := ""
	ips := ""
	ipv6_only := false

	if !is_tor && !is_i2p {
		// TODO: Find alt solution that return lookup IP to `netip.Addr` instead of `net.IP`.
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

		banList, err := r.banList()
		if err != nil {
			return errors.New("Error finding ban list")
		}

		if isBannedIP(banList, hostIps) {
			return errors.New("Cannot add node: host is in our ban list")
		}
	}

	row, err := r.db.Query(`
		SELECT
			id, is_archived
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
		var (
			id         int
			isArchived int
		)
		if err := row.Scan(&id, &isArchived); err != nil {
			return err
		}

		if isArchived == 1 {
			_, err := r.db.Exec(`UPDATE tbl_node SET is_archived = 0 WHERE id = ?`, id)
			if err != nil {
				return errors.New("failed to update node status")
			}
			return nil
		}
		return errors.New("Node already monitored")
	}

	statusDb, _ := json.Marshal([5]int{2, 2, 2, 2, 2})
	_, err = r.db.Exec(`
		INSERT INTO tbl_node (
			protocol,
			hostname,
			port,
			is_tor,
			is_i2p,
			nettype,
			ip_addr,
			lat,
			lon,
			date_entered,
			submitter_iphash,
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
			?,
			?,
			?
		)`,
		protocol,
		hostname,
		port,
		is_tor,
		is_i2p,
		"",
		ipAddr,
		0,
		0,
		time.Now().Unix(),
		hashIPWithSalt(submitterIP, salt),
		0,
		string(statusDb),
		ips,
		ipv6_only)
	if err != nil {
		return err
	}

	return nil
}

// validTorHostname shecks if a given hostname is a valid TOR v3 .onion address
// with optional subdomain
//
// TOR v3 .onion addresses are 56 characters of `base32` followed by ".onion"
func validTorHostname(hostname string) bool {
	return regexp.MustCompile(`^([a-z0-9-]+\.)*[a-z2-7]{56}\.onion$`).MatchString(hostname)
}

// validI2PHostname checks if a given hostname is a valid b32 or naming service
// I2P address
//
// Old b32 addresses are always {52 chars}.b32.i2p and new ones are
// {56+ chars}.b32.i2p. Since I don't know if there is a length limit of new
// b32 addresses, this function allows up to 63 characters.
//
// For naming service, I2P addresses are up to 67 characters, including the
// '.i2p' part. Please note that this naming service validation only validates
// simple length and allowed characters. Advanced validation such as
// internationalized domain name (IDN) is not implemented.
//
// Ref: https://geti2p.net/spec/b32encrypted and https://geti2p.net/en/docs/naming
func validI2PHostname(hostname string) bool {
	// To minimize abuse, I set minimum length of submitted i2p naming service
	// address to 5 characters. If someone have an address of 4 characters or
	// less, let them open an issue or create a pull request.
	return regexp.MustCompile(`^([a-z2-7]{52,63}\.b32|[a-z0-9-]{5,63})\.i2p$`).MatchString(hostname)
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

// Archive node instead of deleteing the records (especially for spy nodes)
// it could be useful somehow in the future.
// https://github.com/ditatompel/xmr-remote-nodes/issues/191#issuecomment-3090599618
func (r *moneroRepo) Archive(id uint) error {
	if _, err := r.db.Exec(`UPDATE tbl_node SET is_archived = ? WHERE id = ?`, 1, id); err != nil {
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

// hashIPWithSalt hashes IP address with salt designed for checksumming, but
// still maintain user privacy, this is NOT cryptographic security.
func hashIPWithSalt(ip, salt string) string {
	hasher := sha256.New()
	hasher.Write([]byte(salt + ip)) // Combine salt and IP
	return hex.EncodeToString(hasher.Sum(nil))
}

// ParseNodeStatuses parses JSONText into [5]int
// Used this to parse last_check_status for templ engine
func ParseNodeStatuses(statuses types.JSONText) [5]int {
	s := [5]int{}
	if err := statuses.Unmarshal(&s); err != nil {
		return [5]int{2, 2, 2, 2, 2}
	}

	return s
}

// ParseCURLGetInfo generates curl command to get node info from given node
//
// Primarily used for Web UI to display example curl command.
func ParseCURLGetInfo(node Node) string {
	d := `'{"jsonrpc":"2.0","id":"0","method":"get_info"}' -H 'Content-Type: application/json'`

	if node.IsI2P {
		return fmt.Sprintf(
			"curl -x socks5h://127.0.0.1:4447 %s://%s:%d/json_rpc -d %s -sL",
			node.Protocol, node.Hostname, node.Port, d,
		)
	}
	if node.IsTor {
		return fmt.Sprintf(
			"curl -x socks5h://127.0.0.1:9050 %s://%s:%d/json_rpc -d %s -sL",
			node.Protocol, node.Hostname, node.Port, d,
		)
	}

	return fmt.Sprintf(
		"curl %s://%s:%d/json_rpc -d %s -sL",
		node.Protocol, node.Hostname, node.Port, d,
	)
}
