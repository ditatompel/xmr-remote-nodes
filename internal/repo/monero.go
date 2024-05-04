package repo

import (
	"encoding/json"
	"errors"
	"net"
	"strings"
	"time"

	"github.com/ditatompel/xmr-nodes/internal/database"
)

type MoneroRepository interface {
	Add(protocol string, host string, port uint) error
}

type MoneroRepo struct {
	db *database.DB
}

func NewMoneroRepo(db *database.DB) MoneroRepository {
	return &MoneroRepo{db}
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
