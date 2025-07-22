package monero

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"strings"
)

// RuckniumNodeData represents a single remote node data from Rucknium's MRL ban list
type RuckniumNodeData struct {
	ID                uint   `json:"id,omitempty" db:"id"`
	Date              string `json:"date" db:"scan_date"`
	ConnectedNodeIP   string `json:"connected_node_ip" db:"connected_node_ip"`
	IsSpyNode         int    `json:"is_spy_node" db:"is_spy_node"` // 0 = no, 1 = yes, 2 = not applied
	RPCDomain         string `json:"rpc_domain" db:"rpc_domain"`
	MRLBanListEnabled int    `json:"mrl_ban_list_enabled" db:"mrl_ban_list_enabled"` // 0 = no, 1 = yes, 2 = not applied
	DNSBanListEnabled int    `json:"dns_ban_list_enabled" db:"dns_ban_list_enabled"` // 0 = no, 1 = yes, 2 = not applied
}

// Get Individual node info from moneronet.info
func (r *moneroRepo) FetchRuckniumNodeData() error {
	req, err := http.NewRequest(http.MethodGet, "https://api.moneronet.info/individual_node_data?date=latest", nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		slog.Error("[MRL] Failed to fetch Rucknium's API")
		return errors.New("Failed to fetch Rucknium's API")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var nodes []RuckniumNodeData
	if err := json.Unmarshal(body, &nodes); err != nil {
		return err
	}

	for _, node := range nodes {
		_, err := r.db.Exec(`INSERT INTO tbl_rucknium_scan (
			scan_date,
			connected_node_ip,
			is_spy_node,
			rpc_domain,
			mrl_ban_list_enabled,
			dns_ban_list_enabled
		) VALUES (
			?, ?, ?, ?, ?, ?
		)
		ON DUPLICATE KEY UPDATE
			is_spy_node = ?,
			rpc_domain = ?,
			mrl_ban_list_enabled = ?,
			dns_ban_list_enabled = ?`,
			node.Date,
			node.ConnectedNodeIP,
			node.IsSpyNode,
			node.RPCDomain,
			node.MRLBanListEnabled,
			node.DNSBanListEnabled,
			node.IsSpyNode,
			node.RPCDomain,
			node.MRLBanListEnabled,
			node.DNSBanListEnabled)
		if err != nil {
			slog.Error(fmt.Sprintf("[MRL] Failed to insert or update Rucknium's node list: %s", err))
		}
	}

	return nil
}

// Get Rucknium node data from local database
func (r *moneroRepo) GetRuckniumData(date string) ([]RuckniumNodeData, error) {
	var nodes []RuckniumNodeData

	query := `
	SELECT
		connected_node_ip,
		is_spy_node,
		rpc_domain,
		mrl_ban_list_enabled,
		dns_ban_list_enabled
	FROM
		tbl_rucknium_scan
	WHERE
		scan_date = ?`
	err := r.db.Select(&nodes, query, date)

	return nodes, err
}

// get latest date from Rucknium data
func (r *moneroRepo) GetLatestRuckniumDate() (string, error) {
	var latestDate string
	query := `SELECT scan_date FROM tbl_rucknium_scan ORDER BY scan_date DESC LIMIT 1`
	err := r.db.Get(&latestDate, query)
	return latestDate, err
}

// Check Rucknium ban list
func (r *moneroRepo) CheckMRLBan() error {
	latestDate, err := r.GetLatestRuckniumDate()
	if err != nil {
		return err
	}

	mrlData, err := r.GetRuckniumData(latestDate)
	if err != nil {
		return err
	}

	if len(mrlData) == 0 {
		return errors.New("No Rucknium data found")
	}

	var nodes []Node

	query := `
		SELECT id, ip_addresses FROM tbl_node
		WHERE nettype = ?
			AND is_archived = ?
			AND is_tor = ?
			AND is_i2p = ?
			AND is_i2p = ?
			AND ipv6_only = ?`
	// For now, Monero Network scan only checks mainnet
	err = r.db.Select(&nodes, query, "mainnet", 0, 0, 0, 0, 0)
	if err != nil {
		return err
	}

	// maps for each property
	spyIPMap := make(map[string]bool)
	mrlBanMap := make(map[string]bool)
	dnsBanMap := make(map[string]bool)

	for _, rn := range mrlData {
		ip := rn.ConnectedNodeIP
		if rn.IsSpyNode == 1 {
			spyIPMap[ip] = true
		}
		if rn.MRLBanListEnabled == 1 {
			mrlBanMap[ip] = true
		}
		if rn.DNSBanListEnabled == 1 {
			dnsBanMap[ip] = true
		}
	}

	for i, node := range nodes {
		ipList := strings.Split(node.IPAddresses, ",")
		nodes[i].IsSpyNode = 0
		nodes[i].MRLBanListEnabled = 1
		nodes[i].DNSBanListEnabled = 1

		for _, ip := range ipList {
			c := net.ParseIP(ip)
			if c != nil && c.To4() == nil {
				slog.Debug(fmt.Sprintf("[MRL] Skipping IPv6 address: %s", ip))
				break
			}

			if spyIPMap[ip] {
				nodes[i].IsSpyNode = 1
			}
			if !mrlBanMap[ip] {
				nodes[i].MRLBanListEnabled = 0
			}
			if !dnsBanMap[ip] {
				nodes[i].DNSBanListEnabled = 0
			}

			// if nodes[i].MRLBanListEnabled == 0 &&
			// 	nodes[i].DNSBanListEnabled == 0 {
			// 	break
			// }
		}

		// Update node MRL columns info in the database
		_, err := r.db.Exec(`UPDATE tbl_node SET
			is_spy_node = ?,
			mrl_ban_list_enabled = ?,
			dns_ban_list_enabled = ?
			WHERE id = ?`,
			nodes[i].IsSpyNode,
			nodes[i].MRLBanListEnabled,
			nodes[i].DNSBanListEnabled,
			nodes[i].ID)
		if err != nil {
			return err
		}

	}

	return nil
}
