package monero

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

// RuckniumNodeData represents a single remote node data from Rucknium's MRL ban list
type RuckniumNodeData struct {
	ID                uint   `json:"id,omitempty" db:"id"`
	Date              string `json:"date" db:"scan_date"`
	ConnectedNodeIP   string `json:"connected_node_ip" db:"connected_node_ip"`
	IsSpyNode         int    `json:"is_spy_node" db:"is_spy_node"`                   // 0 = no, 1 = yes, 2 = not applied
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
		// insert to database with ON DUPLICATE KEY UPDATE

		_, err := r.db.Exec(`INSERT INTO tbl_rucknium_scan (
			scan_date,
			connected_node_ip,
			is_spy_node,
			mrl_ban_list_enabled,
			dns_ban_list_enabled
		) VALUES (
			?, ?, ?, ?, ?
		)
		ON DUPLICATE KEY UPDATE
			is_spy_node = ?,
			mrl_ban_list_enabled = ?,
			dns_ban_list_enabled = ?`,
			node.Date,
			node.ConnectedNodeIP,
			node.IsSpyNode,
			node.MRLBanListEnabled,
			node.DNSBanListEnabled,
			node.IsSpyNode,
			node.MRLBanListEnabled,
			node.DNSBanListEnabled)
		if err != nil {
			slog.Error(fmt.Sprintf("[MRL] Failed to insert or update Rucknium's node list: %s", err))
		}
	}

	return nil
}
